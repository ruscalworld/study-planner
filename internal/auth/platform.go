package auth

import "context"

type Platform[C, T comparable] interface {
	GetAuthenticationConfig(ctx context.Context) (*C, error)
	Authenticate(ctx context.Context, token *T) (*UserInfo, error)
}

type UserInfo struct {
	ExternalID string
	Platform   string
	Name       string
	AvatarURL  string
}
