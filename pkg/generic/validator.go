package generic

import (
	"github.com/silverton-io/gosnowplow/pkg/cache"
	"github.com/silverton-io/gosnowplow/pkg/config"
	"github.com/silverton-io/gosnowplow/pkg/validator"
	"github.com/tidwall/gjson"
)

func validateEvent(event gjson.Result, schemaName string, cache *cache.SchemaCache, conf *config.Generic) (isValid bool, validationError validator.ValidationError, schema []byte) {
	if event.Value() == nil {
		errorType := "payload not present at " + conf.Payload.RootKey + "." + conf.Payload.DataKey
		validationError := validator.ValidationError{
			ErrorType: &errorType,
			Errors:    nil,
		}
		return false, validationError, nil
	}
	if schemaName == "" { // Event does not have schema associated - always invalid.
		errorType := "schema not present at " + conf.Payload.RootKey + "." + conf.Payload.SchemaKey
		validationError := validator.ValidationError{
			ErrorType: &errorType,
			Errors:    nil,
		}
		return false, validationError, nil
	}
	schemaExists, schemaContents := cache.Get(schemaName)
	// FIXME! What happens if the payload key doesn't exist?
	if !schemaExists { // Referenced schema is not present in the cache backend - always invalid
		errorType := "nonexistent schema"
		validationError := validator.ValidationError{
			ErrorType: &errorType,
			Errors:    nil,
		}
		return false, validationError, nil
	} else {
		isValid, validationError := validator.ValidatePayload(event.Value().(map[string]interface{}), schemaContents)
		return isValid, validationError, schemaContents
	}
}
