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
	BindAddr                   string        `envconfig:"BIND_ADDR"`
	GracefulShutdownTimeout    time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval        time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	ComponentTestUseLogFile    bool          `envconfig:"COMPONENT_TEST_USE_LOG_FILE"`
	DatasetAPIURL              string        `envconfig:"DATASET_API_URL"`
	Kafka                      KafkaConfig
	Mongo                      mongo.MongoDriverConfig
	FiltersCollection          string `envconfig:"FILTERS_COLLECTION"`
	FilterOutputsCollection    string `envconfig:"FILTER_OUTPUTS_COLLECTION"`
	EnablePrivateEndpoints     bool   `envconfig:"ENABLE_PRIVATE_ENDPOINTS"`
	EnablePermissionsAuth      bool   `envconfig:"ENABLE_PERMISSIONS_AUTH"`
	ZebedeeURL                 string `envconfig:"ZEBEDEE_URL"`
}

// KafkaConfig contains the config required to connect to Kafka
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
	ExportStartGroup          string   `envconfig:"KAFKA_GROUP_CANTABULAR_EXPORT_START"`
	ExportStartTopic          string   `envconfig:"KAFKA_TOPIC_CANTABULAR_EXPORT_START"`
	CsvCreatedTopic           string   `envconfig:"KAFKA_TOPIC_CSV_CREATED"`
}

var cfg *Config

// Get returns the default config with any modifications through environment
// variables
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
		BindAddr:                   ":27100",
		GracefulShutdownTimeout:    5 * time.Second,
		HealthCheckInterval:        30 * time.Second,
		HealthCheckCriticalTimeout: 90 * time.Second,
		ComponentTestUseLogFile:    false,
		DatasetAPIURL:              "localhost:8082",
		FiltersCollection:          "censusFilters",
		FilterOutputsCollection:    "censusFilterOutputs",
		EnablePrivateEndpoints:     false,
		EnablePermissionsAuth:      true,
		ZebedeeURL:                 "http://localhost:8082",
		Kafka: KafkaConfig{
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
			ExportStartGroup:          "dp-cantabular-csv-exporter",
			ExportStartTopic:          "cantabular-export-start",
			CsvCreatedTopic:           "cantabular-csv-created",
		},
		Mongo: mongo.MongoDriverConfig{
			ClusterEndpoint: "localhost:27017",
			Username:        "",
			Password:        "",
			Database:        "filters",
			Collections: map[string]string{
				"censusFilters":       "censusFilters",
				"censusFilterOutputs": "censuFilterOutputs",
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
