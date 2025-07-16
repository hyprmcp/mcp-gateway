package oauth

import (
	"context"

	"github.com/lestrrat-go/jwx/v3/jwt"
)

type keyType struct{}

var key keyType

func AddTokenToContext(parent context.Context, token jwt.Token) context.Context {
	return context.WithValue(parent, key, token)
}

func TokenFromContext(ctx context.Context) jwt.Token {
	return ctx.Value(key).(jwt.Token)
}
