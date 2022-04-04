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

	var docs []model.FilterOutput
	sc := bson.M{"id": f.ID}

	var err error
	var num int

	if num, err = c.conn.Collection(col.name).Find(ctx, sc, &docs, mongodb.Limit(1)); err != nil {
		return errors.Wrap(err, "failed to update filter output")
	}
	if num < 1 {
		return &er{
			err:      errors.Errorf("failed to find filter output"),
			notFound: true,
		}
	}

	//a record with stat 'completed' can't be updated further
	if docs[0].State == model.Completed {
		return &er{
			err:       errors.Errorf(`filter output is already in "completed" state`),
			forbidden: true,
		}
	}

	updates := bson.M{"$set": bson.M{"state": f.State, "downloads": f.Downloads}}

	var rec *mongodb.CollectionUpdateResult

	if rec, err = c.conn.Collection(col.name).Update(ctx, sc, updates); err != nil {
		return errors.Wrap(err, "failed to update filter output")
	}

	//This should not happen unless the recored is being removed btween the first find and the update.
	//Check if the update failed sue to search condition.
	//Update returns no error but the result's MatchedCount has to be checked for number of updated records
	if rec.MatchedCount < 1 {
		return &er{
			err:      errors.Errorf("failed to find filter output"),
			notFound: true,
		}
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
