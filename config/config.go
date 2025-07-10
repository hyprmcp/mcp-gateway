package config

import (
	"fmt"
	"io"
	"net/url"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	AuthorizationServers []string `yaml:"authorizationServers"`
	Proxy                []Proxy  `yaml:"proxy"`
}

type Proxy struct {
	Path string     `yaml:"path"`
	Http *ProxyHttp `yaml:"http,omitempty"`
}

type ProxyHttp struct {
	Url *ProxyUrl `yaml:"url"`
}

type ProxyUrl url.URL

func (p *ProxyUrl) UnmarshalYAML(value *yaml.Node) error {
	if parsed, err := url.Parse(value.Value); err != nil {
		return err
	} else {
		*p = ProxyUrl(*parsed)
		return nil
	}
}

func ParseFile(fileName string) (*Config, error) {
	if file, err := os.Open(fileName); err != nil {
		return nil, fmt.Errorf("failed to open config file %s: %w", fileName, err)
	} else {
		defer file.Close()
		return Parse(file)
	}
}

func Parse(r io.Reader) (*Config, error) {
	var config Config
	if err := yaml.NewDecoder(r).Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}
	return &config, nil
}
