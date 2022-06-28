package api

import (
	"fmt"
	"net/http"

	dperrors "github.com/ONSdigital/dp-cantabular-filter-flex-api/errors"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/log.go/v2/log"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

func (api *API) addFilterDimensionOption(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := addFilterDimensionOptionRequest{
		FilterID:  chi.URLParam(r, "id"),
		Dimension: chi.URLParam(r, "dimension"),
		Option:    chi.URLParam(r, "option"),
	}

	logData := log.Data{
		"filter_id": req.FilterID,
		"dimension": req.Dimension,
		"option":    req.Option,
	}

	filter, err := api.store.GetFilter(ctx, req.FilterID)
	if err != nil {
		api.respond.Error(
			ctx,
			w,
			statusCode(err),
			Error{
				err:     errors.Wrap(err, "failed to get filter"),
				message: "failed to add dimension option: failed to get filter",
			},
		)
		return
	}

	// Check dimension exists
	var dimension model.Dimension
	var dimExists bool

	for _, d := range filter.Dimensions {
		if d.Name == req.Dimension {
			dimension = d
			dimExists = true
			break
		}
	}

	if !dimExists {
		api.respond.Error(
			ctx,
			w,
			http.StatusBadRequest,
			Error{
				err:     errors.New("failed to add dimension option: dimension not found in filter"),
				logData: logData,
			},
		)
		return
	}

	// Check if option already exists
	var optExists bool
	for _, o := range dimension.Options {
		if o == req.Option {
			optExists = true
			break
		}
	}

	if optExists {
		api.respond.Error(
			ctx,
			w,
			http.StatusBadRequest,
			Error{
				err:     errors.New("failed to add dimension option: option already added to dimension"),
				logData: logData,
			},
		)
		return
	}

	// Check option is valid
	dReq := cantabular.GetDimensionOptionsRequest{
		Dataset: filter.PopulationType,
	}
	dReq.DimensionNames = append(dReq.DimensionNames, dimension.ID)
	dReq.Filters = append(dReq.Filters, cantabular.Filter{
		Codes:    []string{req.Option},
		Variable: dimension.ID,
	})

	if _, err := api.ctblr.GetDimensionOptions(ctx, dReq); err != nil {
		logData["request"] = dReq
		if api.ctblr.StatusCode(err) >= http.StatusInternalServerError {
			err = Error{
				err:     errors.Wrap(err, "failed to query dimension options from Cantabular"),
				message: "Internal Server Error",
				logData: logData,
			}
		} else {
			err = Error{
				err:     errors.WithStack(err),
				message: "invalid option for filter",
				logData: logData,
			}
		}
		api.respond.Error(ctx, w, api.ctblr.StatusCode(err), err)
		return
	}

	var eTag string
	if reqETag := api.getETag(r); reqETag != eTagAny {
		eTag = reqETag
	}

	// Add option to filter
	dimension.Options = append(dimension.Options, req.Option)
	newETag, err := api.store.UpdateFilterDimension(ctx, req.FilterID, req.Dimension, dimension, eTag)
	if err != nil {
		api.respond.Error(ctx, w, statusCode(err), Error{
			err:     errors.Wrap(err, "failed to uodate dimension with option in store"),
			message: "failed to add dimension option",
			logData: logData,
		})
		return
	}

	resp := addFilterDimensionOptionResponse{
		Option: req.Option,
		Links: filterDimensionOptionLinks{
			Filter: filter.Links.Self,
			Self: model.Link{
				ID: "",
				HREF: fmt.Sprintf(
					"%s/filters/%s/dimensions/%s/options/%s",
					api.cfg.BindAddr,
					filter.ID,
					req.Dimension,
					req.Option,
				),
			},
			Dimension: model.Link{
				ID: dimension.ID,
				HREF: fmt.Sprintf(
					"%s/filters/%s/dimensions/%s",
					api.cfg.BindAddr,
					filter.ID,
					req.Dimension,
				),
			},
		},
	}

	w.Header().Set(eTagHeader, newETag)

	api.respond.JSON(ctx, w, http.StatusCreated, resp)
}

func (api *API) deleteFilterDimensionOption(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := deleteFilterDimensionOptionRequest{
		FilterID:  chi.URLParam(r, "id"),
		Dimension: chi.URLParam(r, "dimension"),
		Option:    chi.URLParam(r, "option"),
	}

	logData := log.Data{
		"filter_id": req.FilterID,
		"dimension": req.Dimension,
		"option":    req.Option,
	}

	filter, err := api.store.GetFilter(ctx, req.FilterID)
	if err != nil {
		status := statusCode(err)
		if dperrors.NotFound(err) {
			status = http.StatusBadRequest
		}
		api.respond.Error(
			ctx,
			w,
			status,
			Error{
				err:     errors.Wrap(err, "failed to get filter"),
				message: "failed to delete dimension option: failed to get filter",
			},
		)
		return
	}

	// Check dimension exists
	var dimension model.Dimension
	var dimExists bool

	var dimIndex int
	for i, d := range filter.Dimensions {
		if d.Name == req.Dimension {
			dimension = d
			dimExists = true
			dimIndex = i
			break
		}
	}

	if !dimExists {
		api.respond.Error(
			ctx,
			w,
			http.StatusBadRequest,
			Error{
				err:     errors.New("failed to delete dimension option: dimension not found in filter"),
				logData: logData,
			},
		)
		return
	}

	// Check option exists
	var optExists bool
	var optIndex int
	for i, o := range dimension.Options {
		if o == req.Option {
			optExists = true
			optIndex = i
			break
		}
	}

	if !optExists {
		api.respond.Error(
			ctx,
			w,
			http.StatusNotFound,
			Error{
				err:     errors.New("failed to delete dimension option: option not found"),
				logData: logData,
			},
		)
		return
	}

	// Remove option from dimension by replacing option with duplicate of last element in slice, then
	// chopping the end of the slice (faster than creating a brand new slice in place)
	opts := filter.Dimensions[dimIndex].Options
	opts[optIndex] = opts[len(opts)-1]
	opts = opts[:len(opts)-1]
	filter.Dimensions[dimIndex].Options = opts

	log.Info(ctx, "DEBUG OPTIONS", log.Data{
		"opts":     opts,
		"dim.opts": filter.Dimensions[dimIndex].Options,
	})

	var eTag string
	if reqETag := api.getETag(r); reqETag != eTagAny {
		eTag = reqETag
	}

	newETag, err := api.store.RemoveFilterDimensionOption(ctx, req.FilterID, req.Dimension, req.Option, eTag)
	if err != nil {
		api.respond.Error(ctx, w, statusCode(err), Error{
			err:     errors.Wrap(err, "failed to uodate dimension with option in store"),
			message: "failed to delete dimension option",
			logData: logData,
		})
		return
	}

	w.Header().Set(eTagHeader, newETag)

	api.respond.StatusCode(w, http.StatusNoContent)
}
