package api

import (
	"fmt"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/schema"
)

// ExportStart provides an avro structure for a Export Start event
type ExportStart struct {
	InstanceID string `avro:"instance_id"`
	DatasetID  string `avro:"dataset_id"`
	Edition    string `avro:"edition"`
	Version    string `avro:"version"`
	Filter_ID  string `avro:"filter_id"`
}

// ProduceCSVCreateEvent creates an event to create a csv when POST filters/{id}/submit is hit
func (api *API) ProduceCSVCreateEvent(e *ExportStart, rowCount int32) error {
	if err := api.Pro.Send(schema.CSVCreated); err != nil {
		return fmt.Errorf("error sending csv-created event: %w", err)
	}
	return nil
}
