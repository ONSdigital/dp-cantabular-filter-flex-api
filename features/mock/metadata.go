package mock

import (
	"context"
	"errors"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/cantabularmetadata"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
)

type CantabularMetadataClient struct {
	ErrStatus                    int
	OptionsHappy                 bool
	DimensionsHappy              bool
	GetDefaultClassificationFunc func(ctx context.Context, req cantabularmetadata.GetDefaultClassificationRequest) (*cantabularmetadata.GetDefaultClassificationResponse, error)
}

func (c *CantabularMetadataClient) Reset() {
	c.ErrStatus = 500
	c.OptionsHappy = true
	c.DimensionsHappy = true
}

func (c *CantabularMetadataClient) StatusCode(_ error) int {
	return c.ErrStatus
}

func (c *CantabularMetadataClient) GetDimensionOptions(_ context.Context, _ cantabular.GetDimensionOptionsRequest) (*cantabular.GetDimensionOptionsResponse, error) {
	if c.OptionsHappy {
		return nil, nil
	}

	return nil, errors.New("invalid dimension options")
}

func (c *CantabularMetadataClient) GetDefaultClassification(ctx context.Context, req cantabularmetadata.GetDefaultClassificationRequest) (*cantabularmetadata.GetDefaultClassificationResponse, error) {
	if c.OptionsHappy {
		return c.GetDefaultClassificationFunc(ctx, req)
	}
	return nil, errors.New("error while retrieving default classification")
}

func (c *CantabularMetadataClient) Checker(_ context.Context, _ *healthcheck.CheckState) error {
	return nil
}

func (c *CantabularMetadataClient) CheckerAPIExt(_ context.Context, _ *healthcheck.CheckState) error {
	return nil
}
