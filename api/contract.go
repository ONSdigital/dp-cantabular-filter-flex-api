package api

import (
	"errors"
	"strings"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	"github.com/google/uuid"
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

	for _, d := range r.Dimensions {
		if len(d.Name) == 0 || len(d.DimensionURL) == 0 {
			return errors.New("missing field: [dimension[%d].name | dimension[%d].dimension_url]")
		}
	}

	return nil
}

// createFilterResponse is the response body for POST /filters
type createFilterResponse struct {
	model.Filter
}

// createFilterRequest is the request body for POST /filters
type createFilterOutputsRequest struct {
	State     string             `bson:"state" json:"state"`
	Downloads model.FilterOutput `bson:"downloads"      json:"downloads"`
}

func isNotFullyPopulated(fi *model.FileInfo) error {

	cutset := " "

	if len(strings.Trim(fi.HREF, cutset)) == 0 {
		return errors.New(`"HREF" is empty in input`)
	}

	if len(strings.Trim(fi.Private, cutset)) == 0 {
		return errors.New(`"Private" is empty in input`)
	}

	if len(strings.Trim(fi.Public, cutset)) == 0 {
		return errors.New(`"Public" is empty in input`)
	}

	if len(strings.Trim(fi.Size, cutset)) == 0 {
		return errors.New(`"Size" is empty in input`)
	}

	return nil
}

func (r *createFilterOutputsRequest) Valid() error {

	if err := isNotFullyPopulated(r.Downloads.CSV); err != nil {
		return err
	}
	if err := isNotFullyPopulated(r.Downloads.CSVW); err != nil {
		return err
	}
	if err := isNotFullyPopulated(r.Downloads.TXT); err != nil {
		return err
	}
	if err := isNotFullyPopulated(r.Downloads.XLS); err != nil {
		return err
	}

	return nil
}

// createFilterResponse is the response body for POST /filters
type createFilterOutputsResponse struct {
	ID        uuid.UUID          `bson:"filter_output_id" json:"filter_output_id"`
	Downloads model.FilterOutput `bson:"downloads" json:"downloads"`
}

// getFilterDimensionsResponse is the response body for GET /filters/{id}/dimensions
type getFilterDimensionsResponse struct {
	Dimensions []model.Dimension `json:"dimensions"`
}
