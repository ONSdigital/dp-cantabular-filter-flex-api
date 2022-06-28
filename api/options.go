package api

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

func (api *API) getFilterDimensionOptions(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	filterID := chi.URLParam(r, "id")
	dimensionName := chi.URLParam(r, "name")

	pageLimit, offset, err := getPaginationParams(r.URL, api.cfg.DefaultMaximumLimit)
	if err != nil {
		api.respond.Error(
			ctx,
			w,
			http.StatusBadRequest,
			errors.Wrap(err, "Bad Request"),
		)

	}

	if pageLimit == 0 {
		// define a reasonable default
		// in light of bad input
		// also slice not work with 0
		pageLimit = DefaultLimit
	}

	options, totalCount, eTag, err := api.store.GetFilterDimensionOptions(
		ctx,
		filterID,
		dimensionName,
		pageLimit,
		offset,
	)
	if err != nil {
		code := statusCode(err)
		if totalCount == -1 {
			code = http.StatusBadRequest
		}

		api.respond.Error(
			ctx,
			w,
			code,
			Error{
				err:     errors.Wrap(err, "failed to get filter dimension options"),
				message: "failed to get filter dimension option",
			},
		)
		return
	}

	resp := GetFilterDimensionOptionsResponse{
		Items: parseFilterDimensionOptions(options, filterID, dimensionName, api.cfg.FilterAPIURL),
		paginationResponse: paginationResponse{
			Limit:      pageLimit,
			Offset:     offset,
			Count:      len(options),
			TotalCount: totalCount,
		},
	}

	w.Header().Set(eTagHeader, eTag)
	api.respond.JSON(ctx, w, http.StatusOK, resp)
}

func parseFilterDimensionOptions(options []string, filterID, dimensionName, address string) []GetFilterDimensionOptionsItem {
	responses := make([]GetFilterDimensionOptionsItem, 0)

	for _, option := range options {
		addOptionResponse := GetFilterDimensionOptionsItem{
			Option: option,
			Self: model.Link{
				HREF: address + "/filters/" + filterID + "/dimensions/" + dimensionName + "/options",
				ID:   option,
			},
			Filter: model.Link{
				HREF: address + "/filters/" + filterID,
				ID:   filterID,
			},
			Dimension: model.Link{
				HREF: address + "/filters/" + filterID + "/dimensions/" + dimensionName,
				ID:   dimensionName,
			},
		}

		responses = append(responses, addOptionResponse)
	}

	return responses
}
