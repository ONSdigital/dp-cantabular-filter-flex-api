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

	// Paginates the nested dimensions field
	dimensionsFacet := bson.A{
		bson.D{{"$replaceRoot", bson.D{{"newRoot", "$dimensions"}}}},
		bson.D{{"$sort", bson.D{{"name", 1}}}},
		bson.D{{"$limit", limit}},
		bson.D{{"$skip", offset}},
	}

	// Adds a total record count for all a filter's dimensions
	paginationFacet := bson.A{
		bson.D{{"$count", "totalCount"}},
	}

	// Adds the filter-found flag
	foundFacet := bson.A{
		bson.D{{"$project", bson.D{{"found", "$foundFilters"}}}},
	}

	facet := bson.D{{
		"$facet",
		bson.D{
			{"dimensions", dimensionsFacet},
			{"pagination", paginationFacet},
			{"foundFilters", foundFacet},
		},
	}}

	// Flatten lists where we know there is only one value, or all values are equal
	flatten := bson.D{
		{"$addFields", bson.D{
			{"pagination", bson.D{{"$arrayElemAt", bson.A{"$pagination", 0}}}},
			{"foundFilters", bson.D{{"$arrayElemAt", bson.A{"$foundFilters", 0}}}},
		}},
	}

	pipeline := mongo.Pipeline{
		bson.D{{"$match", bson.D{{"filter_id", fID}}}},
		// Add a flag to indicate if a filter was set. Since stages run on matched rows,
		// if no row was found this will never be set, and will deserialize into `false`.
		bson.D{{"$addFields", bson.D{{"foundFilters", true}}}},
		bson.D{{"$unwind", "$dimensions"}},
		facet,
		flatten,
	}

	var results []dimensionQueryResult
	if err := c.conn.Collection(col.name).Aggregate(ctx, pipeline, &results); err != nil {
		return nil, 0, &er{
			err: errors.Wrap(err, "failed to get filter dimensions"),
		}
	}

	if len(results) == 0 {
		return nil, 0, errors.Errorf("no results found in aggregate query: %v", results)
	}

	result := results[0]

	if !result.FoundFilters.Found {
		return nil, 0, &er{
			err:      errors.Errorf("failed to find filter with ID (%s)", fID),
			notFound: true,
		}
	}

	return result.Dimensions, result.Pagination.TotalCount, nil
}

// foundFilters contains a bool indicating if there were any row matches.
// We need this to distinguish between zero dimensions (which is not an error)
// and zero matched filters (which is an error).
type foundFilters struct {
	Found bool `bson:"found"`
}

type paginationQueryResult struct {
	TotalCount int `bson:"totalCount"`
}

type dimensionQueryResult struct {
	Dimensions   []model.Dimension     `bson:"dimensions"`
	Pagination   paginationQueryResult `bson:"pagination"`
	FoundFilters foundFilters          `bson:"foundFilters"`
}
