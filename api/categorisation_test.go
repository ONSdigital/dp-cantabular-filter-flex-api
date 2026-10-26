package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"slices"
	"strings"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular/gql"
	"github.com/ONSdigital/dp-api-clients-go/v2/cantabularmetadata"
	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/api/mock"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/config"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	dphttp "github.com/ONSdigital/dp-net/v3/http"
	"github.com/ONSdigital/dp-net/v3/request"
	dpresponder "github.com/ONSdigital/dp-net/v3/responder"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

type mockRequestResponses struct {
	datasetName            string
	dimension              *model.Dimension
	categorisationRequest  cantabular.GetCategorisationsRequest
	categorisationResponse *cantabular.GetCategorisationsResponse
	classificationRequest  cantabularmetadata.GetDefaultClassificationRequest
	classificationResponse *cantabularmetadata.GetDefaultClassificationResponse
}

func TestMultivariateMetadataBackendRequestWithTracing(t *testing.T) {
	type backendRequest struct {
		method, path, query string
		header              http.Header
		body                []byte
	}
	received := make(chan backendRequest, 1)
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("read metadata request: %v", err)
			http.Error(w, "unable to read request", http.StatusBadRequest)
			return
		}
		received <- backendRequest{r.Method, r.URL.Path, r.URL.RawQuery, r.Header.Clone(), body}
		http.NotFound(w, r)
	}))
	t.Cleanup(backend.Close)

	backendURL, err := url.Parse(backend.URL)
	if err != nil {
		t.Fatal(err)
	}
	cfg := &config.Config{ServiceAuthToken: "test-service-token"}
	fixture := getMocksValid()
	ctrl := gomock.NewController(t)
	store := mock.NewMockdatastore(ctrl)
	datasets := mock.NewMockdatasetAPIClient(ctrl)
	ctblr := mock.NewMockcantabularClient(ctrl)
	store.EXPECT().GetFilter(gomock.Any(), "filter-id").Return(&model.Filter{
		Dataset:        model.Dataset{ID: "test-dataset", Edition: "2021", Version: 1},
		PopulationType: fixture.datasetName,
	}, nil)
	datasets.EXPECT().GetVersion(gomock.Any(), "", cfg.ServiceAuthToken, "", "", "test-dataset", "2021", "1").Return(dataset.Version{
		State:     "published",
		IsBasedOn: &dataset.IsBasedOn{Type: "cantabular_multivariate_table"},
	}, nil)
	ctblr.EXPECT().GetDimensionsByName(context.Background(), cantabular.GetDimensionsByNameRequest{
		Dataset: fixture.datasetName, DimensionNames: []string{fixture.dimension.Name},
	}).Return(&cantabular.GetDimensionsResponse{Dataset: gql.Dataset{Variables: gql.Variables{
		Edges: []gql.Edge{{Node: gql.Node{Label: fixture.dimension.Label}}},
	}}}, nil)
	ctblr.EXPECT().GetCategorisations(context.Background(), fixture.categorisationRequest).Return(fixture.categorisationResponse, nil)

	router := chi.NewRouter()
	New(context.Background(), cfg, router, nil, dpresponder.New(), nil, store, datasets, nil, ctblr,
		cantabularmetadata.NewClient(cantabularmetadata.Config{Host: backend.URL}, dphttp.NewClient()), nil, false,
		backendURL, backendURL, backendURL)
	spans := tracetest.NewSpanRecorder()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(spans))
	t.Cleanup(func() {
		if err := provider.Shutdown(context.Background()); err != nil {
			t.Errorf("shutdown tracer provider: %v", err)
		}
	})
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	handler := request.HandlerRequestID(16)(otelhttp.NewHandler(router, "/",
		otelhttp.WithTracerProvider(provider), otelhttp.WithPropagators(propagator)))

	const traceparent = "00-0123456789abcdef0123456789abcdef-0123456789abcdef-01"
	const requestID = "caller-request-12345678"
	incoming := httptest.NewRequest(http.MethodPost, "/filters/filter-id/dimensions",
		strings.NewReader(`{"name":"dimension-name","is_area_type":false}`))
	incoming.Header.Set("Content-Type", "application/json")
	incoming.Header.Set("X-Request-Id", requestID)
	incoming.Header.Set("traceparent", traceparent)
	incoming.Header.Set("Authorization", "Bearer caller-token")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, incoming)
	if response.Code != http.StatusInternalServerError {
		t.Fatalf("response status: got %d, want %d", response.Code, http.StatusInternalServerError)
	}

	serverSpans := spans.Ended()
	parent := trace.SpanContextFromContext(propagator.Extract(context.Background(), propagation.HeaderCarrier(incoming.Header)))
	if len(serverSpans) != 1 || !parent.IsValid() {
		t.Fatalf("inbound tracing: got %d spans and a valid parent: %t", len(serverSpans), parent.IsValid())
	}
	if serverSpans[0].SpanContext().TraceID() != parent.TraceID() || serverSpans[0].SpanContext().SpanID() == parent.SpanID() {
		t.Errorf("inbound span did not continue the trace: got %s/%s, parent %s/%s",
			serverSpans[0].SpanContext().TraceID(), serverSpans[0].SpanContext().SpanID(), parent.TraceID(), parent.SpanID())
	}

	var got backendRequest
	select {
	case got = <-received:
	default:
		t.Fatal("metadata backend received no request")
	}
	if got.method != http.MethodPost || got.path != "/graphql" || got.query != "" {
		t.Errorf("backend method and URL: got %s %s?%s, want POST /graphql", got.method, got.path, got.query)
	}
	if got.header.Get("Content-Type") != "application/json" {
		t.Errorf("backend content type: got %q, want application/json", got.header.Get("Content-Type"))
	}
	if id := got.header.Get("X-Request-Id"); len(id) != 20 || strings.HasPrefix(id, requestID+",") {
		t.Errorf("backend request ID: got %q, want a new 20-character ID", id)
	}
	for _, name := range []string{"traceparent", "Authorization", "Cookie"} {
		if value := got.header.Get(name); value != "" {
			t.Errorf("backend %s: got %q, want no copied header", name, value)
		}
	}

	var payload struct {
		Query     string `json:"query"`
		Variables struct {
			Dataset   string   `json:"dataset"`
			Variables []string `json:"variables"`
		} `json:"variables"`
	}
	if err := json.Unmarshal(got.body, &payload); err != nil {
		t.Fatalf("decode metadata request: %v", err)
	}
	if !strings.Contains(payload.Query, "Default_Classification_Flag") ||
		payload.Variables.Dataset != fixture.datasetName ||
		!slices.Equal(payload.Variables.Variables, fixture.classificationRequest.Variables) {
		t.Errorf("backend classification query: got %+v, want %q and %v", payload, fixture.datasetName, fixture.classificationRequest.Variables)
	}
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
