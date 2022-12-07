package api

import (
	"context"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/cantabularmetadata"
	"github.com/pkg/errors"
)

func (api *API) CheckDefaultCategorisation(dimName string, datasetName string) (string, string, error) {
	ctx := context.Background()
	labelMap := make(map[string]string)
	cats, err := api.ctblr.GetCategorisations(ctx, cantabular.GetCategorisationsRequest{
		Dataset:  datasetName,
		Variable: dimName,
	})
	if err != nil {
		return "", "", errors.Wrap(err, "failed to check default categorisation")
	}

	names := make([]string, 0)
	for _, edge := range cats.Dataset.Variables.Edges {
		if len(edge.Node.MapFrom) > 0 {
			for _, mapFrom := range edge.Node.MapFrom {
				for _, _ = range mapFrom.Edges {
					for _, mappedSource := range mapFrom.Edges {
						for _, mappedSourceEdge := range mappedSource.Node.IsSourceOf.Edges {
							names = append(names, mappedSourceEdge.Node.Name)
							labelMap[mappedSourceEdge.Node.Name] = mappedSourceEdge.Node.Label
						}

					}
				}
			}

		} else if len(edge.Node.IsSourceOf.Edges) > 0 {
			for _, sourceOf := range edge.Node.IsSourceOf.Edges {
				names = append(names, sourceOf.Node.Name)
				labelMap[sourceOf.Node.Name] = sourceOf.Node.Label
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
		return "", "", errors.Wrap(err, "failed to check default categorisation")
	}

	return defaultCat.Variable, labelMap[defaultCat.Variable], nil
}
