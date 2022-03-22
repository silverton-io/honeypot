package sink

import (
	"context"
	"encoding/json"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/envelope"
)

type FileSink struct {
	validFile   string
	invalidFile string
}

func (s *FileSink) Initialize(conf config.Sink) {
	log.Debug().Msg("initializing file sink")
	s.validFile = conf.ValidFile
	s.invalidFile = conf.InvalidFile
}

func (s *FileSink) batchPublish(ctx context.Context, filePath string, envelopes []envelope.Envelope) {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	if err != nil {
		log.Error().Stack().Err(err).Msg("could not open file")
	}
	for _, envelope := range envelopes {
		log.Debug().Msg("writing envelope to file " + filePath)
		b, _ := json.Marshal(envelope)
		newline := []byte("\n")
		b = append(b, newline...)
		f.Write(b)
	}
}

func (s *FileSink) BatchPublishValid(ctx context.Context, envelopes []envelope.Envelope) {
	s.batchPublish(ctx, s.validFile, envelopes)
}

func (s *FileSink) BatchPublishInvalid(ctx context.Context, envelopes []envelope.Envelope) {
	s.batchPublish(ctx, s.invalidFile, envelopes)
}

func (s *FileSink) Close() {
	log.Debug().Msg("closing file sink")
}