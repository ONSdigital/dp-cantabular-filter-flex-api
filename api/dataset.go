package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	"github.com/ONSdigital/log.go/v2/log"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

type optionsMap map[string]map[string]dataset.Option

type datasetParams struct {
	id                string
	edition           string
	version           string
	basedOn           string
	datasetLink       dataset.Link
	versionLink       dataset.Link
	metadataLink      dataset.Link
	geoDimensions     []string
	datasetDimensions []string
	sortedDimensions  []string
	options           optionsMap // dimension -> option -> option item
}

func (api *API) getDatasetJSONHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	params, err := api.getDatasetParams(ctx, r)
	if err != nil {
		api.respond.Error(
			ctx,
			w,
			statusCode(err),
			errors.Wrap(err, "failed to get dataset params"),
		)
		return
	}

	logData := log.Data{
		"id":      params.id,
		"edition": params.edition,
		"version": params.version,
	}

	resp, err := api.getDatasetJSON(ctx, r, params)
	if err != nil {
		api.respond.Error(
			ctx,
			w,
			statusCode(err),
			Error{
				err:     err,
				logData: logData,
			},
		)
	}

	api.respond.JSON(ctx, w, http.StatusOK, resp)
}

func (api *API) getDatasetObservationHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	params, err := api.getDatasetParams(ctx, r)
	if err != nil {
		api.respond.Error(
			ctx,
			w,
			statusCode(err),
			errors.Wrap(err, "failed to get dataset params"),
		)
		return
	}

	logData := log.Data{
		"id":      params.id,
		"edition": params.edition,
		"version": params.version,
	}

	response, err := api.getDatasetObservations(ctx, r, params)
	if err != nil {
		api.respond.Error(
			ctx,
			w,
			statusCode(err),
			Error{
				err:     err,
				logData: logData,
			},
		)
	} else {
		api.respond.JSON(ctx, w, http.StatusOK, response)
	}
}

func (api *API) getDatasetJSON(ctx context.Context, r *http.Request, params *datasetParams) (*GetDatasetJSONResponse, error) {
	dReq := cantabular.StaticDatasetQueryRequest{
		Dataset:   params.basedOn,
		Variables: params.sortedDimensions,
	}

	if r.URL.Query().Get("area-type") != "" {
		variables, filters, err := api.validateGeography(ctx, r, params, dReq)
		if err != nil {
			return nil, errors.Wrap(err, "validateGeography failed")
		}
		dReq.Variables = variables
		dReq.Filters = filters
	}

	qRes, err := api.ctblr.StaticDatasetQuery(ctx, dReq)
	if err != nil {
		return nil, errors.Wrap(err, "failed to run query")
	}

	resp, err := api.toGetDatasetJsonResponse(params, qRes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate response")
	}

	return resp, nil
}

func (api *API) getDatasetObservations(ctx context.Context, r *http.Request, params *datasetParams) (*GetObservationsResponse, error) {
	dReq := cantabular.StaticDatasetQueryRequest{
		Dataset:   params.basedOn,
		Variables: params.sortedDimensions,
	}

	if r.URL.Query().Get("area-type") != "" {
		variables, filters, err := api.validateGeography(ctx, r, params, dReq)
		if err != nil {
			return nil, errors.Wrap(err, "validateGeography failed")
		}
		dReq.Variables = variables
		dReq.Filters = filters
	}

	qRes, err := api.ctblr.StaticDatasetQuery(ctx, dReq)
	if err != nil {
		return nil, errors.Wrap(err, "failed to run query")
	}

	resp, err := api.toGetDatasetObservationsResponse(params, qRes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate response")
	}

	return resp, nil
}

func (api *API) validateGeography(ctx context.Context, r *http.Request, params *datasetParams, datasetRequest cantabular.StaticDatasetQueryRequest) ([]string, []cantabular.Filter, error) {
	geographyQuery := strings.Split(r.URL.Query().Get("area-type"), ",")

	if len(geographyQuery) < 1 {
		return nil, nil, errors.New("unable to locate area-type")
	}

	geography := geographyQuery[0]
	geographyOptions := geographyQuery[1:]

	foundGeography := false
	for _, d := range params.geoDimensions {
		if strings.EqualFold(d, geography) {
			foundGeography = true
			params.sortedDimensions[0] = geographyQuery[0]

			datasetRequest.Variables = params.sortedDimensions
			break
		}
	}

	if len(geographyOptions) > 0 {
		// Check the area code is in the area type
		for _, a := range geographyOptions {
			cReq := cantabular.GetAreaRequest{
				Dataset:  params.basedOn,
				Variable: geography,
				Category: a,
			}

			_, err := api.ctblr.GetArea(ctx, cReq)
			if err != nil {
				return nil, nil, errors.New("unable to locate area")
			}

		}
		foundGeography = true
		datasetRequest.Filters = []cantabular.Filter{{Variable: geography, Codes: geographyOptions}}
	}

	if !foundGeography {
		return nil, nil, errors.New("unable to locate area or area-type")
	}

	return datasetRequest.Variables, datasetRequest.Filters, nil
}

func (api *API) getGeographyFilters(r *http.Request, params *datasetParams) ([]cantabular.Filter, error) {
	geographyQuery := strings.Split(r.URL.Query().Get("area-type"), ",")

	if len(geographyQuery) < 2 {
		return nil, errors.New("unable to locate geography")
	}

	geography := geographyQuery[0]
	geographyOptions := geographyQuery[1:]

	foundGeography := false
	for _, d := range params.geoDimensions {
		if strings.EqualFold(d, geography) {
			foundGeography = true
			break
		}
	}

	if !foundGeography {
		return nil, errors.Errorf("unable to validate geography %s", geography)
	}

	var geographyCodes []string
	for _, o := range geographyOptions {
		opt, ok := params.options[geography][o]
		if !ok {
			return nil, errors.Errorf("unable to validate geography option %s", o)
		}
		geographyCodes = append(geographyCodes, opt.Option)
	}

	dimension := r.URL.Query().Get("dimension")
	if dimension == "" {
		return nil, errors.New("unable to locate dimension")
	}

	foundDimension := false
	for _, d := range params.datasetDimensions {
		if d == dimension {
			foundDimension = true
			break
		}
	}
	if !foundDimension {
		return nil, errors.Errorf("unable to validate dimension %s", dimension)
	}

	options := strings.Split(r.URL.Query().Get("options"), ",")
	if len(options) < 1 || options[0] == "" {
		return nil, errors.Errorf("invalid options length or options is empty")
	}

	var dimensionCodes []string
	for _, o := range options {
		opt, ok := params.options[dimension][o]
		if !ok {
			return nil, errors.Errorf("unable to locate dimension option %s", o)
		}
		dimensionCodes = append(dimensionCodes, opt.Option)
	}

	return []cantabular.Filter{{Variable: geography, Codes: geographyCodes}, {Variable: dimension, Codes: dimensionCodes}}, nil

}

func (api *API) toGetDatasetJsonResponse(params *datasetParams, query *cantabular.StaticDatasetQuery) (*GetDatasetJSONResponse, error) {
	var dimensions []DatasetJSONDimension

	for _, dimension := range query.Dataset.Table.Dimensions {
		var options []model.Link

		for _, option := range dimension.Categories {
			options = append(options, model.Link{
				ID:    option.Code,
				Label: option.Label,
			})
		}
		dimensions = append(dimensions, DatasetJSONDimension{
			DimensionName: dimension.Variable.Name,
			Options:       options,
		})
	}

	datasetLinks := DatasetJSONLinks{
		DatasetMetadata: model.Link{
			HREF: params.metadataLink.URL,
			ID:   params.metadataLink.ID,
		},
		Self: model.Link{
			HREF: params.datasetLink.URL,
			ID:   params.datasetLink.ID,
		},
		Version: model.Link{
			HREF: params.versionLink.URL,
			ID:   params.versionLink.ID,
		},
	}

	getDatasetJsonResponse := GetDatasetJSONResponse{
		Dimensions:        dimensions,
		Links:             datasetLinks,
		Observations:      query.Dataset.Table.Values,
		TotalObservations: len(query.Dataset.Table.Values),
	}

	return &getDatasetJsonResponse, nil
}

func (api *API) toGetDatasetObservationsResponse(params *datasetParams, query *cantabular.StaticDatasetQuery) (*GetObservationsResponse, error) {
	var observations []GetObservationResponse

	dimLengths := make([]int, 0)
	dimIndices := make([]int, 0)

	for _, d := range query.Dataset.Table.Dimensions {
		dimLengths = append(dimLengths, len(d.Categories))
		dimIndices = append(dimIndices, 0)
	}

	for _, v := range query.Dataset.Table.Values {
		dimensions := getDimensionRow(query, dimIndices)
		observations = append(observations, GetObservationResponse{
			Dimensions:  dimensions,
			Observation: v,
		})

		l := len(dimIndices) - 1
		for l >= 0 {
			dimIndices[l] += 1
			if dimIndices[l] < dimLengths[l] {
				break
			}
			dimIndices[l] = 0
			l -= 1
		}
	}

	resp := GetObservationsResponse{
		Observations:      observations,
		TotalObservations: len(query.Dataset.Table.Values),
		Links: DatasetJSONLinks{
			DatasetMetadata: model.Link{
				HREF: params.metadataLink.URL,
				ID:   params.metadataLink.ID,
			},
			Self: model.Link{
				HREF: params.datasetLink.URL,
				ID:   params.datasetLink.ID,
			},
			Version: model.Link{
				HREF: params.versionLink.URL,
				ID:   params.versionLink.ID,
			},
		},
	}

	return &resp, nil
}

func getDimensionRow(query *cantabular.StaticDatasetQuery, catIndices []int) []ObservationDimension {
	var dims []ObservationDimension

	for i, index := range catIndices {
		dim := query.Dataset.Table.Dimensions[i]

		dims = append(dims, ObservationDimension{
			Dimension:   dim.Variable.Label,
			DimensionID: dim.Variable.Name,
			Option:      dim.Categories[index].Label,
			OptionID:    dim.Categories[index].Code,
		})
	}

	return dims
}

func (api *API) getDatasetParams(ctx context.Context, r *http.Request) (*datasetParams, error) {
	params := &datasetParams{
		id:      chi.URLParam(r, "dataset_id"),
		edition: chi.URLParam(r, "edition"),
		version: chi.URLParam(r, "version"),
		options: make(optionsMap),
	}

	if params.id == "" {
		return nil, errors.New("invalid dataset id")
	}

	if params.edition == "" {
		return nil, errors.New("invalid edition")
	}

	if params.version == "" {
		return nil, errors.New("invalid version")
	}

	// The following GetVersion() call will only return a 'published' version for an unauthorised caller, i.e. public caller
	// We are therefore guaranteed that the if a version is returned, it is 'published' and the BasedOn, Dimension, and Links.datasetLink/versionLink attributes are complete
	versionItem, err := api.datasets.GetVersion(ctx, "", api.cfg.ServiceAuthToken, "", "", params.id, params.edition, params.version)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get version")
	}

	if versionItem.IsBasedOn.Type != cantabularFlexibleTable && versionItem.IsBasedOn.Type != cantabularMultivariateTable {
		return nil, errors.New("invalid dataset type")
	}

	params.datasetLink = versionItem.Links.Dataset
	params.versionLink = versionItem.Links.Self
	params.basedOn = versionItem.IsBasedOn.ID

	params.metadataLink.URL = api.datasets.GetMetadataURL(params.id, params.edition, params.version)

	if len(versionItem.Dimensions) == 0 {
		return nil, errors.New("invalid dimensions length of zero")
	}

	// Used for improved performance doing checks against extra dimensions
	dimMap := make(map[string]struct{})

	for _, dimension := range versionItem.Dimensions {
		dimMap[dimension.Name] = struct{}{}

		options, err := api.datasets.GetOptionsInBatches(ctx, "", api.cfg.ServiceAuthToken, "", params.id, params.edition, params.version, dimension.Name, api.cfg.DatasetOptionsBatchSize, api.cfg.DatasetOptionsWorkers)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get options")
		}

		params.options[dimension.ID] = make(map[string]dataset.Option)

		for _, option := range options.Items {
			params.options[dimension.ID][option.Label] = option
		}

		params.datasetDimensions = append(params.datasetDimensions, dimension.ID)
	}

	dq := r.URL.Query().Get("dimensions")
	if extraDims := strings.Split(dq, ","); dq != "" && len(extraDims) > 0 {
		if versionItem.IsBasedOn.Type != cantabularMultivariateTable {
			return nil, &Error{
				err:        errors.New("invalid dataset type for custom dimensions"),
				badRequest: true,
			}
		}

		catMap := make(map[string]struct{})

		for _, dim := range versionItem.Dimensions {
			if dim.IsAreaType != nil && *dim.IsAreaType {
				continue
			}

			res, err := api.population.GetCategorisations(ctx, population.GetCategorisationsInput{
				PopulationType: params.basedOn,
				Dimension:      dim.ID,
				PaginationParams: population.PaginationParams{
					Limit: 99999,
				},
				AuthTokens: population.AuthTokens{
					ServiceAuthToken: api.cfg.ServiceAuthToken,
				},
			})
			if err != nil {
				return nil, errors.Wrap(err, "failed to get categorisations")
			}
			for _, dim := range res.Items {
				catMap[dim.ID] = struct{}{}
			}
		}

		for _, ed := range extraDims {
			if _, found := dimMap[ed]; found {
				return nil, errors.Errorf("dimension already in dataset: %s", ed)
			}
			if _, found := catMap[ed]; found {
				return nil, errors.Errorf("categorisation of dimension already in dataset: %s", ed)
			}
		}

		res, err := api.ctblr.GetDimensionsByName(ctx, cantabular.GetDimensionsByNameRequest{
			Dataset:          params.basedOn,
			DimensionNames:   extraDims,
			ExcludeGeography: true,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to get dimensions")
		}

		if len(res.Dataset.Variables.Edges) < len(extraDims) {
			return nil, errors.New("geography variable added in 'dimensions' (use 'area-type')")
		}

		params.datasetDimensions = append(params.datasetDimensions, extraDims...)
	}

	params.geoDimensions, err = api.getGeographyTypes(ctx, params.basedOn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get geography types")
	}

	params.sortedDimensions = api.sortGeography(params.geoDimensions, params.datasetDimensions)

	return params, nil
}

const (
	batchSize     = 100
	numberWorkers = 10
)

func (api *API) getGeographyTypes(ctx context.Context, datasetId string) ([]string, error) {
	var geoDimensions []string

	res, err := api.ctblr.GetGeographyDimensionsInBatches(ctx, datasetId, batchSize, numberWorkers)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get Geography Dimensions")
	}

	for _, d := range res.Variables.Edges {
		geoDimensions = append(geoDimensions, d.Node.Name)
	}

	return geoDimensions, nil
}

func (api *API) sortGeography(geoDimensions []string, dimensions []string) []string {
	foundGeography := false
	var sortedDimensions []string
	var nonGeoDimensions []string

	for _, item := range dimensions {
		isGeography := false

		for _, geo := range geoDimensions {
			if strings.EqualFold(geo, item) {
				isGeography = true

				if !foundGeography {
					foundGeography = true
					sortedDimensions = append(sortedDimensions, item)
				}
			}
		}

		if !isGeography {
			nonGeoDimensions = append(nonGeoDimensions, item)
		}
	}

	sortedDimensions = append(sortedDimensions, nonGeoDimensions...)

	return sortedDimensions
}
