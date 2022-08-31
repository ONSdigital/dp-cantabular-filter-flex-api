package mock

import (
	"context"
	"errors"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
)

type CantabularClient struct {
	ErrStatus                  int
	OptionsHappy               bool
	SearchDimensionsFunc       func(ctx context.Context, req cantabular.SearchDimensionsRequest) (*cantabular.GetDimensionsResponse, error)
	GetGeographyDimensionsFunc func(context.Context, cantabular.GetGeographyDimensionsRequest) (*cantabular.GetGeographyDimensionsResponse, error)
	StaticDatasetQueryFunc     func(context.Context, cantabular.StaticDatasetQueryRequest) (*cantabular.StaticDatasetQuery, error)
}

func (c *CantabularClient) Reset() {
	c.ErrStatus = 500
	c.OptionsHappy = true
}

func (c *CantabularClient) StatusCode(_ error) int {
	return c.ErrStatus
}

func (c *CantabularClient) GetDimensionOptions(_ context.Context, _ cantabular.GetDimensionOptionsRequest) (*cantabular.GetDimensionOptionsResponse, error) {
	if c.OptionsHappy {
		return nil, nil
	}
	return nil, errors.New("invalid dimension options")
}

func (c *CantabularClient) StaticDatasetQuery(ctx context.Context, req cantabular.StaticDatasetQueryRequest) (*cantabular.StaticDatasetQuery, error) {
	if c.OptionsHappy {
		return c.StaticDatasetQueryFunc(ctx, req)
	}
	return nil, errors.New("error while executing dataset query")
}

func (c *CantabularClient) GetGeographyDimensions(ctx context.Context, req cantabular.GetGeographyDimensionsRequest) (*cantabular.GetGeographyDimensionsResponse, error) {
	if c.OptionsHappy {
		return c.GetGeographyDimensionsFunc(ctx, req)
	}
	return nil, errors.New("error while getting geography dimensions")
}

func (c *CantabularClient) SearchDimensions(ctx context.Context, req cantabular.SearchDimensionsRequest) (*cantabular.GetDimensionsResponse, error) {
	if c.OptionsHappy {
		return c.SearchDimensionsFunc(ctx, req)
	}
	return nil, errors.New("error while searching dimensions")
}

func (c *CantabularClient) Checker(_ context.Context, _ *healthcheck.CheckState) error {
	return nil
}

func (c *CantabularClient) CheckerAPIExt(_ context.Context, _ *healthcheck.CheckState) error {
	return nil
}
