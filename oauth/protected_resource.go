package oauth

import (
	"encoding/json"
	"net/http"

	"github.com/jetski-sh/mcp-gateway/config"
	"github.com/jetski-sh/mcp-gateway/log"
)

const ProtectedResourcePath = "/.well-known/oauth-protected-resource"

type ProtectedResourceMetadata struct {
	Resource             string   `json:"resource"`
	AuthorizationServers []string `json:"authorization_servers"`
}

// NewWellKnownHandler implement the OAuth 2.0 Protected Resource Metadata (RFC9728) specification to indicate the
// locations of authorization servers.
//
// Should be used to create a handler for the /.well-known/oauth-protected-resource endpoint.
func NewProtectedResourceHandler(config *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")

			response := ProtectedResourceMetadata{Resource: config.Host.String()}

			if config.Authorization.ServerMetadataProxyEnabled {
				response.AuthorizationServers = []string{config.Host.String()}
			} else {
				response.AuthorizationServers = []string{config.Authorization.Server}
			}

			log.Get(r.Context()).Info("Protected resource metadata", "response", response)

			if err := json.NewEncoder(w).Encode(response); err != nil {
				log.Get(r.Context()).Error(err, "failed to encode protected resource response")
			}
		}
	})
}
