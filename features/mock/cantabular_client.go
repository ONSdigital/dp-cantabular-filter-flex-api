package mock

import (
	"context"
	"errors"
	"fmt"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
)

type CantabularClient struct {
	ErrStatus               int
	OptionsHappy            bool
	DimensionsHappy         bool
	GetDimensionsByNameFunc func(context.Context, cantabular.GetDimensionsByNameRequest) (*cantabular.GetDimensionsResponse, error)
}

func (c *CantabularClient) Reset() {
	c.ErrStatus = 500
	c.OptionsHappy = true
	c.DimensionsHappy = true
}

func (c *CantabularClient) StatusCode(_ error) int {
	return c.ErrStatus
}

func (c *CantabularClient) GetDimensionOptions(_ context.Context, req cantabular.GetDimensionOptionsRequest) (*cantabular.GetDimensionOptionsResponse, error) {
	if c.OptionsHappy {
		fmt.Printf("DEBUG HAPPY OPTIONS CALLED: %v\n", req)
		return nil, nil
	}
	fmt.Printf("DEBUG UNHAPPY OPTIONS CALLED: %v\n", req)

	return nil, errors.New("invalid dimension options")
}

func (c *CantabularClient) StaticDatasetQuery(context.Context, cantabular.StaticDatasetQueryRequest) (*cantabular.StaticDatasetQuery, error) {
	return nil, errors.New("invalid dataset query")
}

func (c *CantabularClient) GetGeographyDimensions(context.Context, cantabular.GetGeographyDimensionsRequest) (*cantabular.GetGeographyDimensionsResponse, error) {
	return nil, errors.New("invalid geography query")
}

func (c *CantabularClient) GetDimensionsByName(ctx context.Context, req cantabular.GetDimensionsByNameRequest) (*cantabular.GetDimensionsResponse, error) {
	if c.DimensionsHappy {
		return c.GetDimensionsByNameFunc(ctx, req)
	}
	return nil, errors.New("error while searching dimensions")
}

func (c *CantabularClient) Checker(_ context.Context, _ *healthcheck.CheckState) error {
	return nil
}

func (c *CantabularClient) CheckerAPIExt(_ context.Context, _ *healthcheck.CheckState) error {
	return nil
}
