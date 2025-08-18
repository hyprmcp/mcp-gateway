package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/jetski-sh/mcp-gateway/config"
)

func NewProxyHandler(config *config.Proxy) http.Handler {
	url := (*url.URL)(config.Http.Url)

	return &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.Out.URL.Scheme = url.Scheme
			r.Out.URL.Host = url.Host
			r.Out.URL.Path = url.Path
			r.Out.URL.RawPath = url.RawPath
			if r.Out.URL.RawQuery == "" || url.RawQuery == "" {
				r.Out.URL.RawQuery = r.Out.URL.RawQuery + url.RawQuery
			} else {
				r.Out.URL.RawQuery = url.RawQuery + "&" + r.Out.URL.RawQuery
			}
			r.Out.Host = ""
		},
		Transport: &mcpAwareTransport{
			config: config.Webhook,
		},
	}
}
