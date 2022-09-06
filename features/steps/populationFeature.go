package steps

import (
	"testing"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/config"
	"github.com/cucumber/godog"
	"github.com/maxcnunes/httpfake"
)

type PopulationFeature struct {
	mockPopulationServer *httpfake.HTTPFake
}

func NewPopulationFeature(t *testing.T, cfg *config.Config) *PopulationFeature {
	df := &PopulationFeature{mockPopulationServer: httpfake.New(httpfake.WithTesting(t))}
	cfg.PopulationTypesAPIURL = df.mockPopulationServer.ResolveURL("")

	return df
}

func (cf *PopulationFeature) Reset() { cf.mockPopulationServer.Reset() }
func (cf *PopulationFeature) Close() { cf.mockPopulationServer.Close() }

func (cf *PopulationFeature) RegisterSteps(ctx *godog.ScenarioContext) {

	ctx.Step(
		"^ the population types api returns these categorisations",
		cf.PopulationTypesReturnsTheseCategorisations,
	)
	ctx.Step(
		"^ the population types api returns an error",
		cf.PopulationTypesReturnsAnError,
	)
}

func (cf *PopulationFeature) PopulationTypesReturnsTheseCategorisations(input *godog.DocString) error {
	return nil
}

func (cf *PopulationFeature) PopulationTypesReturnsAnError() error {
	return nil
}
