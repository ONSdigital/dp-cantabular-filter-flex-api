package api

import (
	"context"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/cantabularmetadata"
	"github.com/pkg/errors"
)

/*
CheckDefaultCategorisation checks the default categorisation of a given dimension
so that we store the correct parent dimension for a given set of options?
*/
func (api *API) CheckDefaultCategorisation(dimName string, datasetName string) (string, string, error) {

	ctx := context.Background()
	labelMap := make(map[string]string)
	cats, err := api.ctblr.GetCategorisations(ctx, cantabular.GetCategorisationsRequest{
		Dataset:  datasetName,
		Variable: dimName,
	})
	if err != nil {
		return "", "", err
	}

	names := make([]string, 0)

	if len(cats.Dataset.Variables.Edges) > 0 {
		for _, edge := range cats.Dataset.Variables.Edges {
			if len(edge.Node.MapFrom) > 0 {
				for _, mapFrom := range edge.Node.MapFrom {
					if len(mapFrom.Edges) > 0 {
						for _, mappedEdge := range mapFrom.Edges {
							if len(mappedEdge.Node.IsSourceOf.Edges) > 0 {
								for _, FINALLY := range mappedEdge.Node.IsSourceOf.Edges {
									names = append(names, FINALLY.Node.Name)
									labelMap[FINALLY.Node.Name] = FINALLY.Node.Label
								}
							}
						}
					}
				}
			} else if len(edge.Node.IsSourceOf.Edges) > 0 {
				for _, sourceOf := range edge.Node.IsSourceOf.Edges {
					names = append(names, sourceOf.Node.Name)
				}
			}
		}
	}

	if len(names) == 0 {
		return "", "", errors.New("no categorisations recieved for variable")
	}

	defaultCat, err := api.metadata.GetDefaultClassification(ctx, cantabularmetadata.GetDefaultClassificationRequest{
		Dataset:   datasetName,
		Variables: names,
	})
	if err != nil {

		return "", "", err
	}

	return defaultCat.Variable, labelMap[defaultCat.Variable], nil

}
