package mongodb

import (
	"context"
	"fmt"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

// CreateFilter creates a new Filter in the CantabularFilters colllection
func (c *Client) CreateFilter(ctx context.Context, f *model.Filter) error {
	var err error

	if f.ID, err = c.generate.UUID(); err != nil {
		return errors.Wrap(err, "failed to generate UUID: %w")
	}

	if f.ETag, err = f.Hash(nil); err != nil {
		return errors.Wrap(err, "failed to generate eTag: %w")
	}

	f.Links.Self = model.Link{
		HREF: fmt.Sprintf("%s/flex/filters/%s", c.cfg.FilterFlexAPIURL, f.ID),
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