package steps

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular/gql"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/event"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/schema"
	"github.com/ONSdigital/log.go/v2/log"

	"github.com/cucumber/godog"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"github.com/rdumont/assistdog"
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
		`^I have this filter with an ETag of "([^"]*)":$`,
		c.iHaveThisFilterWithETag,
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
	ctx.Step(`^the ETag is a hash of the filter "([^"]*)"`,
		c.theETagIsAHashOfTheFilter,
	)
	ctx.Step(
		`^I provide If-Match header "([^"]*)"$`,
		c.iProvideIfMatchHeader,
	)
	ctx.Step(
		`^the client for the dataset API failed and is returning errors$`,
		c.theClientForTheDatasetAPIFailedAndIsReturningErrors,
	)
	ctx.Step(`^the following Export Start events are produced:$`,
		c.theFollowingExportStartEventsAreProduced,
	)
	ctx.Step(`^I have these filter outputs:$`,
		c.iHaveTheseFilterOutputs,
	)
	ctx.Step(
		`^I should receive an errors array`,
		c.iShouldReceiveAnErrorsArray,
	)

	ctx.Step(
		`^Cantabular returns these dimensions for the dataset "([^"]*)" and search term "([^"]*)":$`,
		c.cantabularSearchReturnsTheseDimensions,
	)

	ctx.Step(
		`^Cantabular responds with an error$`,
		c.cantabularRespondsWithAnError,
	)

	ctx.Step(`^the filter output with the id "([^"]*)" is in the datastore`,
		c.filterOutputIsInDatastore)

}
func (c *Component) filterOutputIsInDatastore(id string) error {

	_, err := c.store.GetFilterOutput(c.ctx, id)
	if err != nil {
		return fmt.Errorf("Error encountered while retrieving filter output.")
	}

	return nil

}

// iShouldReceiveAnErrorsArray checks that the response body can be deserialized into
// an error response, and contains at least one error.
func (c *Component) iShouldReceiveAnErrorsArray() error {
	responseBody := c.ApiFeature.HttpResponse.Body

	var errorResponse struct {
		Errors []string `json:"errors"`
	}

	if err := json.NewDecoder(responseBody).Decode(&errorResponse); err != nil {
		return fmt.Errorf("failed to decode error response from body: %w", err)
	}

	if len(errorResponse.Errors) == 0 {
		return errors.New("expected at least one error in response")
	}

	return nil
}

func (c *Component) theFollowingExportStartEventsAreProduced(events *godog.Table) error {
	expected, err := assistdog.NewDefault().CreateSlice(new(event.ExportStart), events)
	if err != nil {
		return fmt.Errorf("failed to create slice from godog table: %w", err)
	}

	consumer, err := GenerateKafkaConsumer(c.ctx)
	if err != nil {
		return fmt.Errorf("failed to generate kafka consumer: %w", err)
	}

	var got []*event.ExportStart
	listen := true
	for listen {
		select {
		case <-time.After(c.waitEventTimeout):
			listen = false
		case <-consumer.Channels().Closer:
			return errors.New("closer channel closed")
		case msg, ok := <-consumer.Channels().Upstream:
			if !ok {
				return errors.New("upstream channel closed")
			}

			var e event.ExportStart
			var s = schema.ExportStart

			if err := s.Unmarshal(msg.GetData(), &e); err != nil {
				msg.Commit()
				msg.Release()
				return fmt.Errorf("error unmarshalling message: %w", err)
			}

			msg.Commit()
			msg.Release()

			got = append(got, &e)
		}

	}

	if err := consumer.Close(c.ctx); err != nil {
		// just log the error, but do not fail the test
		// as it is not relevant to this test.
		log.Error(c.ctx, "error closing kafka consumer", err)
	}

	if diff := cmp.Diff(got, expected); diff != "" {
		return fmt.Errorf("-got +expected)\n%s\n", diff)
	}
	return nil
}

func (c *Component) anETagIsReturned() error {
	eTag := c.ApiFeature.HttpResponse.Header.Get("ETag")
	if eTag == "" {
		return fmt.Errorf("no 'ETag' header returned")
	}
	return nil
}

// theETagIsAHashOfTheFilter checks that the returned ETag header (and stored ETag field)
// are a hash of a filter. Used to validate that the ETag was updated after a mutation.
func (c *Component) theETagIsAHashOfTheFilter(filterID string) error {
	eTag := c.ApiFeature.HttpResponse.Header.Get("ETag")
	if eTag == "" {
		return errors.New("expected ETag")
	}

	ctx := context.Background()
	col := c.cfg.FiltersCollection

	var response model.Filter
	if err := c.store.Conn().Collection(col).FindOne(ctx, bson.M{"filter_id": filterID}, &response); err != nil {
		return errors.Wrap(err, "failed to retrieve filter")
	}

	hash, err := response.Hash(nil)
	if err != nil {
		return errors.Wrap(err, "unable to hash stored filter")
	}

	if eTag != hash {
		return errors.Wrapf(err, "ETag header did not match, expected %s, got %s", hash, eTag)
	}

	if eTag != response.ETag {
		return fmt.Errorf("ETag on stored filter did not match, expected %s, got %s", hash, eTag)
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
	var filters []model.Filter

	err := json.Unmarshal([]byte(docs.Content), &filters)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshall")
	}

	if err := c.insertFilters(filters); err != nil {
		return errors.Wrap(err, "error inserting filters")
	}

	return nil
}

// iHaveThisFilterWithETag inserts the provided filter into the database,
// setting the ETag to the provided stub value.
func (c *Component) iHaveThisFilterWithETag(eTag string, docs *godog.DocString) error {
	var filter model.Filter

	err := json.Unmarshal([]byte(docs.Content), &filter)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal")
	}

	filter.ETag = eTag

	if err := c.insertFilters([]model.Filter{filter}); err != nil {
		return errors.Wrap(err, "failed to insert filter with ETag")
	}

	return nil
}

// cantabularSearchReturnsTheseDimensions sets up a stub response for the `SearchDimensions` method.
func (c *Component) cantabularSearchReturnsTheseDimensions(datasetID, dimension string, docs *godog.DocString) error {
	var response cantabular.GetDimensionsResponse
	if err := json.Unmarshal([]byte(docs.Content), &response); err != nil {
		return errors.Wrap(err, "unable to unmarshal cantabular search response")
	}

	c.MockCantabularClient.SearchDimensionsFunc = func(ctx context.Context, req cantabular.SearchDimensionsRequest) (*cantabular.GetDimensionsResponse, error) {
		if req.Dataset == datasetID && req.Text == dimension {
			return &response, nil
		}

		return &cantabular.GetDimensionsResponse{
			Dataset: gql.DatasetVariables{
				Variables: gql.Variables{
					Search: gql.Search{
						Edges: []gql.Edge{},
					},
				},
			},
		}, nil
	}

	return nil
}

// cantabularSearchRespondsWithAnError sets up a generic error response for the
func (c *Component) cantabularRespondsWithAnError() {
	c.MockCantabularClient.OptionsHappy = false
}

// insertFilters loops through the provided filters and inserts them into the database.
func (c *Component) insertFilters(filters []model.Filter) error {
	ctx := context.Background()
	store := c.store
	col := c.cfg.FiltersCollection

	for _, filter := range filters {
		if _, err := store.Conn().Collection(col).UpsertById(ctx, filter.ID, bson.M{"$set": filter}); err != nil {
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
