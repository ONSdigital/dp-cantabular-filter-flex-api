package api

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	dperrors "github.com/ONSdigital/dp-cantabular-filter-flex-api/errors"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	dprequest "github.com/ONSdigital/dp-net/v2/request"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

func (api *API) addFilterDimension(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fID := chi.URLParam(r, "id")

	var req addFilterDimensionRequest

	if err := api.ParseRequest(r.Body, &req); err != nil {
		api.respond.Error(
			ctx,
			w,
			http.StatusBadRequest,
			errors.Wrap(err, "failed to parse request"),
		)
		return
	}

	logData := log.Data{
		"request": req,
		"id":      fID,
	}

	filter, err := api.store.GetFilter(ctx, fID)
	if err != nil {
		api.respond.Error(
			ctx,
			w,
			statusCode(err),
			Error{
				err:     errors.Wrap(err, "failed to get filter"),
				message: "failed to get filter",
				logData: logData,
			},
		)
		return
	}

	v, err := api.datasets.GetVersion(
		ctx,
		"",
		api.cfg.ServiceAuthToken,
		"",
		"",
		filter.Dataset.ID,
		filter.Dataset.Edition,
		strconv.Itoa(filter.Dataset.Version),
	)
	if err != nil {
		api.respond.Error(
			ctx,
			w,
			statusCode(err),
			Error{
				err:     errors.Wrap(err, "failed to get existing Version"),
				message: "failed to get existing dataset information",
				logData: logData,
			},
		)
		return
	}

	if v.State != published && !dprequest.IsCallerPresent(ctx) {
		api.respond.Error(
			ctx,
			w,
			http.StatusNotFound,
			Error{
				err:     errors.New("unauthenticated request on unpublished dataset"),
				message: "dataset not found",
				logData: logData,
			},
		)
		return
	}

	if err := api.isValidDatasetDimensions(ctx, v, []model.Dimension{req.Dimension}, filter.PopulationType); err != nil {
		api.respond.Error(ctx, w, statusCode(err), err)
		return
	}

	h, err := filter.HashDimensions()
	if err != nil {
		api.respond.Error(
			ctx,
			w,
			http.StatusInternalServerError,
			Error{
				err:     errors.Wrap(err, "failed to hash existing filter dimensions"),
				logData: logData,
			},
		)
		return
	}

	if eTag := api.getETag(r); eTag != eTagAny && !strings.Contains(eTag, h) {
		api.respond.Error(
			ctx,
			w,
			http.StatusConflict,
			Error{
				err:     errors.Wrap(err, "ETag does not match"),
				logData: logData,
			},
		)
		return
	}

	dim := hydrateDimensions([]model.Dimension{req.Dimension}, v.Dimensions)[0]

	if err := api.store.AddFilterDimension(ctx, fID, dim); err != nil {
		api.respond.Error(
			ctx,
			w,
			statusCode(err),
			Error{
				err:     errors.Wrap(err, "failed to add filter dimension"),
				logData: logData,
			},
		)
		return
	}

	filter, err = api.store.GetFilter(ctx, fID)
	if err != nil {
		api.respond.Error(
			ctx,
			w,
			statusCode(err),
			Error{
				err:     errors.Wrap(err, "failed to get updated filter"),
				message: "failed to get updated filter",
				logData: logData,
			},
		)
		return
	}

	b, err := filter.HashDimensions()
	if err != nil {
		api.respond.Error(
			ctx,
			w,
			http.StatusInternalServerError,
			Error{
				err:     errors.Wrap(err, "failed to hash filter dimensions"),
				logData: logData,
			},
		)
		return
	}

	var resp addFilterDimensionResponse
	resp.dimensionItem.fromDimension(
		dim,
		api.cfg.FilterAPIURL,
		fID,
	)

	w.Header().Set(eTagHeader, b)
	api.respond.JSON(ctx, w, http.StatusCreated, resp)
}

func (api *API) updateFilterDimension(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	filterID := chi.URLParam(r, "id")
	dimensionName := chi.URLParam(r, "dimension")

	logData := log.Data{
		"filter_id":      filterID,
		"dimension_name": dimensionName,
	}

	req := updateFilterDimensionRequest{
		Dimension: model.Dimension{
			Name:    dimensionName,
			Options: []string{},
		},
	}

	if err := api.ParseRequest(r.Body, &req); err != nil {
		api.respond.Error(
			ctx,
			w,
			http.StatusBadRequest,
			Error{
				err:     errors.Wrap(err, "failed to parse update filter request"),
				logData: logData,
			},
		)
		return
	}

	// The new dimension won't be present on the dataset (i.e. only `City` will be present, not `Country`),
	// so instead of validating it against the existing Version, we check to see if the dimension exists in Cantabular.
	// TODO: this function gets the dimension via a search, not guaranteed to be correct dimension
	node, err := api.getCantabularDimension(ctx, filterID, req.ID)
	if err != nil {
		api.respond.Error(
			ctx,
			w,
			statusCode(err),
			Error{
				err:     errors.Wrap(err, "error searching for dimension"),
				logData: logData,
			},
		)
		return
	}

	// ID/name is provided by the request, but the label is provided by Cantabular.
	req.Label = node.Label

	var eTag string
	if reqETag := api.getETag(r); reqETag != eTagAny {
		eTag = reqETag
	}

	newETag, err := api.store.UpdateFilterDimension(ctx, filterID, dimensionName, req.Dimension, eTag)
	if err != nil {
		api.respond.Error(
			ctx,
			w,
			statusCode(err),
			Error{
				err:     errors.Wrap(err, "unable to update filter dimension"),
				logData: logData,
			},
		)
		return
	}

	resp := updateFilterDimensionResponse{}
	resp.fromDimension(req.Dimension, api.cfg.FilterAPIURL, filterID)

	w.Header().Set(eTagHeader, newETag)

	api.respond.JSON(ctx, w, http.StatusOK, resp)
}

func (api *API) getFilterDimension(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fID := chi.URLParam(r, "id")
	dim := chi.URLParam(r, "dimension")

	logData := log.Data{
		"id":        fID,
		"dimension": dim,
	}

	// We decode the dimension name since currently dimensions are stored using their pretty name, e.g.
	// `Number of Siblings`, and passed in the URL as encoded (e.g. `Number+of+Siblings`). Until this is
	// changed we need to unescape the dimension before querying.
	dimName, err := url.QueryUnescape(dim)
	if err != nil {
		api.respond.Error(
			ctx,
			w,
			statusCode(err),
			Error{
				err:     errors.Wrap(err, "failed to decode dimension name"),
				message: "failed to decode dimension name",
				logData: logData,
			},
		)
		return
	}

	// Check the filter exists, so we can return a different status code
	// from if the dimension doesn't exist.
	if _, err := api.store.GetFilter(ctx, fID); err != nil {
		status := statusCode(err)
		if dperrors.NotFound(err) {
			status = http.StatusBadRequest
		}

		api.respond.Error(
			ctx,
			w,
			status,
			Error{
				err:     errors.Wrap(err, "failed to get filter from store"),
				message: "failed to get filter",
				logData: logData,
			},
		)
		return
	}

	filterDim, err := api.store.GetFilterDimension(ctx, fID, dimName)
	if err != nil {
		api.respond.Error(
			ctx,
			w,
			statusCode(err),
			Error{
				err:     errors.Wrap(err, "failed to get filter dimension from store"),
				message: "failed to get filter dimension",
				logData: logData,
			},
		)
		return
	}

	var resp getFilterDimensionResponse
	resp.dimensionItem.fromDimension(filterDim, api.cfg.FilterAPIURL, fID)
	resp.IsAreaType = filterDim.IsAreaType

	api.respond.JSON(ctx, w, http.StatusOK, resp)
}

func (api *API) getFilterDimensions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fID := chi.URLParam(r, "id")

	logData := log.Data{"id": fID}

	limit, offset, err := getPaginationParams(r.URL, api.cfg.DefaultMaximumLimit)
	if err != nil {
		api.respond.Error(ctx, w, http.StatusBadRequest, &Error{
			err:     err,
			logData: logData,
		})
		return
	}

	logData["limit"] = limit
	logData["offset"] = offset

	dims, totalCount, err := api.store.GetFilterDimensions(ctx, fID, limit, offset)
	if err != nil {
		api.respond.Error(
			ctx,
			w,
			statusCode(err),
			Error{
				err:     errors.Wrap(err, "failed to get filter dimensions"),
				message: "failed to get filter dimensions",
				logData: logData,
			},
		)
		return
	}

	var items dimensionItems
	items.fromDimensions(dims, api.cfg.FilterAPIURL, fID)

	resp := getFilterDimensionsResponse{
		Items: items,
		paginationResponse: paginationResponse{
			Limit:      limit,
			Offset:     offset,
			Count:      len(dims),
			TotalCount: totalCount,
		},
	}

	api.respond.JSON(ctx, w, http.StatusOK, resp)
}
