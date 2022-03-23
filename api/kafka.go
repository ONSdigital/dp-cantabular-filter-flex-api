package api

import (
	"fmt"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/schema"
)

// ProduceCSVCreateEvent creates an event to create a csv when POST filters/{id}/submit is hit
func (api *API) ProduceCSVCreateEvent(e *event.ExportStart, rowCount int32) error {
	if err := api.producer.Send(schema.CSVCreated, &event.CSVCreated{
		InstanceID: e.InstanceID,
		DatasetID:  e.DatasetID,
		Edition:    e.Edition,
		Version:    e.Version,
		RowCount:   rowCount,
	}); err != nil {
		return fmt.Errorf("error sending csv-created event: %w", err)
	}
	return nil
}
