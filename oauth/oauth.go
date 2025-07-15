package oauth

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

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
			mux.Handle(
				DynamicClientRegistrationPath,
				NewDynamicClientRegistrationHandler(config),
			)
		}

		mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
			if token == "" {
				metadataURL, _ := url.Parse(r.URL.String())
				metadataURL.Path = ProtectedResourcePath
				w.Header().Set(
					"WWW-Authenticate",
					fmt.Sprintf(`Bearer resource_metadata="%s"`, metadataURL.String()),
				)
				w.WriteHeader(http.StatusUnauthorized)
			} else {
				next.ServeHTTP(w, r)
			}
		}))

		return mux
	}
}
