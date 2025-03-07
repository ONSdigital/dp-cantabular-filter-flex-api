package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular/gql"
	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/api/mock"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/config"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	dpresponder "github.com/ONSdigital/dp-net/v2/responder"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	optionsBatch  = 20
	optionsWorker = 2
)

type testParams struct {
	ctx       context.Context
	request   *http.Request
	response  *httptest.ResponseRecorder
	datasetID string
	edition   string
	version   string
}

func TestMissingURLParams(t *testing.T) {
	api := API{}

	Convey("When getDatasetJSON is called with no dataset id", t, func() {
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("dataset_id", "")
		ctx.URLParams.Add("version", "1")
		ctx.URLParams.Add("edition", "1")

		request := httptest.NewRequest("GET", "/dataset/edition/1/version/1", http.NoBody)
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))

		ret, err := api.getDatasetParams(request.Context(), request)

		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "invalid dataset id")
		So(ret, ShouldBeNil)
	})

	Convey("When getDatasetJSON is called with no version id", t, func() {
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("dataset_id", "1")
		ctx.URLParams.Add("version", "")
		ctx.URLParams.Add("edition", "1")

		request := httptest.NewRequest("GET", "/dataset/1/version//edition/1", http.NoBody)
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))

		ret, err := api.getDatasetParams(request.Context(), request)

		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "invalid version")
		So(ret, ShouldBeNil)
	})

	Convey("When getDatasetJSON is called with no edition id", t, func() {
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("dataset_id", "1")
		ctx.URLParams.Add("version", "1")
		ctx.URLParams.Add("edition", "")

		request := httptest.NewRequest("GET", "/dataset/1/version/1/edition/", http.NoBody)
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))

		ret, err := api.getDatasetParams(request.Context(), request)

		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "invalid edition")
		So(ret, ShouldBeNil)
	})
}

func TestInvalidGeography(t *testing.T) {
	Convey("When geography cannot be accessed an error is returned", t, func() {
		api, ctrl, ctblrMock, datasetAPIMock := initMocks(t)
		defer ctrl.Finish()

		expectedError := errors.New(uuid.NewString())
		p := getTestParams()

		versionResponse := getValidVersionResponse()
		gomock.InOrder(
			datasetAPIMock.EXPECT().GetVersion(p.ctx, "", "", "", "", p.datasetID, p.edition, p.version).Return(versionResponse, nil).Times(1),
			datasetAPIMock.EXPECT().GetMetadataURL(p.datasetID, p.edition, p.version).Return(getValidMetadataURLResponse()).Times(1),
			datasetAPIMock.EXPECT().GetOptionsInBatches(p.ctx, "", "", "", p.datasetID, p.edition, p.version, versionResponse.Dimensions[0].Name, optionsBatch, optionsWorker).Return(getValidOptionsResponse(), nil).Times(1),
			ctblrMock.EXPECT().GetGeographyDimensionsInBatches(p.ctx, versionResponse.IsBasedOn.ID, batchSize, numberWorkers).Return(nil, expectedError).Times(1),
		)

		ret, err := api.getDatasetParams(p.ctx, p.request)
		So(ret, ShouldBeNil)
		So(err.Error(), ShouldEqual, "failed to get geography types: failed to get Geography Dimensions: "+expectedError.Error())
	})
}

func TestGetVersion(t *testing.T) {
	Convey("When GetVersion is called with an invalid dataset id then an error is returned", t, func() {
		api, ctrl, _, datasetAPIMock := initMocks(t)
		defer ctrl.Finish()

		expectedError := errors.New(uuid.NewString())
		p := getTestParams()

		gomock.InOrder(
			datasetAPIMock.EXPECT().GetVersion(p.ctx, "", "", "", "", p.datasetID, p.edition, p.version).Return(dataset.Version{}, expectedError).Times(1),
		)

		ret, err := api.getDatasetParams(p.ctx, p.request)

		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "failed to get version: "+expectedError.Error())
		So(ret, ShouldBeNil)
	})
}

func TestGetOptions(t *testing.T) {
	Convey("When GetOptions is called with an invalid id an error is returned", t, func() {
		api, ctrl, _, datasetAPIMock := initMocks(t)
		defer ctrl.Finish()

		p := getTestParams()
		expectedError := errors.New(uuid.NewString())
		versionResponse := getValidVersionResponse()

		gomock.InOrder(
			datasetAPIMock.EXPECT().GetVersion(p.ctx, "", "", "", "", p.datasetID, p.edition, p.version).Return(versionResponse, nil).Times(1),
			datasetAPIMock.EXPECT().GetMetadataURL(p.datasetID, p.edition, p.version).Return(getValidMetadataURLResponse()).Times(1),
			datasetAPIMock.EXPECT().GetOptionsInBatches(p.ctx, "", "", "", p.datasetID, p.edition, p.version, versionResponse.Dimensions[0].Name, optionsBatch, optionsWorker).Return(dataset.Options{}, expectedError).Times(1),
		)

		ret, err := api.getDatasetParams(p.ctx, p.request)

		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "failed to get options: "+expectedError.Error())
		So(ret, ShouldBeNil)
	})
}

func TestStaticDatasetQuery(t *testing.T) {
	Convey("When StaticDatasetQuery is called with an invalid query it returns an error", t, func() {
		api, ctrl, ctblrMock, datasetAPIMock := initMocks(t)
		defer ctrl.Finish()

		p := getTestParams()
		expectedError := errors.New(uuid.NewString())

		versionResponse := getValidVersionResponse()

		datasetRequest := cantabular.StaticDatasetQueryRequest{
			Dataset:   versionResponse.IsBasedOn.ID,
			Variables: []string{versionResponse.Dimensions[0].ID},
		}

		gomock.InOrder(
			datasetAPIMock.EXPECT().GetVersion(p.ctx, "", "", "", "", p.datasetID, p.edition, p.version).Return(versionResponse, nil).Times(1),
			datasetAPIMock.EXPECT().GetMetadataURL(p.datasetID, p.edition, p.version).Return(getValidMetadataURLResponse()).Times(1),
			datasetAPIMock.EXPECT().GetOptionsInBatches(p.ctx, "", "", "", p.datasetID, p.edition, p.version, versionResponse.Dimensions[0].Name, optionsBatch, optionsWorker).Return(getValidOptionsResponse(), nil).Times(1),
			ctblrMock.EXPECT().GetGeographyDimensionsInBatches(p.ctx, versionResponse.IsBasedOn.ID, batchSize, numberWorkers).Return(getValidGeoResponse(), nil).Times(1),
			ctblrMock.EXPECT().StaticDatasetQuery(p.ctx, datasetRequest).Return(nil, expectedError).Times(1),
		)

		params, err := api.getDatasetParams(p.ctx, p.request)
		So(err, ShouldBeNil)
		So(params, ShouldNotBeNil)

		result, err := api.getDatasetJSON(p.ctx, p.request, params)

		So(result, ShouldBeNil)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "failed to run query: "+expectedError.Error())
	})
}

func TestGeographySort(t *testing.T) {
	api := API{}

	Convey("When a geography dimension matches it should be placed at the start of the dimension list", t, func() {
		geoDimensions := []string{"REGION"}

		unsortedDimensions := []string{"AGE", "REGION"}

		result := api.sortGeography(geoDimensions, unsortedDimensions)

		So(result, ShouldResemble, []string{"REGION", "AGE"})
	})

	Convey("When a geography dimension does not match the dimension list should be unchanged", t, func() {
		geoDimensions := []string{"REGION"}

		unsortedDimensions := []string{"AGE", "AGE2"}

		result := api.sortGeography(geoDimensions, unsortedDimensions)

		So(result, ShouldResemble, []string{"AGE", "AGE2"})
	})

	Convey("When additional geography dimensions match they should be ignored", t, func() {
		geoDimensions := []string{"REGION", "REGION2"}

		unsortedDimensions := []string{"AGE", "REGION", "REGION2"}

		result := api.sortGeography(geoDimensions, unsortedDimensions)

		So(result, ShouldResemble, []string{"REGION", "AGE"})
	})
}

func TestToGetJsonResponse(t *testing.T) {
	Convey("When TestToGetJsonResponse is called a valid response should be returned", t, func() {
		api, ctrl, ctblrMock, datasetAPIMock := initMocks(t)
		defer ctrl.Finish()
		p := getTestParams()

		versionResponse := getValidVersionResponse()
		metadataResponse := getValidMetadataURLResponse()
		optionsResponse := getValidOptionsResponse()

		cantabularResponse := getValidCantabularResponse(versionResponse.Dimensions[0].ID, optionsResponse.Items[0].Label)

		gomock.InOrder(
			datasetAPIMock.EXPECT().GetVersion(p.ctx, "", "", "", "", p.datasetID, p.edition, p.version).Return(versionResponse, nil).Times(1),
			datasetAPIMock.EXPECT().GetMetadataURL(p.datasetID, p.edition, p.version).Return(metadataResponse).Times(1),
			datasetAPIMock.EXPECT().GetOptionsInBatches(p.ctx, "", "", "", p.datasetID, p.edition, p.version, versionResponse.Dimensions[0].Name, optionsBatch, optionsWorker).Return(optionsResponse, nil).Times(1),
			ctblrMock.EXPECT().GetGeographyDimensionsInBatches(p.ctx, versionResponse.IsBasedOn.ID, batchSize, numberWorkers).Return(getValidGeoResponse(), nil).Times(1),
		)

		params, err := api.getDatasetParams(p.ctx, p.request)
		So(err, ShouldBeNil)

		result, err := api.toGetDatasetJSONResponse(p.request, params, &cantabularResponse)

		So(err, ShouldBeNil)
		So(result, ShouldNotBeNil)
		So(result.Observations, ShouldResemble, cantabularResponse.Dataset.Table.Values)
		So(result.TotalObservations, ShouldEqual, len(cantabularResponse.Dataset.Table.Values))
		So(len(result.Dimensions), ShouldEqual, 1)
		So(result.Dimensions[0].DimensionName, ShouldEqual, versionResponse.Dimensions[0].ID)
		So(len(result.Dimensions[0].Options), ShouldEqual, 1)
		So(result.Dimensions[0].Options[0].ID, ShouldEqual, optionsResponse.Items[0].Links.Code.ID)
		So(result.Links.DatasetMetadata.HREF, ShouldEqual, metadataResponse)
		So(result.Links.Self.HREF, ShouldEqual, versionResponse.Links.Dataset.URL)
		So(result.Links.Self.ID, ShouldEqual, versionResponse.Links.Dataset.ID)
		So(result.Links.Version.HREF, ShouldEqual, versionResponse.Links.Self.URL)
		So(result.Links.Version.ID, ShouldEqual, versionResponse.Links.Self.ID)
	})
}

func TestToGetJsonResponseWithURLRewriting(t *testing.T) {
	Convey("When TestToGetJsonResponse is called and URL rewriting is enabled, a valid response should be returned", t, func() {
		api, ctrl, ctblrMock, datasetAPIMock := initMocksURLRewriting(t)
		defer ctrl.Finish()
		p := getTestParams()
		p.request.Header.Add("X-Forwarded-Host", "api.example")
		p.request.Header.Add("X-Forwarded-Path-Prefix", "v1")

		versionResponse := getValidURLRewrittenVersionResponse()
		metadataResponse := "https://api.example/v1/metadataURL%2070b53334-eea5-42ac-a456-93136f477865"
		optionsResponse := getValidOptionsResponse()

		cantabularResponse := getValidCantabularResponse(versionResponse.Dimensions[0].ID, optionsResponse.Items[0].Label)

		gomock.InOrder(
			datasetAPIMock.EXPECT().GetVersion(p.ctx, "", "", "", "", p.datasetID, p.edition, p.version).Return(versionResponse, nil).Times(1),
			datasetAPIMock.EXPECT().GetMetadataURL(p.datasetID, p.edition, p.version).Return(metadataResponse).Times(1),
			datasetAPIMock.EXPECT().GetOptionsInBatches(p.ctx, "", "", "", p.datasetID, p.edition, p.version, versionResponse.Dimensions[0].Name, optionsBatch, optionsWorker).Return(optionsResponse, nil).Times(1),
			ctblrMock.EXPECT().GetGeographyDimensionsInBatches(p.ctx, versionResponse.IsBasedOn.ID, batchSize, numberWorkers).Return(getValidGeoResponse(), nil).Times(1),
		)

		params, err := api.getDatasetParams(p.ctx, p.request)
		So(err, ShouldBeNil)

		result, err := api.toGetDatasetJSONResponse(p.request, params, &cantabularResponse)

		So(err, ShouldBeNil)
		So(result, ShouldNotBeNil)
		So(result.Observations, ShouldResemble, cantabularResponse.Dataset.Table.Values)
		So(result.TotalObservations, ShouldEqual, len(cantabularResponse.Dataset.Table.Values))
		So(len(result.Dimensions), ShouldEqual, 1)
		So(result.Dimensions[0].DimensionName, ShouldEqual, versionResponse.Dimensions[0].ID)
		So(len(result.Dimensions[0].Options), ShouldEqual, 1)
		So(result.Dimensions[0].Options[0].ID, ShouldEqual, optionsResponse.Items[0].Links.Code.ID)
		So(result.Links.DatasetMetadata.HREF, ShouldEqual, metadataResponse)
		So(result.Links.Self.HREF, ShouldEqual, versionResponse.Links.Dataset.URL)
		So(result.Links.Self.ID, ShouldEqual, versionResponse.Links.Dataset.ID)
		So(result.Links.Version.HREF, ShouldEqual, versionResponse.Links.Self.URL)
		So(result.Links.Version.ID, ShouldEqual, versionResponse.Links.Self.ID)
	})
}

func TestToGetObservationsResponse(t *testing.T) {
	Convey("When TestToGetObservationsResponse is called a valid response should be returned", t, func() {
		api, ctrl, ctblrMock, datasetAPIMock := initMocks(t)
		defer ctrl.Finish()
		p := getTestParams()

		versionResponse := getValidVersionResponse()
		metadataResponse := getValidMetadataURLResponse()
		optionsResponse := getValidOptionsResponse()

		observationsResponse := getValidObservationsResponse(versionResponse.Dimensions[0].ID)

		cantabularResponse := getValidCantabularResponse(versionResponse.Dimensions[0].ID, optionsResponse.Items[0].Label)

		gomock.InOrder(
			datasetAPIMock.EXPECT().GetVersion(p.ctx, "", "", "", "", p.datasetID, p.edition, p.version).Return(versionResponse, nil).Times(1),
			datasetAPIMock.EXPECT().GetMetadataURL(p.datasetID, p.edition, p.version).Return(metadataResponse).Times(1),
			datasetAPIMock.EXPECT().GetOptionsInBatches(p.ctx, "", "", "", p.datasetID, p.edition, p.version, versionResponse.Dimensions[0].Name, optionsBatch, optionsWorker).Return(optionsResponse, nil).Times(1),
			ctblrMock.EXPECT().GetGeographyDimensionsInBatches(p.ctx, versionResponse.IsBasedOn.ID, batchSize, numberWorkers).Return(getValidGeoResponse(), nil).Times(1),
		)

		params, err := api.getDatasetParams(p.ctx, p.request)
		So(err, ShouldBeNil)

		result, err := api.toGetDatasetObservationsResponse(params, &cantabularResponse)

		So(err, ShouldBeNil)
		So(result, ShouldNotBeNil)
		So(result.TotalObservations, ShouldEqual, len(cantabularResponse.Dataset.Table.Values))
		So(len(result.Observations), ShouldEqual, 1)
		So(result.Observations[0].Dimensions[0].DimensionID, ShouldEqual, observationsResponse.Observations[0].Dimensions[0].DimensionID)
		So(len(result.Observations[0].Dimensions), ShouldEqual, 1)
		So(result.Observations[0].Dimensions[0].OptionID, ShouldEqual, observationsResponse.Observations[0].Dimensions[0].OptionID)
		So(result.Links.DatasetMetadata.HREF, ShouldEqual, metadataResponse)
		So(result.Links.Self.HREF, ShouldEqual, versionResponse.Links.Dataset.URL)
		So(result.Links.Self.ID, ShouldEqual, versionResponse.Links.Dataset.ID)
		So(result.Links.Version.HREF, ShouldEqual, versionResponse.Links.Self.URL)
		So(result.Links.Version.ID, ShouldEqual, versionResponse.Links.Self.ID)
	})
}

func TestGetGeographyFiltersGeoInputs(t *testing.T) {
	api := API{}

	Convey("WHEN getGeography is called with blank geography THEN an error is returned", t, func() {
		request := httptest.NewRequest("GET", "/dataset/edition/1/version/1", http.NoBody)
		result, err := api.getGeographyFilters(request, nil)

		So(result, ShouldBeNil)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "unable to locate geography")
	})

	Convey("WHEN getGeographyFilters is called with an invalid geo query string THEN an error is returned", t, func() {
		request := httptest.NewRequest("GET", "/dataset?area-type=ABC", http.NoBody)
		result, err := api.getGeographyFilters(request, nil)

		So(result, ShouldBeNil)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "unable to locate geography")
	})

	Convey("WHEN getGeographyFilters is called with a geo dimension which does not exist THEN an error is returned", t, func() {
		params := &datasetParams{
			geoDimensions: []string{},
		}
		request := httptest.NewRequest("GET", "/dataset?area-type=ABC,DEF", http.NoBody)
		result, err := api.getGeographyFilters(request, params)

		So(result, ShouldBeNil)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "unable to validate geography ABC")
	})

	Convey("WHEN getGeographyFilters is called with a geo dimension which does not exist THEN an error is returned", t, func() {
		region := "REGION"
		params := &datasetParams{
			geoDimensions: []string{region},
			options:       optionsMap{region: nil},
		}
		request := httptest.NewRequest("GET", "/dataset?area-type=REGION,DEF", http.NoBody)

		result, err := api.getGeographyFilters(request, params)

		So(result, ShouldBeNil)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "unable to validate geography option DEF")
	})
}

func TestGetGeographyFiltersDimensionInput(t *testing.T) {
	api := API{}

	Convey("WHEN getGeographyFilters is called with no dimension THEN an error is returned", t, func() {
		region := uuid.NewString()
		area := uuid.NewString()

		params := &datasetParams{
			geoDimensions: []string{region},
			options:       optionsMap{region: map[string]dataset.Option{area: {}}},
		}
		request := httptest.NewRequest("GET", fmt.Sprintf("/dataset?area-type=%s,%s", region, area), http.NoBody)
		result, err := api.getGeographyFilters(request, params)

		So(result, ShouldBeNil)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "unable to locate dimension")
	})

	Convey("WHEN getGeographyFilters is called with an invalid dimension THEN an error is returned", t, func() {
		region := uuid.NewString()
		area := uuid.NewString()
		dimension := uuid.NewString()

		params := &datasetParams{
			geoDimensions: []string{region},
			options:       optionsMap{region: map[string]dataset.Option{area: {}}},
		}
		request := httptest.NewRequest("GET", fmt.Sprintf("/dataset?area-type=%s,%s&dimension=%s", region, area, dimension), http.NoBody)
		result, err := api.getGeographyFilters(request, params)

		So(result, ShouldBeNil)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "unable to validate dimension "+dimension)
	})
}

func TestGetGeographyFiltersOptions(t *testing.T) {
	api := API{}

	Convey("WHEN getGeographyFilters is called with no options THEN an error is returned", t, func() {
		region := uuid.NewString()
		area := uuid.NewString()
		dimension := uuid.NewString()

		params := &datasetParams{
			geoDimensions:     []string{region},
			options:           optionsMap{region: map[string]dataset.Option{area: {}}},
			datasetDimensions: []string{dimension},
		}
		request := httptest.NewRequest("GET", fmt.Sprintf("/dataset?area-type=%s,%s&dimension=%s", region, area, dimension), http.NoBody)
		result, err := api.getGeographyFilters(request, params)

		So(result, ShouldBeNil)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "invalid options length or options is empty")
	})

	Convey("WHEN getGeographyFilters is called with an invalid option THEN an error is returned", t, func() {
		region := uuid.NewString()
		area := uuid.NewString()
		dimension := uuid.NewString()
		optionValid := uuid.NewString()
		optionInvalid := uuid.NewString()

		testOptionsMap := make(optionsMap)
		testOptionsMap[region] = map[string]dataset.Option{area: {}}
		testOptionsMap[dimension] = map[string]dataset.Option{optionValid: {}}

		params := &datasetParams{
			geoDimensions:     []string{region},
			options:           testOptionsMap,
			datasetDimensions: []string{dimension},
		}
		request := httptest.NewRequest("GET", fmt.Sprintf("/dataset?area-type=%s,%s&dimension=%s&options=%s,%s", region, area, dimension, optionValid, optionInvalid), http.NoBody)
		result, err := api.getGeographyFilters(request, params)

		So(result, ShouldBeNil)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "unable to locate dimension option "+optionInvalid)
	})
}

func TestHandlerErrors(t *testing.T) {
	Convey("When getDatasetParams raises an error it is returned", t, func() {
		api, ctrl, _, _ := initMocks(t)
		defer ctrl.Finish()

		response := httptest.NewRecorder()
		api.getDatasetJSONHandler(response, httptest.NewRequest("GET", "/dataset", http.NoBody))

		So(response.Result().StatusCode, ShouldEqual, http.StatusInternalServerError)
	})

	Convey("When getGeographyFilters raises an error it is returned", t, func() {
		api, ctrl, ctblrMock, datasetAPIMock := initMocks(t)
		defer ctrl.Finish()
		p := getTestParams()

		versionResponse := getValidVersionResponse()

		gomock.InOrder(
			datasetAPIMock.EXPECT().GetVersion(p.ctx, "", "", "", "", p.datasetID, p.edition, p.version).Return(versionResponse, nil).Times(1),
			datasetAPIMock.EXPECT().GetMetadataURL(p.datasetID, p.edition, p.version).Return(getValidMetadataURLResponse()).Times(1),
			datasetAPIMock.EXPECT().GetOptionsInBatches(p.ctx, "", "", "", p.datasetID, p.edition, p.version, versionResponse.Dimensions[0].Name, optionsBatch, optionsWorker).Return(getValidOptionsResponse(), nil).Times(1),
			ctblrMock.EXPECT().GetGeographyDimensionsInBatches(p.ctx, versionResponse.IsBasedOn.ID, batchSize, numberWorkers).Return(getValidGeoResponse(), nil).Times(1),
		)

		request := httptest.NewRequest("GET", "/dataset?area-type=test", http.NoBody)
		request = request.WithContext(p.request.Context())

		api.getDatasetJSONHandler(p.response, request)

		So(p.response.Result().StatusCode, ShouldEqual, http.StatusInternalServerError)
	})

	Convey("When StaticDatasetQuery raises an error it is returned", t, func() {
		api, ctrl, ctblrMock, datasetAPIMock := initMocks(t)
		defer ctrl.Finish()
		p := getTestParams()

		versionResponse := getValidVersionResponse()
		metadataResponse := getValidMetadataURLResponse()
		optionsResponse := getValidOptionsResponse()

		datasetRequest := cantabular.StaticDatasetQueryRequest{
			Dataset:   versionResponse.IsBasedOn.ID,
			Variables: []string{versionResponse.Dimensions[0].ID},
		}

		gomock.InOrder(
			datasetAPIMock.EXPECT().GetVersion(p.ctx, "", "", "", "", p.datasetID, p.edition, p.version).Return(versionResponse, nil).Times(1),
			datasetAPIMock.EXPECT().GetMetadataURL(p.datasetID, p.edition, p.version).Return(metadataResponse).Times(1),
			datasetAPIMock.EXPECT().GetOptionsInBatches(p.ctx, "", "", "", p.datasetID, p.edition, p.version, versionResponse.Dimensions[0].Name, optionsBatch, optionsWorker).Return(optionsResponse, nil).Times(1),
			ctblrMock.EXPECT().GetGeographyDimensionsInBatches(p.ctx, versionResponse.IsBasedOn.ID, batchSize, numberWorkers).Return(getValidGeoResponse(), nil).Times(1),
			ctblrMock.EXPECT().StaticDatasetQuery(p.ctx, datasetRequest).Return(nil, errors.New(uuid.NewString())).Times(1),
		)

		api.getDatasetJSONHandler(p.response, p.request)

		So(p.response.Result().StatusCode, ShouldEqual, http.StatusInternalServerError)
	})
}

func initMocks(t *testing.T) (API, *gomock.Controller, *mock.MockcantabularClient, *mock.MockdatasetAPIClient) {
	var api API

	ctrl := gomock.NewController(t)

	ctblrMock := mock.NewMockcantabularClient(ctrl)
	api.ctblr = ctblrMock
	datasetAPIMock := mock.NewMockdatasetAPIClient(ctrl)
	api.datasets = datasetAPIMock
	api.respond = dpresponder.New()
	api.datasetAPIURL = &url.URL{Scheme: "http", Host: "localhost:27100"}
	api.cfg = &config.Config{
		DatasetOptionsBatchSize: optionsBatch,
		DatasetOptionsWorkers:   optionsWorker,
	}

	return api, ctrl, ctblrMock, datasetAPIMock
}

func initMocksURLRewriting(t *testing.T) (API, *gomock.Controller, *mock.MockcantabularClient, *mock.MockdatasetAPIClient) {
	var api API

	ctrl := gomock.NewController(t)

	ctblrMock := mock.NewMockcantabularClient(ctrl)
	api.ctblr = ctblrMock
	datasetAPIMock := mock.NewMockdatasetAPIClient(ctrl)
	api.datasets = datasetAPIMock
	api.respond = dpresponder.New()
	api.datasetAPIURL = &url.URL{Scheme: "http", Host: "localhost:27100"}
	api.cfg = &config.Config{
		DatasetOptionsBatchSize: optionsBatch,
		DatasetOptionsWorkers:   optionsWorker,
		EnableURLRewriting:      true,
	}

	return api, ctrl, ctblrMock, datasetAPIMock
}

func getTestParams() *testParams {
	datasetID := "datasetID-" + uuid.NewString()
	edition := "edition-" + uuid.NewString()
	version := "version-" + uuid.NewString()

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("dataset_id", datasetID)
	ctx.URLParams.Add("version", version)
	ctx.URLParams.Add("edition", edition)

	request := httptest.NewRequest("GET", "/dataset", http.NoBody)
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))

	return &testParams{
		request.Context(),
		request,
		httptest.NewRecorder(),
		datasetID,
		edition,
		version,
	}
}

func getValidVersionResponse() dataset.Version {
	isBasedOn := &dataset.IsBasedOn{
		ID:   "basedOn" + uuid.NewString(),
		Type: "cantabular_flexible_table",
	}

	return dataset.Version{
		IsBasedOn: isBasedOn,
		Dimensions: []dataset.VersionDimension{
			{
				ID:   "dimensionId " + uuid.NewString(),
				Name: "dimensionName " + uuid.NewString(),
				URL:  "codelistURL " + uuid.NewString(),
			},
		},
		Links: dataset.Links{
			Dataset: dataset.Link{
				URL: "datasetURL " + uuid.NewString(),
				ID:  "datasetID " + uuid.NewString(),
			},
			Self: dataset.Link{
				URL: "versionURL " + uuid.NewString(),
				ID:  "versionId " + uuid.NewString(),
			},
		},
	}
}

func getValidURLRewrittenVersionResponse() dataset.Version {
	isBasedOn := &dataset.IsBasedOn{
		ID:   "basedOn" + uuid.NewString(),
		Type: "cantabular_flexible_table",
	}

	return dataset.Version{
		IsBasedOn: isBasedOn,
		Dimensions: []dataset.VersionDimension{
			{
				ID:   "dimensionId " + uuid.NewString(),
				Name: "dimensionName " + uuid.NewString(),
				URL:  "codelistURL " + uuid.NewString(),
			},
		},
		Links: dataset.Links{
			Dataset: dataset.Link{
				URL: "https://api.example/v1/" + uuid.NewString(),
				ID:  "datasetID " + uuid.NewString(),
			},
			Self: dataset.Link{
				URL: "https://api.example/v1/" + uuid.NewString(),
				ID:  "versionId " + uuid.NewString(),
			},
		},
	}
}

func getValidMetadataURLResponse() string {
	return "metadataURL " + uuid.NewString()
}

func getValidOptionsResponse() dataset.Options {
	return dataset.Options{
		Items: []dataset.Option{
			{
				Label: "optionLabel " + uuid.NewString(),
				Links: dataset.Links{
					Code: dataset.Link{
						URL: "codeURL " + uuid.NewString(),
						ID:  "1",
					},
				},
			},
		},
	}
}

func getValidObservationsResponse(dimensionID string) GetObservationsResponse {
	var obsDimension []ObservationDimension
	obsDimension = append(obsDimension, ObservationDimension{
		Dimension:   "test",
		DimensionID: dimensionID,
		Option:      "test",
		OptionID:    "1",
	})

	var getObservationResponse []GetObservationResponse
	getObservationResponse = append(getObservationResponse, GetObservationResponse{
		Dimensions:  obsDimension,
		Observation: 2,
	})

	return GetObservationsResponse{
		Observations:      getObservationResponse,
		TotalObservations: 1,
		Links: DatasetJSONLinks{
			Version: model.Link{
				ID: "1",
			},
			DatasetMetadata: model.Link{
				ID: "testing",
			},
		},
	}
}

func getValidGeoResponse() *gql.Dataset {
	return &gql.Dataset{
		RuleBase: gql.RuleBase{
			IsSourceOf: gql.Variables{
				Edges: []gql.Edge{
					{
						Node: gql.Node{
							Name: "NodeName " + uuid.NewString(),
						},
					},
				},
			},
		},
	}
}

func getValidCantabularResponse(name, label string) cantabular.StaticDatasetQuery {
	return cantabular.StaticDatasetQuery{
		Dataset: cantabular.StaticDataset{
			Table: cantabular.Table{
				Dimensions: []cantabular.Dimension{
					{
						Categories: []cantabular.Category{
							{
								Label: label,
								Code:  "1",
							},
						},
						Variable: cantabular.VariableBase{
							Name: name,
						},
					},
				},
				Values: []float32{100.0},
			},
		},
	}
}
