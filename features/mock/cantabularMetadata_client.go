package mock

import (
	"context"
	"errors"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/cantabularmetadata"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
)

type CantabularMetadataClient struct {
	ErrStatus                    int
	OptionsHappy                 bool
	DimensionsHappy              bool
	GetDefaultClassificationFunc func(ctx context.Context, req cantabularmetadata.GetDefaultClassificationRequest) (*cantabularmetadata.GetDefaultClassificationResponse, error)
}

func (c *CantabularMetadataClient) Reset() {
	c.ErrStatus = 500
	c.OptionsHappy = true
	c.DimensionsHappy = true
}

func (c *CantabularMetadataClient) StatusCode(_ error) int {
	return c.ErrStatus
}

func (c *CantabularMetadataClient) GetDimensionOptions(_ context.Context, _ cantabular.GetDimensionOptionsRequest) (*cantabular.GetDimensionOptionsResponse, error) {
	if c.OptionsHappy {
		return nil, nil
	}

	return nil, errors.New("invalid dimension options")
}

func (c *CantabularMetadataClient) GetDefaultClassification(ctx context.Context, req cantabularmetadata.GetDefaultClassificationRequest) (*cantabularmetadata.GetDefaultClassificationResponse, error) {
	if c.OptionsHappy {
		return c.GetDefaultClassificationFunc(ctx, req)
	}
	return nil, errors.New("error while retrieving default classification")
}

func (c *CantabularMetadataClient) Checker(_ context.Context, _ *healthcheck.CheckState) error {
	return nil
}

func (c *CantabularMetadataClient) CheckerAPIExt(_ context.Context, _ *healthcheck.CheckState) error {
	return nil
}

/* type CantabularServer struct {
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
	cs.NewHandler().Post(fmt.Sprintf("/graphql")).Handle(cs.PostResponder())

	cs.Lock()
	defer cs.Unlock()
	cs.responses[crc(request)] = response
}

func (cs *CantabularServer) PostResponder() httpfake.Responder {
	return func(w http.ResponseWriter, r *http.Request, rh *httpfake.Request) {
		buf, err := ioutil.ReadAll(r.Body)
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
	//Remove hard coded newlines and tabs as well
	reduced = bytes.ReplaceAll(reduced, []byte("\\n"), []byte(``))
	reduced = bytes.ReplaceAll(reduced, []byte("\\t"), []byte(``))

	return crc32.ChecksumIEEE(reduced)
}
*/
