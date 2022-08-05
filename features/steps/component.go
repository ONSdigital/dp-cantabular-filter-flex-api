package steps

import (
	"context"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/maxcnunes/httpfake"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/config"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/features/mock"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/service"
	componenttest "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/log.go/v2/log"
)

const (
	ComponentTestGroup    = "test-consumer-group"
	DrainTopicTimeout     = 10 * time.Second // maximum time to wait for a topic to be drained
	DrainTopicMaxMessages = 1000             // maximum number of messages that will be drained from a topic
	WaitEventTimeout      = 10 * time.Second // maximum time that the component test consumer will wait for a
)

var (
	BuildTime = strconv.Itoa(time.Now().Nanosecond())
	GitCommit = "component test commit"
	Version   = "component test version"
)

type Component struct {
	componenttest.ErrorFeature
	AuthServiceInjector *componenttest.AuthorizationFeature
	APIInjector         *componenttest.APIFeature
	MongoInjector       *MongoFeature
	DatasetInjector     *DatasetFeature
	CantabularInjector  *CantabularFeature

	svc *service.Service
}

func NewComponent(t *testing.T) *Component {
	component := &Component{
		ErrorFeature:        componenttest.ErrorFeature{TB: t},
		AuthServiceInjector: componenttest.NewAuthorizationFeature(),
		DatasetInjector:     &DatasetFeature{mockDatasetServer: httpfake.New(httpfake.WithTesting(t))},
		CantabularInjector:  &CantabularFeature{CantabularClient: &mock.CantabularClient{OptionsHappy: true}},
	}
	component.APIInjector = componenttest.NewAPIFeature(component.Router)

	cfg, err := config.Get()
	if err != nil {
		component.ErrorFeature.Fatalf("failed to get config: %s", err)
	}

	g := &mock.Generator{URLHost: "http://mockhost:9999"}
	component.MongoInjector = NewMongoFeature(component.ErrorFeature, g, cfg)

	cfg.ZebedeeURL = component.AuthServiceInjector.FakeAuthService.ResolveURL("")
	cfg.DatasetAPIURL = component.DatasetInjector.mockDatasetServer.ResolveURL("")
	log.Info(context.Background(), "config used by component tests", log.Data{"cfg": cfg})

	setInitialiserMock(component.CantabularInjector, g)
	component.svc = service.New()
	component.svc.Cfg = cfg

	return component
}

// Router initialises the service, returning the service's (server) router for tests
// This delayed initialisation is needed to ensure that any changes to the router (or the service in general)
// as a result of test setup, are picked up
func (c *Component) Router() (http.Handler, error) {
	if err := c.svc.Init(context.Background(), c.svc.Cfg, BuildTime, GitCommit, Version); err != nil {
		c.ErrorFeature.Fatalf("failed to initialise service: %s", err)
	}

	return c.svc.Api.Router, nil
}

// Reset re-initialises the service under test and the api mocks.
func (c *Component) Reset() {
	c.AuthServiceInjector.Reset()
	c.APIInjector.Reset()
	c.DatasetInjector.Reset()
	c.CantabularInjector.Reset()
	c.MongoInjector.Reset()
}

func (c *Component) Close() {
	c.AuthServiceInjector.Close()
	c.DatasetInjector.Close()
	c.MongoInjector.Close()
}

func setInitialiserMock(c *CantabularFeature, g service.Generator) {
	service.GetHTTPServer = func(bindAddr string, router http.Handler) service.HTTPServer {
		return &http.Server{Addr: bindAddr, Handler: router}
	}

	service.GetGenerator = func() service.Generator {
		return g
	}

	service.GetCantabularClient = func(_ *config.Config) service.CantabularClient {
		return c
	}
}
