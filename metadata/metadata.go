package metadata

import (
	"context"
	"slices"
)

type Metadata map[string]any

var _ MetadataSource = Metadata(nil)

const (
	RegistrationEndpointKey  = "registration_endpoint"
	TokenEndpointKey         = "token_endpoint"
	AuthorizationEndpointKey = "authorization_endpoint"
	SupportedScopesKey       = "scopes_supported"
)

func (meta Metadata) GetMetadata(ctx context.Context) (Metadata, error) {
	return meta, nil
}

func (meta Metadata) GetSupportedScopes() []string {
	if scopesSupported, ok := meta[SupportedScopesKey].([]any); ok {
		scopesSupportedStr := make([]string, 0, len(scopesSupported))
		for _, v := range scopesSupported {
			if s, ok := v.(string); ok {
				scopesSupportedStr = append(scopesSupportedStr, s)
			}
		}

		return slices.Clip(scopesSupportedStr)
	}

	return nil
}

func (meta Metadata) GetRegistrationEndpoint() string {
	return meta.GetStr(RegistrationEndpointKey)
}

func (meta Metadata) GetStr(key string) string {
	if val, ok := meta[key].(string); ok {
		return val
	}
	return ""
}
