package sink

import (
	"context"
	"net/url"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/request"
)

type HttpSink struct {
	id         *uuid.UUID
	name       string
	validUrl   url.URL
	invalidUrl url.URL
}

func (s *HttpSink) Id() *uuid.UUID {
	return s.id
}

func (s *HttpSink) Name() string {
	return s.name
}

func (s *HttpSink) Initialize(conf config.Sink) error {
	log.Debug().Msg("initializing http sink")
	vUrl, vErr := url.Parse(conf.ValidUrl)
	invUrl, invErr := url.Parse(conf.InvalidUrl)
	if vErr != nil {
		log.Debug().Stack().Err(vErr).Msg("validUrl is not a valid url")
		return vErr
	}
	if invErr != nil {
		log.Debug().Stack().Err(invErr).Msg("invalidUrl is not a valid url")
		return invErr
	}
	id := uuid.New()
	s.id, s.name = &id, conf.Name
	s.validUrl, s.invalidUrl = *vUrl, *invUrl
	return nil
}

func (s *HttpSink) BatchPublishValid(ctx context.Context, validEnvelopes []envelope.Envelope) {
	request.PostEnvelopes(s.validUrl, validEnvelopes)
}

func (s *HttpSink) BatchPublishInvalid(ctx context.Context, invalidEnvelopes []envelope.Envelope) {
	request.PostEnvelopes(s.invalidUrl, invalidEnvelopes)
}

func (s *HttpSink) Close() {
	log.Debug().Msg("closing http sink") // no-op
}
