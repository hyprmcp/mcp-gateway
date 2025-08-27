package oauth

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hyprmcp/mcp-gateway/config"
)

const AuthorizationPath = "/oauth/authorize"

func NewAuthorizationHandler(config *config.Config, meta map[string]any) (http.Handler, error) {
	if authorizationEndpointStr, ok := meta["authorization_endpoint"].(string); !ok {
		return nil, errors.New("authorization metadata is missing authorization_endpoint field")
	} else if _, err := url.Parse(authorizationEndpointStr); err != nil {
		return nil, fmt.Errorf("could not parse authorization endpoint: %w", err)
	} else {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			redirectURI, _ := url.Parse(authorizationEndpointStr)
			q := r.URL.Query()
			if scope := q.Get("scope"); !strings.Contains(scope, "openid") {
				q.Set("scope", strings.TrimSpace("openid "+scope))
			}
			redirectURI.RawQuery = q.Encode()
			http.Redirect(w, r, redirectURI.String(), http.StatusFound)
		}), nil
	}
}
