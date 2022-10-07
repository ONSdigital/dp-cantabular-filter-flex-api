package steps

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular/gql"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/config"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/features/mock"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/service"

	"github.com/cucumber/godog"
)

type CantabularFeature struct {
	*mock.CantabularClient
}

func NewCantabularFeature() *CantabularFeature {
	return &CantabularFeature{CantabularClient: &mock.CantabularClient{OptionsHappy: true}}
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
		`^Cantabular returns dimensions for the dataset "([^"]*)" for the following search terms:$`,
		cf.cantabularReturnsMultipleDimensions,
	)

}

// cantabylarSearchReturnsOneOfTheseDimensions sets up a stub response for the `SearchDimensions` method.
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

func (cf *CantabularFeature) cantabularRespondsWithAnError() {
	cf.OptionsHappy = false
}

func (cf *CantabularFeature) setInitialiserMock() {
	service.GetCantabularClient = func(cfg *config.Config) service.CantabularClient {
		return cf.CantabularClient
	}
}
