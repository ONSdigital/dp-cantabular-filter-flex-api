package api

import (
	"net/http"

	"github.com/pkg/errors"
)

func (api *API) CreateFilterOutput(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request CreateFilterOutputRequest
	if err := api.ParseRequest(r.Body, &request); err != nil {
		api.respond.Error(ctx, w, http.StatusBadRequest, errors.Wrap(err, "parse request"))
		return
	}
	if err := api.validate.Struct(&request); err != nil {
		api.respond.Error(ctx, w, http.StatusBadRequest, errors.Wrap(err, "validate request fields"))
		return
	}

	if err := api.store.CreateFilterOutput(ctx, request.Downloads); err != nil {
		api.respond.Error(ctx, w, http.StatusBadRequest, errors.Wrap(err, "upsert filter output"))
		return
	}

	api.respond.JSON(ctx, w, http.StatusOK, &CreateFilterOutputResponse{FilterOutput: request.Downloads})
}
