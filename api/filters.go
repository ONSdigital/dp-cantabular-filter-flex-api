package api

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular/gql"
	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	dperrors "github.com/ONSdigital/dp-cantabular-filter-flex-api/errors"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/event"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	dprequest "github.com/ONSdigital/dp-net/v2/request"
	"github.com/ONSdigital/log.go/v2/log"

	"github.com/go-chi/chi/v5"
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
		Links: model.FilterLinks{
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
		Dimensions:        hydrateDimensions(req.Dimensions, v.Dimensions),
		UniqueTimestamp:   api.generate.UniqueTimestamp(),
		LastUpdated:       api.generate.Timestamp(),
		Dataset:           *req.Dataset,
		InstanceID:        v.ID,
		PopulationType:    req.PopulationType,
		Type:              flexible,
		Published:         v.State == published,
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

	// don't return dimensions in response
	f.Dimensions = nil

	resp := createFilterResponse{f}

	api.respond.JSON(ctx, w, http.StatusCreated, resp)
}

func (api *API) submitFilter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	filterID := chi.URLParam(r, "id")
	logData := log.Data{
		"filter_id": filterID,
	}

	filter, err := api.store.GetFilter(ctx, filterID)
	if err != nil {
		api.respond.Error(
			ctx,
			w,
			statusCode(err),
			Error{
				err:     errors.Wrap(err, "failed to get filter"),
				message: "failed to submit filter: failed to get existing filter",
				logData: logData,
			},
		)
		return

	}

	filterOutput := model.FilterOutput{
		// id created by mongo client
		FilterID:  filter.ID,
		State:     model.Submitted,
		Dataset:   filter.Dataset,
		Downloads: model.Downloads{},
		Events:    []model.Event{},
		Links: model.FilterOutputLinks{
			Version: filter.Links.Version,
			Self: model.Link{
				HREF: api.generate.URL(api.cfg.FilterAPIURL, "/filter-outputs"),
				// uuid created by mongo client, will set there.
				ID: "",
			},
			FilterBlueprint: model.Link{
				HREF: api.generate.URL(api.cfg.FilterAPIURL, "/filters"),
				ID:   filter.ID,
			},
		},
		Published:  filter.Published,
		Dimensions: filter.Dimensions,
		Type:       filter.Type,
	}

	if err = api.store.CreateFilterOutput(ctx, &filterOutput); err != nil {
		api.respond.Error(
			ctx,
			w,
			statusCode(err),
			Error{
				err:     errors.Wrap(err, "failed to create filter output"),
				message: "error submitting filter",
				logData: logData,
			},
		)
		return
	}

	// schema mismatch between avro and model type.
	// naively converting for now.
	version := strconv.Itoa(filter.Dataset.Version)

	e := event.ExportStart{
		InstanceID:     filter.InstanceID,
		DatasetID:      filter.Dataset.ID,
		Edition:        filter.Dataset.Edition,
		Version:        version,
		FilterOutputID: filterOutput.ID,
	}

	// send the export event through Kafka
	if err := api.produceExportStartEvent(e); err != nil {
		api.respond.Error(
			ctx,
			w,
			http.StatusInternalServerError,
			Error{
				err:     errors.Wrap(err, "failed to create export start event"),
				message: "error submitting filter",
				logData: logData,
			},
		)
		return
	}

	resp := submitFilterResponse{
		InstanceID:     filter.InstanceID,
		FilterOutputID: filterOutput.ID,
		Dataset:        filter.Dataset,
		Links:          filter.Links,
		PopulationType: filter.PopulationType,
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

	// don't return dimensions in response
	f.Dimensions = nil

	resp := getFilterResponse{*f}

	api.respond.JSON(ctx, w, http.StatusOK, resp)
}

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

	if !api.isValidDatasetDimensions(w, ctx, logData, v, []model.Dimension{req.Dimension}, filter.PopulationType) {
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
	dimensionName := chi.URLParam(r, "name")

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

func (api *API) putFilter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	resp := putFilterResponse{
		Events: []model.Event{
			{
				Timestamp: "2016-07-17T08:38:25.316+000",
				Name:      "cantabular-export-start",
			},
		},
		Dataset: model.Dataset{
			ID:      "string",
			Edition: "string",
			Version: 0,
		},
		PopulationType: "string",
	}

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

// getCantabularDimension checks that dimension exists in Cantabular by searching for it.
// If the dimension doesn't exist, or couldn't be retrieved, an error is returned.
func (api *API) getCantabularDimension(ctx context.Context, filterID, dimensionName string) (*gql.Node, error) {
	filter, err := api.store.GetFilter(ctx, filterID)
	if err != nil {
		return nil, Error{
			err:        errors.Wrap(err, "error retrieving filter"),
			message:    "failed to get filter dimensions",
			badRequest: true,
		}
	}

	foundDimensions, err := api.ctblr.SearchDimensions(ctx, cantabular.SearchDimensionsRequest{
		Dataset: filter.PopulationType,
		Text:    dimensionName,
	})
	if err != nil {
		return nil, errors.Wrap(err, "error in cantabular response")
	}

	if len(foundDimensions.Dataset.Variables.Search.Edges) == 0 {
		return nil, Error{
			err:      errors.New("no dimensions in response"),
			notFound: true,
			logData:  log.Data{"found_dimensions": foundDimensions},
		}
	}

	return &foundDimensions.Dataset.Variables.Search.Edges[0].Node, nil
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
		hydrated = append(hydrated, dim)
	}

	return hydrated
}
