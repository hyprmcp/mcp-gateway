package oauth

import (
	"io"
	"net/http"
	"net/url"

	"github.com/jetski-sh/mcp-proxy/log"
)

const AuthorizationServerMetadataPath = "/.well-known/oauth-authorization-server"
const OIDCMetadataPath = "/.well-known/openid-configuration"

func NewAuthorizationServerMetadataHandler(authorizationServer string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url, err := url.JoinPath(authorizationServer, OIDCMetadataPath)
		if err != nil {
			log.Get(r.Context()).Error(err, "failed to construct authorization server metadata URL", "base", authorizationServer)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		response, err := http.Get(url)
		if err != nil {
			log.Get(r.Context()).Error(err, "failed to fetch authorization server metadata", "url", url)
			w.WriteHeader(http.StatusBadGateway)
		} else if response.StatusCode >= http.StatusBadRequest {
			log.Get(r.Context()).Error(err, "authorization server returned error", "status", response.StatusCode, "url", response.Request.URL)
			w.WriteHeader(response.StatusCode)
			return
		}
		defer response.Body.Close()
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.Copy(w, response.Body)
	})
}
