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
	"github.com/jetski-sh/mcp-proxy/config"
	"github.com/lestrrat-go/httprc/v3"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

type Manager struct {
	jwkSet jwk.Set
	config *config.Config
}

func NewManager(ctx context.Context, config *config.Config) (*Manager, error) {
	if cache, err := jwk.NewCache(ctx, httprc.NewClient()); err != nil {
		return nil, fmt.Errorf("jwk cache creation error: %w", err)
	} else if meta, err := GetMedatata(config.Authorization.Server); err != nil {
		return nil, fmt.Errorf("authorization server metadata error: %w", err)
	} else if jwksURI, ok := meta["jwks_uri"].(string); !ok {
		return nil, errors.New("no jwks_uri")
	} else if err := cache.Register(ctx, jwksURI); err != nil {
		return nil, fmt.Errorf("jwks registration error: %w", err)
	} else if _, err := cache.Refresh(ctx, jwksURI); err != nil {
		return nil, fmt.Errorf("jwks refresh error: %w", err)
	} else if s, err := cache.CachedSet(jwksURI); err != nil {
		return nil, fmt.Errorf("jwks cache set error: %w", err)
	} else {
		return &Manager{jwkSet: s, config: config}, nil
	}
}

func (mgr *Manager) Register(mux *http.ServeMux) error {
	mux.Handle(ProtectedResourcePath, NewProtectedResourceHandler(mgr.config))

	if mgr.config.Authorization.ServerMetadataProxyEnabled {
		mux.Handle(AuthorizationServerMetadataPath, NewAuthorizationServerMetadataHandler(mgr.config))
	}

	if mgr.config.Authorization.DynamicClientRegistrationEnabled {
		if handler, err := NewDynamicClientRegistrationHandler(mgr.config); err != nil {
			return err
		} else {
			rateLimiter := httprate.LimitByRealIP(3, 10*time.Minute)
			mux.Handle(DynamicClientRegistrationPath, rateLimiter(handler))
		}
	}

	return nil
}

func (mgr *Manager) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawToken :=
			strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(r.Header.Get("Authorization")), "Bearer"))

		token, err := jwt.ParseString(rawToken, jwt.WithKeySet(mgr.jwkSet))
		if err != nil {
			metadataURL, _ := url.Parse(mgr.config.Host.String())
			metadataURL.Path = ProtectedResourcePath
			w.Header().Set(
				"WWW-Authenticate",
				fmt.Sprintf(`Bearer resource_metadata="%s"`, metadataURL.String()),
			)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(TokenContext(r.Context(), token, rawToken)))
	})
}
