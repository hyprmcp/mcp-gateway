package oauth

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/hyprmcp/mcp-gateway/config"
)

const AuthorizationPath = "/oauth/authorize"

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
			redirectURI, _ := url.Parse(authorizationEndpointStr)
			q := r.URL.Query()
			scopes := q.Get("scope")
			for _, scope := range requiredScopes {
				if !strings.Contains(scopes, scope) {
					scopes = strings.TrimSpace(scopes + " " + scope)
				}
			}
			q.Set("scope", scopes)
			redirectURI.RawQuery = q.Encode()
			http.Redirect(w, r, redirectURI.String(), http.StatusFound)
		}), nil
	}
}
