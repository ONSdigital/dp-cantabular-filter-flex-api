package mongodb

import (
	"context"

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

	return result[0].Options, result[0].TotalCount, result[0].eTag, nil
}
