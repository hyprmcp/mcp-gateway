package oauth

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/hyprmcp/mcp-gateway/config"
	"github.com/hyprmcp/mcp-gateway/log"
	"github.com/hyprmcp/mcp-gateway/metadata"
)

func NewAuthorizationServerMetadataHandler(config *config.Config, metaSource metadata.MetadataSource) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		meta, err := metaSource.GetMetadata(r.Context())
		if err != nil {
			log.Get(r.Context()).Error(err, "failed to get authorization server metadata from upstream")
			http.Error(w, "Failed to retrieve authorization server metadata", http.StatusInternalServerError)
		}

		if config.Authorization.GetDynamicClientRegistration().Enabled {
			if _, ok := meta["registration_endpoint"]; !ok {
				registrationURI, _ := url.Parse(config.Host.String())
				registrationURI.Path = DynamicClientRegistrationPath
				meta["registration_endpoint"] = registrationURI.String()
				log.Get(r.Context()).Info("Adding registration endpoint to authorization server metadata",
					"url", meta["registration_endpoint"])
			}
		}

		if config.Authorization.AuthorizationProxyEnabled || config.Authorization.ClientSecret != "" {
			authorizationURI, _ := url.Parse(config.Host.String())
			authorizationURI.Path = AuthorizationPath
			meta["authorization_endpoint"] = authorizationURI.String()
			log.Get(r.Context()).Info("Adding authorization endpoint to authorization server metadata",
				"url", meta["authorization_endpoint"])
		}

		if config.Authorization.ClientSecret != "" {
			tokenURI, _ := url.Parse(config.Host.String())
			tokenURI.Path = TokenPath
			meta["token_endpoint"] = tokenURI.String()
			log.Get(r.Context()).Info("Adding token endpoint to authorization server metadata",
				"url", meta["token_endpoint"])
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(meta); err != nil {
			log.Get(r.Context()).Error(err, "failed to encode authorization server metadata response")
		}
	})
}
