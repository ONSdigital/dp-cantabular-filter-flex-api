package api

import (
	"context"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
)

//  ValidateAndReturnDimensions encapsulates the dimension validation paths for filters based on flexible and multivariate tables.
func (api *API) ValidateAndReturnDimensions(v dataset.Version, dimensions []model.Dimension, populationType string) (finalDims []model.Dimension, filterType string, err error) {
	ctx := context.Background()

	if v.IsBasedOn.Type == cantabularFlexible {
		filterType = flexible
		if err = api.isValidDatasetDimensions(ctx, v, dimensions, populationType); err != nil {
			return
		}

		finalDims = hydrateDimensions(dimensions, v.Dimensions)

	} else if v.IsBasedOn.Type == cantabularMultivariate {
		filterType = multivariate
		multivariateDims, err := api.isValidMultivariateDimensions(ctx, dimensions, populationType)
		if err != nil {
			return finalDims, "", err
		}

		finalDims = multivariateDims

	}

	return
}
