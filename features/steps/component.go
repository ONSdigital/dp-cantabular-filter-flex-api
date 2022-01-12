package steps

import (
	"context"
	"fmt"
	"net/http"
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
	producer   kafka.IProducer
	errorChan  chan error
	svc        *service.Service
	cfg        *config.Config
	wg         *sync.WaitGroup
	signals    chan os.Signal
	ctx        context.Context
	HTTPServer *http.Server
}

func NewComponent() *Component {
	return &Component{
		errorChan:  make(chan error),
		wg:         &sync.WaitGroup{},
		ctx:        context.Background(),
		HTTPServer: &http.Server{},
	}
}

// initService initialises the server, the mocks and waits for the dependencies to be ready
func (c *Component) InitService() (http.Handler, error) {
	// register interrupt signals
	c.signals = make(chan os.Signal, 1)
	signal.Notify(c.signals, os.Interrupt, syscall.SIGTERM)

	// Read config
	var err error
	c.cfg, err = config.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}
	log.Info(c.ctx, "config used by component tests", log.Data{"cfg": c.cfg})

	// producer for triggering test events
	if c.producer, err = kafka.NewProducer(
		c.ctx,
		&kafka.ProducerConfig{
			BrokerAddrs:       c.cfg.KafkaConfig.Addr,
			Topic:             c.cfg.KafkaConfig.ExportStartTopic,
			MinBrokersHealthy: &c.cfg.KafkaConfig.ProducerMinBrokersHealthy,
			KafkaVersion:      &c.cfg.KafkaConfig.Version,
			MaxMessageBytes:   &c.cfg.KafkaConfig.MaxBytes,
		},
	); err != nil {
		return nil, fmt.Errorf("error creating kafka producer: %w", err)
	}

	// Create service and initialise it
	c.svc = service.New()
	if err = c.svc.Init(c.ctx, c.cfg, BuildTime, GitCommit, Version); err != nil {
		return nil, fmt.Errorf("unexpected service Init error in NewComponent: %w", err)
	}

	// wait for producer to be initialised
	<-c.producer.Channels().Initialised
	log.Info(c.ctx, "component-test kafka producer initialised")

	return c.HTTPServer.Handler, nil
}

func (c *Component) setInitialiserMock() {
	service.GetHTTPServer = func(bindAddr string, router http.Handler) service.HTTPServer {
		c.HTTPServer.Addr = bindAddr
		c.HTTPServer.Handler = router
		return c.HTTPServer
	}
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
	// TODO we should fix the timeout issue and then uncomment the following line.
	// c.wg.Wait()

	// close producer
	if err := c.producer.Close(c.ctx); err != nil {
		log.Error(c.ctx, "error closing kafka producer", err)
	}
}

// Reset re-initialises the service under test and the api mocks.
// Note that the service under test should not be started yet
// to prevent race conditions if it tries to call un-initialised dependencies (steps)
func (c *Component) Reset() error {
	if _, err := c.InitService(); err != nil {
		return fmt.Errorf("failed to initialise service: %w", err)
	}

	c.setInitialiserMock()

	return nil
}
