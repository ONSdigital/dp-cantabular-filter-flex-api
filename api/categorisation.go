package api

import (
	"context"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/cantabularmetadata"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	"github.com/pkg/errors"
)

// RetrieveDefaultCategorisation takes dimension, returns categorisations, and checks if any are default.
// if so, it returns the relevant information. If there are no default categorisations, it returns empty string
// for default categorisation, and the original dimension name and label to persist instead.
// returns (finalDimension, finalLabel, finalCategorisation, error)
func (api *API) RetrieveDefaultCategorisation(dimension *model.Dimension, datasetName string) (string, string, string, error) {
	ctx := context.Background()
	labelMap := make(map[string]string)
	cats, err := api.ctblr.GetCategorisations(ctx, cantabular.GetCategorisationsRequest{
		Dataset:  datasetName,
		Variable: dimension.Name,
	})
	if err != nil {
		return "", "", "", errors.Wrap(err, "failed to check default categorisation")
	}

	names := make([]string, 0)

	for edgeIndex := range cats.Dataset.Variables.Edges {
		edge := cats.Dataset.Variables.Edges[edgeIndex]

		if len(edge.Node.MapFrom) > 0 {
			for mapFromIndex := range edge.Node.MapFrom {
				mapFrom := edge.Node.MapFrom[mapFromIndex]
				for mapFromEdgeIndex := range mapFrom.Edges {
					mappedSource := mapFrom.Edges[mapFromEdgeIndex]
					for mappedSourceEdgeIndex := range mappedSource.Node.IsSourceOf.Edges {
						mappedSourceEdge := mappedSource.Node.IsSourceOf.Edges[mappedSourceEdgeIndex]
						names = append(names, mappedSourceEdge.Node.Name)
						labelMap[mappedSourceEdge.Node.Name] = mappedSourceEdge.Node.Label
					}
				}
			}
		} else if len(edge.Node.IsSourceOf.Edges) > 0 {
			for i := range edge.Node.IsSourceOf.Edges {
				sourceOf := edge.Node.IsSourceOf.Edges[i]
				names = append(names, sourceOf.Node.Name)
				labelMap[sourceOf.Node.Name] = sourceOf.Node.Label
			}
		}
	}

	if len(names) == 0 {
		return "", "", "", errors.New("no categorisations received for variable")
	}

	defaultCat, err := api.metadata.GetDefaultClassification(ctx, cantabularmetadata.GetDefaultClassificationRequest{
		Dataset:   datasetName,
		Variables: names,
	})
	if err != nil {
		return "", "", "", errors.Wrap(err, "failed to check default categorisation")
	}

	if len(defaultCat.Variables) > 1 {
		return "", "", "", errors.New("more than 1 categorisation returned")
	}

	if len(defaultCat.Variables) == 0 {
		return dimension.Name, dimension.Label, "", nil
	}

	return defaultCat.Variables[0], labelMap[defaultCat.Variables[0]], defaultCat.Variables[0], nil
}
