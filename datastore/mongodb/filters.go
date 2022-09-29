package mongodb

import (
	"context"
	"fmt"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	"github.com/ONSdigital/dp-mongodb/v3/mongodb"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
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
	f.Links.Dimensions = model.Link{
		HREF: f.Links.Self.HREF + "/dimensions",
	}

	col := c.collections.filters

	// lockID, err := col.lock(ctx, f.ID)
	// if err != nil {
	//	return err
	// }
	// defer col.unlock(ctx, lockID)

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
		return nil, &er{
			err:      errors.Wrap(err, "failed to find filter"),
			notFound: errors.Is(err, mongodb.ErrNoDocumentFound),
		}
	}

	return &f, nil
}
