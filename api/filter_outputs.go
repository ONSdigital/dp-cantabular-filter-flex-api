package api

import (
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
				err:     errors.Wrap(err, "failed to get filter output"),
				message: "failed to get filter output",
				logData: log.Data{
					"id": fID,
				},
			},
		)
		return
	}

	resp := getFilterOutputResponse{
		*filterOutput,
	}

	api.respond.JSON(ctx, w, http.StatusOK, resp)
}

func (api *API) putFilterOutput(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fID := chi.URLParam(r, "filter_output_id")

	var req putFilterOutputRequest

	if err := api.ParseRequest(r.Body, &req); err != nil {
		api.respond.Error(
			ctx,
			w,
			http.StatusBadRequest,
			Error{
				err: errors.Wrap(err, "invalid request body"),
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
			Error{
				err:     errors.Wrap(err, "failed to update filter output"),
				message: "failed to update filter output",
				logData: log.Data{
					"id": fID,
				},
			},
		)
		return
	}

	api.respond.StatusCode(w, http.StatusOK)
}

//createFilterOutputEvent will add a new event to the list of existing events in the filter outputs
func (api *API) addFilterOutputEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fID := chi.URLParam(r, "filter_output_id")

	var req addFilterOutputEventRequest

	if err := api.ParseRequest(r.Body, &req); err != nil {
		api.respond.Error(
			ctx,
			w,
			http.StatusBadRequest,
			Error{
				err: errors.Wrap(err, "invalid request body"),
				logData: log.Data{
					"id": fID,
				},
			},
		)
		return
	}

	if err := api.store.AddFilterOutputEvent(ctx, fID, &req.Event); err != nil {
		api.respond.Error(
			ctx,
			w,
			statusCode(err),
			err,
		)
		return
	}

	api.respond.StatusCode(w, http.StatusCreated)
}
