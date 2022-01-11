package steps

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/config"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/service"
	componenttest "github.com/ONSdigital/dp-component-test"
	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/log.go/v2/log"
)

const (
	DrainTopicTimeout     = 1 * time.Second // maximum time to wait for a topic to be drained
	DrainTopicMaxMessages = 1000            // maximum number of messages that will be drained from a topic
	MinioCheckRetries     = 3               // maximum number of retires to validate that a file is present in minio
)

var (
	BuildTime string = "1625046891"
	GitCommit string = "7434fe334d9f51b7239f978094ea29d10ac33b16"
	Version   string = ""
)

type Component struct {
	componenttest.ErrorFeature
	producer         kafka.IProducer
	errorChan        chan error
	svc              *service.Service
	cfg              *config.Config
	wg               *sync.WaitGroup
	signals          chan os.Signal
	waitEventTimeout time.Duration
	testETag         string
	ctx              context.Context
	apiFeature       *componenttest.APIFeature
}

func NewComponent() *Component {
	return &Component{
		errorChan: make(chan error),
		wg:        &sync.WaitGroup{},
		testETag:  "13c7791bafdbaaf5e6660754feb1a58cd6aaa892",
		ctx:       context.Background(),
	}
}

// initService initialises the server, the mocks and waits for the dependencies to be ready
func (c *Component) initService(ctx context.Context) error {
	// register interrupt signals
	c.signals = make(chan os.Signal, 1)
	signal.Notify(c.signals, os.Interrupt, syscall.SIGTERM)

	// Read config
	cfg, err := config.Get()
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	log.Info(ctx, "config used by component tests", log.Data{"cfg": cfg})

	// producer for triggering test events
	if c.producer, err = kafka.NewProducer(
		ctx,
		&kafka.ProducerConfig{
			BrokerAddrs:       cfg.KafkaConfig.Addr,
			Topic:             cfg.KafkaConfig.ExportStartTopic,
			MinBrokersHealthy: &cfg.KafkaConfig.ProducerMinBrokersHealthy,
			KafkaVersion:      &cfg.KafkaConfig.Version,
			MaxMessageBytes:   &cfg.KafkaConfig.MaxBytes,
		},
	); err != nil {
		return fmt.Errorf("error creating kafka producer: %w", err)
	}

	// Create service and initialise it
	c.svc = service.New()
	if err = c.svc.Init(ctx, cfg, BuildTime, GitCommit, Version); err != nil {
		return fmt.Errorf("unexpected service Init error in NewComponent: %w", err)
	}

	c.cfg = cfg

	// wait for producer to be initialised
	<-c.producer.Channels().Initialised
	log.Info(ctx, "component-test kafka producer initialised")

	return nil
}

// startService starts the service under test and blocks until an error or an os interrupt is received.
// Then it closes the service (graceful shutdown)
func (c *Component) startService(ctx context.Context) {
	defer c.wg.Done()
	c.svc.Start(ctx, c.errorChan)

	select {
	case err := <-c.errorChan:
		err = fmt.Errorf("service error received: %w", err)
		c.svc.Close(ctx)
		panic(fmt.Errorf("unexpected error received from errorChan: %w", err))
	case sig := <-c.signals:
		log.Info(ctx, "os signal received", log.Data{"signal": sig})
	}
	if err := c.svc.Close(ctx); err != nil {
		panic(fmt.Errorf("unexpected error during service graceful shutdown: %w", err))
	}
}

// Close kills the application under test, and then it shuts down the testing producer.
func (c *Component) Close() {
	// kill application
	c.signals <- os.Interrupt

	// wait for graceful shutdown to finish (or timeout)
	c.wg.Wait()

	// close producer
	if err := c.producer.Close(c.ctx); err != nil {
		log.Error(c.ctx, "error closing kafka producer", err)
	}
}

// Reset re-initialises the service under test and the api mocks.
// Note that the service under test should not be started yet
// to prevent race conditions if it tries to call un-initialised dependencies (steps)
func (c *Component) Reset() error {
	if err := c.initService(c.ctx); err != nil {
		return fmt.Errorf("failed to initialise service: %w", err)
	}

	return nil
}
