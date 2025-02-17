package api

import (
	"net/http"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	"github.com/ONSdigital/dp-net/v2/links"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

func (api *API) getFilterOutput(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fID := chi.URLParam(r, "id")

	var filterOutput *model.FilterOutput

	filterFlexLinksBuilder := links.FromHeadersOrDefault(&r.Header, api.cantabularFilterFlexAPIURL)

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

	if api.cfg.EnableURLRewriting {
		filterOutput.Links.Version.HREF, err = filterFlexLinksBuilder.BuildLink(filterOutput.Links.Version.HREF)
		if err != nil {
			api.respond.Error(
				ctx,
				w,
				statusCode(err),
				Error{
					err:     errors.Wrap(err, "failed to build version link"),
					message: "failed to build version link",
					logData: log.Data{
						"id":   fID,
						"href": filterOutput.Links.Version.HREF,
					},
				},
			)
			return
		}

		filterOutput.Links.Self.HREF, err = filterFlexLinksBuilder.BuildLink(filterOutput.Links.Self.HREF)
		if err != nil {
			api.respond.Error(
				ctx,
				w,
				statusCode(err),
				Error{
					err:     errors.Wrap(err, "failed to build self link"),
					message: "failed to build self link",
					logData: log.Data{
						"id":   fID,
						"href": filterOutput.Links.Self.HREF,
					},
				},
			)
			return
		}
		filterOutput.Links.FilterBlueprint.HREF, err = filterFlexLinksBuilder.BuildLink(filterOutput.Links.FilterBlueprint.HREF)
		if err != nil {
			api.respond.Error(
				ctx,
				w,
				statusCode(err),
				Error{
					err:     errors.Wrap(err, "failed to build filter blueprint link"),
					message: "failed to build filter blueprint link",
					logData: log.Data{
						"id": fID,
					},
				},
			)
			return
		}
	}

	resp := getFilterOutputResponse{
		*filterOutput,
	}

	api.respond.JSON(ctx, w, http.StatusOK, resp)
}

func (api *API) putFilterOutput(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fID := chi.URLParam(r, "id")

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
		ID:              fID,
		State:           req.State,
		Downloads:       req.Downloads,
		LastUpdated:     api.generate.Timestamp(),
		UniqueTimestamp: api.generate.UniqueTimestamp(),
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

// createFilterOutputEvent will add a new event to the list of existing events in the filter outputs
func (api *API) addFilterOutputEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fID := chi.URLParam(r, "id")

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
