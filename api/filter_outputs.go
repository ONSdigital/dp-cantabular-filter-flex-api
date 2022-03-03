package api

import (
	"net/http"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	dprequest "github.com/ONSdigital/dp-net/v2/request"
	"github.com/pkg/errors"
)

func (api *API) createFilterOutput(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if !dprequest.IsCallerPresent(ctx) {
		api.respond.Error(
			ctx,
			w,
			http.StatusNotFound,
			Error{
				err:     errors.New("unauthenticated request on unpublished dataset"),
				message: "caller not found",
			},
		)
		return
	}

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

	f := model.FilterOutputResponse{
		Download: req.Downloads,
	}
	if err := api.store.CreateFilterOutput(ctx, &f); err != nil {
		api.respond.Error(
			ctx,
			w,
			statusCode(err),
			errors.Wrap(err, "failed to create filter"),
		)
		return
	}

	resp := createFilterOutputsResponse{
		ID:        f.ID,
		Downloads: f.Download,
	}

	api.respond.JSON(ctx, w, http.StatusCreated, resp)
}
