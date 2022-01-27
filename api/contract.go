package api

import (
	"errors"
	"fmt"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"

	"github.com/google/uuid"
)

// createFilterRequest is the request body for POST /filters
type createFilterRequest struct{
	InstanceID     *uuid.UUID        `bson:"instance_id"     json:"instance_id"`
	DatasetID      string            `bson:"dataset_id"      json:"dataset_id"`
	Edition        string            `bson:"edition"         json:"edition"`
	Version        int               `bson:"version"         json:"version"`
	CantabularBlob string            `bson:"cantabular_blob" json:"cantabular_blob"`
	Dimensions     []model.Dimension `bson:"dimensions"      json:"dimensions"`
}

func (r *createFilterRequest) Valid() error {
	if r.DatasetID == "" || r.Edition == "" || r.CantabularBlob == "" {
		return errors.New("missing field: [dataset_id | edition | cantabular_blob]")
	}

	if r.InstanceID == nil || r.Version == 0 {
		return errors.New("missing/invalid field: [instance_id | version]")
	}

	if len(r.Dimensions) < 2 {
		return errors.New("missing/invalid field: 'dimensions' must contain at least 2 values")
	}

	for i, d := range r.Dimensions{
		if err := d.Valid(); err != nil{
			return fmt.Errorf("invalid field: 'dimension [%d]': %w", i, err)
		}
	}

	return nil
}

// createFilterResponse is the response body for POST /filters
type createFilterResponse struct{
	model.Filter
}
