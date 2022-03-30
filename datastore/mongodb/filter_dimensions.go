package mongodb

import (
	"context"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type dimensionQueryResult struct {
	Dimensions []model.Dimension `bson:"dimensions"`
	TotalCount int               `bson:"totalCount"`
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
		return nil, 0, errors.Wrap(err, "failed to get filter dimensions")
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

func (c *Client) AddFilterDimension(ctx context.Context, fID string, dim model.Dimension) error {
	col := c.collections.filters

	q := map[string]model.Dimension{
		"dimensions": dim,
	}

	if _, err := c.conn.Collection(col.name).Update(ctx, bson.M{"filter_id": fID}, bson.M{"$push": q}); err != nil {
		return errors.Wrap(err, "failed to add dimension to filter")
	}

	return nil
}
