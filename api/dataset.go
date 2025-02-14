package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	"github.com/ONSdigital/dp-net/v2/links"
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

	if resp != nil {
		if resp.TotalObservations > api.cfg.MaxRowsReturned {
			api.respond.Error(
				ctx,
				w,
				403,
				Error{
					message: "Too many rows returned, please refine your query by requesting specific areas or reducing the number of categories returned.  For further information please visit https://developer.ons.gov.uk/createyourowndataset/",
				},
			)
			return
		}
	}

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

	cancelContext, cancel := context.WithTimeout(ctx, time.Second*300)
	defer cancel()

	filterFlexLinksBuilder := links.FromHeadersOrDefault(&r.Header, api.cantabularFilterFlexAPIURL)

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

	cReq := cantabular.StaticDatasetQueryRequest{
		Dataset:   params.basedOn,
		Variables: params.datasetDimensions,
	}

	if r.URL.Query().Get("area-type") != "" {
		variables, filters, err := api.validateGeography(ctx, r, params, cReq)
		if err != nil {
			api.respond.Error(
				ctx,
				w,
				statusCode(err),
				errors.Wrap(err, "validateGeography failed"),
			)
			return
		}
		cReq.Variables = variables
		cReq.Filters = filters
	}

	logData := log.Data{
		"id":      params.id,
		"edition": params.edition,
		"version": params.version,
	}

	countcheck, err := api.ctblr.CheckQueryCount(ctx, cReq)
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
		return
	}

	if countcheck > api.cfg.MaxRowsReturned {
		api.respond.Error(
			ctx,
			w,
			403,
			Error{
				message: "Too many rows returned, please refine your query by requesting specific areas or reducing the number of categories returned.  For further information please visit https://developer.ons.gov.uk/createyourowndataset/",
			},
		)
		return
	}

	var consume = func(ctx context.Context, file io.Reader) error {
		if file == nil {
			return errors.New("no file content has been provided")
		}

		response, err := api.processObservationsResponse(cancelContext, file, w)
		if err != nil {
			api.respond.Error(
				ctx,
				w,
				http.StatusUnprocessableEntity,
				Error{
					message: err.Error(),
				},
			)
		}
		if response == "" {
			api.respond.Error(
				ctx,
				w,
				http.StatusNotFound,
				Error{
					message: "No results found",
				},
			)
		}
		return nil
	}

	qRes, err := api.ctblr.StaticDatasetQueryStreamJSON(cancelContext, cReq, consume)
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
		return
	}

	qRes.TotalObservations = countcheck

	if api.cfg.EnableURLRewriting {
		params.metadataLink.URL, err = filterFlexLinksBuilder.BuildLink(params.metadataLink.URL)
		if err != nil {
			api.respond.Error(
				ctx,
				w,
				statusCode(err),
				errors.Wrap(err, "failed to get build metadata link"),
			)
			return
		}
		params.datasetLink.URL, err = filterFlexLinksBuilder.BuildLink(params.datasetLink.URL)
		if err != nil {
			api.respond.Error(
				ctx,
				w,
				statusCode(err),
				errors.Wrap(err, "failed to get build dataset link"),
			)
			return
		}
		params.versionLink.URL, err = filterFlexLinksBuilder.BuildLink(params.versionLink.URL)
		if err != nil {
			api.respond.Error(
				ctx,
				w,
				statusCode(err),
				errors.Wrap(err, "failed to get build version link"),
			)
			return
		}
	}

	datasetLink := &cantabular.Link{
		HREF: params.metadataLink.URL,
		ID:   params.metadataLink.ID,
	}

	qRes.Links.DatasetMetadata = datasetLink

	selfLink := &cantabular.Link{
		HREF: params.datasetLink.URL,
		ID:   params.datasetLink.ID,
	}

	qRes.Links.Self = *selfLink

	versionLink := &cantabular.Link{
		HREF: params.versionLink.URL,
		ID:   params.versionLink.ID,
	}

	qRes.Links.Version = versionLink

	api.respond.JSON(ctx, w, http.StatusOK, qRes)
}

func (api *API) processObservationsResponse(ctx context.Context, r io.Reader, w http.ResponseWriter) (string, error) {
	buf := new(strings.Builder)

	writ, err := io.Copy(buf, r)
	if err != nil {
		logData := log.Data{
			"method":  http.MethodGet,
			"message": err,
		}

		api.respond.Error(
			ctx,
			w,
			http.StatusUnprocessableEntity,
			Error{
				err:     fmt.Errorf("%s", fmt.Sprintf("An error occurred while processing the response, bytes written %d", writ)),
				logData: logData,
			},
		)
		return "", err
	}

	return buf.String(), nil
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

	resp, err := api.toGetDatasetJSONResponse(r, params, qRes)
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

	geographyCodes := make([]string, 0, len(geographyOptions))
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

	dimensionCodes := make([]string, 0, len(options))
	for _, o := range options {
		opt, ok := params.options[dimension][o]
		if !ok {
			return nil, errors.Errorf("unable to locate dimension option %s", o)
		}
		dimensionCodes = append(dimensionCodes, opt.Option)
	}

	return []cantabular.Filter{{Variable: geography, Codes: geographyCodes}, {Variable: dimension, Codes: dimensionCodes}}, nil
}

func (api *API) toGetDatasetJSONResponse(r *http.Request, params *datasetParams, query *cantabular.StaticDatasetQuery) (*GetDatasetJSONResponse, error) {
	dimensions := make([]DatasetJSONDimension, 0, len(query.Dataset.Table.Dimensions))
	var err error

	filterFlexLinksBuilder := links.FromHeadersOrDefault(&r.Header, api.cantabularFilterFlexAPIURL)

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

	if api.cfg.EnableURLRewriting {
		params.metadataLink.URL, err = filterFlexLinksBuilder.BuildLink(params.metadataLink.URL)
		if err != nil {
			return nil, errors.New("failed to build metadata link")
		}
		params.datasetLink.URL, err = filterFlexLinksBuilder.BuildLink(params.datasetLink.URL)
		if err != nil {
			return nil, errors.New("failed to build dataset link")
		}
		params.versionLink.URL, err = filterFlexLinksBuilder.BuildLink(params.versionLink.URL)
		if err != nil {
			return nil, errors.New("failed to build version link")
		}
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

	getDatasetJSONResponse := GetDatasetJSONResponse{
		Dimensions:        dimensions,
		Links:             datasetLinks,
		Observations:      query.Dataset.Table.Values,
		TotalObservations: len(query.Dataset.Table.Values),
		BlockedAreas:      query.Dataset.Table.Rules.Blocked.Count,
		TotalAreas:        query.Dataset.Table.Rules.Total.Count,
		AreasReturned:     query.Dataset.Table.Rules.Passed.Count,
	}

	return &getDatasetJSONResponse, nil
}

func (api *API) toGetDatasetObservationsResponse(params *datasetParams, query *cantabular.StaticDatasetQuery) (*GetObservationsResponse, error) {
	observations := make([]GetObservationResponse, 0, len(query.Dataset.Table.Values))

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
			dimIndices[l]++
			if dimIndices[l] < dimLengths[l] {
				break
			}
			dimIndices[l] = 0
			l--
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
		BlockedAreas:  query.Dataset.Table.Rules.Blocked.Count,
		AreasReturned: query.Dataset.Table.Rules.Passed.Count,
		TotalAreas:    query.Dataset.Table.Rules.Total.Count,
	}

	return &resp, nil
}

func getDimensionRow(query *cantabular.StaticDatasetQuery, catIndices []int) []ObservationDimension {
	dims := make([]ObservationDimension, 0, len(catIndices))

	for i, index := range catIndices {
		dim := &query.Dataset.Table.Dimensions[i]

		dims = append(dims, ObservationDimension{
			Dimension:   dim.Variable.Label,
			DimensionID: dim.Variable.Name,
			Option:      dim.Categories[index].Label,
			OptionID:    dim.Categories[index].Code,
		})
	}

	return dims
}

//nolint:gocognit,gocyclo // should break this function down in future
func (api *API) getDatasetParams(ctx context.Context, r *http.Request) (*datasetParams, error) {
	params := &datasetParams{
		id:      chi.URLParam(r, "dataset_id"),
		edition: chi.URLParam(r, "edition"),
		version: chi.URLParam(r, "version"),
		options: make(optionsMap),
	}

	err := validateBaseParams(params)
	if err != nil {
		return nil, err
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

	for dimensionIndex := range versionItem.Dimensions {
		dimension := &versionItem.Dimensions[dimensionIndex]
		dimMap[dimension.Name] = struct{}{}

		options, err := api.datasets.GetOptionsInBatches(ctx, "", api.cfg.ServiceAuthToken, "", params.id, params.edition, params.version, dimension.Name, api.cfg.DatasetOptionsBatchSize, api.cfg.DatasetOptionsWorkers)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get options")
		}

		params.options[dimension.ID] = make(map[string]dataset.Option)

		for optionIndex := range options.Items {
			option := options.Items[optionIndex]
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

		for dimIndex := range versionItem.Dimensions {
			dim := &versionItem.Dimensions[dimIndex]
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

func validateBaseParams(params *datasetParams) error {
	if params.id == "" {
		return errors.New("invalid dataset id")
	}

	if params.edition == "" {
		return errors.New("invalid edition")
	}

	if params.version == "" {
		return errors.New("invalid version")
	}

	return nil
}

const (
	batchSize     = 100
	numberWorkers = 10
)

func (api *API) getGeographyTypes(ctx context.Context, datasetID string) ([]string, error) {
	res, err := api.ctblr.GetGeographyDimensionsInBatches(ctx, datasetID, batchSize, numberWorkers)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get Geography Dimensions")
	}

	geoDimensions := make([]string, len(res.Variables.Edges))

	for i := range res.Variables.Edges {
		geoDimensions = append(geoDimensions, res.Variables.Edges[i].Node.Name)
	}

	return geoDimensions, nil
}

func (api *API) sortGeography(geoDimensions, dimensions []string) []string {
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
