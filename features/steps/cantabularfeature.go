package steps

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/config"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/features/mock"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/service"

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
	*mock.CantabularClient
}

func NewCantabularFeature() *CantabularFeature {
	return &CantabularFeature{
		CantabularClient: &mock.CantabularClient{
			OptionsHappy:    true,
			DimensionsHappy: true,
		},
	}
}

func (cf *CantabularFeature) Reset() {
	cf.CantabularClient.Reset()
}

func (cf *CantabularFeature) RegisterSteps(ctx *godog.ScenarioContext) {
	ctx.Step(
		`^Cantabular returns these dimensions for the dataset "([^"]*)" and search term "([^"]*)":$`,
		cf.cantabularSearchReturnsTheseDimensions,
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

}

func (cf *CantabularFeature) cantabularReturnsMultipleDimensions(datasetID string, docs *godog.DocString) error {
	cantabularResponses := struct {
		Responses map[string]cantabular.GetDimensionsResponse `json:"responses"`
	}{}

	if err := json.Unmarshal([]byte(docs.Content), &cantabularResponses); err != nil {
		return fmt.Errorf("unable to unmarshal cantabular search response: %w", err)
	}

	cf.CantabularClient.GetDimensionsByNameFunc = func(_ context.Context, req cantabular.GetDimensionsByNameRequest) (*cantabular.GetDimensionsResponse, error) {
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

func (cf *CantabularFeature) cantabularSearchReturnsTheseDimensions(datasetID, dimension string, docs *godog.DocString) error {
	var resp cantabular.GetDimensionsResponse
	if err := json.Unmarshal([]byte(docs.Content), &resp); err != nil {
		return fmt.Errorf("unable to unmarshal cantabular search response: %w", err)
	}

	cf.CantabularClient.GetDimensionsByNameFunc = func(_ context.Context, req cantabular.GetDimensionsByNameRequest) (*cantabular.GetDimensionsResponse, error) {
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

func (cf *CantabularFeature) cantabularRespondsWithAnError() {
	cf.OptionsHappy = false
	cf.DimensionsHappy = false
}

func (cf *CantabularFeature) cantabularGetOptionsRespondsWithAnError() {
	cf.OptionsHappy = false
}

func (cf *CantabularFeature) setInitialiserMock() {
	service.GetCantabularClient = func(cfg *config.Config) service.CantabularClient {
		return cf.CantabularClient
	}
}
