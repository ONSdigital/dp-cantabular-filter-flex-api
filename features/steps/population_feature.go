package steps

import (
	"encoding/json"
	"fmt"

	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/config"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/features/mock"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/service"
	"github.com/cucumber/godog"
)

type PopulationFeature struct {
	client *mock.PopulationTypesAPIClient
}

func (f *PopulationFeature) Reset() {
	f.client.Reset()
}

func NewPopulationFeature() *PopulationFeature {
	return &PopulationFeature{
		client: mock.NewPopulationTypesAPIClient(),
	}
}

func (f *PopulationFeature) RegisterSteps(ctx *godog.ScenarioContext) {
	ctx.Step(
		`^Population Types API returns this GetCategorisations response for the given request:$`,
		f.PopulationTypesReturnsTheseCategorisations,
	)
	ctx.Step(
		`^Population Types API returns this GetDefaultDatasetMetadata response for the given request:$`,
		f.PopulationTypesReturnsThisDefaultDatasetMetadata,
	)
}

func (f *PopulationFeature) setMockedInterface() {
	service.GetPopulationClient = func(cfg *config.Config) (service.PopulationTypesAPIClient, error) {
		return f.client, nil
	}
}

func (f *PopulationFeature) setInitialiserMock() {
	f.setMockedInterface()
}

func (f *PopulationFeature) PopulationTypesReturnsTheseCategorisations(input *godog.DocString) error {
	rr := &struct {
		Req  population.GetCategorisationsInput    `json:"request"`
		Resp population.GetCategorisationsResponse `json:"response"`
	}{}

	if err := json.Unmarshal([]byte(input.Content), &rr); err != nil {
		return fmt.Errorf("unable to unmarshal [request:response] body: %w", err)
	}

	f.client.SetGetCategorisationsResponse(rr.Req, rr.Resp)

	return nil
}

func (f *PopulationFeature) PopulationTypesReturnsThisDefaultDatasetMetadata(input *godog.DocString) error {
	var rr population.GetPopulationTypeMetadataResponse

	if err := json.Unmarshal([]byte(input.Content), &rr); err != nil {
		return fmt.Errorf("unable to unmarshal [request:response] body: %w", err)
	}

	f.client.SetDefaultDatasetMetadata(rr)

	return nil
}
