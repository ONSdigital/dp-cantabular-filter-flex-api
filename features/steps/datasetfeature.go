package steps

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/config"

	"github.com/cucumber/godog"
	"github.com/maxcnunes/httpfake"
)

type DatasetFeature struct {
	datasetServer *httpfake.HTTPFake
}

func NewDatasetFeature(t *testing.T, cfg *config.Config) *DatasetFeature {
	df := &DatasetFeature{datasetServer: httpfake.New(httpfake.WithTesting(t))}
	cfg.DatasetAPIURL = df.datasetServer.ResolveURL("")

	return df
}

func (df *DatasetFeature) Reset() { df.datasetServer.Reset() }
func (df *DatasetFeature) Close() { df.datasetServer.Close() }

func (df *DatasetFeature) RegisterSteps(ctx *godog.ScenarioContext) {
	ctx.Step(
		`^the following dataset document with dataset id "([^"]*)" is available from dp-dataset-api:$`,
		df.theFollowingDatasetDocumentIsAvailable,
	)
	ctx.Step(
		`^the following version document with dataset id "([^"]*)", edition "([^"]*)" and version "([^"]*)" is available from dp-dataset-api:$`,
		df.theFollowingVersionDocumentIsAvailable,
	)
	ctx.Step(
		`^the following metadata document for dataset id "([^"]*)", edition "([^"]*)" and version "([^"]*)" is available from dp-dataset-api:$`,
		df.theFollowingMetadataDocumentIsAvailable,
	)
	ctx.Step(
		`^the following dimensions document for dataset id "([^"]*)", edition "([^"]*)" and version "([^"]*)" is available from dp-dataset-api:$`,
		df.theFollowingDimensionsDocumentIsAvailable,
	)
	ctx.Step(
		`^the following options document for dataset id "([^"]*)", edition "([^"]*)", version "([^"]*)" and dimension "([^"]*)" is available from dp-dataset-api:$`,
		df.theFollowingOptionsDocumentIsAvailable,
	)
	ctx.Step(
		`^the client for the dataset API failed and is returning errors$`,
		df.theClientForTheDatasetAPIFailedAndIsReturningErrors,
	)
}

// theFollowingDatasetDocumentIsAvailable generates a mocked response for dataset API
// GET /datasets/{dataset_id}
func (df *DatasetFeature) theFollowingDatasetDocumentIsAvailable(datasetID string, v *godog.DocString) error {
	url := fmt.Sprintf("/datasets/%s", datasetID)

	df.datasetServer.NewHandler().
		Get(url).
		Reply(http.StatusOK).
		BodyString(v.Content)

	return nil
}

// theFollowingVersionDocumentIsAvailable generates a mocked response for dataset API
// GET /datasets/{dataset_id}/editions/{edition}/versions/{version}
func (df *DatasetFeature) theFollowingVersionDocumentIsAvailable(datasetID, edition, version string, v *godog.DocString) error {
	url := fmt.Sprintf(
		"/datasets/%s/editions/%s/versions/%s",
		datasetID,
		edition,
		version,
	)

	df.datasetServer.NewHandler().
		Get(url).
		Reply(http.StatusOK).
		BodyString(v.Content)

	return nil
}

// theFollowingMetadataDocumentIsAvailable generates a mocked response for dataset API
// GET /datasets/{dataset_id}/editions/{edition}/versions/{version}/metadata
func (df *DatasetFeature) theFollowingMetadataDocumentIsAvailable(datasetID, edition, version string, v *godog.DocString) error {
	url := fmt.Sprintf("/datasets/%s/editions/%s/versions/%s/metadata", datasetID, edition, version)

	df.datasetServer.NewHandler().
		Get(url).
		Reply(http.StatusOK).
		BodyString(v.Content)

	return nil
}

// theFollowingDimensionsDocumentIsAvailable generates a mocked response for dataset API
// GET /datasets/{dataset_id}/editions/{edition}/versions/{version}/dimensions
func (df *DatasetFeature) theFollowingDimensionsDocumentIsAvailable(datasetID, edition, version string, v *godog.DocString) error {
	url := fmt.Sprintf("/datasets/%s/editions/%s/versions/%s/dimensions", datasetID, edition, version)

	df.datasetServer.NewHandler().
		Get(url).
		Reply(http.StatusOK).
		BodyString(v.Content)

	return nil
}

// theFollowingOptionsDocumentIsAvailable generates a mocked response for dataset API
// GET /datasets/{dataset_id}/editions/{edition}/versions/{version}/dimensions/{dimension}/options
func (df *DatasetFeature) theFollowingOptionsDocumentIsAvailable(datasetID, edition, version, dimension string, v *godog.DocString) error {
	url := fmt.Sprintf("/datasets/%s/editions/%s/versions/%s/dimensions/%s/options", datasetID, edition, version, dimension)

	df.datasetServer.NewHandler().
		Get(url).
		Reply(http.StatusOK).
		BodyString(v.Content)

	return nil
}

func (df *DatasetFeature) theClientForTheDatasetAPIFailedAndIsReturningErrors() error {
	df.datasetServer.Reset()
	df.datasetServer.NewHandler().
		Get("/datasets/cantabular-example-1/editions/2021/versions/1").
		Reply(http.StatusInternalServerError)
	return nil
}

func (df *DatasetFeature) setInitialiserMock() {
	// the default initialiser is fine as it will pick up the dataset server url from the config
}
