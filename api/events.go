package api

import (
	"fmt"

	"github.com/ONSdigital/dp-kafka/v3/avro"
)

/*
   Note: schema copied from
   csv export event lifted
   from the repository
*/

// ExportStart provides an avro structure for a Export Start event
type ExportStart struct {
	InstanceID string `avro:"instance_id"`
	DatasetID  string `avro:"dataset_id"`
	Edition    string `avro:"edition"`
	Version    string `avro:"version"`
	Filter_ID  string `avro:"filter_id"`
}

// ProduceCSVCreateEvent creates an event to create a csv when POST filters/{id}/submit is hit
func (api *API) ProduceCSVCreateEvent(e *ExportStart) error {

	var exportStart = `{
  "type": "record",
  "name": "cantabular-export-start",
  "fields": [
    {"name": "instance_id", "type": "string", "default": ""},
    {"name": "dataset_id",  "type": "string", "default": ""},
    {"name": "edition",     "type": "string", "default": ""},
    {"name": "version",     "type": "string", "default": ""},
    {"name": "filter_id",   "type":"string",  "default": ""}
  ]
}`

	var ExportStart = &avro.Schema{
		Definition: exportStart,
	}

	if err := api.Producer.Send(ExportStart, e); err != nil {
		return fmt.Errorf("error sending csv-created event: %w", err)
	}
	return nil
}
