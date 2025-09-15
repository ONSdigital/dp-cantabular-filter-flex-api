package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/config"
	mongo "github.com/ONSdigital/dp-cantabular-filter-flex-api/datastore/mongodb"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/generator"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/cantabularmetadata"
	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	dphttp "github.com/ONSdigital/dp-net/v3/http"
	"github.com/ONSdigital/dp-net/v3/responder"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// GetHTTPServer creates an http server and sets the Server
var GetHTTPServer = func(bindAddr string, router http.Handler) HTTPServer {
	otelHandler := otelhttp.NewHandler(router, "/")
	s := dphttp.NewServer(bindAddr, otelHandler)
	s.HandleOSSignals = false
	return s
}

// GetResponder gets a http request responder
var GetResponder = func() Responder {
	return responder.New()
}

var GetMongoDB = func(ctx context.Context, cfg *config.Config, g Generator) (Datastore, error) {
	return mongo.NewClient(ctx, g, mongo.Config{
		MongoDriverConfig:       cfg.Mongo,
		FilterAPIURL:            cfg.FilterAPIURL,
		FiltersCollection:       cfg.FiltersCollection,
		FilterOutputsCollection: cfg.FilterOutputsCollection,
	})
}

// GetCantabularClient gets and initialises the Cantabular Client
var GetCantabularClient = func(cfg *config.Config) CantabularClient {
	return cantabular.NewClient(
		cantabular.Config{
			Host:           cfg.CantabularURL,
			ExtApiHost:     cfg.CantabularExtURL,
			GraphQLTimeout: cfg.DefaultRequestTimeout,
		},
		dphttp.NewClient(),
		nil,
	)
}

var GetMetadataClient = func(cfg *config.Config) MetadataClient {
	// NB: Client initialisation is not consistent, see above...
	return cantabularmetadata.NewClient(cantabularmetadata.Config{
		Host:           cfg.MetadataAPIURL,
		GraphQLTimeout: cfg.DefaultRequestTimeout,
	}, dphttp.NewClient())
}

// GetDatasetAPIClient gets and initialises the DatasetAPI Client
var GetDatasetAPIClient = func(cfg *config.Config) DatasetAPIClient {
	return dataset.NewAPIClient(cfg.DatasetAPIURL)
}

// GetPopulationClient gets and initialises the PopultaionTypesAPI Client
var GetPopulationClient = func(cfg *config.Config) (PopulationTypesAPIClient, error) {
	return population.NewClient(cfg.PopulationTypesAPIURL)
}

// GetHealthCheck creates a healthcheck with versionInfo
var GetHealthCheck = func(cfg *config.Config, buildTime, gitCommit, version string) (HealthChecker, error) {
	versionInfo, err := healthcheck.NewVersionInfo(buildTime, gitCommit, version)
	if err != nil {
		return nil, fmt.Errorf("failed to get version info: %w", err)
	}

	hc := healthcheck.New(
		versionInfo,
		cfg.HealthCheckCriticalTimeout,
		cfg.HealthCheckInterval,
	)
	return &hc, nil
}

var GetGenerator = func() Generator {
	return generator.New()
}
