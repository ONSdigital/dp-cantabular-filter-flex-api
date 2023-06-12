package steps

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/api"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/event"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/schema"
	"github.com/ONSdigital/log.go/v2/log"

	"github.com/cucumber/godog"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"github.com/rdumont/assistdog"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func (c *Component) RegisterSteps(ctx *godog.ScenarioContext) {
	c.AuthFeature.RegisterSteps(ctx)
	c.APIFeature.RegisterSteps(ctx)
	c.DatasetFeature.RegisterSteps(ctx)
	c.CantabularFeature.RegisterSteps(ctx)
	c.MetadataFeature.RegisterSteps(ctx)
	c.MongoFeature.RegisterSteps(ctx)
	c.PopulationFeature.RegisterSteps(ctx)

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
	ctx.Step(`^the getDatasetObservationResult result is:$`,
		c.theDatasetObservationResult,
	)
	ctx.Step(`^the getGeographyDatasetJSON result should be:$`,
		c.theGeographyDatasetJSONResult,
	)
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

// we are passing the string array as [xxxx,yyyy,zzz]
// this is required to support array being used in kafka messages
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
		return fmt.Errorf("+got -expected)\n%s", diff)
	}
	return nil
}

func (c *Component) theGeographyDatasetJSONResult(expected *godog.DocString) error {
	var got, expt api.GetDatasetJSONResponse

	b, err := io.ReadAll(c.APIFeature.HttpResponse.Body)
	if err != nil {
		return fmt.Errorf("Component::theGeographyDatasetJSONResult: error reading APIfeature response body: %w", err)
	}
	if err = json.Unmarshal(b, &got); err != nil {
		return fmt.Errorf("Component::theGeographyDatasetJSONResult error unmarshalling APIfeature response body: %w", err)
	}
	if err = json.Unmarshal([]byte(expected.Content), &expt); err != nil {
		return fmt.Errorf("Component::theGeographyDatasetJSONResult error unmarshalling 'expected' parameter: %w", err)
	}

	urlCompare := func(s1, s2 string) bool {
		if !strings.Contains(s1, "localhost:9999") && !strings.Contains(s2, "localhost:9999") {
			return s1 == s2
		}
		s1URL, e := url.Parse(s1)
		if e != nil {
			return false
		}
		s2URL, e := url.Parse(s2)
		if e != nil {
			return false
		}

		return s1URL.Path == s2URL.Path
	}

	assert.Empty(c, cmp.Diff(got, expt, cmp.Comparer(urlCompare)))

	return c.StepError()
}

func (c *Component) theDatasetObservationResult(expected *godog.DocString) error {
	var got, expt api.GetObservationsResponse

	b, err := io.ReadAll(c.APIFeature.HttpResponse.Body)
	if err != nil {
		return fmt.Errorf("Component::theDatasetObservationResult: error reading APIfeature response body: %w", err)
	}
	if err = json.Unmarshal(b, &got); err != nil {
		return fmt.Errorf("Component::theDatasetObservationResult error unmarshalling APIfeature response body: %w", err)
	}
	if err = json.Unmarshal([]byte(expected.Content), &expt); err != nil {
		return fmt.Errorf("Component::theDatasetObservationResult error unmarshalling 'expected' parameter: %w", err)
	}

	urlCompare := func(s1, s2 string) bool {
		if !strings.Contains(s1, "localhost:9999") && !strings.Contains(s2, "localhost:9999") {
			return s1 == s2
		}
		s1URL, e := url.Parse(s1)
		if e != nil {
			return false
		}
		s2URL, e := url.Parse(s2)
		if e != nil {
			return false
		}

		return s1URL.Path == s2URL.Path
	}

	assert.Empty(c, cmp.Diff(got, expt, cmp.Comparer(urlCompare)))

	return c.StepError()
}
