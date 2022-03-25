package service

import (
	"context"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/config"
	kafka "github.com/ONSdigital/dp-kafka/v3"
)

// GetKafkaProducer creates a Kafka producer. Currently just for POST filters/{id}/submit
func GetKafkaProducer(ctx context.Context, cfg *config.Config) (kafka.IProducer, error) {
	pConfig := &kafka.ProducerConfig{
		BrokerAddrs: cfg.KafkaConfig.Addr,
		// TODO: Right default?
		Topic:             cfg.KafkaConfig.ExportStartTopic,
		MinBrokersHealthy: &cfg.KafkaConfig.ProducerMinBrokersHealthy,
		KafkaVersion:      &cfg.KafkaConfig.Version,
		MaxMessageBytes:   &cfg.KafkaConfig.MaxBytes,
	}
	// TODOwhat should the types be here really ?
	// if cfg.KafkaConfig.SecProtocol == config.KafkaConfig.TLSProtocolFlag {
	if false {
		pConfig.SecurityConfig = kafka.GetSecurityConfig(
			cfg.KafkaConfig.SecCACerts,
			cfg.KafkaConfig.SecClientCert,
			cfg.KafkaConfig.SecClientKey,
			cfg.KafkaConfig.SecSkipVerify,
		)
	}
	return kafka.NewProducer(ctx, pConfig)
}
