package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	dprequest "github.com/ONSdigital/dp-net/v2/request"
	"github.com/ONSdigital/log.go/v2/log"

	"github.com/pkg/errors"
)

const (
	flexible  = "flexible"
	published = "published"
)

func (api *API) createFilter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req createFilterRequest

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
	}

	v, err := api.datasets.GetVersion(
		ctx,
		"",
		api.cfg.ServiceAuthToken,
		"",
		"",
		req.Dataset.ID,
		req.Dataset.Edition,
		strconv.Itoa(req.Dataset.Version),
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

	dimIDs, err := api.validateDimensions(ctx, req.Dimensions, v.Dimensions)
	if err != nil {
		api.respond.Error(
			ctx,
			w,
			http.StatusBadRequest,
			Error{
				err:     errors.Wrap(err, "failed to validate request dimensions"),
				logData: logData,
			},
		)
		return
	}

	// Validate dimension options by performing Cantabular query with selections,
	// skip this check if requesting all options
	if err := api.validateDimensionOptions(ctx, req, dimIDs); err != nil {
		api.respond.Error(
			ctx,
			w,
			statusCode(err),
			Error{
				err:     errors.Wrap(err, "failed to validate dimension options"),
				logData: logData,
			},
		)
		return
	}

	f := model.Filter{
		Links: model.Links{
			Version: model.Link{
				HREF: api.generate.URL(
					api.cfg.DatasetAPIURL,
					"/datasets/%s/editions/%s/version/%d",
					req.Dataset.ID,
					req.Dataset.Edition,
					req.Dataset.Version,
				),
				ID: strconv.Itoa(v.Version),
			},
		},
		Dimensions:        req.Dimensions,
		UniqueTimestamp:   api.generate.UniqueTimestamp(),
		LastUpdated:       api.generate.Timestamp(),
		Dataset:           *req.Dataset,
		InstanceID:        v.ID,
		PopulationType:    req.PopulationType,
		Type:              flexible,
		Published:         v.State == published,
		Events:            nil, // TODO: Not sure what to
		DisclosureControl: nil, // populate for these fields yet
	}

	if err := api.store.CreateFilter(ctx, &f); err != nil {
		api.respond.Error(
			ctx,
			w,
			statusCode(err),
			Error{
				err:     errors.Wrap(err, "failed to create filter"),
				logData: logData,
			},
		)
		return
	}

	resp := createFilterResponse{
		Filter: f,
	}

	api.respond.JSON(ctx, w, http.StatusCreated, resp)
}

func (api *API) getFilter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fID := chi.URLParam(r, "id")

	logData := log.Data{
		"filter_id": fID,
	}

	f, err := api.store.GetFilter(ctx, fID)
	if err != nil {
		api.respond.Error(
			ctx,
			w,
			statusCode(err),
			Error{
				err:     errors.Wrap(err, "failed to get filter"),
				message: "failed to get filter",
				logData: log.Data{
					"id": fID,
				},
			},
		)
		return
	}

	if eTag := api.getETag(r); eTag != eTagAny {
		if eTag != f.ETag{
			api.respond.Error(
				ctx,
				w,
				http.StatusConflict,
				Error{
					err:     errors.New("conflict: invalid ETag provided or filter has been updated"),
					logData: log.Data{
						"expected_etag": eTag,
						"actual_etag": f.ETag,
					},
				},
			)
		}
		return
	}

	if !f.Published && !dprequest.IsCallerPresent(ctx) {
		api.respond.Error(
			ctx,
			w,
			http.StatusNotFound,
			Error{
				err:     errors.New("unauthenticated request on unpublished dataset"),
				message: "failed to get filter",
				logData: logData,
			},
		)
		return
	}

	resp := getFilterResponse{*f}

	api.respond.JSON(ctx, w, http.StatusOK, resp)
}

func (api *API) getFilterDimensions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fID := chi.URLParam(r, "id")

	logData := log.Data{"id": fID}

	limit, offset, err := getPaginationParams(r.URL, api.cfg.DefaultMaximumLimit, logData)
	if err != nil {
		api.respond.Error(ctx, w, http.StatusBadRequest, err)
		return
	}

	logData["limit"] = limit
	logData["offset"] = offset

	dimensions, totalCount, err := api.store.GetFilterDimensions(ctx, fID, limit, offset)
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

	resp := getFilterDimensionsResponse{
		Items: dimensions,
		paginationResponse: paginationResponse{
			Limit:      limit,
			Offset:     offset,
			Count:      len(dimensions),
			TotalCount: totalCount,
		},
	}

	api.respond.JSON(ctx, w, http.StatusOK, resp)
}

// validateDimensions validates provided filter dimensions exist within the dataset dimensions provided.
// Returns a map of the dimensions name:id for use in the following validation calls
func (api *API) validateDimensions(ctx context.Context, filterDims []model.Dimension, dims []dataset.VersionDimension) (map[string]string, error) {
	dimensions := make(map[string]string)
	for _, d := range dims {
		dimensions[d.Name] = d.ID
	}

	var incorrect []string
	for _, fd := range filterDims {
		if _, ok := dimensions[fd.Name]; !ok {
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

// validateDimensionOptions validates the dimension options in a createFilterRequest by making a call
// to Cantabular to check they exist. Takes as a second argument a map mapping the dimension names to
// ids, which are required by the Cantabular query
// TODO: Requires rule variable to be first in POST request, acceptable?
func (api *API) validateDimensionOptions(ctx context.Context, req createFilterRequest, dimIDs map[string]string) error {
	dReq := cantabular.GetDimensionOptionsRequest{
		Dataset: req.PopulationType,
	}

	// only validate dimensions with specified options
	for _, d := range req.Dimensions {
		if len(d.Options) > 0 {
			dReq.DimensionNames = append(dReq.DimensionNames, dimIDs[d.Name])
			dReq.Filters = append(dReq.Filters, cantabular.Filter{
				Codes:    d.Options,
				Variable: dimIDs[d.Name],
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
