package authorization

import (
	"crypto/rand"
	"net/url"
	"slices"
	"strings"

	"github.com/hyprmcp/mcp-gateway/metadata"
	"github.com/hyprmcp/mcp-gateway/oauth/callback"
)

type EditQueryFunc func(q url.Values) error

func EditChain(manipulators ...EditQueryFunc) EditQueryFunc {
	return func(q url.Values) error {
		for _, manipulator := range manipulators {
			if err := manipulator(q); err != nil {
				return err
			}
		}
		return nil
	}
}

func RequiredScopes(requiredScopes []string, meta metadata.Metadata) EditQueryFunc {
	actualRequiredScopes := requiredScopes
	if meta != nil {
		supportedScopes := meta.GetSupportedScopes()
		actualRequiredScopes = slices.DeleteFunc(
			requiredScopes,
			func(s string) bool { return !slices.Contains(supportedScopes, s) },
		)
	}

	return func(q url.Values) error {
		scopes := q.Get("scope")
		for _, scope := range actualRequiredScopes {
			if !strings.Contains(scopes, scope) {
				scopes = strings.TrimSpace(scopes + " " + scope)
			}
		}
		q.Set("scope", scopes)
		return nil
	}
}

func RedirectURI(hostURI string, store callback.URIStore) EditQueryFunc {
	return func(q url.Values) error {
		if origRedirectURI := q.Get("redirect_uri"); origRedirectURI != "" {
			state := q.Get("state")
			if state != "" {
				state = rand.Text()
				q.Set("state", state)
			}

			store.Set(state, origRedirectURI)
		}

		overrideRedirectURI, _ := url.Parse(hostURI)
		overrideRedirectURI.Path = callback.CallbackPath
		q.Set("redirect_uri", overrideRedirectURI.String())
		return nil
	}
}
