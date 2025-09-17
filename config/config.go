package config

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Host          *URL          `yaml:"host" json:"host"`
	Authorization Authorization `yaml:"authorization" json:"authorization"`
	Proxy         []Proxy       `yaml:"proxy" json:"proxy"`
}

func ParseFile(fileName string) (*Config, error) {
	if file, err := os.Open(fileName); err != nil {
		return nil, fmt.Errorf("failed to open config file %s: %w", fileName, err)
	} else {
		defer func() { _ = file.Close() }()
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

func (c *Config) YAMLString() (string, error) {
	if data, err := yaml.Marshal(c); err != nil {
		return "", err
	} else {
		return string(data), nil
	}
}

func (c *Config) Validate() error {
	if c.Host == nil {
		return fmt.Errorf("host is required")
	}

	if err := c.Authorization.Validate(); err != nil {
		return fmt.Errorf("authorization is invalid: %w", err)
	}

	return nil
}
