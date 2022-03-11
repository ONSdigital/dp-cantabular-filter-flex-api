package api

import (
	"net/http"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	"github.com/pkg/errors"
)

func (api *API) createFilterOutput(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req createFilterOutputRequest

	if err := api.ParseRequest(r.Body, &req); err != nil {
		api.respond.Error(
			ctx,
			w,
			http.StatusBadRequest,
			errors.Wrap(err, "failed to parse request"),
		)
		return
	}

	f := model.FilterOutput{
		State:     req.State,
		Downloads: req.Downloads,
	}

	if err := api.store.CreateFilterOutput(ctx, &f); err != nil {
		api.respond.Error(
			ctx,
			w,
			statusCode(err),
			errors.Wrap(err, "failed to create filter outputs"),
		)
		return
	}

	resp := createFilterOutputResponse{
		f,
	}

	api.respond.JSON(ctx, w, http.StatusCreated, resp)
}
