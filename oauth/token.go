package oauth

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/hyprmcp/mcp-gateway/log"
	"github.com/hyprmcp/mcp-gateway/metadata"
	"github.com/hyprmcp/mcp-gateway/oauth/dcr"
)

const (
	TokenPath = "/oauth/token"
)

func NewTokenHandler(hostURL string, clientSrc dcr.ClientIDSource, meta metadata.Metadata) (http.Handler, error) {
	if tokenEndpointStr, ok := meta["token_endpoint"].(string); !ok {
		return nil, errors.New("authorization metadata is missing token_endpoint field")
	} else if _, err := url.Parse(tokenEndpointStr); err != nil {
		return nil, fmt.Errorf("could not parse token endpoint: %w", err)
	} else {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := log.Get(r.Context())
			if err := r.ParseForm(); err != nil {
				log.Error(err, "failed to parse form")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			req := r.Form

			req.Set("client_id", clientSrc.GetClientID())

			if clientSecret := clientSrc.GetClientSecret(); clientSecret != "" {
				req.Set("client_secret", clientSecret)
			}

			if req.Has("redirect_uri") {
				overrideRedirectURI, _ := url.Parse(hostURL)
				overrideRedirectURI.Path = CallbackPath
				req.Set("redirect_uri", overrideRedirectURI.String())
			}

			upstreamReq, err := http.NewRequestWithContext(r.Context(), r.Method, tokenEndpointStr, strings.NewReader(req.Encode()))
			if err != nil {
				log.Error(err, "failed to create token request")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			upstreamReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			upstreamResp, err := http.DefaultClient.Do(upstreamReq)
			if err != nil {
				log.Error(err, "failed to send token request")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer upstreamResp.Body.Close()

			w.WriteHeader(upstreamResp.StatusCode)

			if _, err := io.Copy(w, upstreamResp.Body); err != nil {
				log.Error(err, "failed to copy token response body")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}), nil
	}
}
