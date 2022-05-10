package api

import (
	"context"
	"net/http"
)

func (api *API) getDatasetJSON(w http.ResponseWriter, r *http.Request) {
	var response *getDatasetJsonObservationsResponse
	var err error

	ctx := r.Context()

	if r.URL.Query().Get("geography") == "" {
		response, err = api.getDefaultDatasetJSON(ctx, r)
	} else {
		response, err = api.getGeographyDatsetJSON(ctx, r)
	}

	if err != nil {
		api.respond.Error(
			ctx,
			w,
			http.StatusBadRequest,
			err,
		)
		return
	}

	api.respond.JSON(ctx, w, http.StatusOK, response)
}

func (api *API) getDefaultDatasetJSON(ctx context.Context, r *http.Request) (*getDatasetJsonObservationsResponse, error) {
	return &getDatasetJsonObservationsResponse{}, nil
}

func (api *API) getGeographyDatsetJSON(ctx context.Context, r *http.Request) (*getDatasetJsonObservationsResponse, error) {
	return &getDatasetJsonObservationsResponse{}, nil
}
