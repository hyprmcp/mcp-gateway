package oauth

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hyprmcp/mcp-gateway/config"
	"github.com/hyprmcp/mcp-gateway/metadata"
	"github.com/hyprmcp/mcp-gateway/oauth/authorization"
)

const AuthorizationPath = "/oauth/authorize"

var stateMap = make(map[string]string)

func NewAuthorizationHandler(
	config *config.Config,
	meta metadata.Metadata,
	fn authorization.EditQueryFunc,
) (http.Handler, error) {
	if authorizationEndpointStr, ok := meta["authorization_endpoint"].(string); !ok {
		return nil, errors.New("authorization metadata is missing authorization_endpoint field")
	} else if _, err := url.Parse(authorizationEndpointStr); err != nil {
		return nil, fmt.Errorf("could not parse authorization endpoint: %w", err)
	} else {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			q := r.URL.Query()
			if err := fn(q); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			upstreamAuthorizationURI, _ := url.Parse(authorizationEndpointStr)
			upstreamAuthorizationURI.RawQuery = q.Encode()
			http.Redirect(w, r, upstreamAuthorizationURI.String(), http.StatusFound)
		}), nil
	}
}
