package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	dprequest "github.com/ONSdigital/dp-net/v2/request"
	"github.com/ONSdigital/log.go/v2/log"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	flexible  = "flexible"
	published = "published"
)

func (api *API) createFilter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req createFilterRequest

	if err := api.ParseRequest(r.Body, &req); err != nil {
		api.respond.Error(
			ctx,
			w,
			http.StatusBadRequest,
			errors.Wrap(err, "failed to parse request"),
		)
		return
	}

	v, err := api.datasets.GetVersion(
		ctx,
		"",
		api.cfg.ServiceAuthToken,
		"",
		"",
		req.Dataset.ID,
		req.Dataset.Edition,
		strconv.Itoa(req.Dataset.Version),
	)
	if err != nil {
		api.respond.Error(
			ctx,
			w,
			statusCode(err),
			Error{
				err:     errors.Wrap(err, "failed to get existing Version"),
				message: "failed to get existing dataset information",
			},
		)
		return
	}

	if v.State != published && !dprequest.IsCallerPresent(ctx) {
		api.respond.Error(
			ctx,
			w,
			http.StatusNotFound,
			Error{
				err:     errors.New("unauthenticated request on unpublished dataset"),
				message: "dataset not found",
			},
		)
		return
	}

	dimIDs, err := api.validateDimensions(ctx, req.Dimensions, v.Dimensions)
	if err != nil {
		api.respond.Error(
			ctx,
			w,
			http.StatusBadRequest,
			errors.Wrap(err, "failed to validate request dimensions"),
		)
		return
	}

	// Validate dimension options by performing Cantabular query with selections,
	// skip this check if requesting all options
	if err := api.validateDimensionOptions(ctx, req, dimIDs); err != nil {
		api.respond.Error(
			ctx,
			w,
			statusCode(err),
			errors.Wrap(err, "failed to validate dimension options"),
		)
		return
	}

	f := model.Filter{
		Links: model.Links{
			Version: model.Link{
				HREF: api.generate.URL(
					api.cfg.DatasetAPIURL,
					"/datasets/%s/editions/%s/version/%d",
					req.Dataset.ID,
					req.Dataset.Edition,
					req.Dataset.Version,
				),
				ID: strconv.Itoa(v.Version),
			},
		},
		Dimensions:        req.Dimensions,
		UniqueTimestamp:   api.generate.Timestamp(),
		LastUpdated:       api.generate.Timestamp(),
		Dataset:           *req.Dataset,
		PopulationType:    req.PopulationType,
		Type:              flexible,
		Published:         v.State == published,
		Events:            nil, // TODO: Not sure what to
		DisclosureControl: nil, // populate for these fields yet
	}

	if f.InstanceID, err = uuid.Parse(v.ID); err != nil {
		api.respond.Error(
			ctx,
			w,
			http.StatusInternalServerError,
			Error{
				err:     errors.Wrap(err, "failed to parse version instance id"),
				message: "Internal Server Error",
				logData: log.Data{"version": v},
			},
		)
		return
	}

	if err := api.store.CreateFilter(ctx, &f); err != nil {
		api.respond.Error(
			ctx,
			w,
			statusCode(err),
			errors.Wrap(err, "failed to create filter"),
		)
		return
	}

	resp := createFilterResponse{
		Filter: f,
	}

	api.respond.JSON(ctx, w, http.StatusCreated, resp)
}

// validateDimensions validates provided filter dimensions exist within the dataset dimensions provided.
// Returns a map of the dimensions name:id for use in the following validation calls
func (api *API) validateDimensions(ctx context.Context, filterDims []model.Dimension, dims []dataset.VersionDimension) (map[string]string, error) {
	dimensions := make(map[string]string)
	for _, d := range dims {
		dimensions[d.Name] = d.ID
	}

	var incorrect []string
	for _, fd := range filterDims {
		if _, ok := dimensions[fd.Name]; !ok {
			incorrect = append(incorrect, fd.Name)
		}
	}

	if incorrect != nil {
		return nil, errors.Errorf("incorrect dimensions chosen: %v", incorrect)
	}

	return dimensions, nil
}

// validateDimensionOptions validates the dimension options in a createFilterRequest by making a call
// to Cantabular to check they exist. Takes as a second argument a map mapping the dimension names to
// ids, which are required by the Cantabular query
// TODO: Requires rule variable to be first in POST request, acceptable?
func (api *API) validateDimensionOptions(ctx context.Context, req createFilterRequest, dimIDs map[string]string) error {
	dReq := cantabular.GetDimensionOptionsRequest{
		Dataset: req.PopulationType,
	}

	// only validate dimensions with specified options
	for _, d := range req.Dimensions {
		if len(d.Options) > 0 {
			dReq.DimensionNames = append(dReq.DimensionNames, dimIDs[d.Name])
			dReq.Filters = append(dReq.Filters, cantabular.Filter{
				Codes:    d.Options,
				Variable: dimIDs[d.Name],
			})
		}
	}
	if len(dReq.Filters) == 0 {
		return nil
	}

	if _, err := api.ctblr.GetDimensionOptions(ctx, dReq); err != nil {
		if api.ctblr.StatusCode(err) >= http.StatusInternalServerError {
			return Error{
				err:     errors.Wrap(err, "failed to query dimension options from Cantabular"),
				message: "Internal Server Error",
				logData: log.Data{
					"request": dReq,
				},
			}
		}
		return Error{
			err:     errors.WithStack(err),
			message: "failed to validate dimension options for filter",
		}
	}

	return nil
}
