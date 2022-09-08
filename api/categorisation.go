package api

import (
	"context"
)

/*
   CheckDefaultCategorisation checks the default categorisation of a given dimension
   so that we store the correct parent dimension for a given set of options?
*/
func (api *API) CheckDefaultCategorisation(dimName string, datasetName string) (string, error) {

	ctx := context.Background()
	cats, err := api.populations.GetCategorisations(ctx, nil)
	names := make([]string, 0)

	for _, cat := range cats {
		names = append(names, cat.Name)
	}
	defaultCat, err := api.metadata.GetDefaultClassification(ctx, GetDefaultClassificationRequest{
		Dataset:   datasetName,
		Variables: names,
	})
	if err != nil {
		return "", err
	}

	return defaultCat.Variable, nil

}
