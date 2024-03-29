package api

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/event"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/schema"
)

// ProduceCSVCreateEvent creates an event to create a csv when POST filters/{id}/submit is hit
func (api *API) produceExportStartEvent(ctx context.Context, e event.ExportStart) error {
	if err := api.producer.Send(ctx, schema.ExportStart, e); err != nil {
		return errors.Wrap(err, "error sending 'export_start' event")
	}

	return nil
}
