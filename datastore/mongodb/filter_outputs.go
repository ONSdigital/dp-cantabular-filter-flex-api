package mongodb

import (
	"context"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	"github.com/ONSdigital/dp-mongodb/v3/mongodb"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

// FilterOutputs returns the FilterOutputs collection
func (c *Client) FilterOutputs() *mongodb.Collection {
	return c.conn.Collection(c.collections.filterOutputs.name)
}

// CreateFilterOutput upserts a FilterOutput in the FilterOutputs collection
func (c *Client) CreateFilterOutput(ctx context.Context, filterOutput *model.FilterOutput) error {
	if filterOutput.ID == uuid.Nil {
		id, err := c.generate.UUID()
		if err != nil {
			return errors.Wrap(err, "failed to generate uuid")
		}
		filterOutput.ID = id
	}

	col := c.collections.filters
	lockID, err := col.lock(ctx, filterOutput.ID.String())
	if err != nil {
		return err
	}
	defer col.unlock(ctx, lockID)

	if _, err := c.FilterOutputs().UpsertById(ctx, filterOutput.ID, bson.M{"$set": filterOutput}); err != nil {
		return errors.Wrap(err, "failed to upsert filter output")
	}
	return nil
}
