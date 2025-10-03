package oauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/hyprmcp/mcp-gateway/config"
	"github.com/hyprmcp/mcp-gateway/log"
	"github.com/hyprmcp/mcp-gateway/proxy/proxyutil"
	"go.uber.org/multierr"
)

const AuthorizationServerMetadataPath = "/.well-known/oauth-authorization-server"
const OIDCMetadataPath = "/.well-known/openid-configuration"

func NewAuthorizationServerMetadataHandler(config *config.Config) http.Handler {
	if len(config.Proxy) == 1 && !config.Proxy[0].Authentication.Enabled {
		return &httputil.ReverseProxy{
			Rewrite:        proxyutil.RewriteHostFunc((*url.URL)(config.Proxy[0].Http.Url)),
			ModifyResponse: proxyutil.RemoveCORSHeaders,
		}
	} else {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			metadata, err := GetMedatata(config.Authorization.Server)
			if err != nil {
				log.Get(r.Context()).Error(err, "failed to get authorization server metadata from upstream")
				http.Error(w, "Failed to retrieve authorization server metadata", http.StatusInternalServerError)
			}

			if config.Authorization.GetDynamicClientRegistration().Enabled {
				if _, ok := metadata["registration_endpoint"]; !ok {
					registrationURI, _ := url.Parse(config.Host.String())
					registrationURI.Path = DynamicClientRegistrationPath
					metadata["registration_endpoint"] = registrationURI.String()
					log.Get(r.Context()).Info("Adding registration endpoint to authorization server metadata",
						"url", metadata["registration_endpoint"])
				}
			}

			if config.Authorization.AuthorizationProxyEnabled {
				authorizationURI, _ := url.Parse(config.Host.String())
				authorizationURI.Path = AuthorizationPath
				metadata["authorization_endpoint"] = authorizationURI.String()
				log.Get(r.Context()).Info("Adding authorization endpoint to authorization server metadata",
					"url", metadata["authorization_endpoint"])
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(metadata); err != nil {
				log.Get(r.Context()).Error(err, "failed to encode authorization server metadata response")
			}
		})
	}
}

func GetMedatata(server string) (map[string]any, error) {
	uris, err := getMetadataURIs(server)
	if err != nil {
		return nil, err
	}

	getMetatadatFunc := func(u string) (map[string]any, error) {
		resp, err := http.Get(u)
		if err != nil {
			return nil, err
		}

		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode >= http.StatusBadRequest {
			return nil, fmt.Errorf("failed to fetch %v: %v", u, resp.Status)
		}

		var metadata map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&metadata); err != nil {
			return nil, err
		}

		return metadata, nil
	}

	for _, u := range uris {
		if metadata, err1 := getMetatadatFunc(u); err1 != nil {
			multierr.AppendInto(&err, err1)
		} else {
			return metadata, nil
		}
	}

	return nil, err
}

func getMetadataURIs(server string) ([]string, error) {
	var uris []string

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, fmt.Errorf("failed to parse authorization server URL: %w", err)
	}

	serverURL.Path = AuthorizationServerMetadataPath
	uris = append(uris, serverURL.String())
	serverURL.Path = OIDCMetadataPath
	uris = append(uris, serverURL.String())
	return uris, nil
}
