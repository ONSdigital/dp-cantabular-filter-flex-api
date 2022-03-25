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
type updateFilterResponse struct {
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

// createFilterOutputResponse is the response body for POST /filters-output
type createFilterOutputResponse struct {
	model.FilterOutput
}

// filterOutputResponse is the response body for PUT /filters-outputs
type filterOutputResponse struct {
	model.FilterOutput
	model.JobState
	Links model.FilterOutputLinks `json:"links"`
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
