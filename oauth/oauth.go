package oauth

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func NewOAuthMiddleware(authorizationServers []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		mux := http.NewServeMux()
		mux.Handle(
			ProtectedResourcePath,
			NewProtectedResourceHandler(authorizationServers),
		)
		mux.Handle(
			AuthorizationServerMetadataPath,
			NewAuthorizationServerMetadataHandler(authorizationServers[0]),
		)
		mux.Handle("/{asd}/mcp", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
