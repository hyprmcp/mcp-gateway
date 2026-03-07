package oauth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/go-chi/httprate"
	"github.com/hyprmcp/mcp-gateway/config"
	"github.com/hyprmcp/mcp-gateway/htmlresponse"
	"github.com/hyprmcp/mcp-gateway/log"
	"github.com/lestrrat-go/httprc/v3"
	"github.com/lestrrat-go/httprc/v3/errsink"
	"github.com/lestrrat-go/httprc/v3/tracesink"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

type Manager struct {
	jwkSet         jwk.Set
	config         *config.Config
	authServerMeta map[string]any
}

func NewManager(ctx context.Context, config *config.Config) (*Manager, error) {
	log := log.Get(ctx)

	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cache, err := jwk.NewCache(ctx, httprc.NewClient(
		httprc.WithTraceSink(tracesink.Func(func(ctx context.Context, s string) { log.V(1).Info(s) })),
		httprc.WithErrorSink(errsink.NewFunc(func(ctx context.Context, err error) { log.V(1).Error(err, "httprc.NewClient error") })),
	))
	if err != nil {
		return nil, fmt.Errorf("jwk cache creation error: %w", err)
	}

	meta, err := GetMedatata(config.Authorization.Server)
	if err != nil {
		return nil, fmt.Errorf("authorization server metadata error: %w", err)
	}

	jwksURIRaw, ok := meta["jwks_uri"].(string)
	if !ok {
		return nil, errors.New("no jwks_uri")
	}

	// Build internal JWKS URI using the auth server address to avoid a
	// circular dependency when the auth server's issuer is set to the
	// public gateway URL (the metadata would advertise the gateway's own
	// URL as jwks_uri, causing the gateway to fetch from itself).
	baseURL, err := url.Parse(config.Authorization.Server)
	if err != nil {
		return nil, fmt.Errorf("invalid authorization server URL %q: %w", config.Authorization.Server, err)
	}
	jwksURL, err := url.Parse(jwksURIRaw)
	if err != nil {
		return nil, fmt.Errorf("invalid jwks_uri %q: %w", jwksURIRaw, err)
	}
	jwksPath := jwksURL.Path
	if jwksPath == "" || jwksPath == "/" {
		jwksPath = "/keys"
	}
	internalJWKSURI, err := url.JoinPath(baseURL.String(), jwksPath)
	if err != nil {
		return nil, fmt.Errorf("failed to construct internal JWKS URI: %w", err)
	}

	if err := cache.Register(timeoutCtx, internalJWKSURI,
		jwk.WithMinInterval(10*time.Second),
		jwk.WithMaxInterval(5*time.Minute),
	); err != nil {
		return nil, fmt.Errorf("jwks registration error: %w", err)
	}
	if _, err := cache.Refresh(timeoutCtx, internalJWKSURI); err != nil {
		return nil, fmt.Errorf("jwks refresh error: %w", err)
	}
	s, err := cache.CachedSet(internalJWKSURI)
	if err != nil {
		return nil, fmt.Errorf("jwks cache set error: %w", err)
	}

	return &Manager{jwkSet: s, config: config, authServerMeta: meta}, nil
}

func (mgr *Manager) Register(mux *http.ServeMux) error {
	mux.Handle(ProtectedResourcePath, NewProtectedResourceHandler(mgr.config))

	if mgr.config.Authorization.ServerMetadataProxyEnabled {
		mux.Handle(AuthorizationServerMetadataPath, NewAuthorizationServerMetadataHandler(mgr.config))
	}

	if mgr.config.Authorization.GetDynamicClientRegistration().Enabled {
		if handler, err := NewDynamicClientRegistrationHandler(mgr.config, mgr.authServerMeta); err != nil {
			return err
		} else {
			rateLimiter := httprate.LimitByRealIP(3, 10*time.Minute)
			mux.Handle(DynamicClientRegistrationPath, rateLimiter(handler))
		}
	}

	if mgr.config.Authorization.AuthorizationProxyEnabled {
		if handler, err := NewAuthorizationHandler(mgr.config, mgr.authServerMeta); err != nil {
			return err
		} else {
			mux.Handle(AuthorizationPath, handler)
		}

		// Reverse-proxy auth server endpoints so that external clients can
		// reach them through the public gateway URL. Paths are derived from
		// the authorization server metadata (token, jwks, userinfo, etc.).
		authURL, err := url.Parse(mgr.config.Authorization.Server)
		if err != nil {
			return fmt.Errorf("failed to parse auth server URL: %w", err)
		}
		authProxy := &httputil.ReverseProxy{
			Rewrite: func(r *httputil.ProxyRequest) {
				r.Out.URL.Scheme = authURL.Scheme
				r.Out.URL.Host = authURL.Host
				r.Out.Host = ""
			},
		}
		for _, path := range authServerProxyPaths(mgr.authServerMeta) {
			mux.Handle(path, authProxy)
		}
	}

	return nil
}

func (mgr *Manager) Handler(next http.Handler) http.Handler {
	htmlHandler := htmlresponse.NewHandler(mgr.config, true)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawToken :=
			strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(r.Header.Get("Authorization")), "Bearer"))
		if token, err := jwt.ParseString(rawToken, jwt.WithKeySet(mgr.jwkSet)); err != nil {
			htmlHandler.Handler(mgr.unauthorizedHandler()).ServeHTTP(w, r)
		} else {
			next.ServeHTTP(w, r.WithContext(TokenContext(r.Context(), token, rawToken)))
		}
	})
}

func (mgr *Manager) unauthorizedHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Bearer resource_metadata="%s"`, mgr.getMetadataURL(r.URL)))
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func (mgr *Manager) getMetadataURL(u *url.URL) *url.URL {
	metadataURL, _ := url.Parse(mgr.config.Host.String())
	metadataURL.Path = ProtectedResourcePath
	metadataURL = metadataURL.JoinPath(u.Path)
	return metadataURL
}

// authServerProxyPaths returns the HTTP paths that should be reverse-proxied
// to the authorization server. Standard endpoint paths are extracted from the
// authorization server metadata. The authorization endpoint is also registered
// as a prefix pattern (trailing slash) so that connector sub-paths (e.g.
// /auth/github) are proxied as well.
func authServerProxyPaths(meta map[string]any) []string {
	endpointKeys := []string{
		"token_endpoint",
		"jwks_uri",
		"userinfo_endpoint",
		"introspection_endpoint",
		"revocation_endpoint",
		"device_authorization_endpoint",
		"end_session_endpoint",
	}

	seen := make(map[string]bool)
	var paths []string
	add := func(p string) {
		if p != "" && p != "/" && !seen[p] {
			seen[p] = true
			paths = append(paths, p)
		}
	}

	for _, key := range endpointKeys {
		if s, ok := meta[key].(string); ok {
			if u, err := url.Parse(s); err == nil && strings.HasPrefix(u.Path, "/") {
				add(u.Path)
			}
		}
	}

	// Register the authorization endpoint both as an exact match and as a
	// prefix pattern so that connector-specific sub-paths (e.g.
	// /auth/github) are also proxied.
	if s, ok := meta["authorization_endpoint"].(string); ok {
		if u, err := url.Parse(s); err == nil && strings.HasPrefix(u.Path, "/") {
			add(u.Path)
			add(strings.TrimRight(u.Path, "/") + "/")
		}
	}

	// Common auth server paths that are not advertised in metadata but are
	// required for the OAuth flow when the auth server uses external
	// identity provider connectors (e.g. the IdP redirects back to
	// /callback after user authentication).
	add("/callback")
	add("/approval")

	return paths
}
