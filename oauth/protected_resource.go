package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/hyprmcp/mcp-gateway/config"
	"github.com/hyprmcp/mcp-gateway/log"
)

const ProtectedResourcePath = "/.well-known/oauth-protected-resource/"

var upstreamMetadataURLs sync.Map

type ProtectedResourceMetadata struct {
	Resource             string         `json:"resource"`
	AuthorizationServers []string       `json:"authorization_servers"`
	ExtraFields          map[string]any `json:"-"`
}

// NewWellKnownHandler implement the OAuth 2.0 Protected Resource Metadata (RFC9728) specification to indicate the
// locations of authorization servers.
//
// Should be used to create a handler for the /.well-known/oauth-protected-resource endpoint.
func NewProtectedResourceHandler(config *config.Config) http.Handler {
	return http.StripPrefix(
		ProtectedResourcePath,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}

			if strings.Trim(r.URL.Path, "/") == "" && len(config.Proxy) == 1 && !config.Proxy[0].Authentication.Enabled {
				http.NotFound(w, r)
				return
			}

			var response ProtectedResourceMetadata
			if upstream, ok := upstreamMetadataURLs.Load(strings.Trim(r.URL.Path, `/`)); ok {
				fmt.Println("loaded", upstream)
				if upstreamStr, ok := upstream.(string); ok {
					if req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, upstreamStr, nil); err != nil {
						log.Get(r.Context()).Error(err, "failed to create request")
						http.Error(w, "failed to create request", http.StatusBadGateway)
						return
					} else if resp, err := http.DefaultClient.Do(req); err != nil {
						log.Get(r.Context()).Error(err, "failed to fetch metadata")
						http.Error(w, "failed to fetch metadata", http.StatusBadGateway)
						return
					} else if resp.StatusCode != http.StatusOK {
						log.Get(r.Context()).Error(errors.New("upstream returned non-200 status code"), "upstream returned non-200 status code")
						for key, values := range resp.Header {
							for _, val := range values {
								w.Header().Set(key, val)
							}
						}
						w.WriteHeader(resp.StatusCode)
						io.Copy(w, resp.Body)
						return
					} else {
						defer resp.Body.Close()

						if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
							log.Get(r.Context()).Error(err, "failed to decode metadata")
							http.Error(w, "failed to decode metadata", http.StatusInternalServerError)
							return
						}
					}
				} else {
					log.Get(r.Context()).Error(errors.New("upstream is not a string"), "upstream is not a string")
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
			} else {
				if config.Authorization.ServerMetadataProxyEnabled {
					response.AuthorizationServers = []string{config.Host.String()}
				} else {
					response.AuthorizationServers = []string{config.Authorization.Server}
				}
			}

			resourceURL, _ := url.Parse(config.Host.String())
			resourceURL = resourceURL.JoinPath(r.URL.Path)
			response.Resource = resourceURL.String()

			log.Get(r.Context()).Info("Protected resource metadata", "response", response)

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(response); err != nil {
				log.Get(r.Context()).Error(err, "failed to encode protected resource response")
			}
		}),
	)
}
