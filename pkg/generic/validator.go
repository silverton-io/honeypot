package generic

import (
	"github.com/silverton-io/honeypot/pkg/cache"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/validator"
	"github.com/tidwall/gjson"
)

func validateEvent(event gjson.Result, schemaName string, cache *cache.SchemaCache, conf *config.Generic) (isValid bool, validationError validator.ValidationError, schema []byte) {
	if event.Value() == nil {
		validationError := validator.ValidationError{
			ErrorType:       "payload not present at " + conf.Payload.RootKey + "." + conf.Payload.DataKey,
			ErrorResolution: "generic payload configuration and payload path should match",
			Errors:          nil,
		}
		return false, validationError, nil
	}
	if schemaName == "" { // Event does not have schema associated - always invalid.
		validationError := validator.ValidationError{
			ErrorType:       "schema not present at " + conf.Payload.RootKey + "." + conf.Payload.SchemaKey,
			ErrorResolution: "generic schema configuration and schema path should match",
			Errors:          nil,
		}
		return false, validationError, nil
	}
	schemaExists, schemaContents := cache.Get(schemaName)
	if !schemaExists { // Referenced schema is not present in the cache backend - always invalid
		validationError := validator.ValidationError{
			ErrorType:       "nonexistent schema",
			ErrorResolution: "publish the specified schema to the cache backend",
			Errors:          nil,
		}
		return false, validationError, nil
	} else {
		isValid, validationError := validator.ValidatePayload([]byte(event.Raw), schemaContents)
		return isValid, validationError, schemaContents
	}
}
