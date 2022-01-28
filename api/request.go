package api

import (
	"io"
	"fmt"
	"net/http"
	"encoding/json"

	"github.com/ONSdigital/log.go/v2/log"
	"github.com/pkg/errors"
)

// ParseRequest attemts to read unmarshal a request body into a 
// request object, returning an appropriate error on failure.
// 'req' must be a pointer to a struct.
// ParseRequest will also attempt to call the request's Valid()
// function if it has one and will throw an error if it fails
func (api *API) ParseRequest(body io.Reader, req interface{}) error {
	b, err := io.ReadAll(body)
	if err != nil{
		return Error{
			err:     errors.Wrap(err, "failed to read request body"),
			message: "failed to read request body",
		}
	}

	if err := json.Unmarshal(b, &req); err != nil{
		return Error{
			statusCode: http.StatusBadRequest,
			err:        fmt.Errorf("failed to unmarshal request body: %w", err),
			message:    "badly formed request body",
			logData:    log.Data{
				"body": string(b),
			},
		}
	}

	if v, ok := req.(validator); ok {
		if err := v.Valid(); err != nil{
			return Error{
				statusCode: http.StatusBadRequest,
				err:        fmt.Errorf("invalid request: %w", err),
				logData:    log.Data{
					"body":    string(b),
					"request": fmt.Sprintf("%+v", req),
				},
			}
		}
	}

	return nil
}
