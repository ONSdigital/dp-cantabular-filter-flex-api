package api

import (
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
)

func (api *API) createFilter(w http.ResponseWriter, r *http.Request){
	ctx := r.Context()
	var req createFilterRequest

	if err := api.ParseRequest(r.Body, &req); err != nil{
		api.respond.Error(ctx, w, err)
		return
	}
	defer r.Body.Close()

	f := model.Filter{
		Links: model.Links{
			Version: model.Link{
				HREF: fmt.Sprintf(
					"%s/datasets/%s/editions/%s/version/%d",
					api.cfg.DatasetAPIURL,
					req.DatasetID,
					req.Edition,
					req.Version,
				),
				ID: req.InstanceID,
			},
		},
		Dimensions:      req.Dimensions,
		UniqueTimestamp: api.generate.Timestamp(),
		LastUpdated:     api.generate.Timestamp(),
		InstanceID:      *req.InstanceID,
		Dataset: model.Dataset{
			ID:      req.DatasetID,
			Edition: req.Edition,
			Version: req.Version,
		},
		Published:         true, // TODO: Not sure what to
		Events:            nil,  // populate for these
		DisclosureControl: nil,  // fields yet
	}

	if err := api.store.CreateFilter(ctx, &f); err != nil{
		api.respond.Error(ctx, w, fmt.Errorf("failed to create filter: %w", err))
		return
	}

	resp := createFilterResponse{
		Filter: f,
	}

	api.respond.JSON(ctx, w, http.StatusCreated, resp)
}
