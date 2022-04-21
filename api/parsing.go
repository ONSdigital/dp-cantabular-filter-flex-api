package api

import "github.com/ONSdigital/dp-cantabular-filter-flex-api/model"

func (api *API) generateFilterOutput(filter *model.Filter) model.FilterOutput {
	filterLinks := model.FilterOutputLinks{
		Version: filter.Links.Version,
		Self: model.Link{
			HREF: api.generate.URL(api.cfg.FilterAPIURL, "/filter-outputs"),
			// uuid created by mongo client, will set there.
			ID: "",
		},
		FilterBlueprint: model.Link{
			HREF: api.generate.URL(api.cfg.FilterAPIURL, "/filters"),
			ID:   filter.ID,
		},
	}

	// Gets amended through PUT as time goes on.
	downloads := model.Downloads{}

	// same as above
	events := []model.Event{}

	filterOutput := model.FilterOutput{
		// id created by mongo client
		FilterID:   filter.ID,
		State:      model.Submitted,
		Dataset:    filter.Dataset,
		Downloads:  downloads,
		Events:     events,
		Links:      filterLinks,
		Published:  filter.Published,
		Dimensions: filter.Dimensions,
	}

	return filterOutput
}
