package mock

import (
	"context"
	"errors"
	"fmt"

	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
)

type PopulationTypesAPIClient struct {
	//GetCategorisationsFunc      func(ctx context.Context, req population.GetCategorisationsInput) (*population.GetCategorisationsResponse, error)
	GetCategorisationsResponses map[population.GetCategorisationsInput]population.GetCategorisationsResponse
	GetCategoristionsHappy      bool
}

func NewPopulationTypesAPIClient() *PopulationTypesAPIClient {
	return &PopulationTypesAPIClient{
		GetCategoristionsHappy:      true,
		GetCategorisationsResponses: make(map[population.GetCategorisationsInput]population.GetCategorisationsResponse),
	}
}

func (c *PopulationTypesAPIClient) Reset() {}

func (c *PopulationTypesAPIClient) Checker(_ context.Context, _ *healthcheck.CheckState) error {
	return nil
}

func (c *PopulationTypesAPIClient) GetCategorisations(ctx context.Context, req population.GetCategorisationsInput) (population.GetCategorisationsResponse, error) {
	if !c.GetCategoristionsHappy {
		return population.GetCategorisationsResponse{}, errors.New("failed to get categorisations")
	}

	if resp, ok := c.GetCategorisationsResponses[req]; ok {
		return resp, nil
	}
	fmt.Println("DEBUGMAP")
	for k, _ := range c.GetCategorisationsResponses {
		fmt.Printf("%+v\n", k)
	}

	fmt.Printf("DEBUGREQ: %+v", req)
	return population.GetCategorisationsResponse{}, errors.New("no response for provided input")
}
