package oauth

import (
	"context"
	"net/url"

	"github.com/lestrrat-go/jwx/v3/jwt"
)

type tokenKey struct{}
type rawTokenKey struct{}
type originalURLKey struct{}

func TokenContext(parent context.Context, token jwt.Token, rawToken string) context.Context {
	parent = context.WithValue(parent, tokenKey{}, token)
	parent = context.WithValue(parent, rawTokenKey{}, rawToken)
	return parent
}

func GetToken(ctx context.Context) jwt.Token {
	if val, ok := ctx.Value(tokenKey{}).(jwt.Token); ok {
		return val
	} else {
		return nil
	}
}

func GetRawToken(ctx context.Context) string {
	if val, ok := ctx.Value(rawTokenKey{}).(string); ok {
		return val
	} else {
		return ""
	}
}

func WithOriginalURL(ctx context.Context, url *url.URL) context.Context {
	return context.WithValue(ctx, originalURLKey{}, url)
}

func GetOriginalURL(ctx context.Context) *url.URL {
	if url, ok := ctx.Value(originalURLKey{}).(*url.URL); ok {
		return url
	}

	return nil
}
