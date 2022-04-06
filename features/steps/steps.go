package steps

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/api"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	"github.com/ONSdigital/dp-kafka/v3/avro"

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

	// this is the same as above, but added for clearer step definition
	ctx.Step(
		`^Mongo datastore is failing`,
		c.MongoDatastoreIsFailing,
	)

	ctx.Step(`^an ETag is returned`,
		c.anETagIsReturned,
	)

	ctx.Step(
		`^I (try to )?add a new dimension to an existing filter$`,
		c.iAddANewDimensionToAnExistingFilter,
	)

	ctx.Step(
		`^I (try to )?add a new dimension with no options to an existing filter$`,
		c.iAddANewDimensionWithoutOptionsToAnExistingFilter,
	)

	ctx.Step(
		`^I receive the dimension\'s body back in the body of the response$`,
		c.iReceiveTheDimensionsBodyBackInTheBodyOfTheResponse,
	)

	ctx.Step(
		`^I receive the dimension\'s body back with an empty \'options\' slice$`,
		c.iReceiveTheDimensionWithEmptyOptionsSliceBackInTheBodyOfTheResponse,
	)

	ctx.Step(
		`^I try to add a malformed dimension to an existing filter$`,
		c.iTryToAddAMalformedDimensionToAnExistingFilter,
	)

	ctx.Step(
		`^I try to add a new dimension to a filter which has dimensions which were modified since I retrieved them$`,
		c.iTryToAddANewDimensionToAFilterWhichHasDimensionsWhichWereModifiedSinceIRetrievedThem,
	)

	ctx.Step(
		`^I try to add a new dimension to a non-existent filter$`,
		c.iTryToAddANewDimensionToANonexistentFilter,
	)

	ctx.Step(
		`^the client for the dataset API failed and is returning errors$`,
		c.theClientForTheDatasetAPIFailedAndIsReturningErrors,
	)

	ctx.Step(`^one event with the following fields are in the produced kafka topic catabular-export-start:$`, c.oneEventWithTheFollowingFieldsAreInTheProducedKafkaTopicCatabularexportstart)

	ctx.Step(`^I should receive the following time ignored JSON response:$`, c.iShouldReceiveTheFollowingTimeIgnoredJSONResponse)

	ctx.Step(`^I have these filter outputs:$`, c.iHaveTheseFilterOutputs)

}

/*
   The POST api now returns a Jobstate with a time component in it.
   This results in failed tests as example is not current.
   This custom JSON response parser is just to ensure time is ignored
   when comparing requests.

   open to suggestions on better way to handle
*/
func (c *Component) iShouldReceiveTheFollowingTimeIgnoredJSONResponse(arg1 *godog.DocString) error {
	now := time.Now()

	actual := new(bytes.Buffer)
	actual.ReadFrom(c.ApiFeature.HttpResponse.Body)

	var actualUpdateResponse api.UpdateFilterResponse
	var expectedUpdateResponse api.UpdateFilterResponse

	if err := json.Unmarshal([]byte(actual.String()), &actualUpdateResponse); err != nil {
		return errors.New("Failed to unmarshal ACTUAL response.")
	}

	if err := json.Unmarshal([]byte(arg1.Content), &expectedUpdateResponse); err != nil {
		return errors.New("Failed to unmarshal EXPECTED")
	}

	// naive setting because as per method, it should always be set.
	actualUpdateResponse.Events[0].Timestamp = now
	expectedUpdateResponse.Events[0].Timestamp = now

	if !reflect.DeepEqual(actualUpdateResponse, expectedUpdateResponse) {
		return errors.New("Structs are not equal")
	}
	return nil
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

		var e api.ExportStart
		var exportStart = `{
  "type": "record",
  "name": "cantabular-export-start",
  "fields": [
    {"name": "instance_id", "type": "string", "default": ""},
    {"name": "dataset_id",  "type": "string", "default": ""},
    {"name": "edition",     "type": "string", "default": ""},
    {"name": "version",     "type": "string", "default": ""},
    {"name": "filter_id",   "type":"string",  "default": ""}
  ]
}`

		var s = &avro.Schema{
			Definition: exportStart,
		}

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

func (c *Component) MongoDatastoreFailsForUpdateFilterOutput() error {
	var err error
	c.store, err = GetFailingMongo(c.ctx, c.cfg, c.g)
	if err != nil {
		return fmt.Errorf("failed to create new mongo mongoClient: %w", err)
	}

	return nil
}

func (c *Component) MongoDatastoreIsFailing() error {
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

const example_dimension = `{
	"name": "Number of siblings (3 mappings)",
	"is_area_type": false,
	"options": ["4-7", "7+"]
}`

func (c *Component) iAddANewDimensionToAnExistingFilter() error {
	c.postedJSON = example_dimension

	return c.ApiFeature.IPostToWithBody(
		"/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions",
		&godog.DocString{Content: c.postedJSON},
	)
}

func (c *Component) iAddANewDimensionWithoutOptionsToAnExistingFilter() error {
	c.postedJSON = `{
		"name": "Number of siblings (3 mappings)"
	}`

	return c.ApiFeature.IPostToWithBody(
		"/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions",
		&godog.DocString{Content: c.postedJSON},
	)
}

func (c *Component) iReceiveTheDimensionsBodyBackInTheBodyOfTheResponse() error {
	return c.ApiFeature.IShouldReceiveTheFollowingJSONResponse(&godog.DocString{Content: c.postedJSON})
}

func (c *Component) iReceiveTheDimensionWithEmptyOptionsSliceBackInTheBodyOfTheResponse() error {
	var want = `{
		"name": "Number of siblings (3 mappings)",
		"is_area_type": false,
		"options": []
	}`
	return c.ApiFeature.IShouldReceiveTheFollowingJSONResponse(&godog.DocString{Content: want})
}

func (c *Component) iTryToAddAMalformedDimensionToAnExistingFilter() error {
	return c.ApiFeature.IPostToWithBody(
		"/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions",
		&godog.DocString{Content: "not valid JSON"},
	)
}

func (c *Component) iTryToAddANewDimensionToAFilterWhichHasDimensionsWhichWereModifiedSinceIRetrievedThem() error {
	staleEtag := "a-stale-etag" // this would normally be a SHA-1 hash representing an older copy of the dimensions slice
	json := example_dimension
	c.ApiFeature.ISetTheHeaderTo("If-Match", staleEtag)
	return c.ApiFeature.IPostToWithBody(
		"/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions",
		&godog.DocString{Content: json},
	)
}

func (c *Component) iTryToAddANewDimensionToANonexistentFilter() error {
	json := example_dimension
	return c.ApiFeature.IPostToWithBody(
		"/filters/94310d8d-72d6-492a-bc30-000000000000/dimensions",
		&godog.DocString{Content: json},
	)
}

func (c *Component) theClientForTheDatasetAPIFailedAndIsReturningErrors() error {
	c.DatasetAPI.Reset()
	c.DatasetAPI.NewHandler().
		Get("/datasets/cantabular-example-1/editions/2021/versions/1").
		Reply(http.StatusInternalServerError)

	json := example_dimension
	return c.ApiFeature.IPostToWithBody(
		"/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions",
		&godog.DocString{Content: json},
	)
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
