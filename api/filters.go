package api

import (
	"fmt"
	"github.com/pkg/errors"
	"net/http"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
)

func (api *API) createFilter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req createFilterRequest

	if err := api.ParseRequest(r.Body, &req); err != nil {
		api.respond.Error(
			ctx,
			w,
			http.StatusBadRequest,
			fmt.Errorf("failed to parse request: %w", err),
		)
		return
	}

	// getVersion
	// validate dimensions

	f := model.Filter{
		Links: model.Links{
			Version: model.Link{
				HREF: fmt.Sprintf(
					"%s/datasets/%s/editions/%s/version/%d",
					api.cfg.DatasetAPIURL,
					req.Dataset.ID,
					req.Dataset.Edition,
					req.Dataset.Version,
				),
				ID: nil, //"",//version.ID,,
			},
		},
		Dimensions:      req.Dimensions,
		UniqueTimestamp: api.generate.Timestamp(),
		LastUpdated:     api.generate.Timestamp(),
		//InstanceID:      "",//version.ID,
		Dataset: model.Dataset{
			ID:      req.Dataset.ID,
			Edition: req.Dataset.Edition,
			Version: req.Dataset.Version,
		},
		Published:         true, // TODO: Not sure what to
		Events:            nil,  // populate for these
		DisclosureControl: nil,  // fields yet
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
