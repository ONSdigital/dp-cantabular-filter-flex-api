package api

import (
	"context"
	"net/http"
	"time"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"

	"github.com/google/uuid"
)

// responder handles responding to http requests
type responder interface {
	JSON(context.Context, http.ResponseWriter, int, interface{})
	Error(context.Context, http.ResponseWriter, int, error)
	Errors(context.Context, http.ResponseWriter, int, []error)
}

type datastore interface {
	CreateFilter(context.Context, *model.Filter) error
}

type validator interface {
	Valid() error
}

type generator interface {
	PSK() ([]byte, error)
	UUID() (uuid.UUID, error)
	Timestamp() time.Time
	URL(host, path string, args ...interface{}) string
}

type cantabularClient interface {
	GetDimensionOptions(context.Context, cantabular.GetDimensionOptionsRequest) (*cantabular.GetDimensionOptionsResponse, error)
	StatusCode(error) int
}

type datasetAPIClient interface {
	GetVersion(ctx context.Context, userAuthToken, svcAuthToken, downloadSvcAuthToken, collectionID, datasetID, edition, version string) (dataset.Version, error)
}

type coder interface {
	Code() int
}
