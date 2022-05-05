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
	id                 string
	edition            string
	version            string
	basedOn            string
	datasetLink        dataset.Link
	versionLink        dataset.Link
	metadataLink       dataset.Link
	unsortedDimensions []string
	options            optionsMap // dimension -> option -> option item
}

func (api *API) getDatasetJSON(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.URL.Query().Get("geography") == "" {
		api.getDefaultDatasetJSON(ctx, w, r)
	} else {
		api.getGeographyDatsetJSON(ctx, r)
	}
}

func (api *API) getDefaultDatasetJSON(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	params, err := api.getDatasetParams(ctx, r)
	if err != nil {
		api.respond.Error(ctx, w, statusCode(err), errors.Wrap(err, "failed to get dataset params"))
		return
	}

	logData := log.Data{
		"id":      params.id,
		"edition": params.edition,
		"version": params.version,
	}

	geoDimensions, err := api.getGeographyTypes(ctx, params.basedOn)
	if err != nil {
		api.respond.Error(
			ctx,
			w,
			statusCode(err),
			Error{
				err:     errors.Wrap(err, "failed to get geography types"),
				logData: logData,
			},
		)
		return
	}

	dimensions, _ := api.sortGeography(geoDimensions, params.unsortedDimensions)

	datasetRequest := cantabular.StaticDatasetQueryRequest{
		Dataset:   params.basedOn,
		Variables: dimensions,
	}

	queryResult, err := api.ctblr.StaticDatasetQuery(ctx, datasetRequest)
	if err != nil {
		api.respond.Error(
			ctx,
			w,
			statusCode(err),
			Error{
				err:     errors.Wrap(err, "failed to run query"),
				logData: logData,
			},
		)
		return
	}

	response, err := api.toGetDatasetJsonResponse(params, queryResult)
	if err != nil {
		api.respond.Error(
			ctx,
			w,
			statusCode(err),
			Error{
				err:     errors.Wrap(err, "failed to generate response"),
				logData: logData,
			},
		)
		return
	}

	api.respond.JSON(ctx, w, http.StatusOK, response)
}

func (api *API) getGeographyDatsetJSON(ctx context.Context, r *http.Request) (*getDatasetJSONResponse, error) {
	return &getDatasetJSONResponse{}, nil
}

func (api *API) toGetDatasetJsonResponse(params *datasetParams, query *cantabular.StaticDatasetQuery) (*getDatasetJSONResponse, error) {
	var dimensions []DatasetJSONDimension

	for _, dimension := range query.Dataset.Table.Dimensions {
		if _, ok := params.options[dimension.Variable.Name]; !ok {
			return nil, errors.New("dimension mismatch")
		}

		var options []model.Link

		for _, option := range dimension.Categories {
			if _, ok := params.options[dimension.Variable.Name][option.Label]; !ok {
				return nil, errors.New("option mismatch")
			}

			option := model.Link{
				HREF: params.options[dimension.Variable.Name][option.Label].Links.Code.URL,
				ID:   option.Label,
			}

			options = append(options, option)
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

	getDatasetJsonResponse := getDatasetJSONResponse{
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

	datasetItem, err := api.datasets.GetDatasetCurrentAndNext(ctx, "", "", "", params.id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get dataset")
	}

	params.datasetLink = datasetItem.Links.Self
	params.basedOn = datasetItem.IsBasedOn.ID

	if datasetItem.DatasetDetails.Type != "cantabular_flexible_table" {
		return nil, errors.New("invalid dataset type")
	}

	versionItem, err := api.datasets.GetVersion(ctx, "", "", "", "", params.id, params.edition, params.version)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get version")
	}

	params.versionLink = versionItem.Links.Self

	metadata, err := api.datasets.GetVersionMetadata(ctx, "", "", "", params.id, params.edition, params.version)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get metadata")
	}

	params.metadataLink = metadata.Version.Links.Self

	dimensions, err := api.datasets.GetVersionDimensions(ctx, "", "", "", params.id, params.edition, params.version)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get dimensions")
	}

	if dimensions.Items.Len() == 0 {
		return nil, errors.New("invalid dimensions length of zero")
	}

	for _, dimension := range dimensions.Items {
		options, err := api.datasets.GetOptionsInBatches(ctx, "", "", "", params.id, params.edition, params.version, dimension.Name, api.cfg.DatasetOptionsBatchSize, api.cfg.DatasetOptionsWorkers)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get options")
		}

		params.options[dimension.Links.CodeList.ID] = make(map[string]dataset.Option)

		for _, option := range options.Items {
			params.options[dimension.Links.CodeList.ID][option.Label] = option
		}

		params.unsortedDimensions = append(params.unsortedDimensions, dimension.Links.CodeList.ID)
	}

	return params, nil
}

func (api *API) getGeographyTypes(ctx context.Context, datasetId string) ([]string, error) {
	var geoDimensions []string

	res, err := api.ctblr.GetGeographyDimensions(ctx, datasetId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get Geography Dimensions")
	}

	for _, dimension := range res.Dataset.RuleBase.IsSourceOf.Edges {
		geoDimensions = append(geoDimensions, dimension.Node.Name)
	}

	return geoDimensions, nil
}

func (api *API) sortGeography(geoDimensions []string, dimensions []string) ([]string, bool) {
	foundGeography := false
	var sortedDimensions []string
	var nonGeoDimensions []string

	for _, item := range dimensions {
		isGeography := false

		for _, geo := range geoDimensions {
			if geo == strings.ToUpper(item) {
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

	return sortedDimensions, foundGeography
}
