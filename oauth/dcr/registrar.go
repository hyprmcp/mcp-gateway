package dcr

import (
	"context"
	"fmt"

	"github.com/hyprmcp/mcp-gateway/config"
)

type ClientRegistrar interface {
	RegisterClient(ctx context.Context, client Client) (*Client, error)
}

func NewRegistrarFromConfig(ctx context.Context, cfg *config.Authorization) (ClientRegistrar, error) {
	switch cfg.Type {

	case config.AuthorizationTypeOAuth2:
		if meta, err := cfg.OAuth2.GetMetadata(ctx); err != nil {
			return nil, err
		} else if reg := meta.GetRegistrationEndpoint(); reg == "" {
			return NewFakeRegistrar(cfg.OAuth2), nil
		} else {
			return nil, nil
		}

	case config.AuthorizationTypeOIDC:
		if meta, err := cfg.OIDC.GetMetadata(ctx); err != nil {
			return nil, err
		} else if reg := meta.GetRegistrationEndpoint(); reg == "" {
			return NewFakeRegistrar(cfg.OIDC), nil
		} else {
			return nil, nil
		}

	case config.AuthorizationTypeGitHub:
		return NewFakeRegistrar(cfg.GitHub), nil

	case config.AuthorizationTypeDex:
		if cfg.Dex.DynamicClientRegistration.Enabled {
			return NewDexRegistrar(cfg.Dex)
		} else {
			return nil, nil
		}

	default:
		return nil, fmt.Errorf("unsupported authorization type: %s", cfg.Type)

	}
}
