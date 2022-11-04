package mock

import (
	"context"
	"fmt"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/config"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/service"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
)

type MockDatasetAPIClient struct {
	client service.DatasetAPIClient
}

func NewMockDatasetAPIClient() *MockDatasetAPIClient {
	cfg, _ := config.Get()
	mc := &MockDatasetAPIClient{}
	mc.client = dataset.NewAPIClient(cfg.DatasetAPIURL)
	return mc
}

func (c *MockDatasetAPIClient) GetVersion(ctx context.Context, userAuthToken, svcAuthToken, downloadSvcAuthToken, collectionID, datasetID, edition, version string) (dataset.Version, error) {
	return c.client.GetVersion(ctx, userAuthToken, svcAuthToken, downloadSvcAuthToken, collectionID, datasetID, edition, version)
}

func (c *MockDatasetAPIClient) GetOptionsInBatches(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, id, edition, version, dimension string, batchSize, maxWorkers int) (dataset.Options, error) {
	return c.client.GetOptionsInBatches(ctx, userAuthToken, serviceAuthToken, collectionID, id, edition, version, dimension, batchSize, maxWorkers)
}

func (c *MockDatasetAPIClient) GetMetadataURL(id, edition, version string) string {
	return fmt.Sprintf("%s/datasets/%s/editions/%s/versions/%s/metadata", "http://localhost:9999", id, edition, version)
}

func (c *MockDatasetAPIClient) Checker(ctx context.Context, state *healthcheck.CheckState) error {
	return c.client.Checker(ctx, state)
}
