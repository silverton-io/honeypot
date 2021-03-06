package envelope

import (
	"github.com/gin-gonic/gin"
	"github.com/silverton-io/honeypot/pkg/util"
)

const HTTP_HEADERS_CONTEXT string = "io.silverton/honeypot/internal/contexts/httpHeaders/v1.0.json"

func buildContextsFromRequest(c *gin.Context) map[string]interface{} {
	headers := util.HttpHeadersToMap(c)
	context := map[string]interface{}{
		HTTP_HEADERS_CONTEXT: headers,
	}
	return context
}
