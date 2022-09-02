package steps

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular/gql"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/config"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/features/mock"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/service"
	dphttp "github.com/ONSdigital/dp-net/v2/http"

	"github.com/cucumber/godog"
)

type CantabularFeature struct {
	cantabularMock   *mock.CantabularClient
	cantabularServer *mock.CantabularServer

	cfg *config.Config
}

func NewCantabularFeature(t *testing.T, cfg *config.Config) *CantabularFeature {
	return &CantabularFeature{
		cantabularMock:   &mock.CantabularClient{OptionsHappy: true},
		cantabularServer: mock.NewCantabularServer(t),
		cfg:              cfg,
	}
}

func (cf *CantabularFeature) Reset() {
	cf.cantabularMock.Reset()
	cf.cantabularServer.Reset()

	cf.setMockedInterface()
}

func (cf *CantabularFeature) RegisterSteps(ctx *godog.ScenarioContext) {
	ctx.Step(
		`^the Cantabular service is a mocked extended Cantabular server`,
		cf.useAMockedExtCantabularServer,
	)
	ctx.Step(
		`^the Cantabular service is a mocked interface$`,
		cf.setMockedInterface,
	)

	ctx.Step(
		`^Cantabular returns these dimensions for the dataset "([^"]*)" and search term "([^"]*)":$`,
		cf.cantabularSearchReturnsTheseDimensions,
	)
	ctx.Step(
		`^Cantabular returns these geography dimensions for the given request:$`,
		cf.cantabularReturnsTheseGeographyDimensionsForTheGivenRequest,
	)
	ctx.Step(
		`^Cantabular returns this static dataset for the given request:$`,
		cf.cantabularReturnsThisStaticDatasetForTheGivenRequest,
	)

	ctx.Step(
		`^Cantabular responds with an error$`,
		cf.cantabularRespondsWithAnError,
	)
}

func (cf *CantabularFeature) useAMockedExtCantabularServer() error {
	cf.cfg.CantabularExtURL = cf.cantabularServer.ResolveURL("")

	cf.setMockedServer()

	return nil
}

func (cf *CantabularFeature) cantabularSearchReturnsTheseDimensions(datasetID, dimension string, docs *godog.DocString) error {
	var response cantabular.GetDimensionsResponse
	if err := json.Unmarshal([]byte(docs.Content), &response); err != nil {
		return fmt.Errorf("unable to unmarshal cantabular search response: %w", err)
	}

	cf.cantabularMock.SearchDimensionsFunc = func(ctx context.Context, req cantabular.SearchDimensionsRequest) (*cantabular.GetDimensionsResponse, error) {
		if req.Dataset == datasetID && req.Text == dimension {
			return &response, nil
		}

		return &cantabular.GetDimensionsResponse{
			Dataset: gql.Dataset{
				Variables: gql.Variables{
					Search: gql.Search{
						Edges: []gql.Edge{},
					},
				},
			},
		}, nil
	}

	return nil
}

func (cf *CantabularFeature) cantabularReturnsTheseGeographyDimensionsForTheGivenRequest(docs *godog.DocString) error {
	request, response, found := strings.Cut(docs.Content, "response:")
	if !found {
		return errors.New("CantabularFeature::cantabularReturnsTheseGeographyDimensionsForTheGivenRequest - request and response were not found")
	}
	request = strings.TrimPrefix(request, "request:")

	cf.cantabularServer.Handle([]byte(request), []byte(response))

	return nil
}

func (cf *CantabularFeature) cantabularReturnsThisStaticDatasetForTheGivenRequest(docs *godog.DocString) error {
	request, response, found := strings.Cut(docs.Content, "response:")
	if !found {
		return errors.New("CantabularFeature::cantabularReturnsThisStaticDatasetForTheGivenRequest - request and response were not found")
	}
	request = strings.TrimPrefix(request, "request:")

	cf.cantabularServer.Handle([]byte(request), []byte(response))

	return nil
}

func (cf *CantabularFeature) cantabularRespondsWithAnError() {
	cf.cantabularMock.OptionsHappy = false
}

func (cf *CantabularFeature) setMockedServer() {
	service.GetCantabularClient = func(cfg *config.Config) service.CantabularClient {
		return cantabular.NewClient(
			cantabular.Config{
				Host:           cfg.CantabularURL,
				ExtApiHost:     cfg.CantabularExtURL,
				GraphQLTimeout: cfg.DefaultRequestTimeout,
			},
			dphttp.NewClient(),
			nil,
		)
	}
}

func (cf *CantabularFeature) setMockedInterface() {
	service.GetCantabularClient = func(cfg *config.Config) service.CantabularClient {
		return cf.cantabularMock
	}
}

func (cf *CantabularFeature) setInitialiserMock() {
	cf.setMockedInterface()
}
