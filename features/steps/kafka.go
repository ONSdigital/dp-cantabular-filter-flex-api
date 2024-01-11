package steps

import (
	"context"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/config"
	kafka "github.com/ONSdigital/dp-kafka/v4"
)

// GenerateKafkaConsumer produces a consumer for specific steps.
// Currently only for PUT filters/{id}/submit
func GenerateKafkaConsumer(ctx context.Context) (*kafka.ConsumerGroup, error) {
	kafkaConfig := config.KafkaConfig{
		Addr:                      []string{"kafka-1:9092"},
		ConsumerMinBrokersHealthy: 1,
		ProducerMinBrokersHealthy: 1,
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
		ExportStartGroup:          ComponentTestGroup,
		TLSProtocolFlag:           false,
	}

	kafkaOffset := kafka.OffsetOldest
	consumer, err := kafka.NewConsumerGroup(
		ctx,
		&kafka.ConsumerGroupConfig{
			BrokerAddrs:       kafkaConfig.Addr,
			Topic:             kafkaConfig.ExportStartTopic,
			GroupName:         kafkaConfig.ExportStartGroup,
			MinBrokersHealthy: &(kafkaConfig.ConsumerMinBrokersHealthy),
			KafkaVersion:      &(kafkaConfig.Version),
			Offset:            &(kafkaOffset),
		},
	)
	if err != nil {
		return nil, err
	}

	if err := consumer.Start(); err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	consumer.LogErrors(ctx)
	return consumer, nil
}
