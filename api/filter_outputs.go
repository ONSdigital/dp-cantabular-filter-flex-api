package api

import (
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

func (api *API) getFilterOutput(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fID := chi.URLParam(r, "filter-output-id")

	var filterOutput *model.FilterOutput

	filterOutput, err := api.store.GetFilterOutput(ctx, fID)
	if err != nil {
		api.respond.Error(
			ctx,
			w,
			statusCode(err),
			Error{
				err: errors.Wrap(err, "failed to get filter output"),
				logData: log.Data{
					"id": fID,
				},
			},
		)
		return
	}

	resp := getFilterOutputResponse{
		model.JobState{
			InstanceID:       filterOutput.InstanceID,
			FilterID:         filterOutput.FilterID,
			DimensionListUrl: fmt.Sprintf("%s/filter-outputs/%s", api.cfg.BindAddr, filterOutput.FilterID),
			Events:           filterOutput.Events,
		},
		filterOutput.Links,
		filterOutput.Downloads,
		filterOutput.State,
	}

	api.respond.JSON(ctx, w, http.StatusOK, resp)

}

func (api *API) updateFilterOutput(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fID := chi.URLParam(r, "filter-output-id")

	var req updateFilterOutputRequest

	if err := api.ParseRequest(r.Body, &req); err != nil {
		api.respond.Error(
			ctx,
			w,
			http.StatusBadRequest,
			Error{
				err: errors.Wrap(err, "failed to parse request"),
				logData: log.Data{
					"id": fID,
				},
			},
		)
		return
	}

	f := model.FilterOutput{
		ID:        fID,
		State:     req.State,
		Downloads: req.Downloads,
	}

	if err := api.store.UpdateFilterOutput(ctx, &f); err != nil {
		api.respond.Error(
			ctx,
			w,
			statusCode(err),
			errors.Wrap(err, "failed to update filter output"),
		)
		return
	}

	api.respond.StatusCode(w, http.StatusOK)
}
