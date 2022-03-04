package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/api"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/config"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/service"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/service/mock"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/dp-kafka/v3/kafkatest"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func newMockService(t *testing.T, store service.Datastore) (*service.Service, error) {
	cfg, err := config.Get()
	if err != nil {
		return nil, err
	}

	service.GetHealthCheck = func(_ *config.Config, _, _, _ string) (service.HealthChecker, error) {
		return &mock.HealthCheckerMock{
			StartFunc: func(contextMoqParam context.Context) {},
			StopFunc:  func() {},
			AddAndGetCheckFunc: func(name string, checker healthcheck.Checker) (*healthcheck.Check, error) {
				return &healthcheck.Check{}, nil
			},
		}, nil
	}
	service.GetKafkaProducer = func(_ context.Context, _ *config.Config) (kafka.IProducer, error) {
		return &kafkatest.IProducerMock{
			CloseFunc:     func(ctx context.Context) error { return nil },
			LogErrorsFunc: func(ctx context.Context) {},
		}, nil
	}
	service.GetMongoDB = func(_ context.Context, _ *config.Config, _ service.Generator) (service.Datastore, error) {
		return store, nil
	}

	s := service.New()
	if err := s.Init(context.Background(), cfg, "BuildTime", "GitCommit", "Version"); err != nil {
		return nil, err
	}
	return s, nil
}

func TestCreateFilterOutputs(t *testing.T) {
	datastoreMock := &mock.DatastoreMock{
		CreateFilterOutputFunc: func(ctx context.Context, filterOutput *model.FilterOutput) error {
			if filterOutput.ID == uuid.Nil {
				id, err := uuid.NewUUID()
				if err != nil {
					return errors.Wrap(err, "failed to generate uuid")
				}
				filterOutput.ID = id
			}
			return nil
		},
	}

	Convey("Run test cases for API /filter-outputs", t, func() {
		s, err := newMockService(t, datastoreMock)
		So(err, ShouldBeNil)
		testCases := []struct {
			Request *api.CreateFilterOutputRequest
			Error   string
		}{
			{
				Request: &api.CreateFilterOutputRequest{
					State: "created",
					Downloads: &model.FilterOutput{
						CSV: &model.FileInfo{
							HREF:    "href",
							Size:    "size",
							Public:  "public",
							Private: "private",
							Skipped: false,
						},
						CSVW: &model.FileInfo{
							HREF:    "href",
							Size:    "size",
							Public:  "public",
							Private: "private",
							Skipped: false,
						},
						TXT: &model.FileInfo{
							HREF:    "href",
							Size:    "size",
							Public:  "public",
							Private: "private",
							Skipped: false,
						},
						XLS: &model.FileInfo{
							HREF:    "href",
							Size:    "size",
							Public:  "public",
							Private: "private",
							Skipped: false,
						},
					},
				},
			},
			{
				Request: &api.CreateFilterOutputRequest{
					State: "",
					Downloads: &model.FilterOutput{
						CSV: &model.FileInfo{
							HREF:    "href",
							Size:    "size",
							Public:  "public",
							Private: "private",
							Skipped: false,
						},
						CSVW: &model.FileInfo{
							HREF:    "href",
							Size:    "size",
							Public:  "public",
							Private: "private",
							Skipped: false,
						},
						TXT: &model.FileInfo{
							HREF:    "href",
							Size:    "size",
							Public:  "public",
							Private: "private",
							Skipped: false,
						},
						XLS: &model.FileInfo{
							HREF:    "href",
							Size:    "size",
							Public:  "public",
							Private: "private",
							Skipped: false,
						},
					},
				},
				Error: "CreateFilterOutputRequest.State",
			},
		}

		for i, tc := range testCases {
			tc := tc

			name := fmt.Sprintf("CreateFilterOutputsTestCase-%d", i)
			t.Run(name, func(t *testing.T) {
				Convey(name, t, func() {
					data, err := json.Marshal(tc.Request)
					So(err, ShouldBeNil)

					r := httptest.NewRequest("POST", "/filter-outputs", bytes.NewBuffer(data))
					w := httptest.NewRecorder()
					s.Api.CreateFilterOutput(w, r)

					code := w.Result().StatusCode
					if tc.Error == "" {
						So(code, ShouldEqual, 200)
					} else {
						So(code, ShouldEqual, 400)
					}

					data, err = ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)

					var res api.CreateFilterOutputResponse
					err = json.Unmarshal(data, &res)
					if tc.Error == "" {
						So(err, ShouldBeNil)
						So(res.FilterOutput.ID, ShouldNotEqual, uuid.Nil)
					} else {
						So(string(data), ShouldContainSubstring, tc.Error)
					}
				})
			})
		}
	})

}
