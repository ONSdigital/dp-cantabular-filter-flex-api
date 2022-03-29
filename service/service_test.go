package service_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"testing"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/config"
	service "github.com/ONSdigital/dp-cantabular-filter-flex-api/service"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/service/mock"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-kafka/v3/kafkatest"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	ctx           = context.Background()
	testBuildTime = "BuildTime"
	testGitCommit = "GitCommit"
	testVersion   = "Version"
	testChecks    = map[string]*healthcheck.Check{ // TODO - adjust when app is more fully implemented
		"Cantabular client": {},
	}
)

var (
	errHealthcheck = fmt.Errorf("healthCheck error")
	errServer      = fmt.Errorf("HTTP Server error")
	errAddCheck    = fmt.Errorf("healthcheck add check error")
)

func TestInit(t *testing.T) {
	Convey("Having a set of mocked dependencies", t, func() {

		cfg, err := config.Get()
		So(err, ShouldBeNil)

		// Mock clients
		subscribedTo := []*healthcheck.Check{}
		serverMock := &mock.HTTPServerMock{}
		datastoreMock := &mock.DatastoreMock{}
		generatorMock := &mock.GeneratorMock{}
		responderMock := &mock.ResponderMock{}

		hcMock := &mock.HealthCheckerMock{
			AddAndGetCheckFunc: func(name string, checker healthcheck.Checker) (*healthcheck.Check, error) {
				return testChecks[name], nil
			},
			SubscribeFunc: func(s healthcheck.Subscriber, checks ...*healthcheck.Check) {
				subscribedTo = append(subscribedTo, checks...)
			},
		}

		// Initialiser functions
		service.GetHealthCheck = func(_ *config.Config, _, _, _ string) (service.HealthChecker, error) {
			return hcMock, nil
		}

		service.GetHTTPServer = func(_ string, _ http.Handler) service.HTTPServer {
			return serverMock
		}

		service.GetMongoDB = func(_ context.Context, _ *config.Config, _ service.Generator) (service.Datastore, error) {
			return datastoreMock, nil
		}

		service.GetGenerator = func() service.Generator {
			return generatorMock
		}

		service.GetResponder = func() service.Responder {
			return responderMock
		}

		// Service
		svc := service.New()

		Convey("Given that initialising healthcheck returns an error", func() {
			service.GetHealthCheck = func(cfg *config.Config, buildTime, gitCommit, version string) (service.HealthChecker, error) {
				return nil, errHealthcheck
			}

			Convey("Then service Init fails with the same error and no further initialisations are attempted", func() {
				err := svc.Init(ctx, cfg, testBuildTime, testGitCommit, testVersion)
				So(errors.Unwrap(err), ShouldResemble, errHealthcheck)
				So(svc.Cfg, ShouldResemble, cfg)

				Convey("And no checkers are registered ", func() {
					So(hcMock.AddAndGetCheckCalls(), ShouldHaveLength, 0)
				})
			})
		})

		Convey("Given that Checkers cannot be registered", func() {
			hcMock.AddAndGetCheckFunc = func(name string, checker healthcheck.Checker) (*healthcheck.Check, error) { return nil, errAddCheck }

			Convey("Then service Init fails with the expected error", func() {
				err := svc.Init(ctx, cfg, testBuildTime, testGitCommit, testVersion)
				So(err, ShouldNotBeNil)
				So(errors.Is(err, errAddCheck), ShouldBeTrue)
				So(svc.Cfg, ShouldResemble, cfg)

				Convey("And other checkers don't try to register", func() {
					So(hcMock.AddAndGetCheckCalls(), ShouldHaveLength, 1)
				})
			})
		})

		Convey("Given that all dependencies are successfully initialised", func() {

			Convey("Then service Init succeeds, all dependencies are initialised", func() {
				err := svc.Init(ctx, cfg, testBuildTime, testGitCommit, testVersion)
				So(err, ShouldBeNil)
				So(svc.Cfg, ShouldResemble, cfg)
				So(svc.Server, ShouldEqual, serverMock)
				So(svc.HealthCheck, ShouldResemble, hcMock)

				Convey("Then all checks are registered", func() {
					So(hcMock.AddAndGetCheckCalls(), ShouldHaveLength, 5)
					So(hcMock.AddAndGetCheckCalls()[0].Name, ShouldResemble, "Cantabular server")
					So(hcMock.AddAndGetCheckCalls()[1].Name, ShouldResemble, "Cantabular API Extension")
					So(hcMock.AddAndGetCheckCalls()[2].Name, ShouldResemble, "Dataset API client")
					So(hcMock.AddAndGetCheckCalls()[3].Name, ShouldResemble, "Datastore")
					So(hcMock.AddAndGetCheckCalls()[4].Name, ShouldResemble, "Zebedee")
				})
			})
		})
	})
}

func TestStart(t *testing.T) {
	Convey("Having a correctly initialised Service with mocked dependencies", t, func() {
		cfg, err := config.Get()
		So(err, ShouldBeNil)

		hcMock := &mock.HealthCheckerMock{
			StartFunc: func(ctx context.Context) {},
		}

		serverWg := &sync.WaitGroup{}
		serverMock := &mock.HTTPServerMock{}

		svc := &service.Service{
			Cfg:         cfg,
			Server:      serverMock,
			HealthCheck: hcMock,
		}

		Convey("When a service with a successful HTTP server is started", func() {
			serverMock.ListenAndServeFunc = func() error {
				serverWg.Done()
				return nil
			}
			serverWg.Add(1)
			svc.Start(ctx, make(chan error, 1))

			Convey("Then healthcheck is started and HTTP server starts listening", func() {
				So(len(hcMock.StartCalls()), ShouldEqual, 1)
				serverWg.Wait() // Wait for HTTP server go-routine to finish
				So(len(serverMock.ListenAndServeCalls()), ShouldEqual, 1)
			})
		})

		Convey("When a service with a failing HTTP server is started", func() {
			serverMock.ListenAndServeFunc = func() error {
				serverWg.Done()
				return errServer
			}
			errChan := make(chan error, 1)
			serverWg.Add(1)
			svc.Start(ctx, errChan)

			Convey("Then HTTP server errors are reported to the provided errors channel", func() {
				rxErr := <-errChan
				So(rxErr.Error(), ShouldResemble, fmt.Sprintf("failure in http listen and serve: %s", errServer.Error()))
			})
		})
	})
}

func TestClose(t *testing.T) {
	Convey("Having a correctly initialised service", t, func() {
		cfg, err := config.Get()
		So(err, ShouldBeNil)

		hcStopped := false

		// healthcheck Stop does not depend on any other service being closed/stopped
		hcMock := &mock.HealthCheckerMock{
			StopFunc: func() { hcStopped = true },
		}

		// server Shutdown will fail if healthcheck is not stopped
		serverMock := &mock.HTTPServerMock{
			ShutdownFunc: func(ctx context.Context) error {
				if !hcStopped {
					return fmt.Errorf("server stopped before healthcheck")
				}
				return nil
			},
		}

		producerMock := &kafkatest.IProducerMock{
			CloseFunc: func(ctx context.Context) error { return nil },
		}

		svc := &service.Service{
			Cfg:         cfg,
			Server:      serverMock,
			HealthCheck: hcMock,
			Producer:    producerMock,
		}
		Convey("Closing the service results in all the dependencies being closed in the expected order", func() {
			err := svc.Close(context.Background())
			So(err, ShouldBeNil)
			So(hcMock.StopCalls(), ShouldHaveLength, 1)
			So(serverMock.ShutdownCalls(), ShouldHaveLength, 1)
		})

		Convey("If services fail to stop, the Close operation tries to close all dependencies and returns an error", func() {
			serverMock.ShutdownFunc = func(ctx context.Context) error {
				return errServer
			}

			err = svc.Close(context.Background())
			So(err, ShouldNotBeNil)
			So(hcMock.StopCalls(), ShouldHaveLength, 1)
			So(serverMock.ShutdownCalls(), ShouldHaveLength, 1)
		})
		/*
			   TODO - figure out if this has been added to the boilerplate code or been removed from the csv exporter,
			   such that, do we need to make this work here

					Convey("If service times out while shutting down, the Close operation fails with the expected error", func() {
						cfg.GracefulShutdownTimeout = 1 * time.Millisecond
						timeoutServerMock := &mock.HTTPServerMock{
							ListenAndServeFunc: func() error { return nil },
							ShutdownFunc: func(ctx context.Context) error {
								time.Sleep(2 * time.Millisecond)
								return nil
							},
						}

						svcList := service.NewServiceList(nil)
						svcList.HealthCheck = true
						svc := service.Service{
							Config:      cfg,
							ServiceList: svcList,
							Server:      timeoutServerMock,
							HealthCheck: hcMock,
						}

						err = svc.Close(context.Background())
						So(err, ShouldNotBeNil)
						So(err.Error(), ShouldResemble, "context deadline exceeded")
						So(len(hcMock.StopCalls()), ShouldEqual, 1)
						So(len(timeoutServerMock.ShutdownCalls()), ShouldEqual, 1)
					})*/
	})
}
