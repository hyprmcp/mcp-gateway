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
	"github.com/lestrrat-go/httprc/v3"
	"github.com/lestrrat-go/httprc/v3/errsink"
	"github.com/lestrrat-go/httprc/v3/tracesink"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

type Manager struct {
	jwkSet         jwk.Set
	config         *config.Config
	authConfig     metadata.MetadataSource
	authServerMeta map[string]any
}

func NewManager(ctx context.Context, cfg *config.Config) (*Manager, error) {
	log := log.Get(ctx)

	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if authConfig := config.GetActualAuthorizationConfig(cfg); authConfig == nil {
		return nil, errors.New("actual auth config must not be nil")
	} else if meta, err := authConfig.GetMetadata(ctx); err != nil {
		return nil, fmt.Errorf("authorization server metadata error: %w", err)
	} else if cache, err := jwk.NewCache(ctx, httprc.NewClient(
		httprc.WithTraceSink(tracesink.Func(func(ctx context.Context, s string) { log.V(1).Info(s) })),
		httprc.WithErrorSink(errsink.NewFunc(func(ctx context.Context, err error) { log.V(1).Error(err, "httprc.NewClient error") })),
	)); err != nil {
		return nil, fmt.Errorf("jwk cache creation error: %w", err)
	} else if jwksURI, ok := meta["jwks_uri"].(string); !ok {
		return nil, errors.New("no jwks_uri")
	} else if err := cache.Register(
		timeoutCtx,
		jwksURI,
		jwk.WithMinInterval(10*time.Second),
		jwk.WithMaxInterval(5*time.Minute),
	); err != nil {
		return nil, fmt.Errorf("jwks registration error: %w", err)
	} else if _, err := cache.Refresh(timeoutCtx, jwksURI); err != nil {
		return nil, fmt.Errorf("jwks refresh error: %w", err)
	} else if s, err := cache.CachedSet(jwksURI); err != nil {
		return nil, fmt.Errorf("jwks cache set error: %w", err)
	} else {
		return &Manager{jwkSet: s, config: cfg, authConfig: authConfig, authServerMeta: meta}, nil
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
	authFns = append(authFns, authorization.RequiredScopes([]string{"openid", "profile", "email"}, mgr.authServerMeta))

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
		if token, err := jwt.ParseString(rawToken, jwt.WithKeySet(mgr.jwkSet)); err != nil {
			htmlHandler.Handler(mgr.unauthorizedHandler()).ServeHTTP(w, r)
		} else {
			next.ServeHTTP(w, r.WithContext(TokenContext(r.Context(), token, rawToken)))
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
