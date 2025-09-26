package oauth

import (
	"context"

	"github.com/hyprmcp/mcp-gateway/oauth/userinfo"
)

type userInfoKey struct{}

func UserInfoContext(parent context.Context, userInfo *userinfo.UserInfo) context.Context {
	parent = context.WithValue(parent, userInfoKey{}, userInfo)
	return parent
}

func GetUserInfo(ctx context.Context) *userinfo.UserInfo {
	if val, ok := ctx.Value(userInfoKey{}).(*userinfo.UserInfo); ok {
		return val
	} else {
		return nil
	}
}
