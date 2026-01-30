package model

type Profile struct {
	UCP     UCPProfile     `json:"ucp"`
	Payment *PaymentConfig `json:"payment,omitempty"`
	Keys    []SigningKey   `json:"signing_keys,omitempty"`
}

type UCPProfile struct {
	Version      string                       `json:"version"`
	Services     map[string]ServiceDefinition `json:"services"`
	Capabilities []Capability                 `json:"capabilities"`
}

type ServiceDefinition struct {
	Version string       `json:"version"`
	Spec    string       `json:"spec"`
	REST    *RESTBinding `json:"rest,omitempty"`
}

type RESTBinding struct {
	Schema   string `json:"schema"`
	Endpoint string `json:"endpoint"`
}

type Capability struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Spec    string `json:"spec"`
	Schema  string `json:"schema"`
	Extends string `json:"extends,omitempty"`
}

type PaymentConfig struct {
	Handlers []PaymentHandler `json:"handlers"`
}

type PaymentHandler struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	Version           string   `json:"version"`
	Spec              string   `json:"spec"`
	ConfigSchema      string   `json:"config_schema"`
	InstrumentSchemas []string `json:"instrument_schemas"`
	Config            any      `json:"config"`
}

type SigningKey struct {
	KID string `json:"kid"`
	KTY string `json:"kty"`
	CRV string `json:"crv"`
	X   string `json:"x"`
	Y   string `json:"y"`
	Use string `json:"use"`
	Alg string `json:"alg"`
}
