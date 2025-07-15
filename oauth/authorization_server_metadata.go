package oauth

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/jetski-sh/mcp-proxy/config"
	"github.com/jetski-sh/mcp-proxy/log"
)

const AuthorizationServerMetadataPath = "/.well-known/oauth-authorization-server"
const OIDCMetadataPath = "/.well-known/openid-configuration"

func NewAuthorizationServerMetadataHandler(config *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upstream := config.Authorization.Servers[0]

		metadataURL, err := url.JoinPath(upstream, OIDCMetadataPath)
		if err != nil {
			log.Get(r.Context()).Error(err, "failed to construct authorization server metadata URL", "base", upstream)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		upstreamResp, err := http.Get(metadataURL)
		if err != nil {
			log.Get(r.Context()).Error(err, "failed to fetch authorization server metadata", "url", metadataURL)
			w.WriteHeader(http.StatusBadGateway)
		} else if upstreamResp.StatusCode >= http.StatusBadRequest {
			log.Get(r.Context()).Error(err, "authorization server returned error", "status", upstreamResp.StatusCode, "url", upstreamResp.Request.URL)
			w.WriteHeader(upstreamResp.StatusCode)
			return
		}

		defer func() { _ = upstreamResp.Body.Close() }()

		var resp map[string]any
		if err := json.NewDecoder(upstreamResp.Body).Decode(&resp); err != nil {
			log.Get(r.Context()).Error(err, "failed to decode authorization server metadata response", "url", metadataURL)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if _, ok := resp["registration_endpoint"]; !ok {
			resp["registration_endpoint"] = "http://" + r.Host + DynamicClientRegistrationPath
			log.Get(r.Context()).Info("Adding registration endpoint to authorization server metadata",
				"url", resp["registration_endpoint"])
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Get(r.Context()).Error(err, "failed to encode authorization server metadata response")
		}
	})
}
