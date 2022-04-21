package mock

import (
	"context"
	"errors"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
)

type CantabularClient struct {
	ErrStatus    int
	OptionsHappy bool
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

func (c *CantabularClient) StaticDatasetQuery(context.Context, cantabular.StaticDatasetQueryRequest) (*cantabular.StaticDatasetQuery, error) {
	return nil, errors.New("invalid dataset query")
}

func (c *CantabularClient) GetGeographyDimensions(context.Context, string) (*cantabular.GetGeographyDimensionsResponse, error) {
	return nil, errors.New("invalid geography query")
}

func (c *CantabularClient) SearchDimensions(_ context.Context, _ cantabular.SearchDimensionsRequest) (*cantabular.GetDimensionsResponse, error) {
	if c.OptionsHappy {
		return nil, nil
	}
	return nil, errors.New("error searching dimensions")
}

func (c *CantabularClient) Checker(_ context.Context, _ *healthcheck.CheckState) error {
	return nil
}

func (c *CantabularClient) CheckerAPIExt(_ context.Context, _ *healthcheck.CheckState) error {
	return nil
}
