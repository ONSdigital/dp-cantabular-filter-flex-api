package api

import (
	"context"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/cantabularmetadata"
)

/*
   CheckDefaultCategorisation checks the default categorisation of a given dimension
   so that we store the correct parent dimension for a given set of options?
*/
func (api *API) CheckDefaultCategorisation(dimName string, datasetName string) (string, error) {

	ctx := context.Background()
	cats, err := api.ctblr.GetCategorisations(ctx, cantabular.GetCategorisationsRequest{
		Dataset:  datasetName,
		Variable: dimName,
	})
	if err != nil {
		return "", err
	}

	names := make([]string, 0)
	for _, edge := range cats.Dataset.Variables.Search.Edges {
		names = append(names, edge.Node.Name)

	}

	defaultCat, err := api.metadata.GetDefaultClassification(ctx, cantabularmetadata.GetDefaultClassificationRequest{
		Dataset:   datasetName,
		Variables: names,
	})
	if err != nil {

		return "", err
	}

	return defaultCat.Variable, nil

}
