package sink

import (
	"context"
	"sync"
	"sync/atomic"

	awsconf "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/firehose"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/gosnowplow/pkg/config"
	"github.com/silverton-io/gosnowplow/pkg/input"
	"github.com/silverton-io/gosnowplow/pkg/tele"
)

type KinesisFirehoseSink struct {
	client              *firehose.Client
	validEventsStream   string
	invalidEventsStream string
}

func (s *KinesisFirehoseSink) Initialize(conf config.Sink) {
	ctx := context.Background()
	cfg, _ := awsconf.LoadDefaultConfig(ctx)
	client := firehose.NewFromConfig(cfg)
	s.client, s.validEventsStream, s.invalidEventsStream = client, conf.ValidEventTopic, conf.InvalidEventTopic
}

func (s *KinesisFirehoseSink) batchPublish(ctx context.Context, stream string, events []interface{}) {
	var wg sync.WaitGroup
}

func (s *KinesisFirehoseSink) batchPublishValid(ctx context.Context, events []interface{}) {
	s.batchPublish(ctx, s.validEventsStream, events)
}

func (s *KinesisFirehoseSink) batchPublishInvalid(ctx context.Context, events []interface{}) {
	s.batchPublish(ctx, s.invalidEventsStream, events)
}

func (s *KinesisFirehoseSink) BatchPublishValidAndInvalid(ctx context.Context, inputType string, validEvents []interface{}, invalidEvents []interface{}, meta *tele.Meta) {
	var validCounter *int64
	var invalidCounter *int64
	switch inputType {
	case input.GENERIC_INPUT:
		validCounter = &meta.ValidGenericEventsProcessed
		invalidCounter = &meta.InvalidGenericEventsProcessed
	case input.CLOUDEVENTS_INPUT:
		validCounter = &meta.ValidCloudEventsProcessed
		invalidCounter = &meta.InvalidCloudEventsProcessed
	default:
		validCounter = &meta.ValidSnowplowEventsProcessed
		invalidCounter = &meta.InvalidSnowplowEventsProcessed
	}
	// Publish
	s.batchPublishValid(ctx, validEvents)
	s.batchPublishInvalid(ctx, invalidEvents)
	// Increment global metadata counters
	atomic.AddInt64(validCounter, int64(len(validEvents)))
	atomic.AddInt64(invalidCounter, int64(len(invalidEvents)))
}

func (s *KinesisFirehoseSink) Close() {
	log.Debug().Msg("closing kinesis firehose sink client")
}