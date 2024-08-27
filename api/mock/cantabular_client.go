package mock

import (
	"context"
	"errors"
	"fmt"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular/gql"
	stream "github.com/ONSdigital/dp-api-clients-go/v2/stream"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
)

type CantabularClient struct {
	ErrStatus               int
	OptionsHappy            bool
	ExpectedFilterDimension string
	GetObservationsResponse *cantabular.GetObservationsResponse
}

func (c *CantabularClient) StatusCode(_ error) int {
	return c.ErrStatus
}

func (c *CantabularClient) GetDimensionOptions(_ context.Context, req cantabular.GetDimensionOptionsRequest) (*cantabular.GetDimensionOptionsResponse, error) {
	if !c.OptionsHappy {
		return nil, errors.New("invalid dimension options")
	}

	if c.ExpectedFilterDimension == "" {
		return nil, nil
	}

	for _, f := range req.Filters {
		if f.Variable == c.ExpectedFilterDimension {
			return nil, nil
		}
	}

	return nil, fmt.Errorf(
		"expected dimension not found in request (expected: %s found: %v)",
		c.ExpectedFilterDimension,
		req.DimensionNames,
	)
}

func (c *CantabularClient) StaticDatasetQuery(context.Context, cantabular.StaticDatasetQueryRequest) (*cantabular.StaticDatasetQuery, error) {
	return nil, errors.New("invalid dataset query")
}

func (c *CantabularClient) GetGeographyDimensionsInBatches(context.Context, string, int, int) (*gql.Dataset, error) {
	return nil, errors.New("invalid geography query")
}

func (c *CantabularClient) GetDimensionsByName(_ context.Context, _ cantabular.GetDimensionsByNameRequest) (*cantabular.GetDimensionsResponse, error) {
	return nil, nil
}

func (c *CantabularClient) GetArea(context.Context, cantabular.GetAreaRequest) (*cantabular.GetAreaResponse, error) {
	return nil, errors.New("invalid area query")
}

func (c *CantabularClient) Checker(_ context.Context, _ *healthcheck.CheckState) error {
	return nil
}

func (c *CantabularClient) CheckerAPIExt(_ context.Context, _ *healthcheck.CheckState) error {
	return nil
}

func (c *CantabularClient) GetCategorisations(context.Context, cantabular.GetCategorisationsRequest) (*cantabular.GetCategorisationsResponse, error) {
	if c.OptionsHappy {
		return nil, nil
	}
	return nil, errors.New("invalid Categorisation request")
}

func (c *CantabularClient) CheckQueryCount(_ context.Context, _ cantabular.StaticDatasetQueryRequest) (int, error) {
	return 0, nil
}

func (c *CantabularClient) StaticDatasetQueryStreamJSON(_ context.Context, _ cantabular.StaticDatasetQueryRequest, _ stream.Consumer) (cantabular.GetObservationsResponse, error) {
	return *c.GetObservationsResponse, nil
}
