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

	return nil
}

//AddFilterOutputEvent will append to the existing list of events in the filter output
func (c *Client) AddFilterOutputEvent(ctx context.Context, id string, f *model.Event) error {
	col := c.collections.filterOutputs

	var fo model.FilterOutput
	//find the model output for a given id
	if err := c.conn.Collection(col.name).FindOne(ctx, bson.M{"id": id}, &fo); err != nil {
		err := &er{
			err: errors.Wrap(err, "failed to find filter"),
		}
		if errors.Is(err, mongodb.ErrNoDocumentFound) {
			err.notFound = true
		}
		return err
	}

	//append the filteroutput events with the new event
	ev := fo.Events
	ev = append(ev, *f)

	searchCondition := bson.M{"id": id}
	update := bson.M{"events": ev}

	if _, err := c.conn.Collection(col.name).Update(ctx, searchCondition, update); err != nil {
		return errors.Wrap(err, "failed to add dimension to filter")
	}
	return nil
}
