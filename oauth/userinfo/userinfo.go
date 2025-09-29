package userinfo

import "context"

type UserInfo struct {
	Subject string
	Email   string
	Token   string
}

type TokenValidator interface {
	ValidateToken(ctx context.Context, token string) (*UserInfo, error)
}

type TokenValidatorFunc func(ctx context.Context, token string) (*UserInfo, error)

func (f TokenValidatorFunc) ValidateToken(ctx context.Context, token string) (*UserInfo, error) {
	return f(ctx, token)
}
