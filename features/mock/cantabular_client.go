package mock

import (
	"bytes"
	"context"
	"errors"
	"hash/crc32"
	"io"
	"net/http"
	"sync"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular/gql"
	"github.com/ONSdigital/dp-api-clients-go/v2/stream"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"

	"github.com/maxcnunes/httpfake"
)

type CantabularClient struct {
	ErrStatus                           int
	OptionsHappy                        bool
	DimensionsHappy                     bool
	ResponseTooLarge                    bool
	GetObservationsResponse             *cantabular.GetObservationsResponse
	GetDimensionsByNameFunc             func(context.Context, cantabular.GetDimensionsByNameRequest) (*cantabular.GetDimensionsResponse, error)
	SearchDimensionsFunc                func(ctx context.Context, req cantabular.SearchDimensionsRequest) (*cantabular.GetDimensionsResponse, error)
	GetGeographyDimensionsInBatchesFunc func(ctx context.Context, datasetID string, batchSize, maxWorkers int) (*gql.Dataset, error)
	GetAreaFunc                         func(context.Context, cantabular.GetAreaRequest) (*cantabular.GetAreaResponse, error)
	StaticDatasetQueryFunc              func(context.Context, cantabular.StaticDatasetQueryRequest) (*cantabular.StaticDatasetQuery, error)
	GetCategorisationsFunc              func(ctx context.Context, req cantabular.GetCategorisationsRequest) (*cantabular.GetCategorisationsResponse, error)
	StaticDatasetQueryStreamJSONFunc    func(context.Context, cantabular.StaticDatasetQueryRequest, stream.Consumer) (cantabular.GetObservationsResponse, error)
	CheckQueryCountFunc                 func(context.Context, cantabular.StaticDatasetQueryRequest) (int, error)
}

func (c *CantabularClient) Reset() {
	c.ErrStatus = 500
	c.OptionsHappy = true
	c.DimensionsHappy = true
}

func (c *CantabularClient) StatusCode(_ error) int {
	return c.ErrStatus
}

func (c *CantabularClient) GetDimensionOptions(_ context.Context, _ cantabular.GetDimensionOptionsRequest) (*cantabular.GetDimensionOptionsResponse, error) {
	if c.OptionsHappy {
		return nil, nil
	}

	return nil, errors.New("invalid dimension options")
}

func (c *CantabularClient) GetCategorisations(ctx context.Context, req cantabular.GetCategorisationsRequest) (*cantabular.GetCategorisationsResponse, error) {
	if c.OptionsHappy {
		return c.GetCategorisationsFunc(ctx, req)
	}
	return nil, errors.New("error while retrieving categorisations")
}

func (c *CantabularClient) StaticDatasetQuery(ctx context.Context, req cantabular.StaticDatasetQueryRequest) (*cantabular.StaticDatasetQuery, error) {
	if c.OptionsHappy {
		return c.StaticDatasetQueryFunc(ctx, req)
	}
	return nil, errors.New("error while executing dataset query")
}

func (c *CantabularClient) StaticDatasetQueryStreamJSON(ctx context.Context, req cantabular.StaticDatasetQueryRequest, str stream.Consumer) (cantabular.GetObservationsResponse, error) {
	if c.OptionsHappy {
		return c.StaticDatasetQueryStreamJSONFunc(ctx, req, str)
	}
	return *c.GetObservationsResponse, errors.New("error while executing dataset query")
}

func (c *CantabularClient) CheckQueryCount(ctx context.Context, req cantabular.StaticDatasetQueryRequest) (int, error) {
	if c.OptionsHappy {
		return c.CheckQueryCountFunc(ctx, req)
	}
	if c.ResponseTooLarge {
		return 500000, nil
	}
	return 0, errors.New("error while executing dataset query")
}

func (c *CantabularClient) GetGeographyDimensionsInBatches(ctx context.Context, datasetID string, batchSize, maxWorkers int) (*gql.Dataset, error) {
	if c.OptionsHappy {
		return c.GetGeographyDimensionsInBatchesFunc(ctx, datasetID, batchSize, maxWorkers)
	}
	return nil, errors.New("error while getting geography dimensions")
}

func (c *CantabularClient) GetArea(ctx context.Context, req cantabular.GetAreaRequest) (*cantabular.GetAreaResponse, error) {
	if c.OptionsHappy {
		return c.GetAreaFunc(ctx, req)
	}
	return nil, errors.New("error while getting area dimensions")
}

func (c *CantabularClient) GetDimensionsByName(ctx context.Context, req cantabular.GetDimensionsByNameRequest) (*cantabular.GetDimensionsResponse, error) {
	if c.DimensionsHappy {
		return c.GetDimensionsByNameFunc(ctx, req)
	}
	return nil, errors.New("error while searching dimensions")
}

func (c *CantabularClient) Checker(_ context.Context, _ *healthcheck.CheckState) error {
	return nil
}

func (c *CantabularClient) CheckerAPIExt(_ context.Context, _ *healthcheck.CheckState) error {
	return nil
}

//nolint:gocritic // embedded mutex only used in mock client
type CantabularServer struct {
	*httpfake.HTTPFake

	sync.RWMutex
	responses map[uint32][]byte
}

func NewCantabularServer(t *testing.T) *CantabularServer {
	return &CantabularServer{
		HTTPFake:  httpfake.New(httpfake.WithTesting(t)),
		responses: map[uint32][]byte{},
	}
}

func (cs *CantabularServer) Reset() {
	cs.HTTPFake.Reset()
	cs.responses = map[uint32][]byte{}
}

func (cs *CantabularServer) Handle(request, response []byte) {
	cs.NewHandler().Post("/graphql").Handle(cs.PostResponder())

	cs.Lock()
	defer cs.Unlock()
	cs.responses[crc(request)] = response
}

func (cs *CantabularServer) PostResponder() httpfake.Responder {
	return func(w http.ResponseWriter, r *http.Request, rh *httpfake.Request) {
		buf, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		cs.RLock()
		defer cs.RUnlock()
		if resp, ok := cs.responses[crc(buf)]; ok {
			_, err = w.Write(resp)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		w.WriteHeader(http.StatusBadRequest)
	}
}

func crc(request []byte) uint32 {
	reduced := bytes.Map(func(r rune) rune {
		switch r {
		case '\n', '\t', ' ':
			return -1
		default:
			return r
		}
	}, request)
	// Remove hard coded newlines and tabs as well
	reduced = bytes.ReplaceAll(reduced, []byte("\\n"), []byte(``))
	reduced = bytes.ReplaceAll(reduced, []byte("\\t"), []byte(``))

	return crc32.ChecksumIEEE(reduced)
}
