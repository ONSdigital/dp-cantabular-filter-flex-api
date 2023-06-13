package api

import (
	"context"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular/gql"
	"github.com/ONSdigital/dp-api-clients-go/v2/cantabularmetadata"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/api/mock"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

type mockRequestResponses struct {
	datasetName            string
	dimension              *model.Dimension
	categorisationRequest  cantabular.GetCategorisationsRequest
	categorisationResponse *cantabular.GetCategorisationsResponse
	classificationRequest  cantabularmetadata.GetDefaultClassificationRequest
	classificationResponse *cantabularmetadata.GetDefaultClassificationResponse
}

func TestRetrieveDefaultCategorisation(t *testing.T) {
	Convey("When Cantabular returns an error", t, func() {
		api, _, ctblrMock, _ := initCategoryMocks(t)

		ctx := context.Background()
		expectedError := errors.New(uuid.NewString())

		mockRequests := getMocksValid()

		gomock.InOrder(
			ctblrMock.EXPECT().GetCategorisations(ctx, mockRequests.categorisationRequest).Return(nil, expectedError).Times(1),
		)

		finalDimension, finalLabel, finalCategorisation, err := api.RetrieveDefaultCategorisation(mockRequests.dimension, mockRequests.datasetName)

		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "failed to check default categorisation: "+expectedError.Error())
		So(finalDimension, ShouldEqual, "")
		So(finalLabel, ShouldEqual, "")
		So(finalCategorisation, ShouldEqual, "")
	})

	Convey("When Cantabular returns no categorisations", t, func() {
		api, _, ctblrMock, _ := initCategoryMocks(t)

		ctx := context.Background()
		mockRequests := getMocksValid()
		mockRequests.categorisationResponse = getCategorisationResponseWithNoCategorisations()

		gomock.InOrder(
			ctblrMock.EXPECT().GetCategorisations(ctx, mockRequests.categorisationRequest).Return(mockRequests.categorisationResponse, nil).Times(1),
		)

		finalDimension, finalLabel, finalCategorisation, err := api.RetrieveDefaultCategorisation(mockRequests.dimension, mockRequests.datasetName)

		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "no categorisations received for variable")
		So(finalDimension, ShouldEqual, "")
		So(finalLabel, ShouldEqual, "")
		So(finalCategorisation, ShouldEqual, "")
	})

	Convey("When Cantabular metadata service returns an error", t, func() {
		api, _, ctblrMock, metadataMock := initCategoryMocks(t)

		ctx := context.Background()
		expectedError := errors.New(uuid.NewString())

		mockRequests := getMocksValid()

		gomock.InOrder(
			ctblrMock.EXPECT().GetCategorisations(ctx, mockRequests.categorisationRequest).Return(mockRequests.categorisationResponse, nil).Times(1),
			metadataMock.EXPECT().GetDefaultClassification(ctx, mockRequests.classificationRequest).Return(nil, expectedError).Times(1),
		)

		finalDimension, finalLabel, finalCategorisation, err := api.RetrieveDefaultCategorisation(mockRequests.dimension, mockRequests.datasetName)

		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "failed to check default categorisation: "+expectedError.Error())
		So(finalDimension, ShouldEqual, "")
		So(finalLabel, ShouldEqual, "")
		So(finalCategorisation, ShouldEqual, "")
	})

	Convey("When Cantabular returns a categorisation with Mapped Source edges", t, func() {
		api, _, ctblrMock, metadataMock := initCategoryMocks(t)

		ctx := context.Background()

		mockRequests := getMocksValid()

		gomock.InOrder(
			ctblrMock.EXPECT().GetCategorisations(ctx, mockRequests.categorisationRequest).Return(mockRequests.categorisationResponse, nil).Times(1),
			metadataMock.EXPECT().GetDefaultClassification(ctx, mockRequests.classificationRequest).Return(mockRequests.classificationResponse, nil).Times(1),
		)

		finalDimension, finalLabel, finalCategorisation, err := api.RetrieveDefaultCategorisation(mockRequests.dimension, mockRequests.datasetName)

		So(err, ShouldBeNil)
		So(finalDimension, ShouldEqual, "test-name-3")
		So(finalLabel, ShouldEqual, "test label 3")
		So(finalCategorisation, ShouldEqual, "test-name-3")
	})

	Convey("When Cantabular returns a categorisation without Mapped Source edges", t, func() {
		api, _, ctblrMock, metadataMock := initCategoryMocks(t)

		ctx := context.Background()

		mockRequests := getMocksValid()
		mockRequests.categorisationResponse = &cantabular.GetCategorisationsResponse{
			Dataset: gql.Dataset{
				Variables: gql.Variables{
					Edges: []gql.Edge{
						// edge
						{
							Node: gql.Node{
								// mapFrom
								IsSourceOf: gql.Variables{
									Edges: []gql.Edge{
										{
											Node: gql.Node{
												Name:  "test-name-1",
												Label: "test label 1",
											},
										},
										{
											Node: gql.Node{
												Name:  "test-name-2",
												Label: "test label 2",
											},
										},
										{
											Node: gql.Node{
												Name:  "test-name-3",
												Label: "test label 3",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		gomock.InOrder(
			ctblrMock.EXPECT().GetCategorisations(ctx, mockRequests.categorisationRequest).Return(mockRequests.categorisationResponse, nil).Times(1),
			metadataMock.EXPECT().GetDefaultClassification(ctx, mockRequests.classificationRequest).Return(mockRequests.classificationResponse, nil).Times(1),
		)

		finalDimension, finalLabel, finalCategorisation, err := api.RetrieveDefaultCategorisation(mockRequests.dimension, mockRequests.datasetName)

		So(err, ShouldBeNil)
		So(finalDimension, ShouldEqual, "test-name-3")
		So(finalLabel, ShouldEqual, "test label 3")
		So(finalCategorisation, ShouldEqual, "test-name-3")
	})

	Convey("When more than 1 categorisation is returned", t, func() {
		api, _, ctblrMock, metadataMock := initCategoryMocks(t)

		ctx := context.Background()

		mockRequests := getMocksValid()
		mockRequests.classificationResponse = &cantabularmetadata.GetDefaultClassificationResponse{
			Variables: []string{"test-name-3", "test-name-2"},
		}

		gomock.InOrder(
			ctblrMock.EXPECT().GetCategorisations(ctx, mockRequests.categorisationRequest).Return(mockRequests.categorisationResponse, nil).Times(1),
			metadataMock.EXPECT().GetDefaultClassification(ctx, mockRequests.classificationRequest).Return(mockRequests.classificationResponse, nil).Times(1),
		)

		finalDimension, finalLabel, finalCategorisation, err := api.RetrieveDefaultCategorisation(mockRequests.dimension, mockRequests.datasetName)

		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "more than 1 categorisation returned")
		So(finalDimension, ShouldEqual, "")
		So(finalLabel, ShouldEqual, "")
		So(finalCategorisation, ShouldEqual, "")
	})

	Convey("When no categorisation is returned", t, func() {
		api, _, ctblrMock, metadataMock := initCategoryMocks(t)

		ctx := context.Background()

		mockRequests := getMocksValid()
		mockRequests.classificationResponse = &cantabularmetadata.GetDefaultClassificationResponse{
			Variables: []string{},
		}

		gomock.InOrder(
			ctblrMock.EXPECT().GetCategorisations(ctx, mockRequests.categorisationRequest).Return(mockRequests.categorisationResponse, nil).Times(1),
			metadataMock.EXPECT().GetDefaultClassification(ctx, mockRequests.classificationRequest).Return(mockRequests.classificationResponse, nil).Times(1),
		)

		finalDimension, finalLabel, finalCategorisation, err := api.RetrieveDefaultCategorisation(mockRequests.dimension, mockRequests.datasetName)

		So(err, ShouldBeNil)
		So(finalDimension, ShouldEqual, mockRequests.dimension.Name)
		So(finalLabel, ShouldEqual, mockRequests.dimension.Label)
		So(finalCategorisation, ShouldEqual, "")
	})
}

func getCategorisationResponseWithNoCategorisations() *cantabular.GetCategorisationsResponse {
	return &cantabular.GetCategorisationsResponse{}
}

func getMocksValid() mockRequestResponses {
	datasetName := "Test Dataset name"

	dimension := &model.Dimension{
		Name:  "dimension-name",
		Label: "Dimension Label",
	}

	categorisationRequest := cantabular.GetCategorisationsRequest{
		Dataset:  datasetName,
		Variable: dimension.Name,
	}

	categorisationResponse := &cantabular.GetCategorisationsResponse{
		Dataset: gql.Dataset{
			Variables: gql.Variables{
				Edges: []gql.Edge{
					// edge
					{
						Node: gql.Node{
							// mapFrom
							MapFrom: []gql.Variables{
								{
									Edges: []gql.Edge{
										// mappedSource
										{
											Node: gql.Node{
												IsSourceOf: gql.Variables{
													// mappedSourceEdge
													Edges: []gql.Edge{
														{
															Node: gql.Node{
																Name:  "test-name-1",
																Label: "test label 1",
															},
														},
														{
															Node: gql.Node{
																Name:  "test-name-2",
																Label: "test label 2",
															},
														},
														{
															Node: gql.Node{
																Name:  "test-name-3",
																Label: "test label 3",
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	classificationRequest := cantabularmetadata.GetDefaultClassificationRequest{
		Dataset:   datasetName,
		Variables: []string{"test-name-1", "test-name-2", "test-name-3"},
	}

	classificationResponse := &cantabularmetadata.GetDefaultClassificationResponse{
		Variables: []string{"test-name-3"},
	}

	return mockRequestResponses{
		datasetName,
		dimension,
		categorisationRequest,
		categorisationResponse,
		classificationRequest,
		classificationResponse,
	}
}

func initCategoryMocks(t *testing.T) (API, *gomock.Controller, *mock.MockcantabularClient, *mock.MockmetadataAPIClient) {
	var api API

	ctrl := gomock.NewController(t)

	ctblrMock := mock.NewMockcantabularClient(ctrl)
	api.ctblr = ctblrMock
	metadataAPIMock := mock.NewMockmetadataAPIClient(ctrl)
	api.metadata = metadataAPIMock

	return api, ctrl, ctblrMock, metadataAPIMock
}
