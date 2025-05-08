package oauth

import "context"

const (
	GOOGLE_PROVIDER string = "google"
	GITHUB_PROVIDER string = "github"
)

type Provider interface {
	GetAuthUrl() string
	Auth(ctx context.Context, code string) (*User, error)
}
