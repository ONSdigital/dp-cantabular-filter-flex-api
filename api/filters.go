package api

import (
	"context"
	"crypto/sha1"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	dprequest "github.com/ONSdigital/dp-net/v2/request"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"

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

	if !api.isValidDatasetDimensions(w, ctx, logData, v, req.Dimensions, req.PopulationType) {
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
		model.JobState{
			InstanceID:       f.InstanceID,
			DimensionListUrl: fmt.Sprintf("%s/filters/%s/dimensions", api.cfg.BindAddr, f.ID),
			FilterID:         f.ID,
			Events:           f.Events,
		},
		f.Links,
		f.Dataset,
		f.PopulationType,
	}

	api.respond.JSON(ctx, w, http.StatusCreated, resp)
}

func (api *API) postFilter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	filterID := chi.URLParam(r, "id")
	postTime, _ := time.Parse(time.RFC3339, "2016-07-17T08:38:25.316Z")

	filter, error := api.store.GetFilter(ctx, filterID)
	if err != nil {
		// error 500
	}

	error := api.store.CreateFilterOutput(ctx, nil)
	if err != nil {
		// error 500
	}

	// send the export event through Kafka
	// return a mock response for now.
	resp := updateFilterResponse{
		model.JobState{
			InstanceID: "",
			FilterID:   filterID,
			Events: []model.Event{
				{
					Timestamp: postTime,
					Name:      "mock-export-event",
				},
			},
		},

		model.Dataset{
			ID:      "mock-id",
			Edition: "mock-edition",
			Version: 0,
		},
		model.Links{
			Version: model.Link{
				HREF: "",
				ID:   "",
			},
			Self: model.Link{
				HREF: "",
				ID:   "",
			},
		},
	}

	api.respond.JSON(ctx, w, http.StatusAccepted, resp)
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
		if eTag != f.ETag {
			api.respond.Error(
				ctx,
				w,
				http.StatusConflict,
				Error{
					err: errors.New("conflict: invalid ETag provided or filter has been updated"),
					logData: log.Data{
						"expected_etag": eTag,
						"actual_etag":   f.ETag,
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

func (api *API) addFilterDimension(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fID := chi.URLParam(r, "id")

	newDimension := model.Dimension{
		Options: make([]string, 0),
	}

	if err := api.ParseRequest(r.Body, &newDimension); err != nil {
		api.respond.Error(
			ctx,
			w,
			http.StatusBadRequest,
			errors.Wrap(err, "failed to parse request"),
		)
		return
	}

	logData := log.Data{
		"request": newDimension,
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

	hashedFilterDimensions, err := api.hashFilterDimensions(ctx, fID)
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

	ifMatch := r.Header.Get("If-Match") // e.g. `"asdf", "qwer", "1234"`
	if ifMatch != "" && ifMatch != eTagAny && !strings.Contains(ifMatch, hashedFilterDimensions) {
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

	if !api.isValidDatasetDimensions(w, ctx, logData, v, []model.Dimension{newDimension}, filter.PopulationType) {
		return
	}

	if err := api.store.AddFilterDimension(ctx, fID, newDimension); err != nil {
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
	resp := addFilterDimensionResponse{
		Dimension: newDimension,
	}
	fdBytes, err := api.hashFilterDimensions(ctx, filter.ID)
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
	w.Header().Set("ETag", fdBytes)
	api.respond.JSON(ctx, w, http.StatusCreated, resp)
}

func (api *API) hashFilterDimensions(ctx context.Context, fID string) (string, error) {
	filter, err := api.store.GetFilter(ctx, fID)
	if err != nil {
		return "", err
	}

	h := sha1.New()

	dimensions := struct {
		items []model.Dimension
	}{
		items: filter.Dimensions,
	}
	fdBytes, err := bson.Marshal(dimensions)
	if err != nil {
		return "", err
	}
	if _, err := h.Write(fdBytes); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func (api *API) putFilter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	time, _ := time.Parse(time.RFC3339, "2016-07-17T08:38:25.316Z")

	resp := putFilterResponse{
		model.PutFilter{
			Events: []model.Event{
				{
					Timestamp: time,
					Name:      "cantabular-export-start",
				},
			},
			Dataset: model.Dataset{
				ID:      "string",
				Edition: "string",
				Version: 0,
			},
			PopulationType: "string",
		},
	}

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

func (api *API) isValidDatasetDimensions(w http.ResponseWriter, ctx context.Context, logData log.Data, v dataset.Version, d []model.Dimension, pt string) bool {
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
		return false
	}

	dimIDs, err := api.validateDimensions(d, v.Dimensions)
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
		return false
	}

	// Validate dimension options by performing Cantabular query with selections,
	// skip this check if requesting all options
	if err := api.validateDimensionOptions(ctx, d, dimIDs, pt); err != nil {
		api.respond.Error(
			ctx,
			w,
			statusCode(err),
			Error{
				err:     errors.Wrap(err, "failed to validate dimension options"),
				logData: logData,
			},
		)
		return false
	}
	return true
}

// validateDimensions validates provided filter dimensions exist within the dataset dimensions provided.
// Returns a map of the dimensions name:id for use in the following validation calls
func (api *API) validateDimensions(filterDims []model.Dimension, dims []dataset.VersionDimension) (map[string]string, error) {
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

func (api *API) validateDimensionOptions(ctx context.Context, filterDimensions []model.Dimension, dimIDs map[string]string, populationType string) error {
	dReq := cantabular.GetDimensionOptionsRequest{
		Dataset: populationType,
	}
	for _, d := range filterDimensions {
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
