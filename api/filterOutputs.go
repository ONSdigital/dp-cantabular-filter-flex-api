package api

import (
	"net/http"

	"github.com/pkg/errors"
)

func (api *API) createFilterOutputs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req createFilterOutputsRequest

	if err := api.ParseRequest(r.Body, &req); err != nil {
		api.respond.Error(
			ctx,
			w,
			http.StatusBadRequest,
			errors.Wrap(err, "failed to parse request"),
		)
		return
	}

	if err := api.store.CreateFilterOutputs(ctx, &req.Downloads); err != nil {
		api.respond.Error(
			ctx,
			w,
			statusCode(err),
			errors.Wrap(err, "failed to create filter"),
		)
		return
	}

	resp := createFilterOutputsResponse{
		FilterOutput: req.Downloads,
	}

	api.respond.JSON(ctx, w, http.StatusCreated, resp)
}
