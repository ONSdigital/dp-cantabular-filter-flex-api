package mongodb

import (
	"context"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	"github.com/ONSdigital/log.go/v2/log"

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
		{Key: "_id", Value: 0},
		{Key: "totalCount", Value: "$totalCount"},
	}

	// We can't use `$slice` with a value of 0 , so if we know we're not returning results,
	// don't include the `dimensions` projection.
	if limit > 0 {
		projection = append(projection, bson.E{
			Key:   "dimensions",
			Value: bson.D{{Key: "$slice", Value: bson.A{"$dimensions", offset, limit}}},
		})
	}

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{{Key: "filter_id", Value: fID}}}},
		bson.D{{Key: "$unwind", Value: "$dimensions"}},
		bson.D{{Key: "$replaceRoot", Value: bson.D{{Key: "newRoot", Value: "$dimensions"}}}},
		bson.D{{Key: "$sort", Value: bson.D{{Key: "name", Value: 1}}}},
		bson.D{{
			Key: "$group",
			Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "dimensions", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
				{Key: "totalCount", Value: bson.D{{Key: "$sum", Value: 1}}},
			},
		}},
		bson.D{{Key: "$project", Value: projection}},
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

// GetFilterDimension gets a specific dimensions for a Filter in the CantabularFilters collection
func (c *Client) GetFilterDimension(ctx context.Context, fID, dimName string) (model.Dimension, error) {
	col := c.collections.filters

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{{Key: "filter_id", Value: fID}}}},
		bson.D{{Key: "$unwind", Value: "$dimensions"}},
		bson.D{{Key: "$match", Value: bson.D{{Key: "dimensions.name", Value: dimName}}}},
		bson.D{{Key: "$replaceRoot", Value: bson.D{{Key: "newRoot", Value: "$dimensions"}}}},
	}

	var results []model.Dimension
	if err := c.conn.Collection(col.name).Aggregate(ctx, pipeline, &results); err != nil {
		return model.Dimension{}, errors.Wrap(err, "failed to get filter dimension")
	}

	if len(results) == 0 {
		return model.Dimension{}, &er{
			err:      errors.New("failed to find filter dimension"),
			notFound: true,
		}
	}

	return results[0], nil
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

func (c *Client) UpdateFilterDimension(ctx context.Context, filterID string, dimensionName string, dimension model.Dimension, currentETag string) (string, error) {
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

	// The Etag is optional, so we check if it's provided before comparing it.
	if currentETag != "" && currentETag != filter.ETag {
		logData["expected_etag"] = currentETag
		logData["actual_etag"] = filter.ETag

		return "", &er{
			err:      errors.New("conflict: invalid ETag provided or filter has been updated"),
			conflict: true,
			logData:  logData,
		}
	}

	// Find the dimension in order to replace it in-memory with the old one (rather than just
	// in the datastore). We need to do this in order to generate an accurate ETag/hash.
	newDimensionIndex, err := findDimensionIndex(filter, dimensionName)
	if err != nil {
		return "", &er{
			err:      errors.Wrap(err, "failed to find dimension"),
			notFound: true,
			logData:  logData,
		}
	}

	if filter.Type != "multivariate" {
		if at := filter.Dimensions[newDimensionIndex].IsAreaType; at == nil || !(*at) {
			return "", &er{
				err:     errors.New("non geography variable"),
				logData: logData,
			}
		}
	}

	filter.Dimensions[newDimensionIndex] = dimension

	if filter.ETag, err = filter.Hash(nil); err != nil {
		return "", &er{
			err:     errors.Wrap(err, "failed to generate eTag"),
			logData: logData,
		}
	}

	selectFilter := bson.M{
		"filter_id":       filterID,
		"dimensions.name": dimensionName,
	}

	updateFilter := bson.M{
		"$set": bson.M{
			"etag":         filter.ETag,
			"last_updated": c.generate.Timestamp(),
			"dimensions.$": dimension,
		},
	}

	if _, err := c.conn.Collection(col.name).Update(ctx, selectFilter, updateFilter); err != nil {
		return "", &er{
			err:     errors.Wrap(err, "failed to update filter dimension"),
			logData: logData,
		}
	}

	return filter.ETag, nil
}

func (c *Client) RemoveFilterDimensionOption(ctx context.Context, filterID, dimension, option string, currentETag string) (string, error) {
	col := c.collections.filters
	logData := log.Data{
		"filter_id": filterID,
		"dimension": dimension,
		"option":    option,
	}

	filter, err := c.GetFilter(ctx, filterID)
	if err != nil {
		return "", &er{
			err:     errors.Wrap(err, "unable to fetch filter"),
			logData: logData,
		}
	}

	if currentETag != "" && currentETag != filter.ETag {
		logData["expected_etag"] = currentETag
		logData["actual_etag"] = filter.ETag

		return "", &er{
			err:      errors.New("conflict: invalid ETag provided or filter has been updated"),
			conflict: true,
			logData:  logData,
		}
	}

	dimIndex, err := findDimensionIndex(filter, dimension)
	if err != nil {
		return "", &er{
			err:      errors.Wrap(err, "failed to find dimension index"),
			notFound: true,
			logData:  logData,
		}
	}

	opts := filter.Dimensions[dimIndex].Options

	optIndex, err := findDimensionOptionIndex(opts, option)
	if err != nil {
		return "", &er{
			err:      errors.Wrap(err, "failed to find dimension option index"),
			notFound: true,
			logData:  logData,
		}
	}

	// Remove option from dimension.options while preserving order to generate correct eTag
	opts = append(opts[:optIndex], opts[optIndex+1:]...)
	filter.Dimensions[dimIndex].Options = opts

	if filter.ETag, err = filter.Hash(nil); err != nil {
		return "", &er{
			err:     errors.Wrap(err, "failed to generate eTag"),
			logData: logData,
		}
	}

	selectFilter := bson.M{
		"filter_id":       filterID,
		"dimensions.name": dimension,
	}

	updateFilter := bson.M{
		"$set": bson.M{
			"etag":         filter.ETag,
			"last_updated": c.generate.Timestamp(),
		},
		"$pull": bson.M{
			"dimensions.$.options": option,
		},
	}

	if _, err := c.conn.Collection(col.name).Update(ctx, selectFilter, updateFilter); err != nil {
		return "", &er{
			err:     errors.Wrap(err, "failed to update filter dimension"),
			logData: logData,
		}
	}

	return filter.ETag, nil
}

// findDimension loops through a dimension, looking for a dimension by name, and
// returns Dimension.
func findDimension(filter *model.Filter, dimensionName string) (model.Dimension, error) {
	for i := range filter.Dimensions {
		dim := filter.Dimensions[i]
		if dim.Name == dimensionName {
			return dim, nil
		}
	}

	return model.Dimension{}, errors.New("could not find dimension")
}

// findDimensionIndex loops through a dimension, looking for a dimension by name, and
// returns the index of the item in the Dimensions slice.
func findDimensionIndex(filter *model.Filter, dimensionName string) (int, error) {
	for i := range filter.Dimensions {
		if filter.Dimensions[i].Name == dimensionName {
			return i, nil
		}
	}

	return 0, errors.New("could not find dimension")
}

// findDimensionIndex loops through a dimension, looking for a option by name, and
// returns the index of the item in the options slice.
func findDimensionOptionIndex(options []string, optionName string) (int, error) {
	for i, o := range options {
		if o == optionName {
			return i, nil
		}
	}

	return 0, errors.New("could not find dimension option")
}

// DeleteFilterDimension deletes a given dimension in a Filter
// etag checks should haev passed before calling this function
func (c *Client) DeleteFilterDimension(ctx context.Context, filterID, dimensionName string) (string, error) {
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

	if len(filter.Dimensions) < 2 {
		return "", &er{
			err:      errors.New("can't delete dimension as minimum required condition didn't match"),
			conflict: true,
			logData:  logData,
		}
	}
	// 2. Check if the dimension passed for deletion is not area-type i.e. 'is_area_type' is false.
	dimension, err := findDimension(filter, dimensionName)
	if err != nil {
		return "", &er{
			err:      errors.Wrap(err, "failed to find dimension"),
			notFound: true,
			logData:  logData,
		}
	}
	if dimension.IsAreaType != nil && *(dimension.IsAreaType) {
		return "", &er{
			err:      errors.New("dimension with area type as true can't be deleted"),
			conflict: true,
			logData:  logData,
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

	dimensions := filter.Dimensions
	dimensions = append(dimensions[:dimensionIndex], dimensions[dimensionIndex+1:]...)
	filter.Dimensions = dimensions

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
		},
		"$pull": bson.M{
			"dimensions": bson.M{"name": dimensionName},
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
