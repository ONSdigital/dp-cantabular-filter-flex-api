package steps

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/event"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/schema"

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
		`^Mongo datastore fails`,
		c.MongoDatastoreFails,
	)

	ctx.Step(`^an ETag is returned`,
		c.anETagIsReturned,
	)

	ctx.Step(
		`^I provide If-Match header "([^"]*)"$`,
		c.iProvideIfMatchHeader,
	)

	ctx.Step(
		`^the client for the dataset API failed and is returning errors$`,
		c.theClientForTheDatasetAPIFailedAndIsReturningErrors,
	)

	ctx.Step(`^one event with the following fields are in the produced kafka topic catabular-export-start:$`,
		c.oneEventWithTheFollowingFieldsAreInTheProducedKafkaTopicCatabularexportstart,
	)

	ctx.Step(`^I have these filter outputs:$`,
		c.iHaveTheseFilterOutputs,
	)
}

func (c *Component) oneEventWithTheFollowingFieldsAreInTheProducedKafkaTopicCatabularexportstart() error {
	select {
	case <-time.After(c.waitEventTimeout):
		return nil
	case <-c.consumer.Channels().Closer:
		return errors.New("closer channel closed")
	case msg, ok := <-c.consumer.Channels().Upstream:
		if !ok {
			return errors.New("upstream channel closed")
		}

		var e event.ExportStart
		s := schema.ExportStart

		if err := s.Unmarshal(msg.GetData(), &e); err != nil {
			msg.Commit()
			msg.Release()
			return fmt.Errorf("error unmarshalling message: %w", err)
		}

		msg.Commit()
		msg.Release()

		return fmt.Errorf("kafka event received in csv-created topic: %v", e)
	}
}

func (c *Component) anETagIsReturned() error {
	eTag := c.ApiFeature.HttpResponse.Header.Get("ETag")
	if eTag == "" {
		return fmt.Errorf("no 'ETag' header returned")
	}
	return nil
}

func (c *Component) MongoDatastoreFails() error {
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

func (c *Component) iProvideIfMatchHeader(eTag string) error {
	c.ApiFeature.ISetTheHeaderTo("If-Match", eTag)
	return nil
}

func (c *Component) theClientForTheDatasetAPIFailedAndIsReturningErrors() error {
	c.DatasetAPI.Reset()
	c.DatasetAPI.NewHandler().
		Get("/datasets/cantabular-example-1/editions/2021/versions/1").
		Reply(http.StatusInternalServerError)
	return nil
}

func (c *Component) iHaveTheseFilterOutputs(docs *godog.DocString) error {
	ctx := context.Background()
	var filterOutputs []model.FilterOutput

	err := json.Unmarshal([]byte(docs.Content), &filterOutputs)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshall")
	}

	store := c.store
	col := c.cfg.FilterOutputsCollection

	for _, f := range filterOutputs {
		if _, err = store.Conn().Collection(col).UpsertById(ctx, f.ID, bson.M{"$set": f}); err != nil {
			return errors.Wrap(err, "failed to upsert filter output")
		}
	}

	return nil
}
