package steps

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabularmetadata"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/config"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/features/mock"
	"github.com/cucumber/godog"
)

type MetadataFeature struct {
	*mock.MetadataApiClient
}

func NewMetadataFeature(t *testing.T, mfg *config.Config) *MetadataFeature {
	return &MetadataFeature{MetadataApiClient: &mock.MetadataApiClient{OptionsHappy: true}}
}

func (mf *MetadataFeature) Reset() {
	mf.MetadataApiClient.Reset()
}

func (mf *MetadataFeature) RegisterSteps(ctx *godog.ScenarioContext) {

	ctx.Step(
		"^ cantabular returns these default classifications",
		mf.MetadataReturnsTheseDefaults,
	)
	ctx.Step(
		"^ the cantabular categorisations return an error",
		mf.MetadataReturnsAnError,
	)
}

func (mf *MetadataFeature) MetadataReturnsTheseDefaults(input *godog.DocString) error {
	cantabularResponse := &cantabularmetadata.GetDefaultClassificationResponse{}

	if err := json.Unmarshal([]byte(input.Content), &cantabularResponse); err != nil {
		return fmt.Errorf("unable to unmarshal cantabular search response: %w", err)
	}

	mf.MetadataApiClient.GetDefaultCategorisationFunc = func(ctx context.Context, req cantabularmetadata.GetDefaultClassificationRequest) (*cantabularmetadata.GetDefaultClassificationResponse, error) {
		return cantabularResponse, nil
	}

	return nil
}

func (mf *MetadataFeature) MetadataReturnsAnError() error {
	mf.OptionsHappy = false
}
