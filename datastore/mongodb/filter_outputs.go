package mongodb

import (
	"context"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

// CreateFilterOutput creates a new FilterOutputs in the CantabularFilters colllection
func (c *Client) CreateFilterOutput(ctx context.Context, f *model.FilterOutputResponse) error {

	if id, err := c.generate.UUID(); err != nil {
		return errors.Wrap(err, "failed to generate UUID: %w")
	} else {
		f.ID = id.String()
	}

	col := c.collections.filterOutputs

	if _, err := c.conn.Collection(col.name).UpsertById(ctx, f.ID, bson.M{"$set": f}); err != nil {
		return errors.Wrap(err, "failed to upsert filter")
	}

	return nil

}
