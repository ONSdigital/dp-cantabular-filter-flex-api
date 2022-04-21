package api

import (
	"context"
	"net/http"
	"time"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// To mock interfaces in this file
//go:generate mockgen -source=interface.go -destination=mock/mock_interface.go -package=mock github.com/ONSdigital/dp-cantabular-filter-flex-api/api

// responder handles responding to http requests
type responder interface {
	JSON(context.Context, http.ResponseWriter, int, interface{})
	StatusCode(http.ResponseWriter, int)
	Error(context.Context, http.ResponseWriter, int, error)
	Errors(context.Context, http.ResponseWriter, int, []error)
}

type datastore interface {
	CreateFilter(context.Context, *model.Filter) error
	GetFilter(context.Context, string) (*model.Filter, error)
	CreateFilterOutput(context.Context, *model.FilterOutput) error
	GetFilterOutput(context.Context, string) (*model.FilterOutput, error)
	UpdateFilterOutput(context.Context, *model.FilterOutput) error
	AddFilterOutputEvent(context.Context, string, *model.Event) error
	GetFilterDimensions(context.Context, string, int, int) ([]model.Dimension, int, error)
	GetFilterDimension(ctx context.Context, fID, dimName string) (model.Dimension, error)
	AddFilterDimension(context.Context, string, model.Dimension) error
	UpdateFilterDimension(ctx context.Context, filterID string, dimensionName string, dimension model.Dimension, currentETag string) (eTag string, err error)
}

type validator interface {
	Valid() error
}

type generator interface {
	Timestamp() time.Time
	UniqueTimestamp() primitive.Timestamp
	URL(host, path string, args ...interface{}) string
}

type cantabularClient interface {
	GetDimensionOptions(context.Context, cantabular.GetDimensionOptionsRequest) (*cantabular.GetDimensionOptionsResponse, error)
	StaticDatasetQuery(context.Context, cantabular.StaticDatasetQueryRequest) (*cantabular.StaticDatasetQuery, error)
	GetGeographyDimensions(context.Context, string) (*cantabular.GetGeographyDimensionsResponse, error)
	SearchDimensions(ctx context.Context, req cantabular.SearchDimensionsRequest) (*cantabular.GetDimensionsResponse, error)
	StatusCode(error) int
}

type datasetAPIClient interface {
	GetVersion(ctx context.Context, userAuthToken, svcAuthToken, downloadSvcAuthToken, collectionID, datasetID, edition, version string) (dataset.Version, error)
	GetVersionDimensions(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, id, edition, version string) (dataset.VersionDimensions, error)
	GetOptionsInBatches(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, id, edition, version, dimension string, batchSize, maxWorkers int) (dataset.Options, error)
	GetDatasetCurrentAndNext(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, datasetID string) (dataset.Dataset, error)
}

type coder interface {
	Code() int
}
