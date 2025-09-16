package oauth

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/hyprmcp/mcp-gateway/config"
	"github.com/hyprmcp/mcp-gateway/log"
)

const AuthorizationPath = "/oauth/authorize"

var stateMap = make(map[string]string)

func NewAuthorizationHandler(config *config.Config, meta map[string]any) (http.Handler, error) {
	supportedScopes := getSupportedScopes(meta)
	var requiredScopes = slices.DeleteFunc(
		[]string{"openid", "profile", "email"},
		func(s string) bool { return !slices.Contains(supportedScopes, s) },
	)

	if authorizationEndpointStr, ok := meta["authorization_endpoint"].(string); !ok {
		return nil, errors.New("authorization metadata is missing authorization_endpoint field")
	} else if _, err := url.Parse(authorizationEndpointStr); err != nil {
		return nil, fmt.Errorf("could not parse authorization endpoint: %w", err)
	} else {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			upstreamAuthorizationURI, _ := url.Parse(authorizationEndpointStr)

			q := r.URL.Query()

			if config.Authorization.AuthorizationProxyEnabled {
				scopes := q.Get("scope")
				for _, scope := range requiredScopes {
					if !strings.Contains(scopes, scope) {
						scopes = strings.TrimSpace(scopes + " " + scope)
					}
				}
				q.Set("scope", scopes)
			}

			if config.Authorization.ClientSecret != "" {
				if origRedirectURI := q.Get("redirect_uri"); origRedirectURI != "" {
					state := q.Get("state")
					if state != "" {
						state = genRandom()
						q.Set("state", state)
					}
					log.Get(r.Context()).Info("storing redirect uri", "redirect_uri", origRedirectURI, "state", state)
					stateMap[state] = q.Get("redirect_uri")
				}

				overrideRedirectURI, _ := url.Parse(config.Host.String())
				overrideRedirectURI.Path = CallbackPath
				q.Set("redirect_uri", overrideRedirectURI.String())
			}

			upstreamAuthorizationURI.RawQuery = q.Encode()
			http.Redirect(w, r, upstreamAuthorizationURI.String(), http.StatusFound)
		}), nil
	}
}
