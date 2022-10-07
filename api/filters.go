package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/event"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	dprequest "github.com/ONSdigital/dp-net/v2/request"
	"github.com/ONSdigital/log.go/v2/log"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

const (
	flexible               = "flexible"
	multivariate           = "multivariate"
	published              = "published"
	cantabularTable        = "cantabular_table"
	cantabularMultivariate = "cantabular_multivariate_table"
	cantabularFlexible     = "cantabular_flexible_table"
)

func (api *API) createFilter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req createFilterRequest

	var finalDims []model.Dimension
	var filterType string

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

	if v.IsBasedOn == nil || v.IsBasedOn.Type == cantabularTable {
		api.respond.Error(
			ctx,
			w,
			http.StatusBadRequest,
			Error{
				err:     errors.New("dataset is of type cantabular table"),
				message: "dataset is of invalid type",
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

	finalDims, filterType, err = api.ValidateAndReturnDimensions(
		v,
		req.Dimensions,
		req.PopulationType,
	)
	if err != nil {
		api.respond.Error(ctx, w, statusCode(err), err)
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
		Dimensions:        finalDims,
		UniqueTimestamp:   api.generate.UniqueTimestamp(),
		LastUpdated:       api.generate.Timestamp(),
		Dataset:           *req.Dataset,
		InstanceID:        v.ID,
		PopulationType:    req.PopulationType,
		Type:              filterType,
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
			},
			FilterBlueprint: model.Link{
				HREF: api.generate.URL(api.cfg.FilterAPIURL, "/filters"),
				ID:   filter.ID,
			},
		},
		Published:      filter.Published,
		Dimensions:     filter.Dimensions,
		Type:           filter.Type,
		PopulationType: filter.PopulationType,
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

	var dim []string
	for _, d := range filter.Dimensions {
		dim = append(dim, d.Name)
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
		Dimensions:     dim,
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

// TODO: what is this unfinished function?
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

// isValidMultivariateDimensions checks the validity of the supplied dimensions for a multivariate filter.
// Supplied dimensions may not be in original dataset but still valid, and so isValidDatasetDimensions is
// not relevant.
// NOTE: when we hydrate the dimensions, we will be using the name as the id, and filling out the dimensions
// using the same value for both.
func (api *API) isValidMultivariateDimensions(ctx context.Context, dimensions []model.Dimension, pType string) ([]model.Dimension, error) {
	hydratedDimensions := make([]model.Dimension, 0)

	for _, d := range dimensions {
		dim, err := api.getCantabularDimension(ctx, pType, d.Name)
		if err != nil {
			return nil, errors.Wrap(err, "error in cantabular response")
		}

		dim.IsAreaType = d.IsAreaType
		hydratedDimensions = append(hydratedDimensions, *dim)
	}

	return hydratedDimensions, nil
}

func (api *API) isValidDatasetDimensions(ctx context.Context, v dataset.Version, d []model.Dimension, pType string) error {
	dimIDs, err := api.validateDimensions(d, v.Dimensions)
	if err != nil {
		return Error{
			err:      errors.Wrap(err, "failed to validate request dimensions"),
			notFound: true,
		}
	}

	if err := api.validateDimensionOptions(ctx, d, dimIDs, pType); err != nil {
		return errors.Wrap(err, "failed to validate dimension options")
	}

	return nil
}

// getCantabularDimension checks that dimension exists in Cantabular by searching for it.
// If the dimension doesn't exist, or couldn't be retrieved, an error is returned.
func (api *API) getCantabularDimension(ctx context.Context, popType, dimensionName string) (*model.Dimension, error) {
	resp, err := api.ctblr.GetDimensionsByName(ctx, cantabular.GetDimensionsByNameRequest{
		Dataset:        popType,
		DimensionNames: []string{dimensionName},
	})
	if err != nil {
		return nil, errors.Wrap(err, "error in cantabular response")
	}

	if len(resp.Dataset.Variables.Search.Edges) == 0 {
		return nil, Error{
			err:      errors.New("no dimensions in response"),
			notFound: true,
			logData:  log.Data{"response": resp},
		}
	}

	node := resp.Dataset.Variables.Search.Edges[0].Node
	dim := model.Dimension{
		Label: node.Label,
		ID:    node.Name,
		Name:  node.Name,
	}

	return &dim, nil
}

// validateDimensions validates provided filter dimensions exist within the dataset dimensions provided.
// Returns a map of the dimensions name:id for use in the following validation calls
func (api *API) validateDimensions(filterDims []model.Dimension, dims []dataset.VersionDimension) (map[string]string, error) {

	fDims := make(map[string]bool)
	for _, fd := range filterDims {
		if _, ok := fDims[fd.Name]; ok {
			return nil, Error{
				err: errors.Errorf("duplicate dimensions chosen: %v", fd.Name),
				logData: log.Data{
					"duplicate dimensions chosen": fd.Name,
				},
			}
		} else {
			fDims[fd.Name] = true
		}
	}

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

// validateDimensionOptions by performing Cantabular query with selections,
// will be skipped if requesting all options
func (api *API) validateDimensionOptions(ctx context.Context, filterDimensions []model.Dimension, dimIDs map[string]string, populationType string) error {
	dReq := cantabular.GetDimensionOptionsRequest{
		Dataset: populationType,
	}
	for _, d := range filterDimensions {
		if len(d.Options) > 0 {
			dReq.DimensionNames = append(dReq.DimensionNames, dimIDs[d.Name])
			dReq.Filters = append(dReq.Filters, cantabular.Filter{
				Codes:    d.Options,
				Variable: getFilterVariable(dimIDs, d),
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

func getFilterVariable(dimIDs map[string]string, d model.Dimension) string {
	fVariable := dimIDs[d.Name]
	if len(d.FilterByParent) != 0 {
		fVariable = d.FilterByParent
	}
	return fVariable
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
		if dim.Options == nil {
			dim.Options = []string{}
		}
		hydrated = append(hydrated, dim)
	}

	return hydrated
}
