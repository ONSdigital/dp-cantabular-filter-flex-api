package api

import (
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"

	"github.com/google/uuid"
)

// createFilterRequest is the request body for POST /filters
type createFilterRequest struct{
	InstanceID uuid.UUID         `bson:"instance_id" json:"instance_id"`
	DatasetID  string            `bson:"dataset_id"  json:"datset_id"`
	Edition    string            `bson:"edition"     json:"edition"`
	Version    int               `bson:"version"     json:"version"`
	Dimensions []model.Dimension `bson:"dimensions"  json:"dimensions"`
}

// createFilterResponse is the response body for POST /filters
type createFilterResponse struct{
	model.Filter
}
