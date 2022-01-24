package mongodb

import (
	"fmt"
	"context"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"

	"go.mongodb.org/mongo-driver/bson"
	"github.com/google/uuid"
)

// CreateFilter creates a new Filter in the CantabularFilters colllection
func (c *Client) CreateFilter(ctx context.Context, f *model.Filter) error {
	var err error

	if f.ID, err = uuid.NewRandom(); err != nil{
		return fmt.Errorf("failed to generate uuid: %w", err)
	}

	f.Links.Self = model.Link{
		HREF: fmt.Sprintf("%s/flex/filters/%s", c.cfg.FilterFlexAPIURL, f.ID),
	}

	if _, err = c.conn.Collection(filtersCollection).UpsertById(ctx, f.ID, bson.M{"$set": f}); err != nil{
		return fmt.Errorf("failed to upsert filter: %w", err)
	}

	return nil
}
