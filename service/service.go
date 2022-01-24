package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/api"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/config"
	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

// Service contains all the configs, server and clients to run the event handler service
type Service struct {
	Cfg         *config.Config
	Server      HTTPServer
	HealthCheck HealthChecker
	Producer    kafka.IProducer
	Api         *api.API
	responder   Responder
	store       Datastore
}

func New() *Service {
	return &Service{}
}

// Init initialises the service and it's dependencies
func (svc *Service) Init(ctx context.Context, cfg *config.Config, buildTime, gitCommit, version string) error {
	var err error

	if cfg == nil {
		return errors.New("nil config passed to service init")
	}

	svc.Cfg = cfg

	if svc.Producer, err = GetKafkaProducer(ctx, cfg); err != nil {
		return fmt.Errorf("failed to create kafka producer: %w", err)
	}

	// Get HealthCheck
	if svc.HealthCheck, err = GetHealthCheck(cfg, buildTime, gitCommit, version); err != nil {
		return fmt.Errorf("could not instantiate healthcheck: %w", err)
	}

	if err := svc.registerCheckers(); err != nil {
		return fmt.Errorf("error initialising checkers: %w", err)
	}

	svc.responder = GetResponder()
	if svc.store, err = GetMongoDB(ctx, cfg); err != nil{
		return fmt.Errorf("failed to initialise mongodb store: %w", err)
	}

	r := mux.NewRouter()
	r.StrictSlash(true).Path("/health").HandlerFunc(svc.HealthCheck.Handler)
	// TODO: Add other(s) to serviceList here

	// Setup the API
	svc.Api = api.New(ctx, cfg, r, svc.responder, svc.store)
	svc.Server = GetHTTPServer(cfg.BindAddr, r)

	return nil
}

// Start the service
func (svc *Service) Start(ctx context.Context, svcErrors chan error) {
	log.Info(ctx, "starting service")

	// Always start healthcheck.
	svc.HealthCheck.Start(ctx)

	// Run the http server in a new go-routine
	go func() {
		if err := svc.Server.ListenAndServe(); err != nil {
			svcErrors <- fmt.Errorf("failure in http listen and serve: %w", err)
		}
	}()
}

// Close gracefully shuts the service down in the required order, with timeout
func (svc *Service) Close(ctx context.Context) error {
	timeout := svc.Cfg.GracefulShutdownTimeout
	log.Info(ctx, "commencing graceful shutdown", log.Data{"graceful_shutdown_timeout": timeout})
	ctx, cancel := context.WithTimeout(ctx, timeout)
	hasShutdownError := false

	go func() {
		defer cancel()

		// stop healthcheck, as it depends on everything else
		if svc.HealthCheck != nil {
			log.Info(ctx, "stopping health checker")
			svc.HealthCheck.Stop()
			log.Info(ctx, "stopped health checker")
		}

		// stop any incoming requests before closing any outbound connections
		if svc.Server != nil {
			log.Info(ctx, "stopping http server")
			if err := svc.Server.Shutdown(ctx); err != nil {
				log.Error(ctx, "failed to shutdown http server", err)
				hasShutdownError = true
			}
			log.Info(ctx, "stopped http server")
		}

		// TODO: Close other dependencies, in the expected order
	}()

	// wait for shutdown success (via cancel) or failure (timeout)
	<-ctx.Done()

	// timeout expired
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("shutdown timed out: %w", ctx.Err())
	}

	// other error
	if hasShutdownError {
		return fmt.Errorf("failed to shutdown gracefully")
	}

	log.Info(ctx, "graceful shutdown was successful")
	return nil
}

// registerCheckers adds the checkers for the service clients to the health check object.
func (svc *Service) registerCheckers() error {
	if _, err := svc.HealthCheck.AddAndGetCheck("Kafka producer", svc.Producer.Checker); err != nil {
		return fmt.Errorf("error adding check for Kafka producer: %w", err)
	}

	// TODO: add other health checks here, as per dp-upload-service

	return nil
}
