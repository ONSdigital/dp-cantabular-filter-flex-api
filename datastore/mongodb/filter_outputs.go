package mongodb

import (
	"context"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	"github.com/ONSdigital/dp-mongodb/v3/mongodb"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

// CreateFilterOutput creates a new FilterOutputs in the CantabularFilters collection
func (c *Client) CreateFilterOutput(ctx context.Context, f *model.FilterOutput) error {
	id, err := c.generate.UUID()
	if err != nil {
		return errors.Wrap(err, "failed to generate UUID: %w")
	}

	f.ID = id.String()

	col := c.collections.filterOutputs

	if _, err := c.conn.Collection(col.name).UpsertById(ctx, f.ID, bson.M{"$set": f}); err != nil {
		return errors.Wrap(err, "failed to upsert filter")
	}

	return nil
}

// UpdateFilterOutput creates/updates a FilterOutputs in the CantabularFilters collection
func (c *Client) UpdateFilterOutput(ctx context.Context, f *model.FilterOutput) error {

	col := c.collections.filterOutputs

	searchCondition := bson.M{"id": f.ID}
	update := bson.M{"$set": bson.M{"state": f.State, "downloads": f.Downloads}}

	var result *mongodb.CollectionUpdateResult
	var err error
	if result, err = c.conn.Collection(col.name).Update(ctx, searchCondition, update); err != nil {
		return errors.Wrap(err, "failed to update filter outputs")
	}

	//check for the condition when the search fails. In this case there is no error returned by update API
	//but the result's MatchedCount returns the count of matched documents
	if result.MatchedCount < 1 {
		println("Record not found... Searched ID: ", "`", f.ID, "`")
		err := &er{
			err: errors.Wrap(err, "failed to find filter output"),
		}
		err.notFound = true
		return err
	}
	/* Do we need to handle the error case when for any reson there are more than one record with
	same filter output id
	else if result.ModifiedCount > 1 {
		// return errors.
	}*/

	return nil
}
