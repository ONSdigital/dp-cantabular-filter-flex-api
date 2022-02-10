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
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/features/mock"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/service"

	componenttest "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-component-test/utils"
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
	store      service.Datastore
}

func NewComponent(zebedeeURL, mongoAddr string) (*Component, error) {
	cfg, err := config.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	cfg.ZebedeeURL = zebedeeURL
	cfg.Mongo.ClusterEndpoint = mongoAddr

	return &Component{
		errorChan:  make(chan error),
		wg:         &sync.WaitGroup{},
		ctx:        context.Background(),
		HTTPServer: &http.Server{},
		cfg:        cfg,
	}, nil
}

// Init initialises the server, the mocks and waits for the dependencies to be ready
func (c *Component) Init() (http.Handler, error) {
	log.Info(c.ctx, "config used by component tests", log.Data{"cfg": c.cfg})

	c.signals = make(chan os.Signal, 1)
	signal.Notify(c.signals, os.Interrupt, syscall.SIGTERM)

	/*
		// producer for triggering test events
		if c.producer, err = kafka.NewProducer(
			c.ctx,
			&kafka.ProducerConfig{
				BrokerAddrs:       c.cfg.Kafka.Addr,
				Topic:             c.cfg.Kafka.ExportStartTopic,
				MinBrokersHealthy: &c.cfg.Kafka.ProducerMinBrokersHealthy,
				KafkaVersion:      &c.cfg.Kafka.Version,
				MaxMessageBytes:   &c.cfg.Kafka.MaxBytes,
			},
		); err != nil {
			return fmt.Errorf("error creating kafka producer: %w", err)
		}
	*/
	// Create service and initialise it
	c.svc = service.New()
	if err := c.svc.Init(c.ctx, c.cfg, BuildTime, GitCommit, Version); err != nil {
		return nil, fmt.Errorf("failed to initialise service: %w", err)
	}

	// wait for producer to be initialised
	// <-c.producer.Channels().Initialised
	// log.Info(c.ctx, "component-test kafka producer initialised")

	return c.HTTPServer.Handler, nil
}

func (c *Component) setInitialiserMock() {
	service.GetHTTPServer = func(bindAddr string, router http.Handler) service.HTTPServer {
		c.HTTPServer.Addr = bindAddr
		c.HTTPServer.Handler = router
		return c.HTTPServer
	}

	service.GetGenerator = func() service.Generator {
		return &mock.Generator{}
	}

	c.cfg.Mongo.Database = utils.RandomDatabase()
}

// startService starts the service under test and blocks until an error or an os interrupt is received.
// Then it closes the service (graceful shutdown)
func (c *Component) startService(ctx context.Context) {
	defer c.wg.Done()
	c.svc.Start(ctx, c.errorChan)

	select {
	case err := <-c.errorChan:
		err = fmt.Errorf("service error received: %w", err)
		defer func(){
			if err := c.svc.Close(ctx); err != nil{
				log.Error(ctx, "failed to shutdown service gracefully: %s", err)
			}
		}()
		panic(fmt.Errorf("unexpected error received from errorChan: %w", err))
	case sig := <-c.signals:
		log.Info(ctx, "os signal received", log.Data{"signal": sig})
	}
	if err := c.svc.Close(ctx); err != nil {
		panic(fmt.Errorf("failed to shutdiwn gracefully: %w", err))
	}
}

// Close kills the application under test, and then it shuts down the testing producer.
func (c *Component) Close() {
	// kill application
	c.signals <- os.Interrupt

	// wait for graceful shutdown to finish (or timeout)
	// TODO we should fix the timeout issue and then uncomment the following line.
	c.wg.Wait()

	// close producer
	// if err := c.producer.Close(c.ctx); err != nil {
	//     log.Error(c.ctx, "error closing kafka producer", err)
	// }
}

// Reset re-initialises the service under test and the api mocks.
// Note that the service under test should not be started yet
// to prevent race conditions if it tries to call un-initialised dependencies (steps)
func (c *Component) Reset() error {
	c.setInitialiserMock()

	if _, err := c.Init(); err != nil {
		return fmt.Errorf("failed to initialise component: %w", err)
	}

	return nil
}
