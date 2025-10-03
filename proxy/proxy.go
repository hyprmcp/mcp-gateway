package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/hyprmcp/mcp-gateway/config"
	"github.com/hyprmcp/mcp-gateway/oauth"
	"github.com/hyprmcp/mcp-gateway/proxy/proxyutil"
)

func NewProxyHandler(config *config.Proxy, modifyResponse func(*http.Response) error) http.Handler {
	url := (*url.URL)(config.Http.Url)

	return &httputil.ReverseProxy{
		Rewrite: proxyutil.RewriteChain(
			proxyutil.RewriteFullFunc(url),
			oauth.RewriteSetOriginalURL,
		),
		ModifyResponse: proxyutil.ModifyResponseChain(modifyResponse, proxyutil.RemoveCORSHeaders),
		Transport: &mcpAwareTransport{
			config: config,
		},
	}
}
