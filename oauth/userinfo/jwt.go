package userinfo

import (
	"context"
	"fmt"
	"time"

	"github.com/hyprmcp/mcp-gateway/log"
	"github.com/lestrrat-go/httprc/v3"
	"github.com/lestrrat-go/httprc/v3/errsink"
	"github.com/lestrrat-go/httprc/v3/tracesink"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

func ValidateJWT(options ...jwt.ParseOption) TokenValidatorFunc {
	return func(ctx context.Context, rawToken string) (*UserInfo, error) {
		if token, err := jwt.ParseString(rawToken, options...); err != nil {
			return nil, err
		} else {
			userInfo := &UserInfo{Token: rawToken}
			userInfo.Subject, _ = token.Subject()
			_ = token.Get("email", &userInfo.Email)
			return userInfo, nil
		}
	}
}

// ValidateDynamicJWKS creates an instance of jwt.CachedSet which periodically refreshes the JWKS set from the
// given URI until the passed context is canceled.
func ValidateDynamicJWKS(ctx context.Context, jwksURI string) (TokenValidatorFunc, error) {
	log := log.Get(ctx)

	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if cache, err := jwk.NewCache(ctx, httprc.NewClient(
		httprc.WithTraceSink(tracesink.Func(func(ctx context.Context, s string) { log.V(1).Info(s) })),
		httprc.WithErrorSink(errsink.NewFunc(func(ctx context.Context, err error) { log.V(1).Error(err, "httprc.NewClient error") })),
	)); err != nil {
		return nil, fmt.Errorf("jwk cache creation error: %w", err)
	} else if err := cache.Register(
		timeoutCtx,
		jwksURI,
		jwk.WithMinInterval(10*time.Second),
		jwk.WithMaxInterval(5*time.Minute),
	); err != nil {
		return nil, fmt.Errorf("jwks registration error: %w", err)
	} else if _, err := cache.Refresh(timeoutCtx, jwksURI); err != nil {
		return nil, fmt.Errorf("jwks refresh error: %w", err)
	} else if s, err := cache.CachedSet(jwksURI); err != nil {
		return nil, fmt.Errorf("jwks cache set error: %w", err)
	} else {
		return ValidateJWT(jwt.WithKeySet(s)), nil
	}
}
