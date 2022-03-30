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
	var err error
	var records []model.FilterOutput
	var num int

	searchCondition := bson.M{"id": f.ID}

	if num, err = c.conn.Collection(col.name).Find(ctx, searchCondition, &records, mongodb.Limit(1)); err != nil {
		return errors.Wrap(err, "failed to update filter outputs")
	}
	if num < 1 {
		err := &er{
			err: errors.Errorf("failed to find filter output"),
		}
		err.notFound = true
		return err
	}

	//a record with stat 'completed' can't be updated further
	if records[0].State == model.COMPLETED {
		err := &er{
			err: errors.Errorf(`filter output is already in "completed" state`),
		}
		err.forbidden = true
		return err
	}

	update := bson.M{"$set": bson.M{"state": f.State, "downloads": f.Downloads}}

	var result *mongodb.CollectionUpdateResult

	if result, err = c.conn.Collection(col.name).Update(ctx, searchCondition, update); err != nil {
		return errors.Wrap(err, "failed to update filter outputs")
	}

	//this should not happen unless the recored is being removed btween the first find and the update
	//check for the condition if the update the search failed. In this case there is no error returned by update API
	//but the result's MatchedCount returns the count of matched documents
	if result.MatchedCount < 1 {
		println("Record not found... Searched ID: ", "`", f.ID, "`")
		err := &er{
			err: errors.Errorf("failed to find filter output"),
		}
		err.notFound = true
		return err
	}

	return nil
}
