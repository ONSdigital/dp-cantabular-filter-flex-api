package config

import (
	"time"

	mongo "github.com/ONSdigital/dp-mongodb/v3/mongodb"
	"github.com/kelseyhightower/envconfig"
)

// KafkaTLSProtocolFlag informs service to use TLS protocol for kafka
const KafkaTLSProtocolFlag = "TLS"

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
	PopulationTypesAPIURL        string        `envconfig:"POPULATION_TYPES_API_URL"`
	MetadataAPIURL               string        `envconfig:"CANTABULAR_METADATA_API_URL"`
	FilterAPIURL                 string        `envconfig:"FILTER_API_URL"`
	FiltersCollection            string        `envconfig:"FILTERS_COLLECTION"`
	FilterOutputsCollection      string        `envconfig:"FILTER_OUTPUTS_COLLECTION"`
	EnablePrivateEndpoints       bool          `envconfig:"ENABLE_PRIVATE_ENDPOINTS"`
	EnablePermissionsAuth        bool          `envconfig:"ENABLE_PERMISSIONS_AUTH"`
	CantabularHealthcheckEnabled bool          `envconfig:"CANTABULAR_HEALTHCHECK_ENABLED"`
	ServiceAuthToken             string        `envconfig:"SERVICE_AUTH_TOKEN"`
	ZebedeeURL                   string        `envconfig:"ZEBEDEE_URL"`
	DatasetOptionsWorkers        int           `envconfig:"DATASET_OPTIONS_WORKERS"`
	DatasetOptionsBatchSize      int           `envconfig:"DATASET_OPTIONS_BATCH_SIZE"`
	OTExporterOTLPEndpoint       string        `envconfig:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	OTServiceName                string        `envconfig:"OTEL_SERVICE_NAME"`
	OTBatchTimeout               time.Duration `envconfig:"OTEL_BATCH_TIMEOUT"`
	Mongo                        mongo.MongoDriverConfig
	KafkaConfig                  KafkaConfig
}

type KafkaConfig struct {
	Addr                      []string `envconfig:"KAFKA_ADDR"                            json:"-"`
	ConsumerMinBrokersHealthy int      `envconfig:"KAFKA_CONSUMER_MIN_BROKERS_HEALTHY"`
	ProducerMinBrokersHealthy int      `envconfig:"KAFKA_PRODUCER_MIN_BROKERS_HEALTHY"`
	Version                   string   `envconfig:"KAFKA_VERSION"`
	OffsetOldest              bool     `envconfig:"KAFKA_OFFSET_OLDEST"`
	NumWorkers                int      `envconfig:"KAFKA_NUM_WORKERS"`
	MaxBytes                  int      `envconfig:"KAFKA_MAX_BYTES"`
	SecProtocol               string   `envconfig:"KAFKA_SEC_PROTO"`
	SecCACerts                string   `envconfig:"KAFKA_SEC_CA_CERTS"`
	SecClientKey              string   `envconfig:"KAFKA_SEC_CLIENT_KEY"                  json:"-"`
	SecClientCert             string   `envconfig:"KAFKA_SEC_CLIENT_CERT"`
	SecSkipVerify             bool     `envconfig:"KAFKA_SEC_SKIP_VERIFY"`
	ExportStartTopic          string   `envconfig:"KAFKA_TOPIC_CANTABULAR_EXPORT_START"`
	ExportStartGroup          string   `envconfig:"KAFKA_GROUP_CANTABULAR_EXPORT_START"`
	TLSProtocolFlag           bool     `envconfig:"TLS_PROTOCOL_FLAG"`
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
		PopulationTypesAPIURL:        "http://localhost:27300",
		CantabularURL:                "http://localhost:8491",
		CantabularExtURL:             "http://localhost:8492",
		FilterAPIURL:                 "http://localhost:22100",
		MetadataAPIURL:               "http://dp-cantabular-metadata-service:8493",
		FiltersCollection:            "filters",
		FilterOutputsCollection:      "filterOutputs",
		EnablePrivateEndpoints:       false,
		EnablePermissionsAuth:        true,
		CantabularHealthcheckEnabled: false,
		ServiceAuthToken:             "",
		ZebedeeURL:                   "http://localhost:8082",
		DatasetOptionsWorkers:        2,
		DatasetOptionsBatchSize:      20,
		OTExporterOTLPEndpoint:       "localhost:4317",
		OTServiceName:                "dp-cantabular-filter-flex-api",
		OTBatchTimeout:               5 * time.Second,
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
		KafkaConfig: KafkaConfig{
			Addr:                      []string{"localhost:9092", "localhost:9093", "localhost:9094"},
			ConsumerMinBrokersHealthy: 1,
			ProducerMinBrokersHealthy: 2,
			Version:                   "1.0.2",
			OffsetOldest:              true,
			NumWorkers:                1,
			MaxBytes:                  2000000,
			SecProtocol:               "",
			SecCACerts:                "",
			SecClientKey:              "",
			SecClientCert:             "",
			SecSkipVerify:             false,
			ExportStartTopic:          "cantabular-export-start",
			ExportStartGroup:          "cantabular-export-start-group",
			TLSProtocolFlag:           false,
		},
	}

	return cfg, envconfig.Process("", cfg)
}
