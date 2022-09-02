package steps

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/event"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/schema"
	"github.com/ONSdigital/log.go/v2/log"

	"github.com/cucumber/godog"
	"github.com/google/go-cmp/cmp"
	"github.com/rdumont/assistdog"
	"go.mongodb.org/mongo-driver/bson"
)

func (c *Component) RegisterSteps(ctx *godog.ScenarioContext) {
	c.AuthFeature.RegisterSteps(ctx)
	c.APIFeature.RegisterSteps(ctx)
	c.DatasetFeature.RegisterSteps(ctx)
	c.CantabularFeature.RegisterSteps(ctx)
	c.MongoFeature.RegisterSteps(ctx)

	ctx.Step(
		`^private endpoints are enabled$`,
		c.privateEndpointsAreEnabled,
	)
	ctx.Step(
		`^private endpoints are enabled with permissions checking$`,
		c.privateEndpointsAreEnabledWithPermissions,
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
		`^I provide If-Match header "([^"]*)"$`,
		c.iProvideIfMatchHeader,
	)
	ctx.Step(
		`^I should receive an errors array`,
		c.iShouldReceiveAnErrorsArray,
	)
	ctx.Step(`^an ETag is returned`,
		c.anETagIsReturned,
	)
	ctx.Step(`^the ETag is a hash of the filter "([^"]*)"`,
		c.theETagIsAHashOfTheFilter,
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
		`^Cantabular returns dimensions for the dataset "([^"]*)" for the following search terms:$`,
		c.cantabularReturnsMultipleDimensions,
	)
	ctx.Step(
		`^Cantabular responds with an error$`,
		c.cantabularRespondsWithAnError,
	)
	ctx.Step(`^the filter output with the following structure is in the datastore:$`,
		c.filterOutputIsInDatastore,
	)

}
func (c *Component) filterOutputIsInDatastore(expectedOutput *godog.DocString) error {
	var expected model.FilterOutput

	err := json.Unmarshal([]byte(expectedOutput.Content), &expected)
	if err != nil {
		return fmt.Errorf("failed to unmarshall provided filterOutput: %w", err)
	}

	actual, err := c.store.GetFilterOutput(c.ctx, expected.ID)
	if err != nil {
		return fmt.Errorf("Error encountered while retrieving filter output: %w", err)
	}

	if diff := cmp.Diff(actual, &expected); diff != "" {
		return fmt.Errorf("-got +expected)\n%s\n", diff)
	}
	return nil
}

// iShouldReceiveAnErrorsArray checks that the response body can be deserialized into
// an error response, and contains at least one error.
func (c *Component) iShouldReceiveAnErrorsArray() error {
	responseBody := c.APIFeature.HttpResponse.Body

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

//we are passing the string array as [xxxx,yyyy,zzz]
//this is required to support array being used in kafka messages
func arrayParser(raw string) (interface{}, error) {
	//remove the starting and trailing brackets
	str := strings.Trim(raw, "[]")
	if str == "" {
		return []string{}, nil
	}

	strArray := strings.Split(str, ",")
	return strArray, nil
}

func (c *Component) anETagIsReturned() error {
	eTag := c.APIFeature.HttpResponse.Header.Get("ETag")
	if eTag == "" {
		return fmt.Errorf("no 'ETag' header returned")
	}
	return nil
}

// theETagIsAHashOfTheFilter checks that the returned ETag header (and stored ETag field)
// are a hash of a filter. Used to validate that the ETag was updated after a mutation.
func (c *Component) theETagIsAHashOfTheFilter(filterID string) error {
	eTag := c.APIFeature.HttpResponse.Header.Get("ETag")
	if eTag == "" {
		return errors.New("expected ETag")
	}

	col := c.svc.Cfg.FiltersCollection
	db := c.svc.Cfg.Mongo.Database

	var response model.Filter
	if err := c.MongoFeature.Client.Database(db).Collection(col).FindOne(context.Background(), bson.M{"filter_id": filterID}).Decode(&response); err != nil {
		return fmt.Errorf("failed to retrieve filter: %w", err)
	}

	hash, err := response.Hash(nil)
	if err != nil {
		return fmt.Errorf("unable to hash stored filter: %w", err)
	}

	if eTag != hash {
		return fmt.Errorf("ETag header did not match, expected %s, got %s", hash, eTag)
	}

	if eTag != response.ETag {
		return fmt.Errorf("ETag on stored filter did not match, expected %s, got %s", hash, eTag)
	}

	return nil
}

func (c *Component) privateEndpointsAreEnabled() error {
	c.svc.Cfg.EnablePrivateEndpoints = true
	return nil
}

func (c *Component) privateEndpointsAreEnabledWithPermissions() error {
	c.svc.Cfg.EnablePrivateEndpoints = true
	c.svc.Cfg.EnablePermissionsAuth = true
	return nil
}

func (c *Component) privateEndpointsAreNotEnabled() error {
	c.svc.Cfg.EnablePrivateEndpoints = false
	c.svc.Cfg.EnablePermissionsAuth = false
	return nil
}

func (c *Component) permissionsCheckingIsNotEnabled() error {
	c.svc.Cfg.EnablePermissionsAuth = false
	return nil
}

func (c *Component) theMaximumLimitIsSetTo(val int) error {
	c.svc.Cfg.DefaultMaximumLimit = val
	return nil
}

func (c *Component) iProvideIfMatchHeader(eTag string) error {
	return c.APIFeature.ISetTheHeaderTo("If-Match", eTag)
}

func (c *Component) theFollowingExportStartEventsAreProduced(events *godog.Table) error {
	assist := assistdog.NewDefault()
	assist.RegisterParser([]string{}, arrayParser)
	expected, err := assist.CreateSlice(new(event.ExportStart), events)
	if err != nil {
		return fmt.Errorf("failed to create slice from godog table: %w", err)
	}

	ctx := context.Background()
	consumer, err := GenerateKafkaConsumer(ctx)
	if err != nil {
		return fmt.Errorf("failed to generate kafka consumer: %w", err)
	}

	var got []*event.ExportStart
	listen := true
	for listen {
		select {
		case <-time.After(WaitEventTimeout):
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

	if err := consumer.Close(ctx); err != nil {
		// just log the error, but do not fail the test
		// as it is not relevant to this test.
		log.Error(ctx, "error closing kafka consumer", err)
	}
	if diff := cmp.Diff(expected, got); diff != "" {
		return fmt.Errorf("+got -expected)\n%s\n", diff)
	}
	return nil
}
