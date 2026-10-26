package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/api/mock"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/config"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	dphttp "github.com/ONSdigital/dp-net/v3/http"
	"github.com/ONSdigital/dp-net/v3/request"
	dpresponder "github.com/ONSdigital/dp-net/v3/responder"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

func TestCreateCustomFilterPopulationBackendRequest(t *testing.T) {
	received := make(chan *http.Request, 1)
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received <- r.Clone(context.Background())
		http.NotFound(w, r)
	}))
	t.Cleanup(backend.Close)

	backendURL, err := url.Parse(backend.URL)
	if err != nil {
		t.Fatal(err)
	}
	populationClient, err := population.NewClient(backend.URL)
	if err != nil {
		t.Fatal(err)
	}
	cfg := &config.Config{ServiceAuthToken: "test-service-token"}
	router := chi.NewRouter()
	New(context.Background(), cfg, router, nil, dpresponder.New(), nil, nil, nil, populationClient,
		nil, nil, nil, false, backendURL, backendURL, backendURL)
	handler := request.HandlerRequestID(16)(otelhttp.NewHandler(router, "/", otelhttp.WithTracerProvider(trace.NewNoopTracerProvider())))

	const requestID = "caller-request-12345678"
	incoming := httptest.NewRequest(http.MethodPost, "/custom/filters?tag=one&tag=two",
		strings.NewReader(`{"population_type":"test-population"}`))
	incoming.Header.Set("Content-Type", "application/json")
	incoming.Header.Set("Authorization", "Bearer caller-token")
	incoming.Header.Set("Cookie", "session=caller-session")
	incoming.Header.Set("X-Request-Id", requestID)
	incoming.Header.Set("traceparent", "00-0123456789abcdef0123456789abcdef-0123456789abcdef-01")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, incoming)
	if response.Code != http.StatusNotFound {
		t.Fatalf("response status: got %d, want %d", response.Code, http.StatusNotFound)
	}

	var got *http.Request
	select {
	case got = <-received:
	default:
		t.Fatal("population backend received no request")
	}
	if got.Method != http.MethodGet || got.URL.Path != "/population-types/test-population/metadata" || got.URL.RawQuery != "" {
		t.Errorf("backend method and URL: got %s %s, want GET /population-types/test-population/metadata", got.Method, got.URL.RequestURI())
	}
	if got.Header.Get("Authorization") != "Bearer test-service-token" {
		t.Errorf("backend authorization: got %q, want service token", got.Header.Get("Authorization"))
	}
	if id := got.Header.Get("X-Request-Id"); !strings.HasPrefix(id, requestID+",") || len(strings.TrimPrefix(id, requestID+",")) != len(requestID)/2 {
		t.Errorf("backend request ID: got %q, want %q followed by an appended ID", id, requestID+",")
	}
	for _, name := range []string{"traceparent", "Cookie", "Content-Type"} {
		if value := got.Header.Get(name); value != "" {
			t.Errorf("backend %s: got %q, want no copied header", name, value)
		}
	}
	if got.ContentLength != 0 {
		t.Errorf("backend body length: got %d, want 0", got.ContentLength)
	}
}

func TestCreateFilterCantabularBackendRequest(t *testing.T) {
	type backendRequest struct {
		method, path, query string
		header              http.Header
		body                []byte
	}
	received := make(chan backendRequest, 1)
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("read Cantabular request: %v", err)
			http.Error(w, "unable to read request", http.StatusBadRequest)
			return
		}
		received <- backendRequest{r.Method, r.URL.Path, r.URL.RawQuery, r.Header.Clone(), body}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"data":{"dataset":{"table":{"dimensions":[]}}}}`)
	}))
	t.Cleanup(backend.Close)

	backendURL, err := url.Parse(backend.URL)
	if err != nil {
		t.Fatal(err)
	}
	ctrl := gomock.NewController(t)
	datasets := mock.NewMockdatasetAPIClient(ctrl)
	cfg := &config.Config{ServiceAuthToken: "test-service-token"}
	version := dataset.Version{
		State:      "published",
		IsBasedOn:  &dataset.IsBasedOn{Type: "cantabular_flexible_table"},
		Dimensions: []dataset.VersionDimension{{Name: "age", ID: "age"}},
	}
	datasets.EXPECT().GetVersion(gomock.Any(), "", cfg.ServiceAuthToken, "", "", "test-dataset", "2021", "1").Return(version, nil)
	datasets.EXPECT().GetDatasetCurrentAndNext(gomock.Any(), "", cfg.ServiceAuthToken, "", "test-dataset").Return(dataset.Dataset{}, errors.New("stop after Cantabular request"))
	router := chi.NewRouter()
	New(context.Background(), cfg, router, nil, dpresponder.New(), nil, nil, datasets, nil,
		cantabular.NewClient(cantabular.Config{ExtApiHost: backend.URL}, dphttp.NewClient(), nil), nil, nil, false,
		backendURL, backendURL, backendURL)
	handler := request.HandlerRequestID(16)(otelhttp.NewHandler(router, "/", otelhttp.WithTracerProvider(trace.NewNoopTracerProvider())))

	const requestID = "caller-request-12345678"
	incoming := httptest.NewRequest(http.MethodPost, "/filters?tag=one&tag=two", strings.NewReader(
		`{"dataset":{"id":"test-dataset","edition":"2021","version":1},"population_type":"test-population","dimensions":[{"name":"age","options":["age-code"]}]}`))
	incoming.Header.Set("Content-Type", "application/vnd.api+json")
	incoming.Header.Set("Authorization", "Bearer caller-token")
	incoming.Header.Set("Cookie", "session=caller-session")
	incoming.Header.Set("X-Request-Id", requestID)
	incoming.Header.Set("traceparent", "00-0123456789abcdef0123456789abcdef-0123456789abcdef-01")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, incoming)
	if response.Code != http.StatusInternalServerError {
		t.Fatalf("response status: got %d, want %d", response.Code, http.StatusInternalServerError)
	}

	var got backendRequest
	select {
	case got = <-received:
	default:
		t.Fatal("Cantabular backend received no request")
	}
	if got.method != http.MethodPost || got.path != "/graphql" || got.query != "" {
		t.Errorf("backend method and URL: got %s %s?%s, want POST /graphql", got.method, got.path, got.query)
	}
	if got.header.Get("Content-Type") != "application/json" {
		t.Errorf("backend content type: got %q, want application/json", got.header.Get("Content-Type"))
	}
	if id := got.header.Get("X-Request-Id"); !strings.HasPrefix(id, requestID+",") || len(strings.TrimPrefix(id, requestID+",")) != len(requestID)/2 {
		t.Errorf("backend request ID: got %q, want %q followed by an appended ID", id, requestID+",")
	}
	for _, name := range []string{"Authorization", "Cookie", "traceparent"} {
		if value := got.header.Get(name); value != "" {
			t.Errorf("backend %s: got %q, want no copied header", name, value)
		}
	}

	var payload struct {
		Query     string `json:"query"`
		Variables struct {
			Dataset   string   `json:"dataset"`
			Variables []string `json:"variables"`
			Filters   []struct {
				Variable string   `json:"variable"`
				Codes    []string `json:"codes"`
			} `json:"filters"`
		} `json:"variables"`
	}
	if err := json.Unmarshal(got.body, &payload); err != nil {
		t.Fatalf("decode Cantabular request: %v", err)
	}
	if !strings.Contains(payload.Query, "table(variables: $variables, filters: $filters)") ||
		payload.Variables.Dataset != "test-population" || len(payload.Variables.Variables) != 1 || payload.Variables.Variables[0] != "age" ||
		len(payload.Variables.Filters) != 1 || payload.Variables.Filters[0].Variable != "age" ||
		len(payload.Variables.Filters[0].Codes) != 1 || payload.Variables.Filters[0].Codes[0] != "age-code" {
		t.Errorf("backend GraphQL query and variables: got %+v, want selected age-code for test-population", payload)
	}
}

func Test_getFilterVariable(t *testing.T) {
	Convey("Given an array of dimension IDs exists", t, func() {
		Convey("When I have a dimension with no FilterByParent", func() {
			d := model.Dimension{
				Name:           "city",
				FilterByParent: "",
			}

			Convey("Then the 'city' should be returned as the filter variable", func() {
				got := getFilterVariable(d)
				So(got, ShouldResemble, "city")
			})
		})

		Convey("When I have a dimension with a FilterByParent", func() {
			d := model.Dimension{
				Name:           "city",
				FilterByParent: "region",
			}

			Convey("Then 'region' should be returned as the filter variable", func() {
				got := getFilterVariable(d)
				So(got, ShouldResemble, "region")
			})
		})
	})
}

func TestValidateDimensions(t *testing.T) {
	var api API

	Convey("Given a set of filter dimensions which match existing dimensions found in Version doc", t, func() {
		existingDims := []dataset.VersionDimension{
			{
				Name: "foo",
				ID:   "foo01",
			},
			{
				Name: "bar",
				ID:   "bar01",
			},
		}

		filterDims := []model.Dimension{
			{
				Name: "foo",
				Options: []string{
					"foo_1",
					"foo_2",
				},
			},
			{
				Name: "bar",
				Options: []string{
					"bar_1",
					"bar_2",
				},
			},
		}

		Convey("When validateDimensions is called", func() {
			err := api.validateDimensionsFromVersion(filterDims, existingDims)
			So(err, ShouldBeNil)
		})
	})

	Convey("Given a set of filter dimensions which include dimensions not found in Version doc", t, func() {
		existingDims := []dataset.VersionDimension{
			{
				Name: "foo",
				ID:   "foo01",
			},
			{
				Name: "bar",
				ID:   "bar01",
			},
		}

		filterDims := []model.Dimension{
			{
				Name: "foo",
				Options: []string{
					"foo_1",
					"foo_2",
				},
			},
			{
				Name: "bar",
				Options: []string{
					"bar_1",
					"bar_2",
				},
			},
			{
				Name: "alice",
				Options: []string{
					"alice_1",
					"alice_2",
				},
			},
		}

		Convey("When validateDimensions is called", func() {
			err := api.validateDimensionsFromVersion(filterDims, existingDims)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Given a set of filter dimensions are duplicate", t, func() {
		existingDims := []dataset.VersionDimension{
			{
				Name: "foo",
				ID:   "foo01",
			},
			{
				Name: "foo",
				ID:   "bar01",
			},
		}

		filterDims := []model.Dimension{
			{
				Name: "foo",
				Options: []string{
					"foo_1",
					"foo_2",
				},
			},
			{
				Name: "foo",
				Options: []string{
					"bar_1",
					"bar_2",
				},
			},
		}

		Convey("When validateDimensions is called", func() {
			err := api.validateDimensionsFromVersion(filterDims, existingDims)
			So(err, ShouldNotBeNil)
		})
	})
}

func TestValidateDimensionOptions(t *testing.T) {
	ctx := context.Background()
	var api API

	Convey("Given a Cantabular Client which cannot find provided dimension options", t, func() {
		api.ctblr = &mock.CantabularClient{
			ErrStatus:    http.StatusInternalServerError,
			OptionsHappy: false,
		}

		Convey("When validateDimenOptions is called", func() {
			req := createFilterRequest{
				PopulationType: "Example",
				Dimensions: []model.Dimension{
					{
						Name: "foo",
						Options: []string{
							"foo_1",
							"foo_2",
						},
					},
					{
						Name: "bar",
						Options: []string{
							"bar_1",
							"bar_2",
						},
					},
				},
			}

			err := api.validateDimensionOptions(ctx, req.Dimensions, req.PopulationType)
			So(err, ShouldNotBeNil)
		})

		Convey("When validateDimenOptions is called but there are no options selected", func() {
			req := createFilterRequest{
				PopulationType: "Example",
				Dimensions: []model.Dimension{
					{
						Name: "foo",
					},
					{
						Name: "bar",
					},
				},
			}

			err := api.validateDimensionOptions(ctx, req.Dimensions, req.PopulationType)
			So(err, ShouldBeNil)
		})
	})

	Convey("Given a Cantabular Client which can find provided dimension options", t, func() {
		api.ctblr = &mock.CantabularClient{
			ErrStatus:    http.StatusInternalServerError,
			OptionsHappy: true,
		}

		Convey("When validateDimenOptions is called", func() {
			req := createFilterRequest{
				PopulationType: "Example",
				Dimensions: []model.Dimension{
					{
						Name: "foo",
						Options: []string{
							"foo_1",
							"foo_2",
						},
					},
					{
						Name: "bar",
						Options: []string{
							"bar_1",
							"bar_2",
						},
					},
				},
			}

			err := api.validateDimensionOptions(ctx, req.Dimensions, req.PopulationType)
			So(err, ShouldBeNil)
		})
	})

	Convey("Given a Cantabular Client which is expecting the value present in FilterByParent variable", t, func() {
		parentDim := "ParentFoo"

		api.ctblr = &mock.CantabularClient{
			ErrStatus:               http.StatusInternalServerError,
			OptionsHappy:            true,
			ExpectedFilterDimension: parentDim,
		}

		Convey("When validateDimenOptions is called with value as the parent of a variable", func() {
			req := createFilterRequest{
				PopulationType: "Example",
				Dimensions: []model.Dimension{
					{
						Name: "foo",
						Options: []string{
							"foo_1",
							"foo_2",
						},
						FilterByParent: parentDim,
					},
					{
						Name: "bar",
						Options: []string{
							"bar_1",
							"bar_2",
						},
					},
				},
			}

			err := api.validateDimensionOptions(ctx, req.Dimensions, req.PopulationType)
			So(err, ShouldBeNil)
		})

		Convey("When validateDimenOptions is called with the correct value but not as the parent", func() {
			req := createFilterRequest{
				PopulationType: "Example",
				Dimensions: []model.Dimension{
					{
						Name: parentDim,
						Options: []string{
							"foo_1",
							"foo_2",
						},
						FilterByParent: "incorrect",
					},
					{
						Name: "bar",
						Options: []string{
							"bar_1",
							"bar_2",
						},
					},
				},
			}

			err := api.validateDimensionOptions(ctx, req.Dimensions, req.PopulationType)
			So(err, ShouldNotBeNil)
		})

		Convey("When validateDimenOptions is called without the correct value", func() {
			req := createFilterRequest{
				PopulationType: "Example",
				Dimensions: []model.Dimension{
					{
						Name: "foo",
						Options: []string{
							"foo_1",
							"foo_2",
						},
						FilterByParent: "incorrect",
					},
					{
						Name: "bar",
						Options: []string{
							"bar_1",
							"bar_2",
						},
					},
				},
			}

			err := api.validateDimensionOptions(ctx, req.Dimensions, req.PopulationType)
			So(err, ShouldNotBeNil)
		})
	})
}
