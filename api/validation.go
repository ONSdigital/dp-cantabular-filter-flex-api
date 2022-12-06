package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/pkg/errors"
)

// ValidateAndHydrateDimensions performs validation against the provided dimensions and hydrates missing fields (id/name/label).
// Flexible table types will validate by checking the exisisting variables in the dataset version. Multivariate tables
// will use Cantabular to check the new dimensions exist.
func (api *API) validateAndHydrateDimensions(v dataset.Version, dims []model.Dimension, pType string) ([]model.Dimension, error) {
	ctx := context.Background()

	if len(dims) < 1 {
		return nil, errors.New("no dimensions given")
	}

	if v.IsBasedOn.Type == cantabularFlexibleTable {
		var geodim model.Dimension
		var areaCount int
		for _, d := range dims {
			if d.IsAreaType != nil && *d.IsAreaType {
				areaCount++
				geodims, err := api.getCantabularDimensions(ctx, []model.Dimension{d}, pType)
				if err != nil {
					return nil, errors.Wrap(err, "failed to get geography dimension from Cantabular")
				}
				geodim = geodims[0]
			}
		}

		if areaCount > 1 {
			return nil, &Error{
				err:        errors.New("multiple geography dimensions not permitted"),
				badRequest: true,
				logData: log.Data{
					"dimensions": dims,
				},
			}
		}

		if err := api.validateDimensionsFromVersion(dims, v.Dimensions); err != nil {
			return nil, errors.Wrap(err, "failed to validate dataset dimensions")
		}

		hydrated := hydrateDimensionsFromVersion(dims, v.Dimensions)
		// insert Geography dimension as first in list if present
		if areaCount > 0 {
			hydrated = append([]model.Dimension{geodim}, hydrated...)
		}

		return hydrated, nil
	}

	if v.IsBasedOn.Type == cantabularMultivariateTable {
		resp, err := api.getCantabularDimensions(ctx, dims, pType)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get dimensions from Cantabular")
		}
		return resp, nil
	}

	return nil, &Error{
		err:        errors.New("unexpected IsBasedOn type"),
		logData:    log.Data{"is_based_on.type": v.IsBasedOn.Type},
		badRequest: true,
	}
}

// hydrateDimensions adds additional data (id/label) to a model.Dimension, using values provided by the dataset.
func hydrateDimensions(filterDims []model.Dimension, dims []dataset.VersionDimension) []model.Dimension {
	type record struct{ id, label string }

	lookup := make(map[string]record)
	for _, dim := range dims {
		lookup[dim.Name] = record{id: dim.ID, label: dim.Label}
	}

	var hydrated []model.Dimension
	for _, dim := range filterDims {
		dim.ID = lookup[dim.Name].id
		dim.Label = lookup[dim.Name].label
		if dim.Options == nil {
			dim.Options = []string{}
		}
		hydrated = append(hydrated, dim)
	}

	return hydrated
}

// validateDimensions validates provided filter dimensions exist within the dataset dimensions provided.
// Returns a map of the dimensions name:id for use in the following validation calls
func (api *API) validateDimensions(filterDims []model.Dimension, dims []dataset.VersionDimension, datasetType string) (map[string]string, error) {

	fDims := make(map[string]bool)
	for _, fd := range filterDims {
		if _, ok := fDims[fd.Name]; ok {
			return nil, Error{
				err: errors.Errorf("duplicate dimensions chosen: %v", fd.Name),
				logData: log.Data{
					"duplicate dimensions chosen": fd.Name,
				},
			}
		} else {
			fDims[fd.Name] = true
		}
	}

	dimensions := make(map[string]string)
	for _, d := range dims {
		dimensions[d.Name] = d.ID

	}
	var incorrect []string
	for _, fd := range filterDims {
		if _, ok := dimensions[fd.Name]; !ok {
			// if this is a cantabular multivariate table, then you
			// should be free to create whatever dimension that you
			// want to add, provided that they exist. which is checked
			// in dimension option validation.
			if datasetType == "cantabular_multivariate_table" {
				continue
			}
			incorrect = append(incorrect, fd.Name)
		}
	}

	if incorrect != nil {
		return nil, Error{
			err: errors.Errorf("incorrect dimensions chosen: %v", incorrect),
			logData: log.Data{
				"available_dimensions": dimensions,
			},
		}
	}

	return dimensions, nil
}

func (api *API) isValidDatasetDimensions(ctx context.Context, v dataset.Version, d []model.Dimension, pType string) error {
	/*	dimIDs, err := api.validateDimensions(d, v.Dimensions, pType)
		if err != nil {
			return Error{
				err:      errors.Wrap(err, "failed to validate request dimensions"),
				notFound: true,
			}
		}

		if err := api.validateDimensionOptionsNew(ctx, d, dimIDs, pType); err != nil {
			return errors.Wrap(err, "failed to validate dimension options")
		}
	*/
	if err := api.validateDimensionOptions(ctx, d, pType); err != nil {
		return errors.Wrap(err, "failed to validate dimension options")
	}
	return nil
}

// ValidateAndReturnDimensions encapsulates the dimension validation paths for filters based on flexible and multivariate tables.
func (api *API) ValidateAndReturnDimensions(v dataset.Version, dimensions []model.Dimension, populationType string, postDimension bool) (finalDims []model.Dimension, filterType string, err error) {
	ctx := context.Background()

	if v.IsBasedOn.Type == cantabularFlexibleTable {
		filterType = flexible
		if err = api.isValidDatasetDimensions(ctx, v, dimensions, populationType); err != nil {
			return
		}

		finalDims = hydrateDimensions(dimensions, v.Dimensions)

	} else if v.IsBasedOn.Type == cantabularMultivariateTable {
		println("CREATING FILTER")
		filterType = multivariate
		multivariateDims, err := api.isValidMultivariateDimensions(ctx, dimensions, populationType, postDimension)
		if err != nil {
			return finalDims, "", err
		}

		finalDims = multivariateDims

	}

	return
}

// getCantabularDimensions pulls full dimension information from Cantabular using the names of the provided
// dimensions.
// NOTE: when we hydrate the dimensions, we will be using the name as the id, and filling out the dimensions
// using the same value for both.
func (api *API) getCantabularDimensions(ctx context.Context, dimensions []model.Dimension, pType string) ([]model.Dimension, error) {
	hydratedDimensions := make([]model.Dimension, 0)

	for _, d := range dimensions {
		dim, err := api.getCantabularDimension(ctx, pType, d.Name)
		if err != nil {
			return nil, Error{
				err:     errors.Wrap(err, "failed to get dimension"),
				message: "failed to find dimension: " + d.Name,
				logData: log.Data{
					"dimension": d.Name,
				},
			}
		}

		dim.IsAreaType = d.IsAreaType
		dim.FilterByParent = d.FilterByParent
		dim.Options = d.Options
		if dim.Options == nil {
			dim.Options = []string{}
		}
		hydratedDimensions = append(hydratedDimensions, *dim)
	}

	return hydratedDimensions, nil
}

// getCantabularDimension checks that dimension exists in Cantabular by searching for it.
// If the dimension doesn't exist, or couldn't be retrieved, an error is returned.
func (api *API) getCantabularDimension(ctx context.Context, popType, dimensionName string) (*model.Dimension, error) {
	resp, err := api.ctblr.GetDimensionsByName(ctx, cantabular.GetDimensionsByNameRequest{
		Dataset:        popType,
		DimensionNames: []string{dimensionName},
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get dimension by name")
	}

	if len(resp.Dataset.Variables.Edges) == 0 {
		return nil, Error{
			err:      errors.New("no dimensions in response"),
			notFound: true,
			logData:  log.Data{"response": resp},
		}
	}

	node := resp.Dataset.Variables.Edges[0].Node
	dim := model.Dimension{
		Label: node.Label,
		ID:    dimensionName,
		Name:  dimensionName,
	}

	return &dim, nil
}

// validateDimensionsFromVersion validates provided dimensions exist within the version dimensions provided.
func (api *API) validateDimensionsFromVersion(dims []model.Dimension, versionDims []dataset.VersionDimension) error {
	fDims := make(map[string]bool)

	for _, d := range dims {
		if _, ok := fDims[d.Name]; ok {
			return Error{
				err:        errors.New("duplicate dimensions chosen"),
				message:    "duplicate dimension chosen: " + d.Name,
				badRequest: true,
				logData: log.Data{
					"duplicate_dimension": d.Name,
				},
			}
		} else {
			fDims[d.Name] = true
		}
	}

	dimensions := make(map[string]string)
	for _, vd := range versionDims {
		dimensions[vd.Name] = vd.ID
	}

	var incorrect []string
	for _, d := range dims {
		// allow geography dimensions other than default
		if d.IsAreaType != nil && *d.IsAreaType {
			continue
		}
		if _, ok := dimensions[d.Name]; !ok {
			incorrect = append(incorrect, d.Name)
			continue
		}
	}

	if incorrect != nil {
		return Error{
			err:      errors.New("incorrect dimensions chosen"),
			message:  fmt.Sprintf("incorrect dimension chosen: %s", incorrect),
			notFound: true,
			logData: log.Data{
				"available_dimensions": dimensions,
			},
		}
	}

	return nil
}

// validateDimensionOptions by performing Cantabular query with selections,
// will be skipped if requesting all options
func (api *API) validateDimensionOptions(ctx context.Context, filterDimensions []model.Dimension, populationType string) error {
	dReq := cantabular.GetDimensionOptionsRequest{
		Dataset: populationType,
	}
	for _, d := range filterDimensions {
		if len(d.Options) > 0 {
			dReq.DimensionNames = append(dReq.DimensionNames, d.Name)
			dReq.Filters = append(dReq.Filters, cantabular.Filter{
				Codes:    d.Options,
				Variable: getFilterVariable(d),
			})
		}
	}
	if len(dReq.Filters) == 0 {
		return nil
	}

	if _, err := api.ctblr.GetDimensionOptions(ctx, dReq); err != nil {
		if api.ctblr.StatusCode(err) >= http.StatusInternalServerError {
			return Error{
				err:     errors.Wrap(err, "failed to query dimension options from Cantabular"),
				message: "Internal Server Error",
				logData: log.Data{
					"request": dReq,
				},
			}
		}
		return Error{
			err:     errors.WithStack(err),
			message: "failed to validate dimension options for filter",
		}
	}

	return nil
}

/* // validateDimensionOptions by performing Cantabular query with selections,
// will be skipped if requesting all options
func (api *API) validateDimensionOptionsNew(ctx context.Context, filterDimensions []model.Dimension, dimIDs map[string]string, populationType string) error {
	dReq := cantabular.GetDimensionOptionsRequest{
		Dataset: populationType,
	}
	for _, d := range filterDimensions {
		if len(d.Options) > 0 {
			dReq.DimensionNames = append(dReq.DimensionNames, dimIDs[d.Name])
			dReq.Filters = append(dReq.Filters, cantabular.Filter{
				Codes:    d.Options,
				Variable: getFilterVariable(d),
			})
		}
	}
	if len(dReq.Filters) == 0 {
		return nil
	}

	if _, err := api.ctblr.GetDimensionOptions(ctx, dReq); err != nil {
		if api.ctblr.StatusCode(err) >= http.StatusInternalServerError {
			return Error{
				err:     errors.Wrap(err, "failed to query dimension options from Cantabular"),
				message: "Internal Server Error",
				logData: log.Data{
					"request": dReq,
				},
			}
		}
		return Error{
			err:     errors.WithStack(err),
			message: "failed to validate dimension options for filter",
		}
	}

	return nil
} */

func getFilterVariable(d model.Dimension) string {
	if len(d.FilterByParent) != 0 {
		return d.FilterByParent
	}
	return d.Name
}

// hydrateDimensionsFromDataset adds additional data (id/label) to a model.Dimension, using values provided by the dataset.
func hydrateDimensionsFromVersion(filterDims []model.Dimension, dims []dataset.VersionDimension) []model.Dimension {
	type record struct{ id, label string }

	lookup := make(map[string]record)
	for _, d := range dims {
		// geography dimension gets hydrated from Cantabular
		if d.IsAreaType != nil && *d.IsAreaType {
			continue
		}
		lookup[d.Name] = record{id: d.ID, label: d.Label}
	}

	var hydrated []model.Dimension
	for _, d := range filterDims {
		// geography dimension gets hydrated from Cantabular
		if d.IsAreaType != nil && *d.IsAreaType {
			continue
		}
		d.ID = lookup[d.Name].id
		d.Label = lookup[d.Name].label
		if d.Options == nil {
			d.Options = []string{}
		}
		hydrated = append(hydrated, d)
	}

	return hydrated
}

/*
isValidMultivariateDimensions checks the validity of the supplied dimensions for a multivariate filter.
Supplied dimensions may not be in original dataset but still valid, and so isValidDatasetDimensions is
not relevant.
NOTE: when we hydrate the dimensions, we will be using the name as the id, and filling out the dimensions
using the same value for both.
*/
func (api *API) isValidMultivariateDimensions(ctx context.Context, dimensions []model.Dimension, pType string, postDimension bool) ([]model.Dimension, error) {
	hydratedDimensions := make([]model.Dimension, 0)
	var finalLabel string

	for _, dim := range dimensions {
		finalDimension := dim.Name
		node, err := api.getCantabularDimension(ctx, pType, dim.Name)
		if err != nil {
			return nil, errors.Wrap(err, "error in cantabular response")
		}

		if postDimension {

			finalDimension, finalLabel, err = api.CheckDefaultCategorisation(node.Name, pType)
			if err != nil {
				return nil, err
			}

		}

		if dim.Options == nil {
			dim.Options = []string{}
		}

		hydratedDimensions = append(hydratedDimensions, model.Dimension{
			Name:           finalDimension,
			ID:             finalDimension,
			Label:          finalLabel,
			Options:        dim.Options,
			IsAreaType:     dim.IsAreaType,
			FilterByParent: dim.FilterByParent,
		})

	}

	return hydratedDimensions, nil
}
