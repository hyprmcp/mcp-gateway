package oauth

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func RewriteSetOriginalURL(r *httputil.ProxyRequest) {
	r.Out = r.Out.WithContext(WithOriginalURL(r.Out.Context(), r.In.URL))
}

func (mgr *Manager) UpdateWWWAuthenticateHeader(resp *http.Response) error {
	if resp.StatusCode == http.StatusUnauthorized {
		realRequestURL := GetOriginalURL(resp.Request.Context())
		upstreamMetaURL, _ := url.Parse(resp.Request.URL.String())
		upstreamMetaURL.Path, _ = url.JoinPath(ProtectedResourcePath, upstreamMetaURL.Path)
		upstreamMetaURLStr := upstreamMetaURL.String()

		if value := resp.Header.Get("WWW-Authenticate"); value != "" {
			valueParts := strings.Split(value, " ")
			for i, part := range valueParts {
				if after, ok := strings.CutPrefix(part, "resource_metadata="); ok {
					if resourceMetadataOrig := strings.Trim(after, `"`); resourceMetadataOrig != "" {
						upstreamMetaURLStr = resourceMetadataOrig
						valueParts[i] = fmt.Sprintf(`resource_metadata="%s"`, mgr.getMetadataURL(realRequestURL))
						resp.Header.Set("WWW-Authenticate", strings.Join(valueParts, " "))
						break
					}
				}
			}
		}

		// We also store the metadata URL if it was not found in the WWW-Authenticate header.
		// This is necessary to ensure that we don't return the wrong metadata later when the
		// client calls our metadata endpoint because of some heuristic.
		upstreamMetadataURLs.Store(strings.Trim(realRequestURL.Path, `/`), upstreamMetaURLStr)
	}

	return nil
}
