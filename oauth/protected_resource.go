package oauth

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/hyprmcp/mcp-gateway/config"
	"github.com/hyprmcp/mcp-gateway/log"
)

const ProtectedResourcePath = "/.well-known/oauth-protected-resource/"

type ProtectedResourceMetadata struct {
	Resource             string   `json:"resource"`
	AuthorizationServers []string `json:"authorization_servers"`
}

// NewWellKnownHandler implement the OAuth 2.0 Protected Resource Metadata (RFC9728) specification to indicate the
// locations of authorization servers.
//
// Should be used to create a handler for the /.well-known/oauth-protected-resource endpoint.
func NewProtectedResourceHandler(config *config.Config, authServerURI string) http.Handler {
	return http.StripPrefix(
		ProtectedResourcePath,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")

				resourceURL, _ := url.Parse(config.Host.String())
				resourceURL = resourceURL.JoinPath(r.URL.Path)
				response := ProtectedResourceMetadata{
					Resource:             resourceURL.String(),
					AuthorizationServers: []string{authServerURI},
				}

				log.Get(r.Context()).Info("Protected resource metadata", "response", response)

				if err := json.NewEncoder(w).Encode(response); err != nil {
					log.Get(r.Context()).Error(err, "failed to encode protected resource response")
				}
			}
		}),
	)
}
