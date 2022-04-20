package cache

import (
	"io/ioutil"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
)

type FilesystemCacheBackend struct {
	path string
}

func (b *FilesystemCacheBackend) Initialize(conf config.Backend) {
	log.Debug().Msg("initializing filesystem schema cache backend")
	b.path = conf.Path
	// No-op
}

func (b *FilesystemCacheBackend) GetRemote(schema string) (contents []byte, err error) {
	schemaLocation := filepath.Join(b.path, schema)
	content, err := ioutil.ReadFile(schemaLocation)
	if err != nil {
		log.Error().Stack().Err(err).Msg("could not get schema from filesystem schema cache backend: " + schemaLocation)
		return nil, err
	}
	return content, nil
}

func (b *FilesystemCacheBackend) Close() {
	log.Debug().Msg("closing filesystem schema cache backend")
	// No-op
}
