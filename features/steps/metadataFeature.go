package steps

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabularmetadata"
	api "github.com/ONSdigital/dp-cantabular-filter-flex-api/api/mock"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/config"
	"github.com/cucumber/godog"
)

type MetadataFeature struct {
	*api.MockmetadataAPIClient
}

func NewMetadataFeature(t *testing.T, mfg *config.Config) *MetadataFeature {
	return &MetadataFeature{MockmetadataAPIClient: &api.MockmetadataAPIClient{}}
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

	return nil
}

func (mf *MetadataFeature) MetadataReturnsAnError() error {
	//	mf.OptionsHappy = false
	return nil
}
