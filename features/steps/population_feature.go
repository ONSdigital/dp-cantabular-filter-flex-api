package steps

import (
	"encoding/json"
	"fmt"
	"testing"

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

func NewPopulationFeature(t *testing.T, cfg *config.Config) *PopulationFeature {
	return &PopulationFeature{
		client: mock.NewPopulationTypesAPIClient(),
	}
}

func (f *PopulationFeature) RegisterSteps(ctx *godog.ScenarioContext) {
	ctx.Step(
		`^Population Types API returns this GetCategorisations response for the given request:$`,
		f.PopulationTypesReturnsTheseCategorisations,
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
		return fmt.Errorf("unable to unmarshal request:response body %w", err)
	}

	f.client.GetCategorisationsResponses[rr.Req] = rr.Resp

	return nil
}
