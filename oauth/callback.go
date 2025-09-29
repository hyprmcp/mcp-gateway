package oauth

import (
	"net/http"

	"github.com/hyprmcp/mcp-gateway/config"
	"github.com/hyprmcp/mcp-gateway/oauth/callback"
)

func NewCallbackHandler(config *config.Config, store callback.URIStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if p, err := store.Get(r.FormValue("state")); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			p.RawQuery = r.URL.Query().Encode()
			http.Redirect(w, r, p.String(), http.StatusFound)
		}
	})
}
