package api

import (
	"context"
	"net/http"

	dperrors "github.com/ONSdigital/dp-cantabular-filter-flex-api/errors"
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
		pageLimit = 20
	}

	var eTag string
	if reqETag := api.getETag(r); reqETag != eTagAny {
		eTag = reqETag
	}

	options, totalCount, err := api.store.GetFilterDimensionOptions(
		ctx,
		filterID,
		dimensionName,
		pageLimit,
		offset,
	)
	if err != nil {
		status := statusCode(err)
		finalErr := errors.New("internal server error")

		if dperrors.NotFound(err) {
			status = http.StatusNotFound
			finalErr = err
		}
		if dperrors.BadRequest(err) {
			status = http.StatusBadRequest
			finalErr = err
		}

		api.respond.Error(
			ctx,
			w,
			status,
			errors.Wrap(finalErr, "Failed to get options"),
		)
		return
	}

	resp := GetFilterDimensionOptionsResponse{
		Items: parseFilterDimensionOptions(options, filterID, dimensionName),
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

func parseFilterDimensionOptions(options []string, filterID, dimensionName string) []AddOptionResponse {
	responses := make([]AddOptionResponse, 0)

	for _, option := range options {
		addOptionResponse := AddOptionResponse{
			Option: option,
			Self: model.Link{
				HREF: "/filters/" + filterID + "/dimensions/" + dimensionName + "/options",
				ID:   option,
			},
			Filter: model.Link{
				HREF: "/filters/" + filterID,
				ID:   filterID,
			},
			Dimension: model.Link{
				HREF: "/filters/" + filterID + "/dimensions/" + dimensionName,
				ID:   dimensionName,
			},
		}

		responses = append(responses, addOptionResponse)
	}

	return responses
}
