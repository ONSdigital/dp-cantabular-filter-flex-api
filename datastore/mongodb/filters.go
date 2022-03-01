package mongodb

import (
	"context"
	"fmt"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	"github.com/ONSdigital/dp-mongodb/v3/mongodb"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

// Filters returns the Filters collection
func (c *Client) Filters() *mongodb.Collection {
	return c.conn.Collection(c.collections.filters.name)
}

// CreateFilter creates a new Filter in the CantabularFilters collection
func (c *Client) CreateFilter(ctx context.Context, f *model.Filter) error {
	var err error

	if f.ID, err = c.generate.UUID(); err != nil {
		return errors.Wrap(err, "failed to generate UUID: %w")
	}

	if f.ETag, err = f.Hash(nil); err != nil {
		return errors.Wrap(err, "failed to generate eTag: %w")
	}

	f.Links.Self = model.Link{
		HREF: fmt.Sprintf("%s/filters/%s", c.cfg.FilterFlexAPIURL, f.ID),
	}

	col := c.collections.filters

	lockID, err := col.lock(ctx, f.ID.String())
	if err != nil {
		return err
	}
	defer col.unlock(ctx, lockID)

	if _, err = c.conn.Collection(col.name).UpsertById(ctx, f.ID, bson.M{"$set": f}); err != nil {
		return errors.Wrap(err, "failed to upsert filter")
	}

	return nil
}

// UpdateFilter updates a filter in Filters collection
func (c *Client) UpdateFilter(ctx context.Context, filter *model.Filter) error {
	col := c.collections.filters

	lockID, err := col.lock(ctx, filter.ID.String())
	if err != nil {
		return err
	}
	defer col.unlock(ctx, lockID)

	if _, err := c.Filters().UpdateById(ctx, filter.ID, bson.M{"$set": filter}); err != nil {
		return err
	}
	return nil
}

// GetFilter returns a special filter by id
func (c *Client) GetFilter(ctx context.Context, id string) (*model.Filter, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	var filter model.Filter
	if err := c.Filters().FindOne(ctx, bson.M{"filter_id": uuid}, &filter); err != nil {
		return nil, err
	}

	return &filter, nil
}

// GetFilterDimensions gets the dimensions for a Filter in the CantabularFilters collection
func (c *Client) GetFilterDimensions(ctx context.Context, fID string) ([]model.Dimension, error) {
	var err error

	col := c.collections.filters

	var f model.Filter

	v, err := uuid.Parse(fID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse filter uuid")
	}

	if err = c.conn.Collection(col.name).FindOne(ctx, bson.M{"filter_id": v}, &f); err != nil {
		err := &er{
			err: errors.Wrap(err, "failed to get filter dimensions"),
		}
		if errors.Is(err, mongodb.ErrNoDocumentFound) {
			err.notFound = true
		}
		return nil, err
	}

	return f.Dimensions, nil
}
