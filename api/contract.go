package api

import (
	"errors"
	"fmt"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
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

// createFilterOutputResponse is the response body for POST /filters-output
type createFilterOutputResponse struct {
	model.FilterOutput
}

// filterOutputResponse is the response body for PUT /filters-outputs
type filterOutputResponse struct {
	model.FilterOutput
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

type dimensionItem struct{
	Name string              `json:"name"`
	Links dimensionItemLinks `json:"links"`
}

func (d *dimensionItem) fromDimension(dim model.Dimension, host, filterID string) {
	filterURL := fmt.Sprintf("%s/filters/%s", host, filterID)
	dimURL := fmt.Sprintf("%s/dimensions/%s", filterURL, dim.Name)

	d.Name  = dim.Name
	d.Links = dimensionItemLinks{
		Self:    model.Link{
			HREF: dimURL,
			ID: dim.Name,
		},
		Filter:  model.Link{
			HREF: filterURL,
			ID: filterID,
		},
		Options: model.Link{
			HREF: dimURL + "/options",
		},
	}

}

type dimensionItems []dimensionItem

func (items *dimensionItems) fromDimensions(dims []model.Dimension, host, filterID string) {
	if len(dims) == 0{
		*items = dimensionItems{}
	}
	for _, dim := range dims{
		var item dimensionItem
		item.fromDimension(dim, host, filterID)
		*items = append(*items, item)
	}
}

type dimensionItemLinks struct{
	Filter  model.Link `json:"filter"`
	Options model.Link `json:"options"`
	Self    model.Link `json:"self"`
}
