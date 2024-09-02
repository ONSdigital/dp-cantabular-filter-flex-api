package steps

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/config"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/features/mock"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/service"
	dphttp "github.com/ONSdigital/dp-net/v2/http"

	"github.com/cucumber/godog"
)

type er struct {
	err        error
	statusCode int
}

func (e *er) Error() string {
	if e.err == nil {
		return "nil"
	}
	return e.err.Error()
}

func (e *er) Code() int {
	return e.statusCode
}

type CantabularFeature struct {
	cantabularClient *mock.CantabularClient
	cantabularServer *mock.CantabularServer

	cfg *config.Config
}

func NewCantabularFeature(t *testing.T, cfg *config.Config) *CantabularFeature {
	return &CantabularFeature{
		cantabularClient: &mock.CantabularClient{
			OptionsHappy:     true,
			DimensionsHappy:  true,
			ResponseTooLarge: false,
		},
		cantabularServer: mock.NewCantabularServer(t),
		cfg:              cfg,
	}
}

func (cf *CantabularFeature) Reset() {
	cf.cantabularClient.Reset()
	cf.cantabularServer.Reset()

	cf.setMockedInterface()
}

func (cf *CantabularFeature) RegisterSteps(ctx *godog.ScenarioContext) {
	ctx.Step(
		`^the Cantabular service is a mocked extended Cantabular server$`,
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
		`^Cantabular returns these categorisations for the dataset "([^"]*)" and search term "([^"]*)":$`,
		cf.cantabularSearchReturnsTheseCategories,
	)

	ctx.Step(
		`^Cantabular returns these geography dimensions for the given request:$`,
		cf.cantabularReturnsTheseGeographyDimensionsForTheGivenRequest,
	)
	ctx.Step(
		`^Cantabular returns this area for the given request:$`,
		cf.cantabularReturnsThisAreaForTheGivenRequest,
	)
	ctx.Step(
		`^Cantabular returns this static dataset for the given request:$`,
		cf.cantabularReturnsThisStaticDatasetForTheGivenRequest,
	)
	ctx.Step(
		`^Cantabular responds with an error$`,
		cf.cantabularRespondsWithAnError,
	)
	ctx.Step(
		`^Cantabular GetOptions responds with an error$`,
		cf.cantabularGetOptionsRespondsWithAnError,
	)

	ctx.Step(
		`^Cantabular returns dimensions for the dataset "([^"]*)" for the following search terms:$`,
		cf.cantabularReturnsMultipleDimensions,
	)

	ctx.Step(
		`^Cantabular returns this response for the given request:$`,
		cf.cantabularReturnsThisResponseForTheGivenRequest,
	)

	ctx.Step(`^the count check returns that the count is too large$`,
		cf.countIsTooLarge,
	)

	ctx.Step(
		`^the query sent to cantabular for the check count is:$`,
		cf.theQuerySentToCantabularForTheCheckCountIs,
	)
}

func (cf *CantabularFeature) cantabularSearchReturnsTheseCategories(datasetID, dimension string, docs *godog.DocString) error {
	var resp cantabular.GetCategorisationsResponse
	if err := json.Unmarshal([]byte(docs.Content), &resp); err != nil {
		return fmt.Errorf("unable to unmarshal cantabular search response: %w", err)
	}

	cf.cantabularClient.GetCategorisationsFunc = func(_ context.Context, req cantabular.GetCategorisationsRequest) (*cantabular.GetCategorisationsResponse, error) {
		if req.Dataset == "" {
			return nil, errors.New("no dataset provided in request")
		}
		if req.Dataset == datasetID && req.Variable == dimension {
			return &resp, nil
		}

		return nil, &er{
			err:        errors.New("variable at position 1 does not exist"),
			statusCode: http.StatusNotFound,
		}
	}

	return nil
}

// cantabularReturnsMultipleDimensions sets up a stub response for the `GetDimensionsByName` method.
func (cf *CantabularFeature) cantabularReturnsMultipleDimensions(_ string, docs *godog.DocString) error {
	cantabularResponses := struct {
		Responses map[string]cantabular.GetDimensionsResponse `json:"responses"`
	}{}

	if err := json.Unmarshal([]byte(docs.Content), &cantabularResponses); err != nil {
		return fmt.Errorf("unable to unmarshal cantabular search response: %w", err)
	}

	cf.cantabularClient.GetDimensionsByNameFunc = func(_ context.Context, req cantabular.GetDimensionsByNameRequest) (*cantabular.GetDimensionsResponse, error) {
		if len(req.DimensionNames) == 0 {
			return nil, errors.New("no dimension provided in request")
		}
		if resp, ok := cantabularResponses.Responses[req.DimensionNames[0]]; ok {
			return &resp, nil
		}

		return nil, &er{
			err:        errors.New("variable at position 1 does not exist"),
			statusCode: http.StatusNotFound,
		}
	}

	return nil
}

func (cf *CantabularFeature) useAMockedExtCantabularServer() error {
	cf.cfg.CantabularExtURL = cf.cantabularServer.ResolveURL("")

	cf.setMockedServer()

	return nil
}

func (cf *CantabularFeature) cantabularSearchReturnsTheseDimensions(datasetID, dimension string, docs *godog.DocString) error {
	var resp cantabular.GetDimensionsResponse
	if err := json.Unmarshal([]byte(docs.Content), &resp); err != nil {
		return fmt.Errorf("unable to unmarshal cantabular search response: %w", err)
	}

	cf.cantabularClient.GetDimensionsByNameFunc = func(_ context.Context, req cantabular.GetDimensionsByNameRequest) (*cantabular.GetDimensionsResponse, error) {
		if len(req.DimensionNames) == 0 {
			return nil, errors.New("no dimension provided in request")
		}
		if req.Dataset == datasetID && req.DimensionNames[0] == dimension {
			return &resp, nil
		}

		return nil, &er{
			err:        errors.New("variable at position 1 does not exist"),
			statusCode: http.StatusNotFound,
		}
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

func (cf *CantabularFeature) cantabularReturnsThisAreaForTheGivenRequest(docs *godog.DocString) error {
	request, response, found := strings.Cut(docs.Content, "response:")
	if !found {
		return errors.New("CantabularFeature::cantabularReturnsThisAreaForTheGivenRequest - request and response were not found")
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

func (cf *CantabularFeature) theQuerySentToCantabularForTheCheckCountIs(docs *godog.DocString) error {
	request, response, found := strings.Cut(docs.Content, "response:")
	if !found {
		return errors.New("CantabularFeature::theQuerySentToCantabularForTheCheckCountIs - request and response were not found")
	}
	request = strings.TrimPrefix(request, "request:")

	cf.cantabularServer.Handle([]byte(request), []byte(response))
	cf.cantabularClient.ResponseTooLarge = true

	return nil
}

func (cf *CantabularFeature) cantabularReturnsThisResponseForTheGivenRequest(docs *godog.DocString) error {
	request, response, found := strings.Cut(docs.Content, "response:")
	if !found {
		return errors.New("badly formed [request:response] body")
	}
	request = strings.TrimPrefix(request, "request:")

	cf.cantabularServer.Handle([]byte(request), []byte(response))

	return nil
}

func (cf *CantabularFeature) cantabularRespondsWithAnError() {
	cf.cantabularClient.OptionsHappy = false
	cf.cantabularClient.DimensionsHappy = false
}

func (cf *CantabularFeature) cantabularGetOptionsRespondsWithAnError() {
	cf.cantabularClient.OptionsHappy = false
	cf.cantabularClient.OptionsHappy = false
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
		return cf.cantabularClient
	}
}

func (cf *CantabularFeature) setInitialiserMock() {
	cf.setMockedInterface()
}

func (cf *CantabularFeature) countIsTooLarge() error {
	cf.cantabularClient.ResponseTooLarge = true
	cf.cantabularClient.CheckQueryCountFunc = func(_ context.Context, req cantabular.StaticDatasetQueryRequest) (int, error) {
		return 180000, nil
	}

	return nil
}
