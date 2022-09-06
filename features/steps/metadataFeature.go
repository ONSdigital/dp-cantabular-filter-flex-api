package steps

import (
	"testing"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/config"
	"github.com/cucumber/godog"
	"github.com/maxcnunes/httpfake"
)

type MetadataFeature struct {
	mockMetadataServer *httpfake.HTTPFake
}

func NewMetadataFeature(t *testing.T, cfg *config.Config) *MetadataFeature {
	df := &MetadataFeature{mockMetadataServer: httpfake.New(httpfake.WithTesting(t))}
	cfg.MetadataAPIURL = df.mockMetadataServer.ResolveURL("")

	return df
}

func (df *MetadataFeature) Reset() { df.mockMetadataServer.Reset() }
func (df *MetadataFeature) Close() { df.mockMetadataServer.Close() }

func (cf *MetadataFeature) RegisterSteps(ctx *godog.ScenarioContext) {

	ctx.Step(
		"^ the population types api returns these categorisations",
		cf.MetadataReturnsTheseDefaults,
	)
	ctx.Step(
		"^ the population types api returns an error",
		cf.MetadataReturnsAnError,
	)
}

func (cf *MetadataFeature) MetadataReturnsTheseDefaults(input *godog.DocString) error {
	return nil
}

func (cf *MetadataFeature) MetadataReturnsAnError() error {
	return nil
}
