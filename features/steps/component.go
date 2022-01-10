package steps

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/config"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/service"
	componenttest "github.com/ONSdigital/dp-component-test"
	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/log.go/v2/log"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/maxcnunes/httpfake"
)

const (
	ComponentTestGroup    = "component-test" // kafka group name for the component test consumer
	DrainTopicTimeout     = 1 * time.Second  // maximum time to wait for a topic to be drained
	DrainTopicMaxMessages = 1000             // maximum number of messages that will be drained from a topic
	MinioCheckRetries     = 3                // maximum number of retires to validate that a file is present in minio
	WaitEventTimeout      = 5 * time.Second  // maximum time that the component test consumer will wait for a kafka event
)

var (
	BuildTime string = "1625046891"
	GitCommit string = "7434fe334d9f51b7239f978094ea29d10ac33b16"
	Version   string = ""
)

type Component struct {
	componenttest.ErrorFeature
	DatasetAPI       *httpfake.HTTPFake
	CantabularSrv    *httpfake.HTTPFake
	CantabularAPIExt *httpfake.HTTPFake
	S3Downloader     *s3manager.Downloader
	producer         kafka.IProducer
	consumer         kafka.IConsumerGroup
	errorChan        chan error
	svc              *service.Service
	cfg              *config.Config
	wg               *sync.WaitGroup
	signals          chan os.Signal
	waitEventTimeout time.Duration
	testETag         string
	ctx              context.Context
}

func NewComponent() *Component {
	return &Component{
		errorChan:        make(chan error),
		DatasetAPI:       httpfake.New(),
		CantabularSrv:    httpfake.New(),
		CantabularAPIExt: httpfake.New(),
		wg:               &sync.WaitGroup{},
		waitEventTimeout: WaitEventTimeout,
		testETag:         "13c7791bafdbaaf5e6660754feb1a58cd6aaa892",
		ctx:              context.Background(),
	}
}

// initService initialises the server, the mocks and waits for the dependencies to be ready
func (c *Component) initService(ctx context.Context) error {
	// register interrupt signals
	c.signals = make(chan os.Signal, 1)
	signal.Notify(c.signals, os.Interrupt, syscall.SIGTERM)

	// Read config
	cfg, err := config.Get()
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	log.Info(ctx, "config used by component tests", log.Data{"cfg": cfg})

	// producer for triggering test events
	if c.producer, err = kafka.NewProducer(
		ctx,
		&kafka.ProducerConfig{
			BrokerAddrs:       cfg.KafkaConfig.Addr,
			Topic:             cfg.KafkaConfig.ExportStartTopic,
			MinBrokersHealthy: &cfg.KafkaConfig.ProducerMinBrokersHealthy,
			KafkaVersion:      &cfg.KafkaConfig.Version,
			MaxMessageBytes:   &cfg.KafkaConfig.MaxBytes,
		},
	); err != nil {
		return fmt.Errorf("error creating kafka producer: %w", err)
	}

	// Create service and initialise it
	c.svc = service.New()
	if err = c.svc.Init(ctx, cfg, BuildTime, GitCommit, Version); err != nil {
		return fmt.Errorf("unexpected service Init error in NewComponent: %w", err)
	}

	c.cfg = cfg

	// wait for producer to be initialised
	<-c.producer.Channels().Initialised
	log.Info(ctx, "component-test kafka producer initialised")

	return nil
}

// startService starts the service under test and blocks until an error or an os interrupt is received.
// Then it closes the service (graceful shutdown)
func (c *Component) startService(ctx context.Context) {
	defer c.wg.Done()
	c.svc.Start(ctx, c.errorChan)

	select {
	case err := <-c.errorChan:
		err = fmt.Errorf("service error received: %w", err)
		c.svc.Close(ctx)
		panic(fmt.Errorf("unexpected error received from errorChan: %w", err))
	case sig := <-c.signals:
		log.Info(ctx, "os signal received", log.Data{"signal": sig})
	}
	if err := c.svc.Close(ctx); err != nil {
		panic(fmt.Errorf("unexpected error during service graceful shutdown: %w", err))
	}
}

// drainTopic drains the provided topic and group of any residual messages between scenarios.
// Prevents future tests failing if previous tests fail unexpectedly and
// leave messages in the queue.
//
// A temporary batch consumer is used, that is created and closed within this func
// A maximum of DrainTopicMaxMessages messages will be drained from the provided topic and group.
//
// This method accepts a waitGroup pionter. If it is not nil, it will wait for the topic to be drained
// in a new go-routine, which will be added to the waitgroup. If it is nil, execution will be blocked
// until the topic is drained (or time out expires)
func (c *Component) drainTopic(ctx context.Context, topic, group string, wg *sync.WaitGroup) error {
	msgs := []kafka.Message{}

	defer func() {
		log.Info(ctx, "drained topic", log.Data{
			"len":      len(msgs),
			"messages": msgs,
			"topic":    topic,
			"group":    group,
		})
	}()

	kafkaOffset := kafka.OffsetOldest
	batchSize := DrainTopicMaxMessages
	batchWaitTime := DrainTopicTimeout
	consumer, err := kafka.NewConsumerGroup(
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
	consumer.RegisterBatchHandler(
		ctx,
		func(ctx context.Context, batch []kafka.Message) error {
			defer close(drained)
			msgs = append(msgs, batch...)
			return nil
		},
	)

	// start consumer group
	consumer.Start()

	// start kafka logging go-routines
	consumer.LogErrors(ctx)

	// waitUntilDrained is a func that will wait until the batch is consumed or the timeout expires
	// (with 100 ms of extra time to allow any in-flight drain)
	waitUntilDrained := func() {
		select {
		case <-time.After(DrainTopicTimeout + 100*time.Millisecond):
		case <-drained:
		}

		consumer.Close(ctx)
		<-consumer.Channels().Closed
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

// Close kills the application under test, and then it shuts down the testing consumer and producer.
func (c *Component) Close() {
	// kill application
	c.signals <- os.Interrupt

	// wait for graceful shutdown to finish (or timeout)
	c.wg.Wait()

	// stop listening to consumer, waiting for any in-flight message to be committed
	c.consumer.StopAndWait()

	// close producer
	if err := c.producer.Close(c.ctx); err != nil {
		log.Error(c.ctx, "error closing kafka producer", err)
	}

	// close consumer
	if err := c.consumer.Close(c.ctx); err != nil {
		log.Error(c.ctx, "error closing kafka consumer", err)
	}

	// drain topics in parallel
	wg := &sync.WaitGroup{}
	if err := c.drainTopic(c.ctx, c.cfg.KafkaConfig.CsvCreatedTopic, ComponentTestGroup, wg); err != nil {
		log.Error(c.ctx, "error draining topic", err, log.Data{
			"topic": c.cfg.KafkaConfig.CsvCreatedTopic,
			"group": ComponentTestGroup,
		})
	}
	if err := c.drainTopic(c.ctx, c.cfg.KafkaConfig.ExportStartTopic, c.cfg.KafkaConfig.ExportStartGroup, wg); err != nil {
		log.Error(c.ctx, "error draining topic", err, log.Data{
			"topic": c.cfg.KafkaConfig.ExportStartTopic,
			"group": c.cfg.KafkaConfig.ExportStartGroup,
		})
	}
	wg.Wait()
}

// Reset re-initialises the service under test and the api mocks.
// Note that the service under test should not be started yet
// to prevent race conditions if it tries to call un-initialised dependencies (steps)
func (c *Component) Reset() error {
	if err := c.initService(c.ctx); err != nil {
		return fmt.Errorf("failed to initialise service: %w", err)
	}

	c.DatasetAPI.Reset()
	c.CantabularSrv.Reset()
	c.CantabularAPIExt.Reset()

	return nil
}
