package api

import (
	"fmt"
	"net/http"
	"strings"

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

		status := http.StatusNotFound
		message := "filter output not found"

		if !strings.HasSuffix(err.Error(), "no documents in result") {
			status = http.StatusInternalServerError
			message = "internal service error"
		}

		api.respond.Error(
			ctx,
			w,
			status,
			Error{
				err: errors.Wrap(err, message),
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

	var req createFilterOutputRequest

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
			errors.Wrap(err, "failed to create filter outputs"),
		)
		return
	}
	api.respond.StatusCode(w, http.StatusOK)
}
