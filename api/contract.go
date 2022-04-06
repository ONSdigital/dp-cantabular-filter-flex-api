package api

import (
	"errors"
	"fmt"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
)

// createFilterRequest is the request body for POST /filters
type createFilterRequest struct {
	PopulationType string            `bson:"population_type" json:"population_type"`
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

	if r.PopulationType == "" {
		return errors.New("missing field: population_type")
	}

	if len(r.Dimensions) < 2 {
		return errors.New("missing/invalid field: 'dimensions' must contain at least 2 values")
	}

	for i, d := range r.Dimensions {
		if len(d.Name) == 0 {
			return fmt.Errorf("missing field: [dimension[%d].name]", i)
		}
	}

	return nil
}

// createFilterResponse is the response body for POST /filters
type createFilterResponse struct {
	model.JobState
	Links          model.Links   `json:"links"`
	Dataset        model.Dataset `json:"dataset"`
	PopulationType string        `json:"population_type"`
}

// updateFilter Response is the response body for POST /filters/{id}/submit
// made public because needed for integration tests.
type UpdateFilterResponse struct {
	model.JobState
	Dataset        model.Dataset     `json:"dataset"`
	Links          model.Links       `json:"links"`
	PopulationType string            `json:"population_type"`
	Dimensions     []model.Dimension `json:"dimensions"`
}

// getFilterDimensionsResponse is the response body for GET /filters/{id}
type getFilterResponse struct {
	model.Filter
}

// putFilterResponse is the response body for PUT /filters/{id}
type putFilterResponse struct {
	model.PutFilter
}

// getFilterResponse is the response body for GET /filter-outputs/{id}
type getFilterOutputResponse struct {
	model.JobState
	Links     model.FilterOutputLinks `json:"links"`
	Downloads model.Downloads         `json:"downloads"`
	State     string                  `json:"state"`
}

// createFilterOutputRequest is the request body for POST /filters
type createFilterOutputRequest struct {
	model.FilterOutput
}

func (r *createFilterOutputRequest) Valid() error {
	if err := r.Downloads.CSV.IsNotFullyPopulated(); err != nil {
		return err
	}
	if err := r.Downloads.CSVW.IsNotFullyPopulated(); err != nil {
		return err
	}
	if err := r.Downloads.TXT.IsNotFullyPopulated(); err != nil {
		return err
	}
	if err := r.Downloads.XLS.IsNotFullyPopulated(); err != nil {
		return err
	}

	return nil
}

// getFilterDimensionsResponse is the response body for GET /filters/{id}/dimensions
type getFilterDimensionsResponse struct {
	Items []model.Dimension `json:"items"`
	paginationResponse
}

// paginationResponse represents pagination data as returned to the client.
type paginationResponse struct {
	Limit      int `json:"limit"`
	Offset     int `json:"offset"`
	Count      int `json:"count"`
	TotalCount int `json:"total_count"`
}

// addFilterDimensionResponse is the response body for POST /filters/{id}/dimensions
type addFilterDimensionResponse struct {
	model.Dimension
}

// getDatasetJsonObservationsResponse is the response body for GET /flex/datasets/{dataset_id}/editions/{edition}/versions/{version}/json
type getDatasetJsonObservationsResponse struct {
	Dimensions        []getDatasetJsonResponseDimension `json:"dimensions"`
	Links             getDatasetJsonResponseLinks       `json:"links"`
	Observations      []int                             `json:"observations"`
	TotalObservations int                               `json:"total_observations"`
}

type getDatasetJsonResponseDimension struct {
	DimensionName string                                  `json:"dimension_name"`
	Options       []getDatasetJsonResponseDimensionOption `json:"options"`
}

type getDatasetJsonResponseLinks struct {
	DatasetMetadata getDatasetJsonResponseLink        `json:"dataset_metadata"`
	Self            getDatasetJsonResponseLink        `json:"self"`
	Version         getDatasetJsonResponseVersionLink `json:"version"`
}

type getDatasetJsonResponseDimensionOption struct {
	Href string `json:"href"`
	Id   string `json:"id"`
}

type getDatasetJsonResponseLink struct {
	Href string `json:"href"`
}

type getDatasetJsonResponseVersionLink struct {
	Href    string `json:"href"`
	Version string `json:"version"`
}
