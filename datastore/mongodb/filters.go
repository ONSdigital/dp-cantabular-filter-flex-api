package mongodb

import (
	"context"
	"fmt"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	"github.com/ONSdigital/dp-mongodb/v3/mongodb"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateFilter creates a new Filter in the CantabularFilters collection
func (c *Client) CreateFilter(ctx context.Context, f *model.Filter) error {
	var err error

	id, err := c.generate.UUID()
	if err != nil {
		return errors.Wrap(err, "failed to generate UUID: %w")
	}

	f.ID = id.String()

	if f.ETag, err = f.Hash(nil); err != nil {
		return errors.Wrap(err, "failed to generate eTag: %w")
	}

	f.Links.Self = model.Link{
		HREF: fmt.Sprintf("%s/filters/%s", c.cfg.FilterFlexAPIURL, f.ID),
	}

	col := c.collections.filters

	lockID, err := col.lock(ctx, f.ID)
	if err != nil {
		return err
	}
	defer col.unlock(ctx, lockID)

	if _, err = c.conn.Collection(col.name).UpsertById(ctx, f.ID, bson.M{"$set": f}); err != nil {
		return errors.Wrap(err, "failed to upsert filter")
	}

	return nil
}

// GetFilter gets a filter doc from the filters collections
func (c *Client) GetFilter(ctx context.Context, fID string) (*model.Filter, error) {
	var err error

	col := c.collections.filters

	var f model.Filter

	if err = c.conn.Collection(col.name).FindOne(ctx, bson.M{"filter_id": fID}, &f); err != nil {
		err := &er{
			err: errors.Wrap(err, "failed to find filter"),
		}
		if errors.Is(err, mongodb.ErrNoDocumentFound) {
			err.notFound = true
		}
		return nil, err
	}

	return &f, nil
}

// GetFilterDimensions gets the dimensions for a Filter in the CantabularFilters collection
func (c *Client) GetFilterDimensions(ctx context.Context, fID string, limit, offset int) ([]model.Dimension, int, error) {
	col := c.collections.filters

	projection := bson.D{
		{"_id", 0},
		{"totalCount", "$totalCount"},
	}

	// We can't use `$slice` with a value of 0 , so if we know we're not returning results,
	// don't include the `dimensions` projection.
	if limit > 0 {
		projection = append(projection, bson.E{
			Key:   "dimensions",
			Value: bson.D{{"$slice", bson.A{"$dimensions", offset, limit}}},
		})
	}

	pipeline := mongo.Pipeline{
		bson.D{{"$match", bson.D{{"filter_id", fID}}}},
		bson.D{{"$unwind", "$dimensions"}},
		bson.D{{"$replaceRoot", bson.D{{"newRoot", "$dimensions"}}}},
		bson.D{{"$sort", bson.D{{"name", 1}}}},
		bson.D{{
			"$group",
			bson.D{
				{"_id", 0},
				{"dimensions", bson.D{{"$push", "$$ROOT"}}},
				{"totalCount", bson.D{{"$sum", 1}}},
			},
		}},
		bson.D{{"$project", projection}},
	}

	var results []dimensionQueryResult
	if err := c.conn.Collection(col.name).Aggregate(ctx, pipeline, &results); err != nil {
		return nil, 0, &er{
			err: errors.Wrap(err, "failed to get filter dimensions"),
		}
	}

	if len(results) == 0 {
		return nil, 0, &er{
			err:      errors.Errorf("failed to find filter with ID (%s)", fID),
			notFound: true,
		}
	}

	result := results[0]

	// Ensure we return an empty slice, so it serializes into `[]`.
	if result.Dimensions == nil {
		result.Dimensions = []model.Dimension{}
	}

	return result.Dimensions, result.TotalCount, nil
}

type dimensionQueryResult struct {
	Dimensions []model.Dimension `bson:"dimensions"`
	TotalCount int               `bson:"totalCount"`
}

func (c *Client) AddFilterDimension(ctx context.Context, fID string, dimension model.Dimension) error {
	col := c.collections.filters

	newDimension := map[string]model.Dimension{
		"dimensions": dimension,
	}

	if _, err := c.conn.Collection(col.name).Update(ctx, bson.M{"filter_id": fID}, bson.M{"$push": newDimension}); err != nil {
		return errors.Wrap(err, "failed to add dimension to filter")
	}
	return nil
}
