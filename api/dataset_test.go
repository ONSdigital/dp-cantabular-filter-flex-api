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

var (
	optionsBatch  = 20
	optionsWorker = 2
)

type testParams struct {
	ctx       context.Context
	request   *http.Request
	response  *httptest.ResponseRecorder
	datasetId string
	edition   string
	version   string
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
		api, ctrl, ctblrMock, datasetAPIMock := initMocks(t)
		defer ctrl.Finish()

		expectedError := errors.New(uuid.NewString())
		p := getTestParams()

		versionResponse := getValidVersionResponse()
		geographDimensionsRequest := cantabular.GetGeographyDimensionsRequest{
			Dataset: versionResponse.IsBasedOn.ID,
		}
		gomock.InOrder(
			datasetAPIMock.EXPECT().GetVersion(p.ctx, "", "", "", "", p.datasetId, p.edition, p.version).Return(versionResponse, nil).Times(1),
			datasetAPIMock.EXPECT().GetVersionMetadata(p.ctx, "", "", "", p.datasetId, p.edition, p.version).Return(getValidMetadataResponse(), nil).Times(1),
			datasetAPIMock.EXPECT().GetOptionsInBatches(p.ctx, "", "", "", p.datasetId, p.edition, p.version, versionResponse.Dimensions[0].Name, optionsBatch, optionsWorker).Return(getValidOptionsResponse(), nil).Times(1),
			ctblrMock.EXPECT().GetGeographyDimensions(p.ctx, geographDimensionsRequest).Return(nil, expectedError).Times(1),
		)

		api.getDatasetJSON(p.response, p.request)

		res := p.response.Result()

		So(p.response, ShouldNotBeNil)
		So(res.StatusCode, ShouldEqual, http.StatusInternalServerError)
	})
}

func TestGetVersion(t *testing.T) {
	Convey("When GetVersion is called with an invalid dataset id then an error is returned", t, func() {
		api, ctrl, _, datasetAPIMock := initMocks(t)
		defer ctrl.Finish()

		expectedError := errors.New(uuid.NewString())
		p := getTestParams()

		gomock.InOrder(
			datasetAPIMock.EXPECT().GetVersion(p.ctx, "", "", "", "", p.datasetId, p.edition, p.version).Return(dataset.Version{}, expectedError).Times(1),
		)

		ret, err := api.getDatasetParams(p.ctx, p.request)

		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "failed to get version: "+expectedError.Error())
		So(ret, ShouldBeNil)
	})
}

func TestGetVersionMetadata(t *testing.T) {
	Convey("When GetVersionMetadata is called with an invalid dataset id then an error is returned", t, func() {
		api, ctrl, _, datasetAPIMock := initMocks(t)
		defer ctrl.Finish()

		expectedError := errors.New(uuid.NewString())
		p := getTestParams()

		gomock.InOrder(
			datasetAPIMock.EXPECT().GetVersion(p.ctx, "", "", "", "", p.datasetId, p.edition, p.version).Return(getValidVersionResponse(), nil).Times(1),
			datasetAPIMock.EXPECT().GetVersionMetadata(p.ctx, "", "", "", p.datasetId, p.edition, p.version).Return(dataset.Metadata{}, expectedError).Times(1),
		)

		ret, err := api.getDatasetParams(p.ctx, p.request)

		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "failed to get metadata: "+expectedError.Error())
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
			datasetAPIMock.EXPECT().GetVersion(p.ctx, "", "", "", "", p.datasetId, p.edition, p.version).Return(versionResponse, nil).Times(1),
			datasetAPIMock.EXPECT().GetVersionMetadata(p.ctx, "", "", "", p.datasetId, p.edition, p.version).Return(getValidMetadataResponse(), nil).Times(1),
			datasetAPIMock.EXPECT().GetOptionsInBatches(p.ctx, "", "", "", p.datasetId, p.edition, p.version, versionResponse.Dimensions[0].Name, optionsBatch, optionsWorker).Return(dataset.Options{}, expectedError).Times(1),
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

		geographDimensionsRequest := cantabular.GetGeographyDimensionsRequest{
			Dataset: versionResponse.IsBasedOn.ID,
		}

		gomock.InOrder(
			datasetAPIMock.EXPECT().GetVersion(p.ctx, "", "", "", "", p.datasetId, p.edition, p.version).Return(versionResponse, nil).Times(1),
			datasetAPIMock.EXPECT().GetVersionMetadata(p.ctx, "", "", "", p.datasetId, p.edition, p.version).Return(getValidMetadataResponse(), nil).Times(1),
			datasetAPIMock.EXPECT().GetOptionsInBatches(p.ctx, "", "", "", p.datasetId, p.edition, p.version, versionResponse.Dimensions[0].Name, optionsBatch, optionsWorker).Return(getValidOptionsResponse(), nil).Times(1),
			ctblrMock.EXPECT().GetGeographyDimensions(p.ctx, geographDimensionsRequest).Return(getValidGeoResponse(), nil).Times(1),
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

func TestToGetJsonResponse(t *testing.T) {
	Convey("When TestToGetJsonResponse is called a valid response should be returned", t, func() {
		api, ctrl, _, datasetAPIMock := initMocks(t)
		defer ctrl.Finish()
		p := getTestParams()

		versionResponse := getValidVersionResponse()
		metadataResponse := getValidMetadataResponse()
		optionsResponse := getValidOptionsResponse()

		cantabularResponse := getValidCantabularResponse(versionResponse.Dimensions[0].ID, optionsResponse.Items[0].Label)

		gomock.InOrder(
			datasetAPIMock.EXPECT().GetVersion(p.ctx, "", "", "", "", p.datasetId, p.edition, p.version).Return(versionResponse, nil).Times(1),
			datasetAPIMock.EXPECT().GetVersionMetadata(p.ctx, "", "", "", p.datasetId, p.edition, p.version).Return(metadataResponse, nil).Times(1),
			datasetAPIMock.EXPECT().GetOptionsInBatches(p.ctx, "", "", "", p.datasetId, p.edition, p.version, versionResponse.Dimensions[0].Name, optionsBatch, optionsWorker).Return(optionsResponse, nil).Times(1),
		)

		params, err := api.getDatasetParams(p.ctx, p.request)
		So(err, ShouldBeNil)

		result, err := api.toGetDatasetJsonResponse(params, &cantabularResponse)

		So(err, ShouldBeNil)
		So(result, ShouldNotBeNil)
		So(result.Observations, ShouldResemble, cantabularResponse.Dataset.Table.Values)
		So(result.TotalObservations, ShouldEqual, len(cantabularResponse.Dataset.Table.Values))
		So(len(result.Dimensions), ShouldEqual, 1)
		So(result.Dimensions[0].DimensionName, ShouldEqual, versionResponse.Dimensions[0].ID)
		So(len(result.Dimensions[0].Options), ShouldEqual, 1)
		So(result.Dimensions[0].Options[0].HREF, ShouldEqual, optionsResponse.Items[0].Links.Code.URL)
		So(result.Dimensions[0].Options[0].ID, ShouldEqual, optionsResponse.Items[0].Links.Code.ID)
		So(result.Links.DatasetMetadata.HREF, ShouldEqual, metadataResponse.Version.Links.Self.URL)
		So(result.Links.DatasetMetadata.ID, ShouldEqual, metadataResponse.Version.Links.Self.ID)
		So(result.Links.Self.HREF, ShouldEqual, versionResponse.Links.Dataset.URL)
		So(result.Links.Self.ID, ShouldEqual, versionResponse.Links.Dataset.ID)
		So(result.Links.Version.HREF, ShouldEqual, versionResponse.Links.Self.URL)
		So(result.Links.Version.ID, ShouldEqual, versionResponse.Links.Self.ID)
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
	api.cfg = &config.Config{
		DatasetOptionsBatchSize: optionsBatch,
		DatasetOptionsWorkers:   optionsWorker,
	}

	return api, ctrl, ctblrMock, datasetAPIMock
}

func getTestParams() *testParams {
	datasetId := "datasetId-" + uuid.NewString()
	edition := "edition-" + uuid.NewString()
	version := "version-" + uuid.NewString()

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("dataset_id", datasetId)
	ctx.URLParams.Add("version", version)
	ctx.URLParams.Add("edition", edition)

	request := httptest.NewRequest("GET", "/dataset", nil)
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))

	return &testParams{
		request.Context(),
		request,
		httptest.NewRecorder(),
		datasetId,
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

func getValidMetadataResponse() dataset.Metadata {
	return dataset.Metadata{
		Version: dataset.Version{
			Links: dataset.Links{
				Self: dataset.Link{
					URL: "metadataURL " + uuid.NewString(),
					ID:  "metadataId " + uuid.NewString(),
				},
			},
		},
	}
}

func getValidOptionsResponse() dataset.Options {
	return dataset.Options{
		Items: []dataset.Option{
			{
				Label: "optionLabel " + uuid.NewString(),
				Links: dataset.Links{
					Code: dataset.Link{
						URL: "codeURL " + uuid.NewString(),
					},
				},
			},
		},
	}
}

func getValidGeoResponse() *cantabular.GetGeographyDimensionsResponse {
	return &cantabular.GetGeographyDimensionsResponse{
		Dataset: gql.Dataset{
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
							},
						},
						Variable: cantabular.VariableBase{
							Name: name,
						},
					},
				},
				Values: []int{100},
			},
		},
	}
}
