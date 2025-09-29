package config

type Proxy struct {
	Path           string              `yaml:"path" json:"path"`
	Http           *ProxyHttp          `yaml:"http,omitempty" json:"http,omitempty"`
	Authentication ProxyAuthentication `yaml:"authentication" json:"authentication"`
	Telemetry      ProxyTelemetry      `yaml:"telemetry" json:"telemetry"`
	Webhook        *Webhook            `yaml:"webhook,omitempty" json:"webhook,omitempty"`
}

type ProxyHttp struct {
	Url *URL `yaml:"url" json:"url"`
}

type ProxyAuthentication struct {
	Enabled bool `yaml:"enabled" json:"enabled"`
}

type ProxyTelemetry struct {
	Enabled bool `yaml:"enabled" json:"enabled"`
}

type Webhook struct {
	Method string `yaml:"method,omitempty" json:"method,omitempty"`
	Url    URL    `yaml:"url" json:"url"`
}
