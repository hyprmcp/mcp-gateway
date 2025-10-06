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
	"github.com/lestrrat-go/httprc/v3"
	"github.com/lestrrat-go/httprc/v3/errsink"
	"github.com/lestrrat-go/httprc/v3/tracesink"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

type Manager struct {
	jwkSet         jwk.Set
	config         *config.Config
	authServerMeta map[string]any
}

func NewManager(ctx context.Context, config *config.Config) (*Manager, error) {
	log := log.Get(ctx)

	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if cache, err := jwk.NewCache(ctx, httprc.NewClient(
		httprc.WithTraceSink(tracesink.Func(func(ctx context.Context, s string) { log.V(1).Info(s) })),
		httprc.WithErrorSink(errsink.NewFunc(func(ctx context.Context, err error) { log.V(1).Error(err, "httprc.NewClient error") })),
	)); err != nil {
		return nil, fmt.Errorf("jwk cache creation error: %w", err)
	} else if meta, err := GetMedatata(config.Authorization.Server); err != nil {
		return nil, fmt.Errorf("authorization server metadata error: %w", err)
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
		return &Manager{jwkSet: s, config: config, authServerMeta: meta}, nil
	}
}

func (mgr *Manager) Register(mux *http.ServeMux) error {
	mux.Handle(ProtectedResourcePath, NewProtectedResourceHandler(mgr.config))

	if mgr.config.Authorization.ServerMetadataProxyEnabled {
		mux.Handle(AuthorizationServerMetadataPath, NewAuthorizationServerMetadataHandler(mgr.config))
	}

	if mgr.config.Authorization.GetDynamicClientRegistration().Enabled {
		if handler, err := NewDynamicClientRegistrationHandler(mgr.config, mgr.authServerMeta); err != nil {
			return err
		} else {
			rateLimiter := httprate.LimitByRealIP(3, 10*time.Minute)
			mux.Handle(DynamicClientRegistrationPath, rateLimiter(handler))
		}
	}

	if mgr.config.Authorization.AuthorizationProxyEnabled {
		if handler, err := NewAuthorizationHandler(mgr.config, mgr.authServerMeta); err != nil {
			return err
		} else {
			mux.Handle(AuthorizationPath, handler)
		}
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
		w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Bearer resource_metadata="%s"`, mgr.getMetadataURL(r.URL)))
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func (mgr *Manager) getMetadataURL(u *url.URL) *url.URL {
	metadataURL, _ := url.Parse(mgr.config.Host.String())
	metadataURL.Path = ProtectedResourcePath
	metadataURL = metadataURL.JoinPath(u.Path)
	return metadataURL
}
