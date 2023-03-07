package steps

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabularmetadata"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/config"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/features/mock"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/service"
	"github.com/cucumber/godog"
)

type MetadataFeature struct {
	metadataClient *mock.CantabularMetadataClient
}

func (mf *MetadataFeature) Reset() {
	mf.metadataClient.Reset()
	mf.setMockedInterface()
}

func NewMetadataFeature(t *testing.T, mfg *config.Config) *MetadataFeature {
	return &MetadataFeature{metadataClient: &mock.CantabularMetadataClient{
		OptionsHappy:    true,
		DimensionsHappy: true,
	}}
}

func (mf *MetadataFeature) RegisterSteps(ctx *godog.ScenarioContext) {

	ctx.Step(
		`^Metadata api returns this response for the dataset "([^"]*)" and search term "([^"]*)":$`,
		mf.MetadataReturnsTheseDefaults,
	)
	ctx.Step(
		`^the metadata API returns an error`,
		mf.MetadataReturnsAnError,
	)
}

func (mf *MetadataFeature) MetadataReturnsTheseDefaults(datasetID, search string, input *godog.DocString) error {

	res := &struct {
		Data   cantabularmetadata.Data       `json:"data"`
		Errors []cantabularmetadata.GQLError `json:"errors,omitempty"`
	}{}

	if err := json.Unmarshal([]byte(input.Content), &res); err != nil {
		return fmt.Errorf("unable to unmarshal cantabular search response: %w", err)
	}

	mf.metadataClient.GetDefaultClassificationFunc = func(ctx context.Context, req cantabularmetadata.GetDefaultClassificationRequest) (*cantabularmetadata.GetDefaultClassificationResponse, error) {
		cantabularResponse := cantabularmetadata.GetDefaultClassificationResponse{
			Variables: []string{res.Data.Dataset.Vars[0].Name},
		}

		return &cantabularResponse, nil
	}

	return nil
}

func (mf *MetadataFeature) MetadataReturnsAnError() error {
	var cantabularResponse cantabularmetadata.GetDefaultClassificationResponse

	mf.metadataClient.GetDefaultClassificationFunc = func(ctx context.Context, req cantabularmetadata.GetDefaultClassificationRequest) (*cantabularmetadata.GetDefaultClassificationResponse, error) {
		return &cantabularResponse, errors.New("Not Found DefaultClassification")
	}

	return nil
}

func (cf *MetadataFeature) setMockedInterface() {
	service.GetMetadataClient = func(cfg *config.Config) service.MetadataClient {
		return cf.metadataClient
	}
}

func (cf *MetadataFeature) setInitialiserMock() {
	cf.setMockedInterface()
}
