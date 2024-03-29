package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ONSdigital/log.go/v2/log"
	"github.com/pkg/errors"
)

const (
	ifMatchHeader = "If-Match"
	eTagHeader    = "ETag"
	eTagAny       = "*"
)

// ParseRequest attemts to read unmarshal a request body into a
// request object, returning an appropriate error on failure.
// 'req' must be a pointer to a struct.
// ParseRequest will also attempt to call the request's Valid()
// function if it has one and will throw an error if it fails

// NOTE: for multivariate tables, it is going to use the name field
// for the id when making the call to cantabular.
func (api *API) ParseRequest(body io.Reader, req interface{}) error {
	b, err := io.ReadAll(body)
	if err != nil {
		return Error{
			err:     errors.Wrap(err, "failed to read request body"),
			message: "failed to read request body",
		}
	}

	if err := json.Unmarshal(b, &req); err != nil {
		return Error{
			err:     fmt.Errorf("failed to unmarshal request body: %w", err),
			message: fmt.Sprintf("badly formed request body: %s", err),
			logData: log.Data{
				"body": string(b),
			},
		}
	}

	if v, ok := req.(validator); ok {
		if err := v.Valid(); err != nil {
			return Error{
				err: errors.Wrap(err, "invalid request"),
				logData: log.Data{
					"body":    string(b),
					"request": fmt.Sprintf("%+v", req),
				},
			}
		}
	}

	return nil
}

func (api *API) getETag(r *http.Request) string {
	if eTag := r.Header.Get(ifMatchHeader); eTag != "" {
		return eTag
	}
	return eTagAny
}
