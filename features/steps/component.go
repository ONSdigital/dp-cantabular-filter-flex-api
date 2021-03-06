package steps

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/maxcnunes/httpfake"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/config"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/datastore/mongodb"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/features/mock"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/service"
	servicemock "github.com/ONSdigital/dp-cantabular-filter-flex-api/service/mock"
	componenttest "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-component-test/utils"
	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/log.go/v2/log"
)

const (
	ComponentTestGroup    = "test-consumer-group"
	DrainTopicTimeout     = 10 * time.Second // maximum time to wait for a topic to be drained
	DrainTopicMaxMessages = 1000             // maximum number of messages that will be drained from a topic
	MinioCheckRetries     = 3                // maximum number of retires to validate that a file is present in minio
	WaitEventTimeout      = 10 * time.Second // maximum time that the component test consumer will wait for a
)

var (
	BuildTime string = "1625046891"
	GitCommit string = "7434fe334d9f51b7239f978094ea29d10ac33b16"
	Version   string = ""
)

type Component struct {
	componenttest.ErrorFeature
	ApiFeature        *componenttest.APIFeature
	errorChan         chan error
	DatasetAPI        *httpfake.HTTPFake
	CantabularClient  *mock.CantabularClient
	svc               *service.Service
	cfg               *config.Config
	wg                *sync.WaitGroup
	signals           chan os.Signal
	ctx               context.Context
	HTTPServer        *http.Server
	store             service.Datastore
	g                 service.Generator
	shutdownInitiated bool
	consumer          kafka.IConsumerGroup
	waitEventTimeout  time.Duration
}

func NewComponent(t *testing.T, zebedeeURL, mongoAddr string) (*Component, error) {
	ctx := context.Background()
	cfg, err := config.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	cfg.ZebedeeURL = zebedeeURL
	cfg.Mongo.ClusterEndpoint = mongoAddr
	cfg.Mongo.Database = utils.RandomDatabase()

	g := &mock.Generator{
		URLHost: "http://mockhost:9999",
	}

	mongoClient, err := GetWorkingMongo(ctx, cfg, g)
	if err != nil {
		return nil, fmt.Errorf("failed to create new mongo mongoClient: %w", err)
	}

	return &Component{
		errorChan:  make(chan error),
		wg:         &sync.WaitGroup{},
		ctx:        ctx,
		HTTPServer: &http.Server{},
		cfg:        cfg,
		DatasetAPI: httpfake.New(httpfake.WithTesting(t)),
		CantabularClient: &mock.CantabularClient{
			OptionsHappy: true,
		},
		store:            mongoClient,
		g:                g,
		waitEventTimeout: WaitEventTimeout,
	}, nil
}

// Init initialises the server, the mocks and waits for the dependencies to be ready
func (c *Component) Init() (http.Handler, error) {
	c.signals = make(chan os.Signal, 1)
	signal.Notify(c.signals, os.Interrupt, syscall.SIGTERM)

	log.Info(c.ctx, "config used by component tests", log.Data{"cfg": c.cfg})

	c.cfg.DatasetAPIURL = c.DatasetAPI.ResolveURL("")

	// Create service and initialise it
	c.svc = service.New()
	if err := c.svc.Init(c.ctx, c.cfg, BuildTime, GitCommit, Version); err != nil {
		return nil, fmt.Errorf("failed to initialise service: %w", err)
	}

	return c.HTTPServer.Handler, nil
}

func (c *Component) setInitialiserMock() {
	service.GetHTTPServer = func(bindAddr string, router http.Handler) service.HTTPServer {
		c.HTTPServer.Addr = bindAddr
		c.HTTPServer.Handler = router
		return c.HTTPServer
	}

	service.GetGenerator = func() service.Generator {
		return &mock.Generator{
			URLHost: "http://mockhost:9999",
		}
	}

	service.GetCantabularClient = func(_ *config.Config) service.CantabularClient {
		return c.CantabularClient
	}

	service.GetMongoDB = func(ctx context.Context, cfg *config.Config, g service.Generator) (service.Datastore, error) {
		return c.store, nil
	}

	c.cfg.Mongo.Database = utils.RandomDatabase()
}

// startService starts the service under test and blocks until an error or an os interrupt is received.
// Then it closes the service (graceful shutdown)
func (c *Component) startService(ctx context.Context) {
	defer c.wg.Done()
	c.svc.Start(ctx, c.errorChan)

	wg := sync.WaitGroup{}
	if err := c.drainTopic(c.ctx, c.cfg.KafkaConfig.ExportStartTopic, ComponentTestGroup, &wg); err != nil {
		log.Error(c.ctx, "error draining topic", err)
	}

	select {
	case err := <-c.errorChan:
		err = fmt.Errorf("service error received: %w", err)
		defer func() {
			if err := c.svc.Close(ctx); err != nil {
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

	if !c.shutdownInitiated {
		c.shutdownInitiated = true
		c.signals <- os.Interrupt

		// wait for graceful shutdown to finish (or timeout)
		// TODO we should fix the timeout issue and then uncomment the following line.
		c.wg.Wait()
	}

}

// Reset re-initialises the service under test and the api mocks.
// Note that the service under test should not be started yet
// to prevent race conditions if it tries to call un-initialised dependencies (steps)
func (c *Component) Reset() error {
	c.setInitialiserMock()
	c.DatasetAPI.Reset()
	c.CantabularClient.Reset()

	if _, err := c.Init(); err != nil {
		return fmt.Errorf("failed to initialise component: %w", err)
	}

	return nil
}

func GetWorkingMongo(ctx context.Context, cfg *config.Config, g service.Generator) (service.Datastore, error) {
	mongoClient, err := mongodb.NewClient(ctx, g, mongodb.Config{
		MongoDriverConfig:       cfg.Mongo,
		FilterFlexAPIURL:        cfg.BindAddr,
		FiltersCollection:       cfg.FiltersCollection,
		FilterOutputsCollection: cfg.FilterOutputsCollection,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create new mongo mongoClient: %w", err)
	}
	return mongoClient, nil
}

//keep adding new handler functions for which the mongo needs to fail
func GetFailingMongo(ctx context.Context, cfg *config.Config, g service.Generator) (service.Datastore, error) {
	mongoClient := servicemock.DatastoreMock{
		UpdateFilterOutputFunc: func(_ context.Context, _ *model.FilterOutput) error {
			return errors.New("failed to upsert filter")
		},
		GetFilterOutputFunc: func(_ context.Context, _ string) (*model.FilterOutput, error) {
			return nil, errors.New("mongo client has failed")
		},
		AddFilterOutputEventFunc: func(_ context.Context, _ string, _ *model.Event) error {
			return errors.New("failed to add event")
		},
		GetFilterDimensionOptionsFunc: func(_ context.Context, _, _ string, _, _ int) ([]string, int, string, error) {
			return nil, 0, "", errors.New("error that should not be returned to user.")
		},
		UpdateFilterDimensionFunc: func(_ context.Context, _, _ string, _ model.Dimension, _ string) (string, error) {
			return "", errors.New("failed to update filter dimension")
		},
		RemoveFilterDimensionOptionFunc: func(_ context.Context, _, _, _, _ string) (string, error) {
			return "", errors.New("failed to remove filter dimension option")
		},
		GetFilterFunc: func(_ context.Context, _ string) (*model.Filter, error) {
			return nil, errors.New("failed to get filter")
		},
	}
	return &mongoClient, nil
}
