package oauth

import (
	"testing"
)

func TestAuthServerProxyPaths(t *testing.T) {
	tests := []struct {
		name     string
		meta     map[string]any
		wantLen  int
		contains []string
		excludes []string
	}{
		{
			name: "extracts standard endpoint paths from metadata",
			meta: map[string]any{
				"issuer":                        "https://example.com",
				"authorization_endpoint":        "https://example.com/auth",
				"token_endpoint":                "https://example.com/token",
				"jwks_uri":                      "https://example.com/keys",
				"userinfo_endpoint":             "https://example.com/userinfo",
				"introspection_endpoint":        "https://example.com/token/introspect",
				"device_authorization_endpoint": "https://example.com/device/code",
			},
			contains: []string{
				"/token",
				"/keys",
				"/userinfo",
				"/token/introspect",
				"/device/code",
				"/auth",
				"/auth/",
				"/callback",
				"/approval",
			},
			excludes: []string{"/"},
		},
		{
			name: "does not duplicate paths",
			meta: map[string]any{
				"token_endpoint":    "https://example.com/token",
				"jwks_uri":         "https://example.com/token", // same path
				"userinfo_endpoint": "https://example.com/userinfo",
			},
			contains: []string{"/token", "/userinfo", "/callback", "/approval"},
		},
		{
			name:     "handles empty metadata",
			meta:     map[string]any{},
			contains: []string{"/callback", "/approval"},
			wantLen:  2,
		},
		{
			name: "skips non-string and unparseable values",
			meta: map[string]any{
				"token_endpoint":    42,
				"jwks_uri":         "https://example.com/keys",
				"userinfo_endpoint": "://invalid",
			},
			contains: []string{"/keys", "/callback", "/approval"},
		},
		{
			name: "registers authorization endpoint prefix for sub-paths",
			meta: map[string]any{
				"authorization_endpoint": "https://example.com/dex/auth",
			},
			contains: []string{"/dex/auth", "/dex/auth/"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paths := authServerProxyPaths(tt.meta)

			if tt.wantLen > 0 && len(paths) != tt.wantLen {
				t.Errorf("got %d paths, want %d: %v", len(paths), tt.wantLen, paths)
			}

			pathSet := make(map[string]bool)
			for _, p := range paths {
				if pathSet[p] {
					t.Errorf("duplicate path: %s", p)
				}
				pathSet[p] = true
			}

			for _, want := range tt.contains {
				if !pathSet[want] {
					t.Errorf("missing expected path %q in %v", want, paths)
				}
			}

			for _, exclude := range tt.excludes {
				if pathSet[exclude] {
					t.Errorf("unexpected path %q in %v", exclude, paths)
				}
			}
		})
	}
}
