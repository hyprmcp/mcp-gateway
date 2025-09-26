package userinfo

import (
	"context"
	"net/http"
	"strconv"

	"github.com/google/go-github/v75/github"
)

func ValidateGitHub(client *http.Client) TokenValidatorFunc {
	githubClient := github.NewClient(client)

	return func(ctx context.Context, token string) (*UserInfo, error) {
		// Passing an empty string to Users.Get returns the authenticated user's information.
		if user, _, err := githubClient.WithAuthToken(token).Users.Get(ctx, ""); err != nil {
			return nil, err
		} else {
			return &UserInfo{
				Subject: strconv.FormatInt(user.GetID(), 10),
				Email:   user.GetEmail(),
				Token:   token,
			}, nil
		}
	}
}
