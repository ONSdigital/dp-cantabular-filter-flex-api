package api

import (
	"fmt"
	"time"
	"net/http"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
)

func (api *API) createFilter(w http.ResponseWriter, r *http.Request){
	ctx := r.Context()
	var req createFilterRequest

	if err := api.UnmarshalRequestBody(r.Body, &req); err != nil{
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
		UniqueTimestamp: time.Now(),
		LastUpdated:     time.Now(),
		InstanceID:      req.InstanceID,
		Published:       true,
		DisclosureControl: nil, // Not sure what to populate from here
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
