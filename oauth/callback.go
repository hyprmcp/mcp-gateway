package oauth

import (
	"net/http"
	"net/url"

	"github.com/hyprmcp/mcp-gateway/config"
)

const (
	CallbackPath = "/oauth/callback"
)

func NewCallbackHandler(config *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state := r.FormValue("state")
		if state == "" {
			http.Error(w, "missing state", http.StatusBadRequest)
			return
		}

		if redirectURIStr, ok := stateMap[state]; !ok {
			http.Error(w, "invalid state (no redirect URI)", http.StatusBadRequest)
			return
		} else {
			p, err := url.Parse(redirectURIStr)
			if err != nil {
				http.Error(w, "invalid redirect URI", http.StatusBadRequest)
				return
			}
			p.RawQuery = r.URL.Query().Encode()
			http.Redirect(w, r, p.String(), http.StatusFound)
		}
	})
}
