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

	r.Header.Set("X-Forwarded-Host", r.Header.Get("X-Forwarded-API-Host"))

	var filterOutput *model.FilterOutput

	filterFlexLinksBuilder := links.FromHeadersOrDefault(&r.Header, api.cantabularFilterFlexAPIURL)
	datasetLinksBuilder := links.FromHeadersOrDefault(&r.Header, api.datasetAPIURL)

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
		filterOutput.Links.Version.HREF, err = datasetLinksBuilder.BuildLink(filterOutput.Links.Version.HREF)
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
						"id":   fID,
						"href": filterOutput.Links.FilterBlueprint.HREF,
					},
				},
			)
			return
		}

		if filterOutput.Downloads.TXT != nil {
			filterOutput.Downloads.TXT.HREF, err = links.BuildDownloadLink(filterOutput.Downloads.TXT.HREF, api.downloadServiceURL)
			if err != nil {
				api.respond.Error(
					ctx,
					w,
					statusCode(err),
					Error{
						err:     errors.Wrap(err, "failed to build txt download link"),
						message: "failed to build txt download link",
						logData: log.Data{
							"id":   fID,
							"href": filterOutput.Downloads.TXT.HREF,
						},
					},
				)
				return
			}
		}

		if filterOutput.Downloads.CSV != nil {
			filterOutput.Downloads.CSV.HREF, err = links.BuildDownloadLink(filterOutput.Downloads.CSV.HREF, api.downloadServiceURL)
			if err != nil {
				api.respond.Error(
					ctx,
					w,
					statusCode(err),
					Error{
						err:     errors.Wrap(err, "failed to build csv download link"),
						message: "failed to build csv download link",
						logData: log.Data{
							"id":   fID,
							"href": filterOutput.Downloads.CSV.HREF,
						},
					},
				)
				return
			}
		}

		if filterOutput.Downloads.CSVW != nil {
			filterOutput.Downloads.CSVW.HREF, err = links.BuildDownloadLink(filterOutput.Downloads.CSVW.HREF, api.downloadServiceURL)
			if err != nil {
				api.respond.Error(
					ctx,
					w,
					statusCode(err),
					Error{
						err:     errors.Wrap(err, "failed to build csvw download link"),
						message: "failed to build csvw download link",
						logData: log.Data{
							"id":   fID,
							"href": filterOutput.Downloads.CSVW.HREF,
						},
					},
				)
				return
			}
		}

		if filterOutput.Downloads.XLS != nil {
			filterOutput.Downloads.XLS.HREF, err = links.BuildDownloadLink(filterOutput.Downloads.XLS.HREF, api.downloadServiceURL)
			if err != nil {
				api.respond.Error(
					ctx,
					w,
					statusCode(err),
					Error{
						err:     errors.Wrap(err, "failed to build xls download link"),
						message: "failed to build xls download link",
						logData: log.Data{
							"id":   fID,
							"href": filterOutput.Downloads.XLS.HREF,
						},
					},
				)
				return
			}
		}

		if filterOutput.Downloads.XLSX != nil {
			filterOutput.Downloads.XLSX.HREF, err = links.BuildDownloadLink(filterOutput.Downloads.XLSX.HREF, api.downloadServiceURL)
			if err != nil {
				api.respond.Error(
					ctx,
					w,
					statusCode(err),
					Error{
						err:     errors.Wrap(err, "failed to build xlsx download link"),
						message: "failed to build xlsx download link",
						logData: log.Data{
							"id":   fID,
							"href": filterOutput.Downloads.XLSX.HREF,
						},
					},
				)
				return
			}
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
