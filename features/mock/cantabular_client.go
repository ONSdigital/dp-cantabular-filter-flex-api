package mock

import (
	"context"
	"errors"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
)

type CantabularClient struct {
	ErrStatus            int
	OptionsHappy         bool
	SearchDimensionsFunc func(ctx context.Context, req cantabular.SearchDimensionsRequest) (*cantabular.GetDimensionsResponse, error)
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
