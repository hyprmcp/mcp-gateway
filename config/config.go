package config

import (
	"fmt"
	"io"
	"net/url"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Host          *YamlUrl       `yaml:"host"`
	Authorization Authorization  `yaml:"authorization"`
	DexGRPCClient *DexGRPCClient `yaml:"dexGRPCClient,omitempty"`
	Proxy         []Proxy        `yaml:"proxy"`
}

type Proxy struct {
	Path string     `yaml:"path"`
	Http *ProxyHttp `yaml:"http,omitempty"`
}

type Authorization struct {
	Server                           string `yaml:"server"`
	ServerMetadataProxyEnabled       bool   `yaml:"serverMetadataProxyEnabled"`
	DynamicClientRegistrationEnabled bool   `yaml:"dynamicClientRegistrationEnabled"`
}

type DexGRPCClient struct {
	Addr string `yaml:"addr"`
}

type ProxyHttp struct {
	Url *YamlUrl `yaml:"url"`
}

type YamlUrl url.URL

func (p *YamlUrl) UnmarshalYAML(value *yaml.Node) error {
	if parsed, err := url.Parse(value.Value); err != nil {
		return err
	} else {
		*p = YamlUrl(*parsed)
		return nil
	}
}

func (p *YamlUrl) String() string {
	return (*url.URL)(p).String()
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
	return &config, config.Validate()
}

func (c *Config) Validate() error {
	if c.Host == nil {
		return fmt.Errorf("host is required")
	}

	if c.Authorization.Server == "" {
		return fmt.Errorf("authorization server is required")
	}

	if c.Authorization.DynamicClientRegistrationEnabled {
		if !c.Authorization.ServerMetadataProxyEnabled {
			return fmt.Errorf("serverMetadataProxyEnabled must be true when dynamicClientRegistrationEnabled is true")
		}

		if c.DexGRPCClient == nil || c.DexGRPCClient.Addr == "" {
			return fmt.Errorf("dexGRPCClient is required when dynamicClientRegistrationEnabled is true")
		}
	}

	return nil
}
