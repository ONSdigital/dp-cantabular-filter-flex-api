package steps

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"

	"github.com/cucumber/godog"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
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
		`^the maximum pagination limit is set to (\d+)$`,
		c.theMaximumLimitIsSetTo,
	)
	ctx.Step(
		`^the document in the database for id "([^"]*)" should be:$`,
		c.theDocumentInTheDatabaseShouldBe,
	)
	ctx.Step(
		`^the following version document with dataset id "([^"]*)", edition "([^"]*)" and version "([^"]*)" is available from dp-dataset-api:$`,
		c.theFollowingVersionDocumentIsAvailable,
	)
	ctx.Step(
		`^I have these filters:$`,
		c.iHaveTheseFilters,
	)

	ctx.Step(
		`^Mongo datastore fails for update filter output`,
		c.MongoDatastoreFailsForUpdateFilterOutput,
	)
}

func (c *Component) MongoDatastoreFailsForUpdateFilterOutput() error {
	var err error
	c.store, err = GetFailingMongo(c.ctx, c.cfg, c.g)
	if err != nil {
		return fmt.Errorf("failed to create new mongo mongoClient: %w", err)
	}

	return nil
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

func (c *Component) theMaximumLimitIsSetTo(val int) error {
	c.cfg.DefaultMaximumLimit = val
	return nil
}

func (c *Component) iHaveTheseFilters(docs *godog.DocString) error {
	ctx := context.Background()
	var filters []model.Filter

	err := json.Unmarshal([]byte(docs.Content), &filters)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshall")
	}

	store := c.store
	col := c.cfg.FiltersCollection

	for _, f := range filters {
		if _, err = store.Conn().Collection(col).UpsertById(ctx, f.ID, bson.M{"$set": f}); err != nil {
			return errors.Wrap(err, "failed to upsert filter")
		}
	}

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
