package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular/gql"
	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/api/mock"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/config"
	dpresponder "github.com/ONSdigital/dp-net/v2/responder"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
)

var optionsBatch int = 20
var optionsWorker int = 2

type testParams struct {
	ctx           context.Context
	request       *http.Request
	response      *httptest.ResponseRecorder
	datasetId     string
	edition       string
	version       string
	dimensionId   string
	dimensionName string
	dimensions    dataset.VersionDimensions
	options       dataset.Options
}

var validIsBasedOn = dataset.IsBasedOn{
	ID: "testTable",
}

var validGetDatasetCurrentAndNextResult = dataset.Dataset{
	DatasetDetails: dataset.DatasetDetails{
		Type:      "cantabular_flexible_table",
		IsBasedOn: &validIsBasedOn,
	},
}

var validGeographyResponse = cantabular.GetGeographyDimensionsResponse{
	Dataset: gql.DatasetRuleBase{
		RuleBase: gql.RuleBase{
			IsSourceOf: gql.Variables{
				Edges: []gql.Edge{{
					Node: gql.Node{
						Name: "uk",
					}},
				},
			},
		},
	},
}

func TestMissingURLParams(t *testing.T) {
	api := API{}

	Convey("When getDatasetInfo is called with no dataset id", t, func() {
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("dataset_id", "")
		ctx.URLParams.Add("version", "1")
		ctx.URLParams.Add("edition", "1")

		request := httptest.NewRequest("GET", "/dataset/edition/1/version/1", nil)
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))

		ret, err := api.getDatasetParams(request.Context(), request)

		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "invalid dataset id")
		So(ret, ShouldBeNil)
	})

	Convey("When getDatasetInfo is called with no version id", t, func() {
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("dataset_id", "1")
		ctx.URLParams.Add("version", "")
		ctx.URLParams.Add("edition", "1")

		request := httptest.NewRequest("GET", "/dataset/1/version//edition/1", nil)
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))

		ret, err := api.getDatasetParams(request.Context(), request)

		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "invalid version")
		So(ret, ShouldBeNil)
	})

	Convey("When getDatasetInfo is called with no edition id", t, func() {
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("dataset_id", "1")
		ctx.URLParams.Add("version", "1")
		ctx.URLParams.Add("edition", "")

		request := httptest.NewRequest("GET", "/dataset/1/version/1/edition/", nil)
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))

		ret, err := api.getDatasetParams(request.Context(), request)

		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "invalid edition")
		So(ret, ShouldBeNil)
	})
}

func TestInvalidGeography(t *testing.T) {
	Convey("When geography cannot be accessed an error is returned", t, func() {
		api, ctrl, ctblrMock, datasetAPIMock := initTest(t)
		defer ctrl.Finish()

		expectedError := errors.New(uuid.NewString())
		p := getValidRequest()

		gomock.InOrder(
			datasetAPIMock.EXPECT().GetDatasetCurrentAndNext(p.ctx, "", "", "", p.datasetId).Return(validGetDatasetCurrentAndNextResult, nil).Times(1),
			datasetAPIMock.EXPECT().GetVersionDimensions(p.ctx, "", "", "", p.datasetId, p.edition, p.version).Return(p.dimensions, nil).Times(1),
			datasetAPIMock.EXPECT().GetOptionsInBatches(p.ctx, "", "", "", p.datasetId, p.edition, p.version, p.dimensionName, optionsBatch, optionsWorker).Return(p.options, nil).Times(1),
			ctblrMock.EXPECT().GetGeographyDimensions(p.ctx, validIsBasedOn.ID).Return(nil, expectedError).Times(1),
		)

		api.getDatasetJSON(p.response, p.request)

		res := p.response.Result()

		So(p.response, ShouldNotBeNil)
		So(res.StatusCode, ShouldEqual, http.StatusInternalServerError)
	})
}

func TestInvalidParams(t *testing.T) {
	Convey("When getDatasetInfo is called with an invalid dataset id then an error is returned", t, func() {
		api, ctrl, _, datasetAPIMock := initTest(t)
		defer ctrl.Finish()

		expectedError := errors.New(uuid.NewString())
		p := getValidRequest()

		datasetAPIMock.EXPECT().GetDatasetCurrentAndNext(p.ctx, "", "", "", p.datasetId).Return(dataset.Dataset{}, expectedError).Times(1)

		ret, err := api.getDatasetParams(p.ctx, p.request)

		So(err, ShouldNotBeNil)
		So(err, ShouldResemble, expectedError)
		So(ret, ShouldBeNil)
	})

	Convey("When GetVersionDimensions is called with an invalid dataset id then  an error is returned", t, func() {
		api, ctrl, _, datasetAPIMock := initTest(t)
		defer ctrl.Finish()

		expectedError := errors.New(uuid.NewString())
		p := getValidRequest()

		gomock.InOrder(
			datasetAPIMock.EXPECT().GetDatasetCurrentAndNext(p.ctx, "", "", "", p.datasetId).Return(validGetDatasetCurrentAndNextResult, nil).Times(1),
			datasetAPIMock.EXPECT().GetVersionDimensions(p.ctx, "", "", "", p.datasetId, p.edition, p.version).Return(dataset.VersionDimensions{}, expectedError).Times(1),
		)

		ret, err := api.getDatasetParams(p.ctx, p.request)

		So(err, ShouldNotBeNil)
		So(err, ShouldResemble, expectedError)
		So(ret, ShouldBeNil)
	})

	Convey("When GetVersionDimensions is called and returns an invalid list then an error is returned", t, func() {
		api, ctrl, _, datasetAPIMock := initTest(t)
		defer ctrl.Finish()

		p := getValidRequest()

		gomock.InOrder(
			datasetAPIMock.EXPECT().GetDatasetCurrentAndNext(p.ctx, "", "", "", p.datasetId).Return(validGetDatasetCurrentAndNextResult, nil).Times(1),
			datasetAPIMock.EXPECT().GetVersionDimensions(p.ctx, "", "", "", p.datasetId, p.edition, p.version).Return(dataset.VersionDimensions{}, nil).Times(1),
		)

		ret, err := api.getDatasetParams(p.ctx, p.request)

		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "invalid dimensions length")
		So(ret, ShouldBeNil)
	})

	Convey("When GetOptions is called with an invalid id an error is returned", t, func() {
		api, ctrl, _, datasetAPIMock := initTest(t)
		defer ctrl.Finish()

		p := getValidRequest()
		expectedError := errors.New(uuid.NewString())

		gomock.InOrder(
			datasetAPIMock.EXPECT().GetDatasetCurrentAndNext(p.ctx, "", "", "", p.datasetId).Return(validGetDatasetCurrentAndNextResult, nil).Times(1),
			datasetAPIMock.EXPECT().GetVersionDimensions(p.ctx, "", "", "", p.datasetId, p.edition, p.version).Return(p.dimensions, nil).Times(1),
			datasetAPIMock.EXPECT().GetOptionsInBatches(p.ctx, "", "", "", p.datasetId, p.edition, p.version, p.dimensionName, optionsBatch, optionsWorker).Return(dataset.Options{}, expectedError).Times(1),
		)

		ret, err := api.getDatasetParams(p.ctx, p.request)

		So(err, ShouldNotBeNil)
		So(err, ShouldResemble, expectedError)
		So(ret, ShouldBeNil)
	})

	Convey("When StaticDatasetQuery is called with an invalid query it returns an error", t, func() {
		api, ctrl, ctblrMock, datasetAPIMock := initTest(t)
		defer ctrl.Finish()

		p := getValidRequest()
		expectedError := errors.New(uuid.NewString())
		dimensionId := uuid.NewString()
		dimensionName := uuid.NewString()
		dimensions := dataset.VersionDimensions{
			Items: dataset.VersionDimensionItems{{
				ID:   dimensionId,
				Name: dimensionName,
				Links: dataset.Links{
					CodeList: dataset.Link{ID: "Test"},
				},
			}},
		}

		datasetRequest := cantabular.StaticDatasetQueryRequest{
			Dataset:   validIsBasedOn.ID,
			Variables: []string{dimensions.Items[0].Links.CodeList.ID},
		}

		gomock.InOrder(
			datasetAPIMock.EXPECT().GetDatasetCurrentAndNext(p.ctx, "", "", "", p.datasetId).Return(validGetDatasetCurrentAndNextResult, nil).Times(1),
			datasetAPIMock.EXPECT().GetVersionDimensions(p.ctx, "", "", "", p.datasetId, p.edition, p.version).Return(dimensions, nil).Times(1),
			datasetAPIMock.EXPECT().GetOptionsInBatches(p.ctx, "", "", "", p.datasetId, p.edition, p.version, dimensionName, optionsBatch, optionsWorker).Return(p.options, nil).Times(1),
			ctblrMock.EXPECT().GetGeographyDimensions(p.ctx, validIsBasedOn.ID).Return(&validGeographyResponse, nil).Times(1),
			ctblrMock.EXPECT().StaticDatasetQuery(p.ctx, datasetRequest).Return(nil, expectedError).Times(1),
		)

		api.getDatasetJSON(p.response, p.request)
		res := p.response.Result()

		So(p.response, ShouldNotBeNil)
		So(res.StatusCode, ShouldEqual, http.StatusInternalServerError)
	})
}

func TestGeographySort(t *testing.T) {
	api := API{}

	Convey("When a geography dimension matches it should be placed at the start of the dimension list", t, func() {
		geoDimensions := []string{"REGION"}

		unsortedDimensions := []string{"AGE", "REGION"}

		result, foundGeo := api.sortGeography(geoDimensions, unsortedDimensions)

		So(result, ShouldResemble, []string{"REGION", "AGE"})
		So(foundGeo, ShouldBeTrue)
	})

	Convey("When a geography dimension does not match the dimension list should be unchanged", t, func() {
		geoDimensions := []string{"REGION"}

		unsortedDimensions := []string{"AGE", "AGE2"}

		result, foundGeo := api.sortGeography(geoDimensions, unsortedDimensions)

		So(result, ShouldResemble, []string{"AGE", "AGE2"})
		So(foundGeo, ShouldBeFalse)
	})

	Convey("When additional geography dimensions match they should be ignored", t, func() {
		geoDimensions := []string{"REGION", "REGION2"}

		unsortedDimensions := []string{"AGE", "REGION", "REGION2"}

		result, foundGeo := api.sortGeography(geoDimensions, unsortedDimensions)

		So(result, ShouldResemble, []string{"REGION", "AGE"})
		So(foundGeo, ShouldBeTrue)
	})
}

func TestToGetDatasetJsonResponse(t *testing.T) {
	Convey("When TestToGetDatasetJsonResponse is called a valid response should be returned", t, func() {
		api, ctrl, _, _ := initTest(t)
		defer ctrl.Finish()

		queryResult := &cantabular.StaticDatasetQuery{
			Dataset: cantabular.StaticDataset{
				Table: cantabular.Table{
					Dimensions: []cantabular.Dimension{
						{
							Categories: []cantabular.Category{{
								Label: "16",
							}},
							Variable: cantabular.VariableBase{
								Name: "AGE",
							},
						},
					},
				},
			},
		}

		optionsMap := make(optionsMap)

		optionsMap["AGE"] = make(map[string]dataset.Option)

		option := dataset.Option{
			Label: "16",
			Links: dataset.Links{
				Code: dataset.Link{
					URL: "not valid",
				},
			},
		}

		optionsMap["AGE"]["16"] = option

		result, err := api.toGetDatasetJsonResponse(optionsMap, queryResult)

		So(err, ShouldBeNil)
		So(result, ShouldNotBeNil)
		So(result.Observations, ShouldResemble, queryResult.Dataset.Table.Values)
		So(result.TotalObservations, ShouldEqual, len(queryResult.Dataset.Table.Values))
		So(len(result.Dimensions), ShouldEqual, 1)
		So(result.Dimensions[0].DimensionName, ShouldEqual, ("AGE"))
		So(len(result.Dimensions[0].Options), ShouldEqual, 1)
		So(result.Dimensions[0].Options[0].Href, ShouldEqual, "not valid")
		So(result.Dimensions[0].Options[0].Id, ShouldEqual, "16")
	})
}

func initTest(t *testing.T) (API, *gomock.Controller, *mock.MockcantabularClient, *mock.MockdatasetAPIClient) {
	var api API

	ctrl := gomock.NewController(t)

	ctblrMock := mock.NewMockcantabularClient(ctrl)
	api.ctblr = ctblrMock
	datasetAPIMock := mock.NewMockdatasetAPIClient(ctrl)
	api.datasets = datasetAPIMock
	api.respond = dpresponder.New()
	api.cfg = &config.Config{
		DatasetOptionsBatchSize: optionsBatch,
		DatasetOptionsWorkers:   optionsWorker,
	}

	return api, ctrl, ctblrMock, datasetAPIMock
}

func getValidRequest() *testParams {
	datasetId := "datasetId-" + uuid.NewString()
	edition := "edition-" + uuid.NewString()
	version := "version-" + uuid.NewString()

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("dataset_id", datasetId)
	ctx.URLParams.Add("version", version)
	ctx.URLParams.Add("edition", edition)

	dimensionId := "dimensionId" + uuid.NewString()
	dimensionName := "dimensionName" + uuid.NewString()
	dimensions := dataset.VersionDimensions{
		Items: dataset.VersionDimensionItems{{ID: dimensionId, Name: dimensionName}},
	}
	options := dataset.Options{Items: []dataset.Option{{}}}

	request := httptest.NewRequest("GET", "/dataset", nil)
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))

	return &testParams{
		request.Context(),
		request,
		httptest.NewRecorder(),
		datasetId,
		edition,
		version,
		dimensionId,
		dimensionName,
		dimensions,
		options,
	}
}
