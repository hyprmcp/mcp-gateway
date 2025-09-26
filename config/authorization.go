package config

import (
	"context"
	"fmt"

	"github.com/hyprmcp/mcp-gateway/metadata"
)

type AuthorizationType string

type AuthorizationServerSource interface {
	GetAuthorizationServer() string
}

const (
	AuthorizationTypeOAuth2 AuthorizationType = "oauth2"
	AuthorizationTypeOIDC   AuthorizationType = "oidc"
	AuthorizationTypeGitHub AuthorizationType = "github"
	AuthorizationTypeDex    AuthorizationType = "dex"
)

func (t AuthorizationType) Validate() error {
	switch t {
	case AuthorizationTypeOAuth2,
		AuthorizationTypeOIDC,
		AuthorizationTypeGitHub,
		AuthorizationTypeDex:
		return nil
	default:
		return fmt.Errorf("invalid authorization type: %s", t)
	}
}

type Authorization struct {
	Type   AuthorizationType    `yaml:"type" json:"type"`
	OAuth2 *AuthorizationOAuth2 `yaml:"oauth2,omitempty" json:"oauth2,omitempty"`
	OIDC   *AuthorizationOIDC   `yaml:"oidc,omitempty" json:"oidc,omitempty"`
	GitHub *AuthorizationGitHub `yaml:"github,omitempty" json:"github,omitempty"`
	Dex    *AuthorizationDex    `yaml:"dex,omitempty" json:"dex,omitempty"`
}

func GetActualAuthorizationConfig(config *Config) metadata.MetadataSource {
	if config == nil {
		return nil
	}

	switch config.Authorization.Type {
	case AuthorizationTypeOAuth2:
		return config.Authorization.OAuth2
	case AuthorizationTypeOIDC:
		return config.Authorization.OIDC
	case AuthorizationTypeGitHub:
		return config.Authorization.GitHub
	case AuthorizationTypeDex:
		return config.Authorization.Dex
	default:
		panic("invalid authorization config")
	}
}

func (a *Authorization) Validate() error {
	if a == nil {
		return fmt.Errorf("authorization is nil")
	}

	if err := a.Type.Validate(); err != nil {
		return err
	}

	switch a.Type {
	case AuthorizationTypeOAuth2:
		if err := a.OAuth2.Validate(); err != nil {
			return fmt.Errorf("oauth2 is invalid: %w", err)
		}
	case AuthorizationTypeOIDC:
		if err := a.OIDC.Validate(); err != nil {
			return fmt.Errorf("oidc is invalid: %w", err)
		}
	case AuthorizationTypeGitHub:
		if err := a.GitHub.Validate(); err != nil {
			return fmt.Errorf("github is invalid: %w", err)
		}
	case AuthorizationTypeDex:
		if err := a.Dex.Validate(); err != nil {
			return fmt.Errorf("dex is invalid: %w", err)
		}
	}

	return nil
}

type AuthorizationOAuth2 struct {
	Metadata     metadata.Metadata `yaml:"metadata" json:"metadata"`
	ClientID     string            `yaml:"clientId" json:"clientId"`
	ClientSecret string            `yaml:"clientSecret" json:"clientSecret"`
}

func (a *AuthorizationOAuth2) Validate() error {
	if a == nil {
		return fmt.Errorf("oauth2 is nil")
	} else if a.ClientID == "" {
		return fmt.Errorf("clientId is required")
	} else if a.ClientSecret == "" {
		return fmt.Errorf("clientSecret is required")
	} else if a.Metadata == nil {
		return fmt.Errorf("metadata is required")
	}

	return nil
}

func (a *AuthorizationOAuth2) GetClientID() string {
	return a.ClientID
}

func (a *AuthorizationOAuth2) GetClientSecret() string {
	return a.ClientSecret
}

func (a *AuthorizationOAuth2) GetMetadata(ctx context.Context) (metadata.Metadata, error) {
	return a.Metadata, nil
}

type AuthorizationOIDC struct {
	IssuerURL    string `yaml:"issuerUrl" json:"issuerUrl"`
	ClientID     string `yaml:"clientId" json:"clientId"`
	ClientSecret string `yaml:"clientSecret" json:"clientSecret"`
}

func (a *AuthorizationOIDC) GetAuthorizationServer() string {
	return a.IssuerURL
}

func (a *AuthorizationOIDC) Validate() error {
	if a == nil {
		return fmt.Errorf("oidc is nil")
	} else if a.IssuerURL == "" {
		return fmt.Errorf("issuerUrl is required")
	} else if a.ClientID == "" {
		return fmt.Errorf("clientId is required")
	} else if a.ClientSecret == "" {
		return fmt.Errorf("clientSecret is required")
	}

	return nil
}

func (a *AuthorizationOIDC) GetClientID() string {
	return a.ClientID
}

func (a *AuthorizationOIDC) GetClientSecret() string {
	return a.ClientSecret
}

func (a *AuthorizationOIDC) GetMetadata(ctx context.Context) (metadata.Metadata, error) {
	return metadata.FetchMetadataDiscovery(ctx, a.IssuerURL)
}

// TODO: GitHub access tokens are not JWTs, so we need to validate them by calling the userinfo endpoint
// Something similar should probably be done for OAuth as well
type AuthorizationGitHub struct {
	ClientID     string `yaml:"clientId" json:"clientId"`
	ClientSecret string `yaml:"clientSecret" json:"clientSecret"`
}

func (a *AuthorizationGitHub) GetClientID() string {
	return a.ClientID
}

func (a *AuthorizationGitHub) GetClientSecret() string {
	return a.ClientSecret
}

func (a *AuthorizationGitHub) GetMetadata(ctx context.Context) (metadata.Metadata, error) {
	return metadata.Metadata{
		// TODO: issuer should be host from the root config
		"issuer":                   "https://github.com",
		"response_types_supported": []string{"code"},
		"authorization_endpoint":   "https://github.com/login/oauth/authorize",
		"token_endpoint":           "https://github.com/login/oauth/access_token",
		"userinfo_endpoint":        "https://api.github.com/user",
		"scopes_supported":         []string{"read:user", "user:email"},
	}, nil
}

func (a *AuthorizationGitHub) Validate() error {
	if a == nil {
		return fmt.Errorf("github is nil")
	} else if a.ClientID == "" {
		return fmt.Errorf("clientId is required")
	} else if a.ClientSecret == "" {
		return fmt.Errorf("clientSecret is required")
	}

	return nil
}

type AuthorizationDex struct {
	Server                    string                     `yaml:"server" json:"server"`
	GRPCAddr                  string                     `yaml:"grpcAddr" json:"grpcAddr"`
	DynamicClientRegistration *DynamicClientRegistration `yaml:"dynamicClientRegistration" json:"dynamicClientRegistration"`
}

func (a *AuthorizationDex) GetAuthorizationServer() string {
	return a.Server
}

func (a *AuthorizationDex) Validate() error {
	if a == nil {
		return fmt.Errorf("dex is invalid")
	} else if a.Server == "" {
		return fmt.Errorf("server is required")
	} else if a.GRPCAddr == "" {
		return fmt.Errorf("grpcAddr is required")
	}

	return nil
}

func (a *AuthorizationDex) GetMetadata(ctx context.Context) (metadata.Metadata, error) {
	return metadata.FetchMetadataDiscovery(ctx, a.Server)
}

type ClientType string

const (
	ClientTypePublic       ClientType = "public"
	ClientTypeConfidential ClientType = "confidential"
)

func (ct ClientType) IsPublic() bool {
	return ct == ClientTypePublic
}

func (ct ClientType) IsConfidential() bool {
	return !ct.IsPublic()
}

func (ct ClientType) Validate() error {
	switch ct {
	case ClientTypePublic, ClientTypeConfidential:
		return nil
	default:
		return fmt.Errorf("invalid client type: %s", ct)
	}
}

type DynamicClientRegistration struct {
	Enabled    bool       `yaml:"enabled" json:"enabled"`
	ClientType ClientType `yaml:"clientType" json:"clientType"`
}

func (dcr *DynamicClientRegistration) Validate() error {
	if dcr != nil && dcr.Enabled {
		if err := dcr.ClientType.Validate(); err != nil {
			return fmt.Errorf("dynamicClientRegistration is invalid: %w", err)
		}
	}

	return nil
}
