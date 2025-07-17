package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/jetski-sh/mcp-proxy/config"
)

func NewProxyHandler(config *config.Proxy) http.Handler {
	url := (*url.URL)(config.Http.Url)

	return &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(url)
		},
		Transport: &mcpAwareTransport{
			config: config.Webhook,
		},
	}
}
