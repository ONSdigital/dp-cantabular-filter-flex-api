package mock

import (
	"context"
	"errors"
	"sync"

	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
)

type PopulationTypesAPIClient struct {
	mu                          *sync.Mutex
	GetCategorisationsResponses map[population.GetCategorisationsInput]population.GetCategorisationsResponse
	GetCategoristionsHappy      bool
	GetMetadataResponse         population.GetPopulationTypeMetadataResponse
}

func NewPopulationTypesAPIClient() *PopulationTypesAPIClient {
	return &PopulationTypesAPIClient{
		mu:                          &sync.Mutex{},
		GetCategoristionsHappy:      true,
		GetCategorisationsResponses: make(map[population.GetCategorisationsInput]population.GetCategorisationsResponse),
	}
}

func (c *PopulationTypesAPIClient) Reset() {
	c.GetCategorisationsResponses = make(map[population.GetCategorisationsInput]population.GetCategorisationsResponse)
}

func (c *PopulationTypesAPIClient) Checker(_ context.Context, _ *healthcheck.CheckState) error {
	return nil
}

func (c *PopulationTypesAPIClient) GetCategorisations(ctx context.Context, req population.GetCategorisationsInput) (population.GetCategorisationsResponse, error) {
	if !c.GetCategoristionsHappy {
		return population.GetCategorisationsResponse{}, errors.New("failed to get categorisations")
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	req.ServiceAuthToken = "testToken"
	if resp, ok := c.GetCategorisationsResponses[req]; ok {
		return resp, nil
	}

	return population.GetCategorisationsResponse{}, errors.New("no response for provided input")
}

func (c *PopulationTypesAPIClient) SetGetCategorisationsResponse(req population.GetCategorisationsInput, res population.GetCategorisationsResponse) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.GetCategorisationsResponses[req] = res
}

func (c *PopulationTypesAPIClient) SetDefaultDatasetMetadata(res population.GetPopulationTypeMetadataResponse) {
	c.GetMetadataResponse = res
}

func (c *PopulationTypesAPIClient) GetPopulationTypeMetadata(context.Context, population.GetPopulationTypeMetadataInput) (population.GetPopulationTypeMetadataResponse, error) {
	return c.GetMetadataResponse, nil
}
