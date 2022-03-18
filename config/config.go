package config

import (
	"time"

	mongo "github.com/ONSdigital/dp-mongodb/v3/mongodb"
	"github.com/kelseyhightower/envconfig"
)

// Config represents service configuration for dp-cantabular-filter-flex-api
type Config struct {
	BindAddr                     string        `envconfig:"BIND_ADDR"`
	GracefulShutdownTimeout      time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval          time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout   time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	DefaultRequestTimeout        time.Duration `envconfig:"DEFAULT_REQUEST_TIMEOUT"`
	DefaultMaximumLimit          int           `envconfig:"DEFAULT_MAXIMUM_LIMIT"`
	ComponentTestUseLogFile      bool          `envconfig:"COMPONENT_TEST_USE_LOG_FILE"`
	CantabularURL                string        `envconfig:"CANTABULAR_URL"`
	CantabularExtURL             string        `envconfig:"CANTABULAR_API_EXT_URL"`
	DatasetAPIURL                string        `envconfig:"DATASET_API_URL"`
	FiltersCollection            string        `envconfig:"FILTERS_COLLECTION"`
	FilterOutputsCollection      string        `envconfig:"FILTER_OUTPUTS_COLLECTION"`
	EnablePrivateEndpoints       bool          `envconfig:"ENABLE_PRIVATE_ENDPOINTS"`
	EnablePermissionsAuth        bool          `envconfig:"ENABLE_PERMISSIONS_AUTH"`
	CantabularHealthcheckEnabled bool          `envconfig:"CANTABULAR_HEALTHCHECK_ENABLED"`
	ServiceAuthToken             string        `envconfig:"SERVICE_AUTH_TOKEN"`
	ZebedeeURL                   string        `envconfig:"ZEBEDEE_URL"`
	Mongo                        mongo.MongoDriverConfig
}

var cfg *Config

// Get returns the default config with any modifications through environment
// variables
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
		BindAddr:                     ":27100",
		GracefulShutdownTimeout:      5 * time.Second,
		HealthCheckInterval:          30 * time.Second,
		HealthCheckCriticalTimeout:   90 * time.Second,
		DefaultRequestTimeout:        10 * time.Second,
		DefaultMaximumLimit:          500,
		ComponentTestUseLogFile:      false,
		DatasetAPIURL:                "http://localhost:22000",
		CantabularURL:                "http://localhost:8491",
		CantabularExtURL:             "http://localhost:8492",
		FiltersCollection:            "filters",
		FilterOutputsCollection:      "filterOutputs",
		EnablePrivateEndpoints:       false,
		EnablePermissionsAuth:        true,
		CantabularHealthcheckEnabled: false,
		ServiceAuthToken:             "",
		ZebedeeURL:                   "http://localhost:8082",
		Mongo: mongo.MongoDriverConfig{
			ClusterEndpoint: "localhost:27017",
			Username:        "",
			Password:        "",
			Database:        "filters",
			Collections: map[string]string{
				"filters":       "filters",
				"filterOutputs": "filterOutputs",
			},
			ReplicaSet:                    "",
			IsStrongReadConcernEnabled:    false,
			IsWriteConcernMajorityEnabled: true,
			ConnectTimeout:                5 * time.Second,
			QueryTimeout:                  15 * time.Second,
			TLSConnectionConfig: mongo.TLSConnectionConfig{
				IsSSL: false,
			},
		},
	}

	return cfg, envconfig.Process("", cfg)
}
