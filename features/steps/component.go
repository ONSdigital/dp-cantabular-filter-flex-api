package steps

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/config"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/features/mock"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/service"
	componenttest "github.com/ONSdigital/dp-component-test"
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
	AuthFeature       *componenttest.AuthorizationFeature
	APIFeature        *componenttest.APIFeature
	MongoFeature      *MongoFeature
	DatasetFeature    *DatasetFeature
	CantabularFeature *CantabularFeature

	PopulationFeatures *PopulationFeature
	MetadataFeature    *MetadataFeature

	svc *service.Service
}

func NewComponent(t *testing.T) *Component {
	g := &mock.Generator{URLHost: "http://mockhost:9999"}
	cfg, err := config.Get()
	if err != nil {
		t.Fatalf("failed to get config: %s", err)
	}

	component := &Component{
		ErrorFeature:       componenttest.ErrorFeature{TB: t},
		AuthFeature:        componenttest.NewAuthorizationFeature(),
		DatasetFeature:     NewDatasetFeature(t, cfg),
		CantabularFeature:  NewCantabularFeature(),
		PopulationFeatures: NewPopulationFeature(t, cfg),
	}
	component.MongoFeature = NewMongoFeature(component.ErrorFeature, g, cfg)
	component.APIFeature = componenttest.NewAPIFeature(component.Router)

	cfg.ZebedeeURL = component.AuthFeature.FakeAuthService.ResolveURL("")

	component.setInitialiserMock(g)
	component.svc = service.New()
	component.svc.Cfg = cfg

	return component
}

// Router initialises the service, returning the service's (server) router for tests
// This delayed initialisation is needed to ensure that any changes to the router (or the service in general)
// as a result of test setup, are picked up
func (c *Component) Router() (http.Handler, error) {
	if err := c.svc.Init(context.Background(), c.svc.Cfg, BuildTime, GitCommit, Version); err != nil {
		return nil, fmt.Errorf("failed to initialise service: %s", err)
	}

	return c.svc.Api.Router, nil
}

// Reset re-initialises the service under test and the api dependencies
func (c *Component) Reset() {
	c.AuthFeature.Reset()
	c.APIFeature.Reset()
	c.DatasetFeature.Reset()
	c.CantabularFeature.Reset()
	c.MongoFeature.Reset()
}

func (c *Component) Close() {
	c.AuthFeature.Close()
	c.DatasetFeature.Close()
	c.MongoFeature.Close()
}

func (c *Component) setInitialiserMock(g service.Generator) {
	c.CantabularFeature.setInitialiserMock()
	c.DatasetFeature.setInitialiserMock()
	c.MongoFeature.setInitialiserMock()

	service.GetHTTPServer = func(bindAddr string, router http.Handler) service.HTTPServer {
		return &http.Server{Addr: bindAddr, Handler: router}
	}

	service.GetGenerator = func() service.Generator {
		return g
	}
}
