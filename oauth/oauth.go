package oauth

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-chi/httprate"
	"github.com/jetski-sh/mcp-proxy/config"
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

		oidcProvider, err := oidc.NewProvider(context.TODO(), config.Authorization.Server)
		if err != nil {
			panic(err) // TODO: handle error properly
		}

		mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
			if token == "" {
				metadataURL, _ := url.Parse(config.Host.String())
				metadataURL.Path = ProtectedResourcePath
				w.Header().Set(
					"WWW-Authenticate",
					fmt.Sprintf(`Bearer resource_metadata="%s"`, metadataURL.String()),
				)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			verifier := oidcProvider.Verifier(&oidc.Config{SkipClientIDCheck: true})
			_, err = verifier.Verify(r.Context(), token)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		}))

		return mux
	}
}
