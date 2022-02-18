package cache

import (
	"net/url"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/gosnowplow/pkg/config"
	h "github.com/silverton-io/gosnowplow/pkg/http"
)

type HttpSchemaCacheBackend struct {
	protocol string
	host     string
	path     string
}

func (b *HttpSchemaCacheBackend) Initialize(conf config.SchemaCacheBackend) {
	log.Debug().Msg("initializing http schema cache backend")
	b.protocol = conf.Type
	b.host = conf.Host // FIXME! String trailing / if it's present (or validate it upstream)
	b.path = conf.Path // FIXME! Strip leading / if it's present (or validate it upstream)
	// Auth? TBD
}

func (b *HttpSchemaCacheBackend) GetRemote(schema string) (contents []byte, err error) {
	schemaLocation, _ := url.Parse(b.protocol + "://" + b.host + "/" + b.path + "/" + schema) // FIXME!! There's gotta be a better way here.
	content, err := h.Get(*schemaLocation)
	if err != nil {
		log.Error().Stack().Err(err).Msg("could not get schema from http schema cache backend")
		return nil, err
	}
	return content, nil
}

func (b *HttpSchemaCacheBackend) Close() {
	log.Debug().Msg("closing http schema cache backend")
	// Knock off auth tokens? TBD
}
