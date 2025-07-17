package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/jetski-sh/mcp-proxy/config"
)

func NewProxyHandler(url *url.URL, webhook *config.Webhook) http.Handler {
	return &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(url)
		},
		Transport: &mcpAwareTransport{
			config: webhook,
		},
	}
}
