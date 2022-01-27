package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/config"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/generator"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/responder"
	mongo "github.com/ONSdigital/dp-cantabular-filter-flex-api/datastore/mongodb"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	kafka "github.com/ONSdigital/dp-kafka/v3"
	dphttp "github.com/ONSdigital/dp-net/http"
)

// GetHTTPServer creates an http server and sets the Server
var GetHTTPServer = func(bindAddr string, router http.Handler) HTTPServer {
	s := dphttp.NewServer(bindAddr, router)
	s.HandleOSSignals = false
	return s
}

// GetResponder gets a http request responder
var GetResponder = func() Responder {
	return responder.New()
}

var GetMongoDB = func(ctx context.Context, cfg *config.Config, g Generator) (Datastore, error) {
	return mongo.NewClient(ctx, g, mongo.Config{
		MongoDriverConfig: cfg.Mongo,
		FilterFlexAPIURL:  cfg.BindAddr,
		FiltersCollection: "filters",
		FilterOutputsCollection: "filterOutputs",
	})
}

// GetKafkaProducer creates a Kafka producer
var GetKafkaProducer = func(ctx context.Context, cfg *config.Config) (kafka.IProducer, error) {
	pConfig := &kafka.ProducerConfig{
		BrokerAddrs:       cfg.Kafka.Addr,
		Topic:             cfg.Kafka.CsvCreatedTopic,
		MinBrokersHealthy: &cfg.Kafka.ProducerMinBrokersHealthy,
		KafkaVersion:      &cfg.Kafka.Version,
		MaxMessageBytes:   &cfg.Kafka.MaxBytes,
	}
	if cfg.Kafka.SecProtocol == config.KafkaTLSProtocolFlag {
		pConfig.SecurityConfig = kafka.GetSecurityConfig(
			cfg.Kafka.SecCACerts,
			cfg.Kafka.SecClientCert,
			cfg.Kafka.SecClientKey,
			cfg.Kafka.SecSkipVerify,
		)
	}
	return kafka.NewProducer(ctx, pConfig)
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
