package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/golang/glog"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/config"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/service"
	dpotelgo "github.com/ONSdigital/dp-otel-go"
	"github.com/ONSdigital/log.go/v2/log"
)

const serviceName = "dp-cantabular-filter-flex-api"

var (
	// BuildTime represents the time in which the service was built
	BuildTime string
	// // GitCommit represents the commit (SHA-1) hash of the service that is running
	GitCommit string
	// // Version represents the version of the service that is running
	Version string

	// TODO: remove below explainer before commiting
	//NOTE: replace the above with the below to run code with for example vscode debugger.
	// BuildTime string = "1601119818"
	// GitCommit string = "6584b786caac36b6214ffe04bf62f058d4021538"
	// Version   string = "v0.1.0"
)

func main() {
	log.Namespace = serviceName
	ctx := context.Background()

	if err := run(ctx); err != nil {
		log.Error(ctx, "fatal runtime error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) (err error) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	svcErrors := make(chan error, 1)

	// Read config
	cfg, err := config.Get()
	if err != nil {
		return fmt.Errorf("unable to retrieve service configuration: %w", err)
	}
	log.Info(ctx, "config on startup", log.Data{"config": cfg, "build_time": BuildTime, "git-commit": GitCommit})

	// Set up OpenTelemetry
	otelConfig := dpotelgo.Config{
		OtelServiceName:          cfg.OTServiceName,
		OtelExporterOtlpEndpoint: cfg.OTExporterOTLPEndpoint,
		OtelBatchTimeout:         cfg.OTBatchTimeout,
	}

	otelShutdown, err := dpotelgo.SetupOTelSDK(ctx, otelConfig)

	if err != nil {
		log.Error(ctx, "error setting up OpenTelemetry - hint: ensure OTEL_EXPORTER_OTLP_ENDPOINT is set", err)
	}
	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	// Make sure that context is cancelled when 'run' finishes its execution.
	// Any remaining go-routine that was not terminated during svc.Close (graceful shutdown) will be terminated by ctx.Done()
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	// Run the service
	svc := service.New()
	if err := svc.Init(ctx, cfg, BuildTime, GitCommit, Version); err != nil {
		return fmt.Errorf("running service failed with error: %w", err)
	}
	svc.Start(ctx, svcErrors)

	// Blocks until an os interrupt or a fatal error occurs
	select {
	case err := <-svcErrors:
		err = fmt.Errorf("service error received: %w", err)
		if errClose := svc.Close(ctx); errClose != nil {
			log.Error(ctx, "service close error during error handling", errClose)
		}
		return err
	case sig := <-signals:
		log.Info(ctx, "os signal received", log.Data{"signal": sig})
	}
	return svc.Close(ctx)
}
