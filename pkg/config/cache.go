package config

type Purge struct {
	Enabled bool   `json:"enabled"`
	Path    string `json:"path"`
}

type Backend struct {
	Type   string `json:"type"`
	Region string `json:"region,omitempty"`
	Bucket string `json:"bucket,omitempty"`
	Host   string `json:"host,omitempty"`
	Path   string `json:"path"`
}

type SchemaDirectory struct {
	Enabled bool `json:"enabled"`
}

type SchemaCache struct {
	Backend         `json:"backend"`
	TtlSeconds      int `json:"ttlSeconds"`
	MaxSizeBytes    int `json:"maxSizeBytes"`
	Purge           `json:"purge"`
	SchemaDirectory `json:"schemaDirectory"`
}
