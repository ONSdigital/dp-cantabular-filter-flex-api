package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	dperrors "github.com/ONSdigital/dp-cantabular-filter-flex-api/errors"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

func (api *API) addFilterDimensionOption(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := addFilterDimensionOptionRequest{
		FilterID:  chi.URLParam(r, "id"),
		Dimension: chi.URLParam(r, "dimension"),
		Option:    chi.URLParam(r, "option"),
	}

	logData := log.Data{
		"filter_id": req.FilterID,
		"dimension": req.Dimension,
		"option":    req.Option,
	}

	filter, err := api.store.GetFilter(ctx, req.FilterID)
	if err != nil {
		api.respond.Error(
			ctx,
			w,
			statusCode(err),
			Error{
				err:     errors.Wrap(err, "failed to get filter"),
				message: "failed to add dimension option: failed to get filter",
				logData: logData,
			},
		)
		return
	}

	// Check dimension exists
	var dimension model.Dimension
	var dimExists bool
	for _, d := range filter.Dimensions {
		if d.Name == req.Dimension {
			dimension = d
			dimExists = true
			break
		}
	}

	if !dimExists {
		api.respond.Error(
			ctx,
			w,
			http.StatusBadRequest,
			Error{
				err:     errors.New("failed to add dimension option: dimension not found in filter"),
				logData: logData,
			},
		)
		return
	}

	// Check if option already exists
	var optExists bool
	for _, o := range dimension.Options {
		if o == req.Option {
			optExists = true
			break
		}
	}

	if optExists {
		api.respond.Error(
			ctx,
			w,
			http.StatusBadRequest,
			Error{
				err:     errors.New("failed to add dimension option: option already added to dimension"),
				logData: logData,
			},
		)
		return
	}

	// Check option is valid
	var dReq cantabular.GetDimensionOptionsRequest
	dReq.DimensionNames = append(dReq.DimensionNames, dimension.ID)
	dReq.Filters = append(dReq.Filters, cantabular.Filter{
		Codes:    []string{req.Option},
		Variable: dimension.ID,
	})

	if _, err := api.ctblr.GetDimensionOptions(ctx, dReq); err != nil {
		if api.ctblr.StatusCode(err) >= http.StatusInternalServerError {
			err = Error{
				err:     errors.Wrap(err, "failed to query dimension options from Cantabular"),
				message: "Internal Server Error",
				logData: log.Data{
					"request": dReq,
				},
			}
		} else {
			err = Error{
				err:     errors.WithStack(err),
				message: "invalid option for filter",
			}
		}
		api.respond.Error(ctx, w, api.ctblr.StatusCode(err), err)
		return
	}

	var eTag string
	if reqETag := api.getETag(r); reqETag != eTagAny {
		eTag = reqETag
	}

	// Add option to filter
	dimension.Options = append(dimension.Options, req.Option)
	newETag, err := api.store.UpdateFilterDimension(ctx, req.FilterID, req.Dimension, dimension, eTag)
	if err != nil {
		api.respond.Error(ctx, w, statusCode(err), Error{
			err:     errors.Wrap(err, "failed to uodate dimension with option in store"),
			message: "failed to add dimension option",
			logData: logData,
		})
		return
	}

	resp := addFilterDimensionOptionResponse{
		Option: req.Option,
		Links: filterDimensionOptionLinks{
			Filter: filter.Links.Self,
			Self: model.Link{
				ID: "",
				HREF: fmt.Sprintf(
					"%s/filters/%s/dimensions/%s/options/%s",
					api.cfg.BindAddr,
					filter.ID,
					req.Dimension,
					req.Option,
				),
			},
			Dimension: model.Link{
				ID: dimension.ID,
				HREF: fmt.Sprintf(
					"%s/filters/%s/dimensions/%s",
					api.cfg.BindAddr,
					filter.ID,
					req.Dimension,
				),
			},
		},
	}

	w.Header().Set(eTagHeader, newETag)

	api.respond.JSON(ctx, w, http.StatusCreated, resp)
}

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

// deleteFilterDimensionOptions deletes all options on a given FilterOutput at once
func (api *API) deleteFilterDimensionOptions(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	filterID := chi.URLParam(r, "id")
	dimensionName := chi.URLParam(r, "name")

	filter, err := api.store.GetFilter(ctx, filterID)
	if err != nil {
		api.respond.Error(
			ctx,
			w,
			http.StatusNotFound,
			Error{
				err:     errors.Wrap(err, "failed to get filter"),
				message: "failed to delete option: filter not found",
				logData: log.Data{
					"id": filterID,
				},
			},
		)
		return
	}

	if eTag := api.getETag(r); eTag != eTagAny {
		if eTag != filter.ETag {
			api.respond.Error(
				ctx,
				w,
				http.StatusConflict,
				Error{
					err: errors.New("conflict: invalid ETag provided or filter has been updated"),
					logData: log.Data{
						"expected_etag": eTag,
						"actual_etag":   filter.ETag,
					},
				},
			)
		}
		return
	}

	eTag, err := api.store.DeleteFilterDimensionOptions(
		ctx,
		filterID,
		dimensionName,
	)
	if err != nil {
		code := statusCode(err)
		if dperrors.NotFound(err) {
			code = http.StatusBadRequest
		}
		api.respond.Error(
			ctx,
			w,
			code,
			errors.Wrap(err, "Failed to delete options"),
		)
		return
	}

	w.Header().Set(eTagHeader, eTag)
	api.respond.JSON(ctx, w, http.StatusNoContent, nil)
}
func parseFilterDimensionOptions(options []string, filterID, dimensionName string, address string) []GetFilterDimensionOptionsItem {
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
