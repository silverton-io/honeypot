package sink

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type MaterializeSink struct {
	id           *uuid.UUID
	name         string
	gormDb       *gorm.DB
	validTable   string
	invalidTable string
}

func (s *MaterializeSink) Id() *uuid.UUID {
	return s.id
}

func (s *MaterializeSink) Name() string {
	return s.name
}

func (s *MaterializeSink) Initialize(conf config.Sink) error {
	log.Debug().Msg("initializing postgres sink")
	id := uuid.New()
	s.id, s.name = &id, conf.Name
	connString := generatePgDsn(conf)
	gormDb, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("could not open materialize connection")
		return err
	}
	s.gormDb = gormDb
	s.validTable, s.invalidTable = conf.ValidTable, conf.InvalidTable
	for _, tbl := range []string{s.validTable, s.invalidTable} {
		tblExists := s.gormDb.Migrator().HasTable(tbl)
		if !tblExists {
			log.Debug().Msg(tbl + " table doesn't exist - ensuring")
			err = s.gormDb.Table(tbl).AutoMigrate(&envelope.Envelope{})
			if err != nil {
				log.Error().Err(err).Msg("could not auto migrate table")
				return err
			}
			// NOTE! This is a hacky workaround so that the same gorm struct tag of "json" can be used, but "jsonb" is used for pg.
			for _, col := range []string{"event_metadata", "validation_error", "payload"} {
				alterStmt := "alter table " + tbl + " rename column " + col + " set data type jsonb using " + col + "::jsonb;"
				log.Debug().Msg("ensuring jsonb columns via: " + alterStmt)
				s.gormDb.Exec(alterStmt)
			}
		} else {
			log.Debug().Msg(tbl + " table already exists - not ensuring")
		}
	}
	return nil
}

func (s *MaterializeSink) BatchPublishValid(ctx context.Context, envelopes []envelope.Envelope) {
	s.gormDb.Table(s.validTable).Create(envelopes)
}

func (s *MaterializeSink) BatchPublishInvalid(ctx context.Context, envelopes []envelope.Envelope) {
	s.gormDb.Table(s.invalidTable).Create(envelopes)
}

func (s *MaterializeSink) Close() {
	log.Debug().Msg("closing postgres sink")
	db, _ := s.gormDb.DB()
	db.Close()
}