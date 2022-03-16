package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/protocol"
	"github.com/silverton-io/honeypot/pkg/response"
	"github.com/silverton-io/honeypot/pkg/validator"
)

func SnowplowHandler(p EventHandlerParams) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ctx := context.Background()
		envelopes := buildSnowplowEnvelopesFromRequest(c, *p.Config)
		validEvents, invalidEvents := validator.BifurcateAndAnnotate(envelopes, p.Cache)
		p.Sink.BatchPublishValidAndInvalid(ctx, protocol.SNOWPLOW, validEvents, invalidEvents, p.Meta)
		if c.Request.Method == http.MethodGet {
			redirectUrl, _ := c.GetQuery("u")
			if redirectUrl != "" && p.Config.Snowplow.OpenRedirectsEnabled {
				log.Info().Msg("redirecting to " + redirectUrl)
				c.Redirect(http.StatusFound, redirectUrl)
			}
		}
		c.JSON(http.StatusOK, response.Ok)
	}
	return gin.HandlerFunc(fn)
}