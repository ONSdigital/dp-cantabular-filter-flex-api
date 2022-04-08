package schema

import (
	"github.com/ONSdigital/dp-kafka/v3/avro"
)

var csvCreated = `{
  "type": "record",
  "name": "cantabular-csv-created",
  "fields": [
    {"name": "instance_id", "type": "string", "default": ""},
    {"name": "dataset_id", "type": "string", "default": ""},
    {"name": "edition", "type": "string", "default": ""},
    {"name": "version", "type": "string", "default": ""},
    {"name": "row_count", "type": "int", "default": 0}
  ]
}`

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

var CSVCreated = &avro.Schema{
	Definition: csvCreated,
}

var ExportStart = &avro.Schema{
	Definition: exportStart,
}
