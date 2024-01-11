package api

import (
	"net/http"
	"strconv"

	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/event"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	dprequest "github.com/ONSdigital/dp-net/v2/request"
	"github.com/ONSdigital/log.go/v2/log"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

const (
	flexible                    = "flexible"
	multivariate                = "multivariate"
	published                   = "published"
	cantabularTable             = "cantabular_table"
	cantabularMultivariateTable = "cantabular_multivariate_table"
	cantabularFlexibleTable     = "cantabular_flexible_table"
)

var (
	filterTypes = map[string]string{
		cantabularFlexibleTable:     flexible,
		cantabularMultivariateTable: multivariate,
	}
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

	if req.Custom && v.IsBasedOn.Type != cantabularMultivariateTable {
		api.respond.Error(
			ctx,
			w,
			http.StatusBadRequest,
			Error{
				err:     errors.New("invalid dataset type for custom filter"),
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

	filterType := filterTypes[v.IsBasedOn.Type]
	finalDims, err := api.validateAndHydrateDimensions(
		v,
		req.Dimensions,
		req.PopulationType,
	)
	if err != nil {
		api.respond.Error(ctx, w, statusCode(err), errors.Wrap(err, "failed to validate dimensions"))
		return
	}

	if err := api.validateDimensionOptions(ctx, finalDims, req.PopulationType); err != nil {
		api.respond.Error(
			ctx,
			w,
			statusCode(err),
			errors.Wrap(err, "failed to validate dimension options"),
		)
		return
	}

	currentNext, err := api.datasets.GetDatasetCurrentAndNext(
		ctx,
		"",
		api.cfg.ServiceAuthToken,
		"",
		req.Dataset.ID,
	)
	if err != nil {
		api.respond.Error(
			ctx,
			w,
			statusCode(err),
			Error{
				err:     errors.Wrap(err, "failed to get current and next dataset"),
				message: "failed to get current and next dataset information",
				logData: logData,
			},
		)
		return
	}

	req.Dataset.ReleaseDate = v.ReleaseDate
	req.Dataset.Title = currentNext.Title

	f := model.Filter{
		Links: model.FilterLinks{
			Version: model.Link{
				HREF: v.Links.Self.URL,
				ID:   strconv.Itoa(v.Version),
			},
		},
		Dimensions:        finalDims,
		UniqueTimestamp:   api.generate.UniqueTimestamp(),
		LastUpdated:       api.generate.Timestamp(),
		Dataset:           *req.Dataset,
		InstanceID:        v.ID,
		PopulationType:    req.PopulationType,
		Type:              filterType,
		Custom:            req.Custom,
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

func (api *API) createCustomFilter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req createCustomFilterRequest

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

	// call the population-types-api to get the dataset ID
	input := population.GetPopulationTypeMetadataInput{
		AuthTokens: population.AuthTokens{
			ServiceAuthToken: api.cfg.ServiceAuthToken,
		},
		PopulationType: req.PopulationType,
	}

	dMetadata, err := api.population.GetPopulationTypeMetadata(ctx, input)
	if err != nil {
		api.respond.Error(
			ctx,
			w,
			statusCode(err),
			Error{
				err:     errors.Wrap(err, "failed to get metadata for "+input.PopulationType),
				message: "failed to get metadata for " + input.PopulationType,
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
		dMetadata.DefaultDatasetID,
		dMetadata.Edition,
		strconv.Itoa(dMetadata.Version),
	)
	if err != nil {
		api.respond.Error(
			ctx,
			w,
			statusCode(err),
			Error{
				err:     errors.Wrap(err, "failed to get existing Version"),
				message: "failed to get existing dataset Version",
				logData: logData,
			},
		)
		return
	}

	if v.IsBasedOn == nil || v.IsBasedOn.Type != cantabularMultivariateTable {
		api.respond.Error(
			ctx,
			w,
			http.StatusBadRequest,
			Error{
				err:     errors.New("default dataset is not of type multivariate table"),
				message: "default dataset is not of type multivariate table",
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

	filterType := filterTypes[v.IsBasedOn.Type]
	dataset := model.Dataset{
		ID:      dMetadata.DefaultDatasetID,
		Edition: dMetadata.Edition,
		Version: dMetadata.Version,
		Title:   "custom",
	}

	// get the first available area type dimension from the default dataset
	var dimension model.Dimension
	for i := range v.Dimensions {
		d := &v.Dimensions[i]
		if *d.IsAreaType {
			dimension = model.Dimension{
				Name:       d.Name,
				Label:      d.Label,
				ID:         d.ID,
				IsAreaType: d.IsAreaType,
				Options:    []string{},
			}
			break
		}
	}

	f := model.Filter{
		Links: model.FilterLinks{
			Version: model.Link{
				HREF: v.Links.Self.URL,
				ID:   strconv.Itoa(v.Version),
			},
		},
		Dimensions:        []model.Dimension{dimension},
		UniqueTimestamp:   api.generate.UniqueTimestamp(),
		LastUpdated:       api.generate.Timestamp(),
		Dataset:           dataset,
		InstanceID:        v.ID,
		PopulationType:    req.PopulationType,
		Type:              filterType,
		Published:         v.State == published,
		DisclosureControl: nil, // populate for these fields yet
		Custom:            true,
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
		Custom:         filter.Custom,
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

	dim := make([]string, 0, len(filter.Dimensions))
	for i := range filter.Dimensions {
		dim = append(dim, filter.Dimensions[i].Name)
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
	if err := api.produceExportStartEvent(ctx, e); err != nil {
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
		Custom:         filter.Custom,
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
