package mongodb

import (
	"context"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	"github.com/pkg/errors"
)

// CreateFilterOutputs creates a new FilterOutputs in the CantabularFilters colllection
func (c *Client) CreateFilterOutputs(ctx context.Context, f *model.FilterOutput) error {
	var err error
	col := c.collections.filterOutputs

	if _, err = c.conn.Collection(col.name).Insert(ctx, f); err != nil {
		return errors.Wrap(err, "failed to insert filter outputs")
	}
	return nil
}
