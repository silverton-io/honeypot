package snowplow

import (
	"context"
	"io/ioutil"
	"sync/atomic"

	"github.com/gin-gonic/gin"
	"github.com/silverton-io/gosnowplow/pkg/cache"
	"github.com/silverton-io/gosnowplow/pkg/event"
	"github.com/silverton-io/gosnowplow/pkg/forwarder"
	"github.com/silverton-io/gosnowplow/pkg/http"
	"github.com/silverton-io/gosnowplow/pkg/response"
	"github.com/silverton-io/gosnowplow/pkg/tele"
	"github.com/tidwall/gjson"
)

func RedirectHandler(forwarder forwarder.Forwarder, cache *cache.SchemaCache, meta *tele.Meta) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ctx := context.Background()
		mappedParams := http.MapParams(c)
		e := BuildEventFromMappedParams(c, mappedParams)
		isValid, validationError, schema := validateEvent(e, cache)
		setEventMetadataFields(&e, schema)
		if isValid {
			forwarder.PublishValid(ctx, e)
			atomic.AddInt64(&meta.ValidSnowplowEventsProcessed, int64(1))
		} else {
			iE := event.InvalidEvent{
				ValidationError: &validationError,
				Event:           &e,
			}
			forwarder.PublishInvalid(ctx, iE)
			atomic.AddInt64(&meta.InvalidSnowplowEventsProcessed, int64(1))
		}
		redirectUrl, _ := c.GetQuery("u")
		c.Redirect(302, redirectUrl)
	}
	return gin.HandlerFunc(fn)
}

func GetHandler(forwarder forwarder.Forwarder, cache *cache.SchemaCache, meta *tele.Meta) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ctx := context.Background()
		mappedParams := http.MapParams(c)
		e := BuildEventFromMappedParams(c, mappedParams)
		isValid, validationError, schema := validateEvent(e, cache)
		setEventMetadataFields(&e, schema)
		if isValid {
			forwarder.PublishValid(ctx, e)
			atomic.AddInt64(&meta.ValidSnowplowEventsProcessed, int64(1))
		} else {
			iE := event.InvalidEvent{
				ValidationError: &validationError,
				Event:           &e,
			}
			forwarder.PublishInvalid(ctx, iE)
			atomic.AddInt64(&meta.InvalidSnowplowEventsProcessed, int64(1))
		}
		c.JSON(200, response.Ok)
	}
	return gin.HandlerFunc(fn)
}

func PostHandler(forwarder forwarder.Forwarder, cache *cache.SchemaCache, meta *tele.Meta) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ctx := context.Background()
		body, _ := ioutil.ReadAll(c.Request.Body) // FIXME! Handle errs here
		payloadData := gjson.GetBytes(body, "data")
		var validEvents []interface{}
		var invalidEvents []interface{}
		for _, e := range payloadData.Array() {
			e := BuildEventFromMappedParams(c, e.Value().(map[string]interface{}))
			isValid, validationError, schema := validateEvent(e, cache)
			setEventMetadataFields(&e, schema)
			if isValid {
				validEvents = append(validEvents, e)
			} else {
				iE := event.InvalidEvent{
					ValidationError: &validationError,
					Event:           &e,
				}
				invalidEvents = append(invalidEvents, iE)
			}
		}

		forwarder.BatchPublishValid(ctx, validEvents)
		atomic.AddInt64(&meta.ValidSnowplowEventsProcessed, int64(len(validEvents)))
		forwarder.BatchPublishInvalid(ctx, invalidEvents)
		atomic.AddInt64(&meta.InvalidSnowplowEventsProcessed, int64(len(invalidEvents)))
		c.JSON(200, response.Ok)
	}
	return gin.HandlerFunc(fn)
}
