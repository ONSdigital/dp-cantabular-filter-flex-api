package service

import (
	"context"
	"net/http"
	"time"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	mongo "github.com/ONSdigital/dp-mongodb/v3/mongodb"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//go:generate moq -out mock/server.go -pkg mock . HTTPServer
//go:generate moq -out mock/datastore.go -pkg mock . Datastore
//go:generate moq -out mock/responder.go -pkg mock . Responder
//go:generate moq -out mock/generator.go -pkg mock . Generator
//go:generate moq -out mock/health_check.go -pkg mock . HealthChecker

// HTTPServer defines the required methods from the HTTP server
type HTTPServer interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

// HealthChecker defines the required methods from Healthcheck
type HealthChecker interface {
	Handler(http.ResponseWriter, *http.Request)
	Start(context.Context)
	Stop()
	AddAndGetCheck(name string, checker healthcheck.Checker) (*healthcheck.Check, error)
	Subscribe(healthcheck.Subscriber, ...*healthcheck.Check)
}

// Responder handles responding to http requests
type Responder interface {
	JSON(context.Context, http.ResponseWriter, int, interface{})
	Error(context.Context, http.ResponseWriter, int, error)
	StatusCode(http.ResponseWriter, int)
	Bytes(context.Context, http.ResponseWriter, int, []byte)
	Errors(context.Context, http.ResponseWriter, int, []error)
}

// Datastore is the interface for interacting with the storage backend
type Datastore interface {
	CreateFilter(context.Context, *model.Filter) error
	GetFilter(context.Context, string) (*model.Filter, error)
	CreateFilterOutput(context.Context, *model.FilterOutput) error
	GetFilterDimensions(context.Context, string) ([]model.Dimension, error)
	Checker(context.Context, *healthcheck.CheckState) error
	Conn() *mongo.MongoConnection
}

// Generator is the interface for generating dynamic tokens and timestamps
type Generator interface {
	PSK() ([]byte, error)
	UUID() (uuid.UUID, error)
	Timestamp() time.Time
	UniqueTimestamp() primitive.Timestamp
	URL(host, path string, args ...interface{}) string
}

type CantabularClient interface {
	GetDimensionOptions(context.Context, cantabular.GetDimensionOptionsRequest) (*cantabular.GetDimensionOptionsResponse, error)
	StatusCode(error) int
	Checker(context.Context, *healthcheck.CheckState) error
	CheckerAPIExt(context.Context, *healthcheck.CheckState) error
}

type DatasetAPIClient interface {
	GetVersion(ctx context.Context, userAuthToken, svcAuthToken, downloadSvcAuthToken, collectionID, datasetID, edition, version string) (dataset.Version, error)
	Checker(context.Context, *healthcheck.CheckState) error
}
