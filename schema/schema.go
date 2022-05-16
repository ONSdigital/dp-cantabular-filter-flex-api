package schema

import (
	"github.com/ONSdigital/dp-kafka/v3/avro"
)

var exportStart = `{
  "type": "record",
  "name": "cantabular-export-start",
  "fields": [
    {"name": "instance_id",      "type": "string", "default": ""},
    {"name": "dataset_id",       "type": "string", "default": ""},
    {"name": "edition",          "type": "string", "default": ""},
    {"name": "version",          "type": "string", "default": ""},
    {"name": "filter_output_id", "type": "string", "default": ""},
    {"name": "dimension_ids",    "type": { "type": "array", "items": "string"}, "default": [] }
  ]
}`

var ExportStart = &avro.Schema{
	Definition: exportStart,
}
