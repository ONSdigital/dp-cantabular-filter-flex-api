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
		for i := range dims {
			d := dims[i]
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

// Used for POST Dimension when the filter type is multivariate.
// Extra logic required to hydrate dimension using the default categorisation of the dimension.
func (api *API) HydrateMultivariateDimensionsPOST(dimensions []model.Dimension, pType string) ([]model.Dimension, error) {
	ctx := context.Background()
	hydratedDimensions := make([]model.Dimension, 0)

	for i := range dimensions {
		dim := &dimensions[i]

		node, err := api.getCantabularDimension(ctx, pType, dim.Name)
		if err != nil {
			return nil, errors.Wrap(err, "error in cantabular response")
		}

		finalDimension, finalLabel, finalCategorisation, err := api.RetrieveDefaultCategorisation(node, pType)
		if err != nil {
			return nil, errors.Wrap(err, "failed to hydrate multivariate dimensions")
		}

		if finalDimension != dim.Name || dim.Options == nil {
			dim.Options = []string{}
		}
		hydratedDimensions = append(hydratedDimensions, model.Dimension{
			Name:                  finalDimension,
			ID:                    finalDimension,
			Label:                 finalLabel,
			DefaultCategorisation: finalCategorisation,
			Options:               dim.Options,
			IsAreaType:            dim.IsAreaType,
			FilterByParent:        dim.FilterByParent,
			QualityStatementText:  node.QualityStatementText,
			QualitySummaryURL:     node.QualitySummaryURL,
		})
	}
	return hydratedDimensions, nil
}

// getCantabularDimensions pulls full dimension information from Cantabular using the names of the provided
// dimensions.
// NOTE: when we hydrate the dimensions, we will be using the name as the id, and filling out the dimensions
// using the same value for both.
func (api *API) getCantabularDimensions(ctx context.Context, dimensions []model.Dimension, pType string) ([]model.Dimension, error) {
	hydratedDimensions := make([]model.Dimension, 0)

	for i := range dimensions {
		d := &dimensions[i]
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
		dim.QualityStatementText = d.QualityStatementText
		dim.QualitySummaryURL = d.QualitySummaryURL
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
		Label:                node.Label,
		ID:                   dimensionName,
		Name:                 dimensionName,
		QualityStatementText: node.Meta.ONSVariable.QualityStatementText,
		QualitySummaryURL:    node.Meta.ONSVariable.QualitySummaryURL,
	}

	return &dim, nil
}

// validateDimensionsFromVersion validates provided dimensions exist within the version dimensions provided.
func (api *API) validateDimensionsFromVersion(dims []model.Dimension, versionDims []dataset.VersionDimension) error {
	fDims := make(map[string]bool)

	for i := range dims {
		d := &dims[i]
		if _, ok := fDims[d.Name]; ok {
			return Error{
				err:        errors.New("duplicate dimensions chosen"),
				message:    "duplicate dimension chosen: " + d.Name,
				badRequest: true,
				logData: log.Data{
					"duplicate_dimension": d.Name,
				},
			}
		}

		fDims[d.Name] = true
	}

	dimensions := make(map[string]string)
	for i := range versionDims {
		vd := &versionDims[i]
		dimensions[vd.Name] = vd.ID
	}

	incorrect := make([]string, 0, len(dims))
	for i := range dims {
		d := &dims[i]
		// allow geography dimensions other than default
		if d.IsAreaType != nil && *d.IsAreaType {
			continue
		}
		if _, ok := dimensions[d.Name]; !ok {
			incorrect = append(incorrect, d.Name)
			continue
		}
	}

	if len(incorrect) > 0 {
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
	for i := range filterDimensions {
		d := filterDimensions[i]
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

func getFilterVariable(d model.Dimension) string {
	if d.FilterByParent != "" {
		return d.FilterByParent
	}
	return d.Name
}

// hydrateDimensionsFromDataset adds additional data (id/label) to a model.Dimension, using values provided by the dataset.
func hydrateDimensionsFromVersion(filterDims []model.Dimension, dims []dataset.VersionDimension) []model.Dimension {
	type record struct{ id, label string }

	lookup := make(map[string]record)
	for i := range dims {
		d := &dims[i]
		// geography dimension gets hydrated from Cantabular
		if d.IsAreaType != nil && *d.IsAreaType {
			continue
		}
		lookup[d.Name] = record{id: d.ID, label: d.Label}
	}

	hydrated := make([]model.Dimension, 0, len(filterDims))
	for i := range filterDims {
		d := filterDims[i]
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
