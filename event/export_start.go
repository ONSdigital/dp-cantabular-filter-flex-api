package event

/*
   Note: schema copied from
   csv export event lifted
   from the repository
*/

// ExportStart provides an avro structure for a Export Start event
type ExportStart struct {
	InstanceID     string `avro:"instance_id"`
	DatasetID      string `avro:"dataset_id"`
	Edition        string `avro:"edition"`
	Version        string `avro:"version"`
	FilterOutputID string `avro:"filter_output_id"`
	DimensionsID   string `avro:"dimensions_id"`
}
