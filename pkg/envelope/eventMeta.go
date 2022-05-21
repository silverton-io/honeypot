package envelope

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/google/uuid"
)

type EventMeta struct {
	Protocol           string    `json:"protocol,omitempty"`
	Uuid               uuid.UUID `json:"uuid,omitempty"`
	Vendor             string    `json:"vendor,omitempty"`
	PrimaryNamespace   string    `json:"primaryNamespace,omitempty"`
	SecondaryNamespace string    `json:"secondaryNamespace,omitempty"`
	TertiaryNamespace  string    `json:"tertiaryNamespace,omitempty"`
	Name               string    `json:"name,omitempty"`
	Version            string    `json:"version,omitempty"`
	Path               string    `json:"path,omitempty"`
}

func (e EventMeta) Value() (driver.Value, error) {
	b, err := json.Marshal(e)
	return string(b), err
}

func (e EventMeta) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), &e)
}
