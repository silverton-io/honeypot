package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/silverton-io/honeypot/pkg/protocol"
	"github.com/silverton-io/honeypot/pkg/response"
	"github.com/silverton-io/honeypot/pkg/validator"
)

func GenericHandler(p EventHandlerParams) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ctx := context.Background()
		envelopes := buildGenericEnvelopesFromRequest(c, *p.Config)
		validEvents, invalidEvents := validator.BifurcateAndAnnotate(envelopes, p.Cache)
		p.Sink.BatchPublishValidAndInvalid(ctx, protocol.GENERIC, validEvents, invalidEvents, p.Meta)
		c.JSON(200, response.Ok)
	}
	return gin.HandlerFunc(fn)
}