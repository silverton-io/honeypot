package protocol

const (
	SNOWPLOW    string = "snowplow"
	GENERIC     string = "generic"
	CLOUDEVENTS string = "cloudevents"
	RELAY       string = "relay"
	WEBHOOK     string = "webhook"
)

func GetIntputProtocols() []string {
	return []string{SNOWPLOW, GENERIC, CLOUDEVENTS, RELAY, WEBHOOK}
}
