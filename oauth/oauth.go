package oauth

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/httprate"
	"github.com/jetski-sh/mcp-proxy/config"
	"github.com/jetski-sh/mcp-proxy/log"
	"github.com/lestrrat-go/httprc/v3"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

func NewOAuthMiddleware(config *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		mux := http.NewServeMux()

		mux.Handle(
			ProtectedResourcePath,
			NewProtectedResourceHandler(config),
		)

		if config.Authorization.ServerMetadataProxyEnabled {
			mux.Handle(
				AuthorizationServerMetadataPath,
				NewAuthorizationServerMetadataHandler(config),
			)
		}

		if config.Authorization.DynamicClientRegistrationEnabled {
			rateLimiter := httprate.LimitByRealIP(3, 10*time.Minute)
			mux.Handle(
				DynamicClientRegistrationPath,
				rateLimiter(NewDynamicClientRegistrationHandler(config)),
			)
		}

		var keySet jwk.Set
		if cache, err := jwk.NewCache(context.TODO(), httprc.NewClient()); err != nil {
			panic(err)
		} else if meta, err := GetMedatata(config.Authorization.Server); err != nil {
			panic(err)
		} else if jwksURI, ok := meta["jwks_uri"].(string); !ok {
			panic("no jwks_uri")
		} else if err := cache.Register(context.TODO(), jwksURI); err != nil {
			panic(err)
		} else if _, err := cache.Refresh(context.TODO(), jwksURI); err != nil {
			panic(err)
		} else if s, err := cache.CachedSet(jwksURI); err != nil {
			panic(err)
		} else {
			keySet = s
			log.Root().Info("got jwk set", "jwks_uri", jwksURI)
		}

		mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := jwt.ParseHeader(r.Header, "Authorization", jwt.WithKeySet(keySet))
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

			next.ServeHTTP(w, r.WithContext(AddTokenToContext(r.Context(), token)))
		}))

		return mux
	}
}
