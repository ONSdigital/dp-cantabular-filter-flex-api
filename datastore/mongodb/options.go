package mongodb

import (
	"context"
	"fmt"

	"github.com/ONSdigital/log.go/v2/log"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type optionQueryResult struct {
	Options    []string `bson:"options"`
	TotalCount int      `bson:"total_options"`
	eTag       string   `bson:"etag"`
}

// GetFilterDimensionOptions gets the options for a dimension that is part of a Filter
func (c *Client) GetFilterDimensionOptions(ctx context.Context, filterID, dimensionName string, limit, offset int) ([]string, int, string, error) {
	col := c.collections.filters

	logData := log.Data{
		"filter_id":      filterID,
		"dimension_name": dimensionName,
	}

	pipeline := mongo.Pipeline{
		bson.D{{"$match", bson.D{{"filter_id", filterID}}}},
		bson.D{{"$project", bson.D{
			{"dimensionsCount", bson.D{{"$size", "$dimensions"}}},
			{"dimension", bson.D{{"$filter", bson.D{
				{"input", "$dimensions"},
				{"as", "d"},
				{"cond", bson.D{{"$eq", bson.A{"$$d.name", dimensionName}}}}}}}},
			{"etag", 1},
		},
		}},
		bson.D{{"$unwind", bson.D{{"path", "$dimension"}, {"preserveNullAndEmptyArrays", true}}}},
		bson.D{{"$project", bson.D{
			{"etag", 1},
			{"total_options", bson.D{{"$cond", bson.D{
				{"if", bson.D{{"$gt", bson.A{"$dimension", nil}}}},
				{"then", bson.D{{"$size", "$dimension.options"}}},
				{"else", -1}}}}},
			{"options", bson.D{{"$cond", bson.D{
				{"if", bson.D{{"$gt", bson.A{"$dimension", nil}}}},
				{"then", bson.D{{"$slice", bson.A{"$dimension.options", offset, limit}}}},
				{"else", bson.A{}}}}}}},
		}},
	}

	var result []optionQueryResult

	if err := c.conn.Collection(col.name).Aggregate(ctx, pipeline, &result); err != nil {
		return nil, 0, "", errors.Wrap(err, "failed to get filter options")
	}
	if len(result) == 0 {
		return nil, 0, "", &er{
			err:      errors.New("failed to find filter"),
			notFound: true,
			logData:  logData,
		}
	}

	options := result[0]

	if options.TotalCount == -1 {
		return nil, options.TotalCount, "", &er{
			err:     errors.New("failed to find dimension"),
			logData: logData,
		}
	}

	return result[0].Options, result[0].TotalCount, "", nil
}

// DeleteFilterDimensionOptions deletes all options for a given dimension in a Filter
// by the time that it gets to here, you should actually have determined that the etag is fine
func (c *Client) DeleteFilterDimensionOptions(ctx context.Context, filterID, dimensionName string) (string, error) {
	col := c.collections.filters

	logData := log.Data{
		"filter_id":      filterID,
		"dimension_name": dimensionName,
	}

	filter, err := c.GetFilter(ctx, filterID)
	if err != nil {
		return "", &er{
			err:     errors.Wrap(err, "unable to fetch filter"),
			logData: logData,
		}
	}

	// Find the dimension in order to replace it in-memory with the old one (rather than just
	// in the datastore). We need to do this in order to generate an accurate ETag/hash.
	dimensionIndex, err := findDimensionIndex(filter, dimensionName)
	if err != nil {
		return "", &er{
			err:      errors.Wrap(err, "failed to find dimension index"),
			notFound: true,
			logData:  logData,
		}
	}

	filter.Dimensions[dimensionIndex].Options = make([]string, 0)

	if filter.ETag, err = filter.Hash(nil); err != nil {
		return "", &er{
			err:     errors.Wrap(err, "failed to generate eTag"),
			logData: logData,
		}
	}

	selectFilter := bson.M{
		"filter_id": filterID,
	}

	updateFilter := bson.M{
		"$set": bson.M{
			"etag":         filter.ETag,
			"last_updated": c.generate.Timestamp(),
			fmt.Sprintf("dimensions.%d.options", dimensionIndex): bson.A{},
		},
	}

	if _, err := c.conn.Collection(col.name).Update(ctx, selectFilter, updateFilter); err != nil {
		return "", &er{
			err:     errors.Wrap(err, "failed to update filter dimension"),
			logData: logData,
		}
	}

	return filter.ETag, err
}
