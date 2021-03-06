package api

import (
	"fmt"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"

	"github.com/pkg/errors"
)

// paginationResponse represents pagination data as returned to the client.
type paginationResponse struct {
	Limit      int `json:"limit"`
	Offset     int `json:"offset"`
	Count      int `json:"count"`
	TotalCount int `json:"total_count"`
}

// createFilterRequest is the request body for POST /filters
type createFilterRequest struct {
	PopulationType string            `json:"population_type"`
	Dimensions     []model.Dimension `json:"dimensions"`
	Dataset        *model.Dataset    `json:"dataset"`
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

		if len(d.ID) != 0 {
			return fmt.Errorf("unexpected field id provided for: %s", d.Name)
		}

		if len(d.Label) != 0 {
			return fmt.Errorf("unexpected field label provided for: %s", d.Name)
		}
	}

	return nil
}

// createFilterResponse is the response body for POST /filters
type createFilterResponse struct {
	model.Filter
}

// getFilterDimensionsResponse is the response body for GET /filters/{id}
type getFilterResponse struct {
	model.Filter
}

// putFilterResponse is the response body for PUT /filters/{id}
type putFilterResponse struct {
	Events         []model.Event `json:"events"`
	Dataset        model.Dataset `json:"dataset"`
	PopulationType string        `json:"population_type"`
}

// putFilterOutputRequest is the request body for PUT /filters
type putFilterOutputRequest struct {
	State     string          `json:"state"`
	Downloads model.Downloads `json:"downloads"`
}

type addFilterOutputEventRequest struct {
	model.Event
}

func (r *putFilterOutputRequest) Valid() error {
	if r.Downloads.CSV != nil {
		if err := r.Downloads.CSV.IsNotFullyPopulated(); err != nil {
			return errors.Wrap(err, "'csv' field not fully populated")
		}
	}
	if r.Downloads.CSVW != nil {
		if err := r.Downloads.CSVW.IsNotFullyPopulated(); err != nil {
			return errors.Wrap(err, "'csvw' field not fully populated")
		}
	}
	if r.Downloads.TXT != nil {
		if err := r.Downloads.TXT.IsNotFullyPopulated(); err != nil {
			return errors.Wrap(err, "'txt' field not fully populated")
		}
	}
	if r.Downloads.XLS != nil {
		if err := r.Downloads.XLS.IsNotFullyPopulated(); err != nil {
			return errors.Wrap(err, "'xls' field not fully populated")
		}
	}

	return nil
}

type GetFilterDimensionOptionsItem struct {
	Option    string     `json:"option"`
	Self      model.Link `json:"self"`
	Filter    model.Link `json:"filter"`
	Dimension model.Link `json:"Dimension"`
}

// getDimensionOptionsResponse is the response body for GET /filters/{id}/dimensions/{name}/options
type GetFilterDimensionOptionsResponse struct {
	Items []GetFilterDimensionOptionsItem `json:"items"`
	paginationResponse
}

// getFilterOutputResponse is the response body for GET/filter-outputs/{id}
type getFilterOutputResponse struct {
	model.FilterOutput
}

// getFilterDimensionsResponse is the response body for GET /filters/{id}/dimensions
type getFilterDimensionsResponse struct {
	Items dimensionItems `json:"items"`
	paginationResponse
}

// addFilterDimensionRequest is the request body for POST /filters/{id}/dimensions
type addFilterDimensionRequest struct {
	model.Dimension
}

// addFilterDimensionResponse is the response body for POST /filters/{id}/dimensions
type addFilterDimensionResponse struct {
	dimensionItem
}

// updateFilterDimensionResponse is the request body for PUT /filters/{id}/dimensions/{name}
type updateFilterDimensionRequest struct {
	model.Dimension
}

func (u *updateFilterDimensionRequest) Valid() error {
	if len(u.ID) == 0 {
		return errors.New("missing field: [id]")
	}

	return nil
}

// updateFilterDimensionResponse is the response body for PUT /filters/{id}/dimensions/{name}
type updateFilterDimensionResponse struct {
	dimensionItem
}

type dimensionItem struct {
	ID    string             `json:"id"`
	Name  string             `json:"name"`
	Label string             `json:"label"`
	Links dimensionItemLinks `json:"links"`
}

func (d *dimensionItem) fromDimension(dim model.Dimension, host, filterID string) {
	filterURL := fmt.Sprintf("%s/filters/%s", host, filterID)
	dimURL := fmt.Sprintf("%s/dimensions/%s", filterURL, dim.Name)

	d.ID = dim.ID
	d.Name = dim.Name
	d.Label = dim.Label
	d.Links = dimensionItemLinks{
		Self: model.Link{
			HREF: dimURL,
			ID:   dim.ID,
		},
		Filter: model.Link{
			HREF: filterURL,
			ID:   filterID,
		},
		Options: model.Link{
			HREF: dimURL + "/options",
		},
	}

}

type dimensionItems []dimensionItem

func (items *dimensionItems) fromDimensions(dims []model.Dimension, host, filterID string) {
	if len(dims) == 0 {
		*items = dimensionItems{}
	}
	for _, dim := range dims {
		var item dimensionItem
		item.fromDimension(dim, host, filterID)
		*items = append(*items, item)
	}
}

type getFilterDimensionResponse struct {
	dimensionItem
	IsAreaType bool `json:"is_area_type"`
}

type dimensionItemLinks struct {
	Filter  model.Link `json:"filter"`
	Options model.Link `json:"options"`
	Self    model.Link `json:"self"`
}

type addFilterDimensionOptionRequest struct {
	FilterID, Dimension, Option string
}

type addFilterDimensionOptionResponse struct {
	Option string                     `json:"option"`
	Links  filterDimensionOptionLinks `json:"links`
}

type filterDimensionOptionLinks struct {
	Self      model.Link `json:"self"`
	Filter    model.Link `json:"filter"`
	Dimension model.Link `json:"dimension"`
}

type deleteFilterDimensionOptionRequest struct {
	FilterID, Dimension, Option string
}

type submitFilterResponse struct {
	InstanceID     string            `json:"instance_id"`
	FilterOutputID string            `json:"filter_output_id"`
	Dataset        model.Dataset     `json:"dataset"`
	Links          model.FilterLinks `json:"links"`
	PopulationType string            `json:"population_type"`
}

// getDatasetJSONResponse is the response body for GET /datasets/{dataset_id}/editions/{edition}/versions/{version}/json
type getDatasetJSONResponse struct {
	Dimensions        []DatasetJSONDimension `json:"dimensions"`
	Links             DatasetJSONLinks       `json:"links"`
	Observations      []int                  `json:"observations"`
	TotalObservations int                    `json:"total_observations"`
}

type DatasetJSONDimension struct {
	DimensionName string       `json:"dimension_name"`
	Options       []model.Link `json:"options"`
}

type DatasetJSONLinks struct {
	DatasetMetadata model.Link `json:"dataset_metadata"`
	Self            model.Link `json:"self"`
	Version         model.Link `json:"version"`
}
