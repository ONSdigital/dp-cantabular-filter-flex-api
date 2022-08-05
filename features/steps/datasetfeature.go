package steps

import (
	"fmt"
	"net/http"

	"github.com/cucumber/godog"
	"github.com/maxcnunes/httpfake"
)

type DatasetFeature struct {
	mockDatasetServer *httpfake.HTTPFake
}

func (df *DatasetFeature) Reset() { df.mockDatasetServer.Reset() }
func (df *DatasetFeature) Close() { df.mockDatasetServer.Close() }

func (df *DatasetFeature) RegisterSteps(ctx *godog.ScenarioContext) {
	ctx.Step(
		`^the following version document with dataset id "([^"]*)", edition "([^"]*)" and version "([^"]*)" is available from dp-dataset-api:$`,
		df.theFollowingVersionDocumentIsAvailable,
	)
	ctx.Step(
		`^the client for the dataset API failed and is returning errors$`,
		df.theClientForTheDatasetAPIFailedAndIsReturningErrors,
	)
}

func (df *DatasetFeature) datasetDoesSomething() error {
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

	df.mockDatasetServer.NewHandler().
		Get(url).
		Reply(http.StatusOK).
		BodyString(v.Content)

	return nil
}

func (df *DatasetFeature) theClientForTheDatasetAPIFailedAndIsReturningErrors() error {
	df.mockDatasetServer.Reset()
	df.mockDatasetServer.NewHandler().
		Get("/datasets/cantabular-example-1/editions/2021/versions/1").
		Reply(http.StatusInternalServerError)
	return nil
}
