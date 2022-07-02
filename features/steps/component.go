package steps

import (
	"context"
	"errors"
	"net/http"
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

const (
	mongoVersion = "4.4.8"
	databaseName = "filters"
)

var (
	BuildTime = "1625046891"
	GitCommit = "7434fe334d9f51b7239f978094ea29d10ac33b16"
	Version   = ""
)

type Component struct {
	componenttest.ErrorFeature
	ApiFeature   *componenttest.APIFeature
	AuthFeature  *componenttest.AuthorizationFeature
	MongoFeature *componenttest.MongoFeature

	DatasetAPI       *httpfake.HTTPFake
	CantabularClient *mock.CantabularClient

	HTTPServer *http.Server

	store    service.Datastore
	consumer kafka.IConsumerGroup
	g        service.Generator

	svc *service.Service

	waitEventTimeout time.Duration
}

func NewComponent(t *testing.T) *Component {
	component := &Component{
		ErrorFeature: componenttest.ErrorFeature{TB: t},
		AuthFeature:  componenttest.NewAuthorizationFeature(),
		MongoFeature: componenttest.NewMongoFeature(componenttest.MongoOptions{
			MongoVersion: mongoVersion,
			DatabaseName: databaseName,
		}),

		HTTPServer: &http.Server{},
		DatasetAPI: httpfake.New(httpfake.WithTesting(t)),
		CantabularClient: &mock.CantabularClient{
			OptionsHappy: true,
		},
		g: &mock.Generator{
			URLHost: "http://mockhost:9999",
		},
		waitEventTimeout: WaitEventTimeout,
	}
	component.ApiFeature = componenttest.NewAPIFeature(component.Router)

	cfg, err := config.Get()
	if err != nil {
		t.Fatalf("failed to get config: %s", err)
	}
	cfg.ZebedeeURL = component.AuthFeature.FakeAuthService.ResolveURL("")
	cfg.DatasetAPIURL = component.DatasetAPI.ResolveURL("")
	cfg.Mongo.ClusterEndpoint = component.MongoFeature.Server.URI()
	cfg.Mongo.Database = utils.RandomDatabase()
	component.setWorkingMongo(cfg)

	log.Info(context.Background(), "config used by component tests", log.Data{"cfg": cfg})

	// Create service and initialise it
	component.setInitialiserMock()
	component.svc = service.New()
	if err := component.svc.Init(context.Background(), cfg, BuildTime, GitCommit, Version); err != nil {
		component.ErrorFeature.Fatalf("failed to initialise service: %s", err)
	}

	return component
}

// Router returns the component's server router
func (c *Component) Router() (http.Handler, error) {
	return c.HTTPServer.Handler, nil
}

// Reset re-initialises the service under test and the api mocks.
func (c *Component) Reset() {
	var err error

	c.svc.Reset()

	c.ApiFeature.Reset()
	c.AuthFeature.Reset()
	c.DatasetAPI.Reset()
	c.CantabularClient.Reset()

	err = c.MongoFeature.Reset()
	if err != nil {
		c.ErrorFeature.Fatalf("failed to reset mongo: %s", err)
	}

	c.svc.Cfg.Mongo.Database = utils.RandomDatabase()
	c.setWorkingMongo(c.svc.Cfg)
}

// Close kills the application under test, and then it shuts down the testing producer.
func (c *Component) Close() {
	c.AuthFeature.Close()
	c.DatasetAPI.Close()

	err := c.MongoFeature.Close()
	if err != nil {
		c.ErrorFeature.Errorf("error closing Mongo: %s", err)
	}
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
}

func (c *Component) setWorkingMongo(cfg *config.Config) {
	var err error

	c.store, err = mongodb.NewClient(context.Background(), c.g, mongodb.Config{
		MongoDriverConfig:       cfg.Mongo,
		FilterFlexAPIURL:        cfg.BindAddr,
		FiltersCollection:       cfg.FiltersCollection,
		FilterOutputsCollection: cfg.FilterOutputsCollection,
	})
	if err != nil {
		c.ErrorFeature.Fatalf("failed to get a working mongo: %s", err)
	}
}

//keep adding new handler functions for which the mongo needs to fail
func (c *Component) setFailingMongo() {
	c.store = &servicemock.DatastoreMock{
		UpdateFilterOutputFunc: func(_ context.Context, _ *model.FilterOutput) error {
			return errors.New("failed to upsert filter")
		},
		GetFilterOutputFunc: func(_ context.Context, s string) (*model.FilterOutput, error) {
			return nil, errors.New("mongo client has failed")
		},
		AddFilterOutputEventFunc: func(_ context.Context, _ string, _ *model.Event) error {
			return errors.New("failed to add event")
		},
		GetFilterDimensionOptionsFunc: func(contextMoqParam context.Context, s1 string, s2 string, n1 int, n2 int) ([]string, int, string, error) {
			return nil, 0, "", errors.New("error that should not be returned to user")
		},
	}
}
