package oauth

import (
	"context"

	"github.com/lestrrat-go/jwx/v3/jwt"
)

type tokenKey struct{}
type rawTokenKey struct{}

func TokenContext(parent context.Context, token jwt.Token, rawToken string) context.Context {
	parent = context.WithValue(parent, tokenKey{}, token)
	parent = context.WithValue(parent, rawTokenKey{}, rawToken)
	return parent
}

func GetToken(ctx context.Context) jwt.Token {
	return ctx.Value(tokenKey{}).(jwt.Token)
}

func GetRawToken(ctx context.Context) string {
	return ctx.Value(rawTokenKey{}).(string)
}
