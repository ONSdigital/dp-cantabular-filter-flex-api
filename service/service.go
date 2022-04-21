package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/api"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/config"
	kafka "github.com/ONSdigital/dp-kafka/v3"

	"github.com/ONSdigital/dp-api-clients-go/v2/identity"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/log.go/v2/log"

	"github.com/go-chi/chi/v5"
)

// Service contains all the configs, server and clients to run the event handler service
type Service struct {
	Cfg              *config.Config
	Server           HTTPServer
	HealthCheck      HealthChecker
	Api              *api.API
	responder        Responder
	store            Datastore
	Producer         kafka.IProducer
	generator        Generator
	cantabularClient CantabularClient
	datasetAPIClient DatasetAPIClient
	identityClient   *identity.Client
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

	svc.identityClient = identity.New(cfg.ZebedeeURL)

	// Get HealthCheck
	if svc.HealthCheck, err = GetHealthCheck(cfg, buildTime, gitCommit, version); err != nil {
		return fmt.Errorf("could not instantiate healthcheck: %w", err)
	}

	if svc.Producer, err = GetKafkaProducer(ctx, cfg); err != nil {
		return fmt.Errorf("Could not initialise Kafka producer: %w", err)
	}

	svc.cantabularClient = GetCantabularClient(cfg)
	svc.datasetAPIClient = GetDatasetAPIClient(cfg)
	svc.generator = GetGenerator()
	svc.responder = GetResponder()

	if svc.store, err = GetMongoDB(ctx, cfg, svc.generator); err != nil {
		return fmt.Errorf("failed to initialise mongodb store: %w", err)
	}

	if err := svc.registerCheckers(); err != nil {
		return fmt.Errorf("error initialising checkers: %w", err)
	}

	r := chi.NewRouter()
	r.Handle("/health", http.HandlerFunc(svc.HealthCheck.Handler))
	// TODO: Add other(s) to serviceList here

	// Setup the API
	svc.Api = api.New(
		ctx,
		cfg,
		r,
		svc.identityClient,
		svc.responder,
		svc.generator,
		svc.store,
		svc.datasetAPIClient,
		svc.cantabularClient,
		svc.Producer,
	)
	svc.Server = GetHTTPServer(cfg.BindAddr, r)

	return nil
}

// Start the service
func (svc *Service) Start(ctx context.Context, svcErrors chan error) {
	log.Info(ctx, "starting service")

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
		if err := svc.Producer.Close(ctx); err != nil {
			log.Info(ctx, "failed to shut down kafka producer")
		}

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
	// TODO - when Cantabular server is deployed to Production, remove this placeholder and the flag,
	// and always use the real Checker instead: svc.cantabularClient.Checker
	cantabularChecker := svc.cantabularClient.Checker
	cantabularAPIExtChecker := svc.cantabularClient.CheckerAPIExt
	if !svc.Cfg.CantabularHealthcheckEnabled {
		cantabularChecker = func(ctx context.Context, state *healthcheck.CheckState) error {
			return state.Update(healthcheck.StatusOK, "Cantabular healthcheck placeholder", http.StatusOK)
		}
		cantabularAPIExtChecker = func(ctx context.Context, state *healthcheck.CheckState) error {
			return state.Update(healthcheck.StatusOK, "Cantabular APIExt healthcheck placeholder", http.StatusOK)
		}
	}

	if _, err := svc.HealthCheck.AddAndGetCheck("Cantabular server", cantabularChecker); err != nil {
		return fmt.Errorf("error adding check for Cantabular server: %w", err)
	}

	if _, err := svc.HealthCheck.AddAndGetCheck("Cantabular API Extension", cantabularAPIExtChecker); err != nil {
		return fmt.Errorf("error adding check for Cantabular api extension: %w", err)
	}

	if _, err := svc.HealthCheck.AddAndGetCheck("Dataset API client", svc.datasetAPIClient.Checker); err != nil {
		return fmt.Errorf("error addig check for dataset API client: %w", err)
	}

	if _, err := svc.HealthCheck.AddAndGetCheck("Datastore", svc.store.Checker); err != nil {
		return fmt.Errorf("error adding check for datastore: %w", err)
	}

	if _, err := svc.HealthCheck.AddAndGetCheck("Zebedee", svc.identityClient.Checker); err != nil {
		return fmt.Errorf("error adding check for datastore: %w", err)
	}
	if _, err := svc.HealthCheck.AddAndGetCheck("Kafka", svc.Producer.Checker); err != nil {
		return fmt.Errorf("error adding check for Kafka producer: %w", err)
	}

	return nil
}
