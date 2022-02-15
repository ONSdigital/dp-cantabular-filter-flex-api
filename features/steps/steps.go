package steps

import (
	"fmt"
	"net/http"

	"github.com/cucumber/godog"
)

func (c *Component) RegisterSteps(ctx *godog.ScenarioContext) {
	ctx.Step(
		`^the service starts`,
		c.theServiceStarts,
	)
	ctx.Step(
		`^private endpoints are enabled`,
		c.privateEndpointsAreEnabled,
	)
	ctx.Step(
		`^private endpoints are not enabled`,
		c.privateEndpointsAreNotEnabled,
	)
	ctx.Step(
		`^the document in the database for id "([^"]*)" should be:$`,
		c.theDocumentInTheDatabaseShouldBe,
	)
	ctx.Step(
		`^the following version document with dataset id "([^"]*)", edition "([^"]*)" and version "([^"]*)" is available from dp-dataset-api:$`,
		c.theFollowingVersionDocumentIsAvailable,
	)
}

// theServiceStarts starts the service under test in a new go-routine
// note that this step should be called only after all dependencies have been setup,
// to prevent any race condition, specially during the first healthcheck iteration.
func (c *Component) theServiceStarts() error {
	c.wg.Add(1)
	go c.startService(c.ctx)
	return nil
}

func (c *Component) privateEndpointsAreEnabled() error {
	c.cfg.EnablePrivateEndpoints = true
	return nil
}

func (c *Component) privateEndpointsAreNotEnabled() error {
	c.cfg.EnablePrivateEndpoints = false
	return nil
}

func (c *Component) theDocumentInTheDatabaseShouldBe(id string, doc *godog.DocString) error {
	// TODO: implement step for verifying documents stored in Mongo. No prior
	// art of this being done properly in ONS yet so save to be done in future ticket
	return nil
}

// theFollowingVersionDocumentIsAvailable generates a mocked response for dataset API
// GET /datasets/{dataset_id}/editions/{edition}/versions/{version}
func (c *Component) theFollowingVersionDocumentIsAvailable(datasetID, edition, version string, v *godog.DocString) error {
	url := fmt.Sprintf(
		"/datasets/%s/editions/%s/versions/%s",
		datasetID,
		edition,
		version,
	)

	c.DatasetAPI.NewHandler().
		Get(url).
		Reply(http.StatusOK).
		BodyString(v.Content)

	return nil
}
