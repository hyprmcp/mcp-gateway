package config

import (
	"encoding/json"
	"net/url"

	"gopkg.in/yaml.v3"
)

type URL url.URL

func (p *URL) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	} else if parsed, err := url.Parse(s); err != nil {
		return err
	} else {
		*p = URL(*parsed)
		return nil
	}
}

func (p URL) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

func (p *URL) UnmarshalYAML(value *yaml.Node) error {
	if parsed, err := url.Parse(value.Value); err != nil {
		return err
	} else {
		*p = URL(*parsed)
		return nil
	}
}

func (p URL) MarshalYAML() (any, error) {
	return p.String(), nil
}

func (p *URL) String() string {
	return (*url.URL)(p).String()
}
