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
	"github.com/jetski-sh/mcp-proxy/log"
	"github.com/lestrrat-go/httprc/v3"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

func NewOAuthMiddleware(ctx context.Context, config *config.Config) (func(http.Handler) http.Handler, error) {
	var keySet jwk.Set
	if cache, err := jwk.NewCache(ctx, httprc.NewClient()); err != nil {
		return nil, err
	} else if meta, err := GetMedatata(config.Authorization.Server); err != nil {
		return nil, err
	} else if jwksURI, ok := meta["jwks_uri"].(string); !ok {
		return nil, errors.New("no jwks_uri")
	} else if err := cache.Register(ctx, jwksURI); err != nil {
		return nil, err
	} else if _, err := cache.Refresh(ctx, jwksURI); err != nil {
		return nil, err
	} else if s, err := cache.CachedSet(jwksURI); err != nil {
		return nil, err
	} else {
		keySet = s
		log.Get(ctx).Info("got jwk set", "jwks_uri", jwksURI)
	}

	var mr muxRegistrations

	mr.Add(ProtectedResourcePath, NewProtectedResourceHandler(config))

	if config.Authorization.ServerMetadataProxyEnabled {
		mr.Add(AuthorizationServerMetadataPath, NewAuthorizationServerMetadataHandler(config))
	}

	if config.Authorization.DynamicClientRegistrationEnabled {
		if handler, err := NewDynamicClientRegistrationHandler(config); err != nil {
			return nil, err
		} else {
			rateLimiter := httprate.LimitByRealIP(3, 10*time.Minute)
			mr.Add(DynamicClientRegistrationPath, rateLimiter(handler))
		}
	}

	return func(next http.Handler) http.Handler {
		mux := http.NewServeMux()

		mr.Register(mux)

		mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rawToken :=
				strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(r.Header.Get("Authorization")), "Bearer"))

			token, err := jwt.ParseString(rawToken, jwt.WithKeySet(keySet))
			if err != nil {
				metadataURL, _ := url.Parse(config.Host.String())
				metadataURL.Path = ProtectedResourcePath
				w.Header().Set(
					"WWW-Authenticate",
					fmt.Sprintf(`Bearer resource_metadata="%s"`, metadataURL.String()),
				)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r.WithContext(TokenContext(r.Context(), token, rawToken)))
		}))

		return mux
	}, nil
}
