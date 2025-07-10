package oauth

import (
	"encoding/json"
	"net/http"

	"github.com/jetski-sh/mcp-proxy/log"
)

const ProtectedResourcePath = "/.well-known/oauth-protected-resource"

// NewWellKnownHandler implement the OAuth 2.0 Protected Resource Metadata (RFC9728) specification to indicate the
// locations of authorization servers.
//
// Should be used to create a handler for the /.well-known/oauth-protected-resource endpoint.
func NewProtectedResourceHandler(authorizationServers []string) http.Handler {
	type response struct {
		Resource             string   `json:"resource"`
		AuthorizationServers []string `json:"authorization_servers"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")

			response := response{
				Resource:             "http://" + r.Host + "/",
				AuthorizationServers: authorizationServers, // replace with []string{"http://" + r.Host + "/"} to use built-in server for metadata
			}

			log.Get(r.Context()).Info("Protected resource metadata", "response", response)

			if err := json.NewEncoder(w).Encode(response); err != nil {
				log.Get(r.Context()).Error(err, "failed to encode protected resource response")
			}
		}
	})
}
