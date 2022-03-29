package api

import (
	"net/http"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

func (api *API) updateFilterOutput(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fID := chi.URLParam(r, "filter_output_id")

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

//createFilterOutputEvent will add a new event to the list of existing events in the filter outputs
func (api *API) addFilterOutputEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fID := chi.URLParam(r, "filter_output_id")

	var req eventRequest

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

	f := model.Event{
		Timestamp: req.Timestamp,
		Name:      req.Name,
	}

	if err := api.store.AddFilterOutputEvent(ctx, fID, &f); err != nil {
		api.respond.Error(
			ctx,
			w,
			statusCode(err),
			errors.Wrap(err, "failed to update filter outputs event"),
		)
		return
	}

	api.respond.JSON(ctx, w, http.StatusOK, nil)
}
