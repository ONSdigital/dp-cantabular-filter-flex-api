package steps

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/recipe"

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
	ctx.Step(`^the following recipe is used to create a dataset based on the given cantabular dataset:$`,
		c.theFollowingRecipeIsUsed,
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

func (c *Component) theFollowingRecipeIsUsed(rec *godog.DocString) error {
	var (
		payload struct {
			R   recipe.Recipe                 `json:"recipe"`
			CDS cantabular.StaticDatasetQuery `json:"cantabular_dataset"`
		}
		ds                 dataset.DatasetDetails
		v                  dataset.Version
		dims               dataset.VersionDimensions
		dsj, vj, mj, dimsj []byte
		dimVars            []string
		err                error
	)

	err = json.Unmarshal([]byte(rec.Content), &payload)
	if err != nil {
		return errors.Wrap(err, "error unmarshalling recipe/dataset in Component::theFollowingRecipeIsUsed")
	}

	ds = c.makeDataset(payload.R)
	dsj, err = json.Marshal(ds)
	if err != nil {
		return errors.Wrap(err, "error marshalling dataset in Component::theFollowingRecipeIsUsed")
	}
	err = c.DatasetFeature.theFollowingDatasetDocumentIsAvailable(ds.ID, &godog.DocString{Content: string(dsj)})
	if err != nil {
		return errors.Wrap(err, "error setting dataset into DatasetFeature in Component::theFollowingRecipeIsUsed")
	}

	for _, d := range payload.R.OutputInstances[0].CodeLists {
		dims.Items = append(dims.Items, dataset.VersionDimension{ID: d.ID, Name: d.ID, Label: d.Name, Variable: d.ID, Links: dataset.Links{CodeList: dataset.Link{URL: d.HRef, ID: d.ID}}})
		dimVars = append(dimVars, d.ID)
	}

	v = c.makeVersion(ds, dims)
	vj, err = json.Marshal(v)
	if err != nil {
		return errors.Wrap(err, "error marshalling version in Component::theFollowingRecipeIsUsed")
	}
	err = c.DatasetFeature.theFollowingVersionDocumentIsAvailable(ds.ID, v.Edition, v.Links.Version.ID, &godog.DocString{Content: string(vj)})
	if err != nil {
		return errors.Wrap(err, "error setting version into DatasetFeature in Component::theFollowingRecipeIsUsed")
	}

	m := c.makeMetadata(v)
	mj, err = json.Marshal(m)
	if err != nil {
		return errors.Wrap(err, "error marshalling metadata in Component::theFollowingRecipeIsUsed")
	}
	err = c.DatasetFeature.theFollowingMetadataDocumentIsAvailable(ds.ID, v.Edition, v.Links.Version.ID, &godog.DocString{Content: string(mj)})
	if err != nil {
		return errors.Wrap(err, "error setting metadata into DatasetFeature in Component::theFollowingRecipeIsUsed")
	}

	dimsj, err = json.Marshal(dims)
	if err != nil {
		return errors.Wrap(err, "error marshalling dimensions in Component::theFollowingRecipeIsUsed")
	}
	err = c.DatasetFeature.theFollowingDimensionsDocumentIsAvailable(ds.ID, v.Edition, v.Links.Version.ID, &godog.DocString{Content: string(dimsj)})
	if err != nil {
		return errors.Wrap(err, "error setting dimensions into DatasetFeature in Component::theFollowingRecipeIsUsed")
	}

	_ = c.makeDimensionOptions(v, payload.CDS.Dataset.Table.Dimensions)

	return nil
}

func (c *Component) makeDataset(recipe recipe.Recipe) dataset.DatasetDetails {
	recipeDetails := recipe.OutputInstances[0]

	return dataset.DatasetDetails{ID: recipeDetails.DatasetID,
		Title:     recipeDetails.Title,
		Type:      recipe.Format,
		IsBasedOn: &dataset.IsBasedOn{ID: recipe.CantabularBlob, Type: recipe.Format},
		State:     "published",
		Links: dataset.Links{
			Self: dataset.Link{
				ID:  recipeDetails.DatasetID,
				URL: fmt.Sprintf("http://hostname/datasets/%s", recipeDetails.DatasetID),
			},
		},
	}
}

func (c *Component) makeVersion(ds dataset.DatasetDetails, dims dataset.VersionDimensions) dataset.Version {
	return dataset.Version{
		ID:         ds.ID + "UUID",
		Edition:    "latest",
		Version:    1,
		State:      "published",
		IsBasedOn:  ds.IsBasedOn,
		Dimensions: dims.Items,
		Links: dataset.Links{
			Dataset: ds.Links.Self,
			Version: dataset.Link{
				ID:  "1",
				URL: fmt.Sprintf("%s/editions/%s/versions/%d", ds.Links.Self.URL, "latest", 1),
			},
			Self: dataset.Link{
				ID:  "1",
				URL: fmt.Sprintf("%s/editions/%s/versions/%d", ds.Links.Self.URL, "latest", 1),
			},
		},
	}
}

func (c *Component) makeMetadata(v dataset.Version) interface{} {
	// logic referenced from dp-dataset-api/metadata/models/metadata.go CreateCantabularMetaDataDoc() (commit hash dd295b15704d87c0cb08d127ca108ce2f8df05a7)
	// no metadata links are being set at present in this function, but it is envisaged the links will be set at some point, so I am setting a basic link in
	//this test. The returned value is an extremely cutdown version of dp-dataset-api/metadata/models/Metadata

	return struct {
		Links interface{} `json:"links,omitempty"`
	}{Links: struct {
		Self interface{} `json:"self,omitempty"`
	}{Self: struct {
		HRef string `json:"href,omitempty"`
		ID   string `json:"id,omitempty"`
	}{HRef: v.Links.Self.URL + "/metadata"}}}
}

func (c *Component) makeDimensionOptions(v dataset.Version, dims []cantabular.Dimension) []dataset.Option {
	if len(v.Dimensions) != len(dims) {
		c.Fatalf("dimension lists  not of same length in Component::makeDimensionOptions")
	}
	sort.Slice(v.Dimensions, func(i, j int) bool { return v.Dimensions[i].ID < v.Dimensions[j].ID })
	sort.Slice(dims, func(i, j int) bool { return dims[i].Variable.Name < dims[j].Variable.Name })

	var dimOpts []dataset.Option
	for i, d := range v.Dimensions {

		var opts []dataset.Option
		for _, o := range dims[i].Categories {
			opts = append(opts, dataset.Option{
				DimensionID: d.ID,
				Label:       o.Label,
				Option:      o.Code,
				Links: dataset.Links{
					Versions: v.Links.Version,
					CodeList: d.Links.CodeList,
					Code: dataset.Link{
						ID:  o.Code,
						URL: fmt.Sprintf("%s/codes/%s", d.Links.CodeList.URL, o.Code),
					},
				},
			})
		}

		optsj, err := json.Marshal(dataset.Options{Items: opts, Count: len(opts), TotalCount: len(opts)})
		if err != nil {
			c.Fatalf("error marshalling dimension options in Component::theFollowingRecipeIsUsed: %v", err)
		}
		if err = c.DatasetFeature.theFollowingOptionsDocumentIsAvailable(v.Links.Dataset.ID, v.Edition, v.Links.Version.ID, d.ID, &godog.DocString{Content: string(optsj)}); err != nil {
			c.Fatalf("error setting dimensions into DatasetFeature in Component::theFollowingRecipeIsUsed: %v", err)
		}

		dimOpts = append(dimOpts, opts...)
	}

	return dimOpts
}
