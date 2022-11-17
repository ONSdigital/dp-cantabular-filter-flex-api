package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
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

	response, err := api.getDatasetJSON(ctx, r, params)
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
	datasetRequest := cantabular.StaticDatasetQueryRequest{
		Dataset:   params.basedOn,
		Variables: params.sortedDimensions,
	}

	if r.URL.Query().Get("area-type") != "" {
		variables, filters, err := api.validateGeography(ctx, r, params, datasetRequest)
		if err != nil {
			return nil, errors.Wrap(err, "validateGeography failed")
		}
		datasetRequest.Variables = variables
		datasetRequest.Filters = filters
	}

	queryResult, err := api.ctblr.StaticDatasetQuery(ctx, datasetRequest)
	if err != nil {
		return nil, errors.Wrap(err, "failed to run query")
	}

	response, err := api.toGetDatasetJsonResponse(params, queryResult)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate response")
	}

	return response, nil
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
		dimensions = append(dimensions, DatasetJSONDimension{
			DimensionName: dimension.Variable.Name,
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

	if versionItem.IsBasedOn.Type != cantabularFlexibleTable {
		return nil, errors.New("invalid dataset type")
	}

	params.datasetLink = versionItem.Links.Dataset
	params.versionLink = versionItem.Links.Self
	params.basedOn = versionItem.IsBasedOn.ID

	//metadata, err := api.datasets.GetVersionMetadata(ctx, "", api.cfg.ServiceAuthToken, "", params.id, params.edition, params.version)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "failed to get metadata")
	// }

	//params.metadataLink = metadata.Version.Links.Self
	params.metadataLink.URL = api.datasets.GetMetadataURL(params.id, params.edition, params.version)

	if len(versionItem.Dimensions) == 0 {
		return nil, errors.New("invalid dimensions length of zero")
	}

	for _, dimension := range versionItem.Dimensions {
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
