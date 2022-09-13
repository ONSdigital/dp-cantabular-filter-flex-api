package api

import (
	"context"
	"fmt"

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
		fmt.Printf(dimName)
		fmt.Printf(datasetName)
		fmt.Printf("%+v\n", err)
		return "", err
	}

	names := make([]string, 0)
	for _, edge := range cats.Dataset.Variables.Search.Edges {
		names = append(names, edge.Node.Name)

	}

	fmt.Printf("%+v\n\n\n\n\n\n", dimName)
	fmt.Printf("%+v\n\n\n\n\n\n", names)
	//	fmt.Printf("%+v\n\n\n\n\n\n", cats)
	defaultCat, err := api.metadata.GetDefaultClassification(ctx, cantabularmetadata.GetDefaultClassificationRequest{
		Dataset:   datasetName,
		Variables: names,
	})
	if err != nil {
		fmt.Println("ERROR FROM GETTING DEFAULT CLASSIFICATION")
		return "", err
	}

	fmt.Println("FOUND ONE FOUND ONE \n\n\n\n\n\n")
	fmt.Println(defaultCat.Variable)

	return defaultCat.Variable, nil

}
