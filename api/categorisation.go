package api

import (
	"context"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
)

/*
   CheckDefaultCategorisation checks the default categorisation of a given dimension
   so that we store the correct parent dimension for a given set of options?
*/
func (api *API) CheckDefaultCategorisation(dimName string) (dimension *model.Dimension, err error) {

	ctx := context.Background()
	cats, err := api.populations.GetCategorisations(ctx, nil)
	names := make([]string, 0)

	for _, cat := range cats {
		names = append(names, cat.Name)
	}
	defaultCat, err := api.metadata.GetDefaultCategorisation(ctx, nil)

	for _, categorisation := range defaultCat {
		if categorisation.isDefaultCat == true {
			dimension = defaultCat
			return
		}
	}
	return nil, err
}
