package steps

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular/gql"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/features/mock"

	"github.com/cucumber/godog"
)

type CantabularFeature struct {
	*mock.CantabularClient
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
}

func (cf *CantabularFeature) cantabularSearchReturnsTheseDimensions(datasetID, dimension string, docs *godog.DocString) error {
	var response cantabular.GetDimensionsResponse
	if err := json.Unmarshal([]byte(docs.Content), &response); err != nil {
		return fmt.Errorf("unable to unmarshal cantabular search response: %w", err)
	}

	cf.CantabularClient.SearchDimensionsFunc = func(ctx context.Context, req cantabular.SearchDimensionsRequest) (*cantabular.GetDimensionsResponse, error) {
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

func (cf *CantabularFeature) cantabularRespondsWithAnError() {
	cf.OptionsHappy = false
}
