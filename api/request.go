package api

import (
	"io"
	"fmt"
	"net/http"
	"encoding/json"

	"github.com/ONSdigital/log.go/v2/log"
)

// UnmarshalRequestBody attemts to read and unmarshal a request body into a 
// request object, returning an appropriate error on failure.
// 'req' must be a pointer to a struct
func (api *API) UnmarshalRequestBody(body io.Reader, req interface{}) error {
	b, err := io.ReadAll(body)
	if err != nil{
		return Error{
			err:     fmt.Errorf("failed to read request body: %w", err),
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

	return err
}