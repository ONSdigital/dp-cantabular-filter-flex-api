package api

import (
	"errors"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
)

// createFilterRequest is the request body for POST /filters
type createFilterRequest struct {
	CantabularBlob string            `bson:"cantabular_blob" json:"cantabular_blob"`
	Dimensions     []model.Dimension `bson:"dimensions"      json:"dimensions"`
	Dataset        *model.Dataset    `bson:"dataset"         json:"dataset"`
}

func (r *createFilterRequest) Valid() error {
	if r.Dataset == nil {
		return errors.New("missing field: dataset")
	}

	if r.Dataset.ID == "" || r.Dataset.Edition == "" || r.Dataset.Version == 0 {
		return errors.New("missing field: [dataset.id | dataset.edition | dataset.version]")
	}

	if r.CantabularBlob == "" {
		return errors.New("missing field: cantabular_blob")
	}

	if len(r.Dimensions) < 2 {
		return errors.New("missing/invalid field: 'dimensions' must contain at least 2 values")
	}

	for _, d := range r.Dimensions{
		if len(d.Name) == 0 || len(d.DimensionURL) == 0{
			return errors.New("missing field: [dimension[%d].name | dimension[%d].dimension_url]")
		}
	}

	return nil
}

// createFilterResponse is the response body for POST /filters
type createFilterResponse struct{
	model.Filter
}
