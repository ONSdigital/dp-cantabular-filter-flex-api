package steps

import (
	"context"
	"fmt"
	"sync"
	"time"

	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/log.go/v2/log"
)

func (c *Component) drainTopic(ctx context.Context, topic, group string, wg *sync.WaitGroup) error {
	msgs := []kafka.Message{}

	kafkaOffset := kafka.OffsetOldest
	batchSize := DrainTopicMaxMessages
	batchWaitTime := DrainTopicTimeout
	drainer, err := kafka.NewConsumerGroup(
		ctx,
		&kafka.ConsumerGroupConfig{
			BrokerAddrs:   c.cfg.KafkaConfig.Addr,
			Topic:         topic,
			GroupName:     group,
			KafkaVersion:  &c.cfg.KafkaConfig.Version,
			Offset:        &kafkaOffset,
			BatchSize:     &batchSize,
			BatchWaitTime: &batchWaitTime,
		},
	)
	if err != nil {
		return fmt.Errorf("error creating kafka consumer to drain topic: %w", err)
	}

	// register batch handler with 'drained channel'
	drained := make(chan struct{})
	if err := drainer.RegisterBatchHandler(
		ctx,
		func(ctx context.Context, batch []kafka.Message) error {
			defer close(drained)
			msgs = append(msgs, batch...)
			return nil
		},
	); err != nil {
		return fmt.Errorf("error creating kafka drainer: %w", err)
	}

	// start drainer consumer group
	if err := drainer.Start(); err != nil {
		log.Error(ctx, "error starting kafka drainer", err)
	}

	// start kafka logging go-routines
	drainer.LogErrors(ctx)

	// waitUntilDrained is a func that will wait until the batch is consumed or the timeout expires
	// (with 100 ms of extra time to allow any in-flight drain)
	waitUntilDrained := func() {
		drainer.StateWait(kafka.Consuming)
		log.Info(ctx, "drainer is consuming", log.Data{"topic": topic, "group": group})

		select {
		case <-time.After(DrainTopicTimeout + 100*time.Millisecond):
			log.Info(ctx, "drain timeout has expired (no messages drained)")
		case <-drained:
			log.Info(ctx, "message(s) have been drained")
		}

		defer func() {
			log.Info(ctx, "drained topic", log.Data{
				"len":      len(msgs),
				"messages": msgs,
				"topic":    topic,
				"group":    group,
			})
		}()

		if err := drainer.Close(ctx); err != nil {
			log.Warn(ctx, "error closing drain consumer", log.Data{"err": err})
		}

		<-drainer.Channels().Closed
		log.Info(ctx, "drainer is closed")
	}

	// sync wait if wg is not provided
	if wg == nil {
		waitUntilDrained()
		return nil
	}

	// async wait if wg is provided
	wg.Add(1)
	go func() {
		defer wg.Done()
		waitUntilDrained()
	}()
	return nil
}
