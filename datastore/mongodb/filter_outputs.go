package mongodb

import (
	"context"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	"github.com/ONSdigital/dp-mongodb/v3/mongodb"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

// GetFilterOutput
func (c *Client) GetFilterOutput(ctx context.Context, filterID string) (*model.FilterOutput, error) {
	var filterOutput model.FilterOutput

	coll := c.collections.filterOutputs

	if err := c.conn.Collection(coll.name).FindOne(ctx, bson.M{"id": filterID}, &filterOutput); err != nil {
		err := &er{
			err: err,
		}
		if errors.Is(err, mongodb.ErrNoDocumentFound) {
			err.notFound = true
		}
		return nil, err
	}

	return &filterOutput, nil
}

// CreateFilterOutput creates a new FilterOutputs in the CantabularFilters collection
func (c *Client) CreateFilterOutput(ctx context.Context, f *model.FilterOutput) error {
	id, err := c.generate.UUID()
	if err != nil {
		return errors.Wrap(err, "failed to generate UUID: %w")
	}

	f.ID = id.String()
	f.Links.Self.ID = id.String()

	col := c.collections.filterOutputs

	if _, err := c.conn.Collection(col.name).UpsertById(ctx, f.ID, bson.M{"$set": f}); err != nil {
		return errors.Wrap(err, "failed to upsert filter")
	}

	return nil
}

// UpdateFilterOutput creates/updates a FilterOutputs in the CantabularFilters collection
func (c *Client) UpdateFilterOutput(ctx context.Context, f *model.FilterOutput) error {
	col := c.collections.filterOutputs

	var existing model.FilterOutput
	queryFilter := bson.M{"id": f.ID}

	if err := c.conn.Collection(col.name).FindOne(ctx, queryFilter, &existing); err != nil {
		return &er{
			err:      errors.Wrap(err, "failed to find filter output"),
			notFound: errors.Is(err, mongodb.ErrNoDocumentFound),
		}
	}

	// a record with state 'completed' can't be updated further
	if existing.State == model.Completed {
		return &er{
			err:       errors.Errorf(`filter output is already in "completed" state`),
			forbidden: true,
		}
	}

	fields := bson.M{"state": f.State}
	if f.Downloads.CSV != nil {
		fields["downloads.csv"] = f.Downloads.CSV
	}
	if f.Downloads.CSVW != nil {
		fields["downloads.csvw"] = f.Downloads.CSVW
	}
	if f.Downloads.TXT != nil {
		fields["downloads.txt"] = f.Downloads.TXT
	}
	if f.Downloads.XLS != nil {
		fields["downloads.xls"] = f.Downloads.XLS
	}

	update := bson.M{"$set": fields}

	lockID, err := col.lock(ctx, f.ID)
	if err != nil {
		return err
	}
	defer col.unlock(ctx, lockID)

	rec, err := c.conn.Collection(col.name).Update(ctx, queryFilter, update)
	if err != nil {
		return errors.Wrap(err, "failed to update filter output")
	}

	// This should not happen unless the recored is being removed bteween the
	// first find and the update. Check if the update failed sue to search condition.
	// Update returns no error but the result's MatchedCount has to be checked
	// for number of updated records
	if rec.MatchedCount < 1 {
		return &er{
			err:      errors.Errorf("failed to find filter output"),
			notFound: true,
		}
	}

	return nil
}

//AddFilterOutputEvent will append to the existing list of events in the filter output
func (c *Client) AddFilterOutputEvent(ctx context.Context, id string, f *model.Event) error {
	col := c.collections.filterOutputs

	ne := map[string]model.Event{
		"events": *f,
	}

	var rec *mongodb.CollectionUpdateResult
	var err error
	if rec, err = c.conn.Collection(col.name).Update(ctx, bson.M{"id": id}, bson.M{"$push": ne}); err != nil {
		return errors.Wrap(err, "failed to add event to filter output")
	}

	if rec.MatchedCount != 1 {
		return &er{
			err:      errors.Errorf("filter output not found"),
			notFound: true,
		}
	}
	return nil
}
