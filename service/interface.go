package service

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/config"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"
)

//go:generate moq -out mock/initialiser.go -pkg mock . Initialiser
//go:generate moq -out mock/server.go -pkg mock . HTTPServer
//go:generate moq -out mock/health_check.go -pkg mock . HealthChecker

// Initialiser defines the methods to initialise external services
type Initialiser interface {
	DoGetHTTPServer(bindAddr string, router http.Handler) HTTPServer
	DoGetHealthCheck(cfg *config.Config, buildTime, gitCommit, version string) (HealthChecker, error)
}

// HTTPServer defines the required methods from the HTTP server
type HTTPServer interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

// HealthChecker defines the required methods from Healthcheck
type HealthChecker interface {
	Handler(w http.ResponseWriter, req *http.Request)
	Start(ctx context.Context)
	Stop()
	AddAndGetCheck(name string, checker healthcheck.Checker) (check *healthcheck.Check, err error)
	Subscribe(s healthcheck.Subscriber, checks ...*healthcheck.Check)
}

// Responder handles responding to http requests
type Responder interface{
	JSON(context.Context,http.ResponseWriter, int, interface{})
	Error(context.Context, http.ResponseWriter, error)
	StatusCode(http.ResponseWriter, int)
	Bytes(context.Context, http.ResponseWriter, int, []byte)
}

// Datastore is the interface for interacting with the storage backend
type Datastore interface{
	CreateFilter(context.Context, *model.Filter) error
}