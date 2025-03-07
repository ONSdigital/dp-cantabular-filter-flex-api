package api

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	dperrors "github.com/ONSdigital/dp-cantabular-filter-flex-api/errors"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	"github.com/ONSdigital/dp-net/v2/links"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
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
			},
		)
		return
	}

	// Check dimension exists
	var dimension model.Dimension
	var dimExists bool

	for i := range filter.Dimensions {
		d := filter.Dimensions[i]
		dName := d.Name
		if d.FilterByParent != "" {
			dName = d.FilterByParent
		}
		if req.Dimension == dName {
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
	dID := dimension.ID
	if dimension.FilterByParent != "" {
		dID = dimension.FilterByParent
	}
	dReq := cantabular.GetDimensionOptionsRequest{
		Dataset:        filter.PopulationType,
		DimensionNames: []string{dID},
		Filters: []cantabular.Filter{
			{
				Codes:    []string{req.Option},
				Variable: dID,
			},
		},
	}

	if _, err := api.ctblr.GetDimensionOptions(ctx, dReq); err != nil {
		logData["request"] = dReq
		if api.ctblr.StatusCode(err) >= http.StatusInternalServerError {
			err = Error{
				err:     errors.Wrap(err, "failed to query dimension options from Cantabular"),
				message: "Internal Server Error",
				logData: logData,
			}
		} else {
			err = Error{
				err:     errors.WithStack(err),
				message: "invalid option for filter",
				logData: logData,
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

func (api *API) deleteFilterDimensionOption(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := deleteFilterDimensionOptionRequest{
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
		status := statusCode(err)
		if dperrors.NotFound(err) {
			status = http.StatusBadRequest
		}
		api.respond.Error(
			ctx,
			w,
			status,
			Error{
				err:     errors.Wrap(err, "failed to get filter"),
				message: "failed to delete dimension option: failed to get filter",
			},
		)
		return
	}

	// Check dimension exists
	var dimension model.Dimension
	var dimExists bool

	for i := range filter.Dimensions {
		d := filter.Dimensions[i]
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
				err:     errors.New("failed to delete dimension option: dimension not found in filter"),
				logData: logData,
			},
		)
		return
	}

	// Check option exists
	var optExists bool
	for _, o := range dimension.Options {
		if o == req.Option {
			optExists = true
			break
		}
	}

	if !optExists {
		api.respond.Error(
			ctx,
			w,
			http.StatusNotFound,
			Error{
				err:     errors.New("failed to delete dimension option: option not found"),
				logData: logData,
			},
		)
		return
	}

	var eTag string
	if reqETag := api.getETag(r); reqETag != eTagAny {
		eTag = reqETag
	}

	log.Info(ctx, "removing option from filter dimension", logData)

	newETag, err := api.store.RemoveFilterDimensionOption(ctx, req.FilterID, req.Dimension, req.Option, eTag)
	if err != nil {
		api.respond.Error(ctx, w, statusCode(err), Error{
			err:     errors.Wrap(err, "failed to uodate dimension with option in store"),
			message: "failed to delete dimension option",
			logData: logData,
		})
		return
	}

	w.Header().Set(eTagHeader, newETag)

	api.respond.StatusCode(w, http.StatusNoContent)
}

func (api *API) getFilterDimensionOptions(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	r.Header.Set("X-Forwarded-Host", r.Header.Get("X-Forwarded-API-Host"))

	filterID := chi.URLParam(r, "id")
	dimensionName := chi.URLParam(r, "dimension")

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
		Items: parseFilterDimensionOptions(r, options, filterID, dimensionName, api.cfg.FilterAPIURL, api.cantabularFilterFlexAPIURL, api.cfg.EnableURLRewriting),
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
	dimensionName := chi.URLParam(r, "dimension")

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
			errors.Wrap(err, "failed to delete options"),
		)
		return
	}

	w.Header().Set(eTagHeader, eTag)
	api.respond.JSON(ctx, w, http.StatusNoContent, nil)
}

func parseFilterDimensionOptions(r *http.Request, options []string, filterID, dimensionName, address string, filterFlexAPIURL *url.URL, enableURLRewriting bool) []GetFilterDimensionOptionsItem {
	var err error
	responses := make([]GetFilterDimensionOptionsItem, 0)

	filterFlexLinksBuilder := links.FromHeadersOrDefault(&r.Header, filterFlexAPIURL)

	for _, option := range options {
		selfURL := fmt.Sprintf("%s/filters/%s/dimensions/%s/options", address, filterID, dimensionName)
		filterURL := fmt.Sprintf("%s/filters/%s", address, filterID)
		dimensionURL := fmt.Sprintf("%s/filters/%s/dimensions/%s", address, filterID, dimensionName)

		if enableURLRewriting {
			selfURL, err = filterFlexLinksBuilder.BuildLink(selfURL)
			if err != nil {
				log.Error(r.Context(), "failed to build self link", err, log.Data{"href": selfURL})
				return nil
			}
			filterURL, err = filterFlexLinksBuilder.BuildLink(filterURL)
			if err != nil {
				log.Error(r.Context(), "failed to build filter link", err, log.Data{"href": filterURL})
				return nil
			}
			dimensionURL, err = filterFlexLinksBuilder.BuildLink(dimensionURL)
			if err != nil {
				log.Error(r.Context(), "failed to build dimension link", err, log.Data{"href": dimensionURL})
				return nil
			}
		}

		addOptionResponse := GetFilterDimensionOptionsItem{
			Option: option,
			Self: model.Link{
				HREF: selfURL,
				ID:   option,
			},
			Filter: model.Link{
				HREF: filterURL,
				ID:   filterID,
			},
			Dimension: model.Link{
				HREF: dimensionURL,
				ID:   dimensionName,
			},
		}

		responses = append(responses, addOptionResponse)
	}

	return responses
}
