package oauth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-chi/httprate"
	"github.com/hyprmcp/mcp-gateway/config"
	"github.com/hyprmcp/mcp-gateway/htmlresponse"
	"github.com/hyprmcp/mcp-gateway/log"
	"github.com/hyprmcp/mcp-gateway/metadata"
	"github.com/hyprmcp/mcp-gateway/oauth/authorization"
	"github.com/hyprmcp/mcp-gateway/oauth/callback"
	"github.com/hyprmcp/mcp-gateway/oauth/dcr"
	"github.com/hyprmcp/mcp-gateway/oauth/userinfo"
)

type Manager struct {
	tokenValidator userinfo.TokenValidator
	config         *config.Config
	authConfig     metadata.MetadataSource
	authServerMeta map[string]any
}

func NewManager(ctx context.Context, cfg *config.Config) (*Manager, error) {
	if authConfig := config.GetActualAuthorizationConfig(cfg); authConfig == nil {
		return nil, errors.New("actual auth config must not be nil")
	} else if meta, err := authConfig.GetMetadata(ctx); err != nil {
		return nil, fmt.Errorf("authorization server metadata error: %w", err)
	} else if _, ok := authConfig.(*config.AuthorizationGitHub); ok {
		return &Manager{
			config:         cfg,
			authConfig:     authConfig,
			authServerMeta: meta,
			tokenValidator: userinfo.ValidateGitHub(nil),
		}, nil
	} else if jwksURI, ok := meta["jwks_uri"].(string); !ok {
		return nil, errors.New("no jwks_uri")
	} else if tokenValidator, err := userinfo.ValidateDynamicJWKS(ctx, jwksURI); err != nil {
		return nil, fmt.Errorf("failed to initialize dynamic JWKS token validator: %w", err)
	} else {
		return &Manager{
			config:         cfg,
			authConfig:     authConfig,
			authServerMeta: meta,
			tokenValidator: tokenValidator,
		}, nil
	}
}

func (mgr *Manager) Register(mux *http.ServeMux) error {
	log := log.Root().V(1)
	log.Info("registering OAuth endpoints")

	store := callback.NewURIStore()
	var authFns []authorization.EditQueryFunc
	var metaFns []metadata.EditMetadataFunc

	mux.Handle(ProtectedResourcePath, NewProtectedResourceHandler(mgr.config, mgr.config.Host.String()))

	if reg, err := dcr.NewRegistrarFromConfig(context.TODO(), &mgr.config.Authorization); err != nil {
		return err
	} else if reg != nil {
		if handler, err := NewDynamicClientRegistrationHandler(reg, mgr.authServerMeta); err != nil {
			return err
		} else {
			rateLimiter := httprate.LimitByRealIP(3, 10*time.Minute)
			mux.Handle(DynamicClientRegistrationPath, rateLimiter(handler))
		}

		regEndpointURI, _ := url.Parse(mgr.config.Host.String())
		regEndpointURI.Path = DynamicClientRegistrationPath
		metaFns = append(metaFns, metadata.ReplaceRegistrationEndpoint(regEndpointURI.String()))

		log.Info("DCR endpoint registered")
	}

	// TODO: Make required scopes optional/configurable
	authFns = append(authFns, authorization.RequiredScopes([]string{"openid", "profile", "email", "user:email"}, mgr.authServerMeta))

	if idSrc, ok := mgr.authConfig.(dcr.ClientIDSource); ok {
		if handler, err := NewTokenHandler(mgr.config.Host.String(), idSrc, mgr.authServerMeta); err != nil {
			return err
		} else {
			mux.Handle(TokenPath, handler)
		}
		mux.Handle(callback.CallbackPath, NewCallbackHandler(mgr.config, store))
		authFns = append(authFns, authorization.RedirectURI(mgr.config.Host.String(), store))

		tokenEndpointURI, _ := url.Parse(mgr.config.Host.String())
		tokenEndpointURI.Path = TokenPath
		metaFns = append(metaFns, metadata.ReplaceTokenEndpoint(tokenEndpointURI.String()))

		log.Info("token endpoint registered")
	}

	if len(authFns) > 0 {
		handler, err := NewAuthorizationHandler(mgr.config, mgr.authServerMeta, authorization.EditChain(authFns...))
		if err != nil {
			return err
		}
		mux.Handle(AuthorizationPath, handler)

		authEndpointURI, _ := url.Parse(mgr.config.Host.String())
		authEndpointURI.Path = AuthorizationPath
		metaFns = append(metaFns, metadata.ReplaceAuthorizationEndpoint(authEndpointURI.String()))

		log.Info("authorization endpoint registered", "authFns", len(authFns))
	}

	if len(metaFns) > 0 {
		mux.Handle(
			metadata.OAuth2MetadataPath,
			NewAuthorizationServerMetadataHandler(mgr.config, mgr.authConfig, metadata.EditChain(metaFns...)),
		)

		log.Info("metadata endpoint registered", "metaFns", len(metaFns))
	}

	return nil
}

func (mgr *Manager) Handler(next http.Handler) http.Handler {
	htmlHandler := htmlresponse.NewHandler(mgr.config, true)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawToken :=
			strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(r.Header.Get("Authorization")), "Bearer"))
		if userInfo, err := mgr.tokenValidator.ValidateToken(r.Context(), rawToken); err != nil {
			htmlHandler.Handler(mgr.unauthorizedHandler()).ServeHTTP(w, r)
		} else {
			next.ServeHTTP(w, r.WithContext(UserInfoContext(r.Context(), userInfo)))
		}
	})
}

func (mgr *Manager) unauthorizedHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metadataURL, _ := url.Parse(mgr.config.Host.String())
		metadataURL.Path = ProtectedResourcePath
		metadataURL = metadataURL.JoinPath(r.URL.Path)
		w.Header().Set(
			"WWW-Authenticate",
			fmt.Sprintf(`Bearer resource_metadata="%s"`, metadataURL.String()),
		)
		w.WriteHeader(http.StatusUnauthorized)
	}
}
