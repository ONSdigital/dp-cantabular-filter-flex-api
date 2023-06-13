package steps

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/config"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/datastore/mongodb"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/service"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/service/mock"
	componenttest "github.com/ONSdigital/dp-component-test"

	"github.com/cucumber/godog"
	"github.com/google/go-cmp/cmp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoVersion = "4.4.8"
)

type MongoFeature struct {
	componenttest.ErrorFeature
	*componenttest.MongoFeature

	g   service.Generator
	cfg *config.Config
}

func NewMongoFeature(ef componenttest.ErrorFeature, g service.Generator, cfg *config.Config) *MongoFeature {
	mf := &MongoFeature{
		ErrorFeature: ef,
		MongoFeature: componenttest.NewMongoFeature(componenttest.MongoOptions{
			MongoVersion: mongoVersion,
		}),
		g:   g,
		cfg: cfg}

	mf.cfg.Mongo.ClusterEndpoint = mf.MongoFeature.Server.URI()

	return mf
}

func (mf *MongoFeature) Reset() {
	if err := mf.Client.Database(mf.cfg.Mongo.Database).Drop(context.Background()); err != nil {
		mf.Fatalf("failed to reset mongo (error dropping database): %v", err)
	}

	mf.setWorkingMongo()
}

func (mf *MongoFeature) Close() {
	if err := mf.MongoFeature.Close(); err != nil {
		mf.Fatalf("failed to close mongo: %v", err)
	}
}

func (mf *MongoFeature) setWorkingMongo() {
	service.GetMongoDB = func(ctx context.Context, cfg *config.Config, g service.Generator) (service.Datastore, error) {
		return mongodb.NewClient(ctx, g, mongodb.Config{
			MongoDriverConfig:       cfg.Mongo,
			FilterAPIURL:            cfg.BindAddr,
			FiltersCollection:       cfg.FiltersCollection,
			FilterOutputsCollection: cfg.FilterOutputsCollection,
		})
	}
}

func (mf *MongoFeature) setFailingMongo() {
	service.GetMongoDB = func(ctx context.Context, cfg *config.Config, g service.Generator) (service.Datastore, error) {
		return &mock.DatastoreMock{
			UpdateFilterOutputFunc: func(_ context.Context, _ *model.FilterOutput) error {
				return errors.New("failed to upsert filter")
			},
			GetFilterOutputFunc: func(_ context.Context, _ string) (*model.FilterOutput, error) {
				return nil, errors.New("mongo client has failed")
			},
			AddFilterOutputEventFunc: func(_ context.Context, _ string, _ *model.Event) error {
				return errors.New("failed to add event")
			},
			GetFilterDimensionOptionsFunc: func(_ context.Context, _, _ string, _, _ int) ([]string, int, string, error) {
				return nil, 0, "", errors.New("error that should not be returned to user")
			},
			UpdateFilterDimensionFunc: func(_ context.Context, _, _ string, _ model.Dimension, _ string) (string, error) {
				return "", errors.New("failed to update filter dimension")
			},
			RemoveFilterDimensionOptionFunc: func(_ context.Context, _, _, _, _ string) (string, error) {
				return "", errors.New("failed to remove filter dimension option")
			},
			GetFilterFunc: func(_ context.Context, _ string) (*model.Filter, error) {
				return nil, errors.New("failed to get filter")
			},
		}, nil
	}
}

func (mf *MongoFeature) RegisterSteps(ctx *godog.ScenarioContext) {
	mf.MongoFeature.RegisterSteps(ctx)

	ctx.Step(
		`^Mongo datastore fails for update filter output`,
		mf.MongoDatastoreFailsForUpdateFilterOutput,
	)

	// this is the same as above, but added for clearer step definition
	ctx.Step(
		`^Mongo datastore is failing`,
		mf.MongoDatastoreIsFailing,
	)
	ctx.Step(
		`^a document in collection "([^"]*)" with key "([^"]*)" value "([^"]*)" should match:$`,
		mf.aDocumentInCollectionWithKeyValueShouldMatch,
	)

	ctx.Step(
		`^a document in collection "([^"]*)" with key "([^"]*)" value "([^"]*)" has empty "([^"]*)" options`,
		mf.theFilterHasEmptyOptions,
	)
	ctx.Step(
		`^I have these filters:$`,
		mf.iHaveTheseFilters,
	)
	ctx.Step(`^I have these filter outputs:$`,
		mf.iHaveTheseFilterOutputs,
	)
	ctx.Step(
		`^I have this filter with an ETag of "([^"]*)":$`,
		mf.iHaveThisFilterWithETag,
	)
	ctx.Step(`^the filter output with the following structure is in the datastore:$`,
		mf.filterOutputIsInDatastore,
	)
}

func (mf *MongoFeature) MongoDatastoreFailsForUpdateFilterOutput() {
	mf.setFailingMongo()
}

func (mf *MongoFeature) MongoDatastoreIsFailing() {
	mf.setFailingMongo()
}

func (mf *MongoFeature) theFilterHasEmptyOptions(col, key, val, dimensionName string) error {
	var filter model.Filter
	if err := mf.Client.Database(mf.cfg.Mongo.Database).Collection(col).FindOne(context.Background(), bson.M{key: val}).Decode(&filter); err != nil {
		return fmt.Errorf("failed to retrieve document: %w", err)
	}

	for i := range filter.Dimensions {
		dimension := filter.Dimensions[i]
		if dimension.Name == dimensionName {
			if len(dimension.Options) == 0 {
				return nil
			}
		}
	}

	return errors.New("option either not found or is not empty")
}

func (mf *MongoFeature) aDocumentInCollectionWithKeyValueShouldMatch(col, key, val string, doc *godog.DocString) error {
	var expected, result interface{}
	if err := json.Unmarshal([]byte(doc.Content), &expected); err != nil {
		return fmt.Errorf("failed to unmarshal document: %w", err)
	}

	var bdoc primitive.D
	if err := mf.Client.Database(mf.cfg.Mongo.Database).Collection(col).FindOne(context.Background(), bson.M{key: val}).Decode(&bdoc); err != nil {
		return fmt.Errorf("failed to retrieve document: %w", err)
	}

	b, err := bson.MarshalExtJSON(bdoc, true, true)
	if err != nil {
		return fmt.Errorf("failed to marshal bson document: %w", err)
	}

	if err := json.Unmarshal(b, &result); err != nil {
		return fmt.Errorf("failed to unmarshal result: %w", err)
	}

	if diff := cmp.Diff(expected, result); diff != "" {
		return fmt.Errorf("-expected +got)\n%s", diff)
	}

	return nil
}

func (mf *MongoFeature) iHaveTheseFilters(docs *godog.DocString) error {
	var filters []model.Filter

	err := json.Unmarshal([]byte(docs.Content), &filters)
	if err != nil {
		return fmt.Errorf("failed to unmarshal filter: %w", err)
	}

	if err := mf.insertFilters(filters); err != nil {
		return fmt.Errorf("error inserting filters: %w", err)
	}

	return nil
}

// iHaveThisFilterWithETag inserts the provided filter into the database,
// setting the ETag to the provided stub value.
func (mf *MongoFeature) iHaveThisFilterWithETag(eTag string, docs *godog.DocString) error {
	var filter model.Filter

	err := json.Unmarshal([]byte(docs.Content), &filter)
	if err != nil {
		return fmt.Errorf("failed to unmarshal filter: %w", err)
	}

	filter.ETag = eTag
	filter.LastUpdated = mf.g.Timestamp()
	filter.UniqueTimestamp = mf.g.UniqueTimestamp()

	if err := mf.insertFilters([]model.Filter{filter}); err != nil {
		return fmt.Errorf("failed to insert filter with ETag: %w", err)
	}

	return nil
}

func (mf *MongoFeature) iHaveTheseFilterOutputs(docs *godog.DocString) error {
	var filterOutputs []model.FilterOutput
	err := json.Unmarshal([]byte(docs.Content), &filterOutputs)
	if err != nil {
		return fmt.Errorf("failed to unmarshal filter output: %w", err)
	}

	db := mf.cfg.Mongo.Database
	col := mf.cfg.FilterOutputsCollection

	upsert := true
	for i := range filterOutputs {
		f := filterOutputs[i]
		if _, err := mf.Client.Database(db).Collection(col).UpdateByID(context.Background(), f.ID, bson.M{"$set": f}, &options.UpdateOptions{Upsert: &upsert}); err != nil {
			return fmt.Errorf("failed to upsert filter output: %w", err)
		}
	}

	return nil
}

func (mf *MongoFeature) filterOutputIsInDatastore(expectedOutput *godog.DocString) error {
	ctx := context.Background()
	col := mf.cfg.FilterOutputsCollection
	db := mf.cfg.Mongo.Database

	var actual, expected model.FilterOutput
	err := json.Unmarshal([]byte(expectedOutput.Content), &expected)
	if err != nil {
		return fmt.Errorf("failed to unmarshall provided filterOutput: %w", err)
	}

	err = mf.Client.Database(db).Collection(col).FindOne(ctx, bson.M{"id": expected.ID}).Decode(&actual)
	if err != nil {
		return fmt.Errorf("error encountered while retrieving filter output: %w", err)
	}

	if diff := cmp.Diff(actual, expected); diff != "" {
		return fmt.Errorf("-got +expected)\n%s", diff)
	}
	return nil
}

func (mf *MongoFeature) insertFilters(filters []model.Filter) error {
	ctx := context.Background()
	db := mf.cfg.Mongo.Database
	col := mf.cfg.FiltersCollection

	upsert := true
	for i := range filters {
		filter := filters[i]
		if _, err := mf.Client.Database(db).Collection(col).UpdateByID(ctx, filter.ID, bson.M{"$set": filter}, &options.UpdateOptions{Upsert: &upsert}); err != nil {
			return fmt.Errorf("failed to upsert filter: %w", err)
		}
	}

	return nil
}

func (mf *MongoFeature) setInitialiserMock() {
	// the default initialiser is fine as it will pick up the mongo server url from the config
}
