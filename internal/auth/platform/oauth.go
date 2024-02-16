package platform

import (
	"context"
	"fmt"

	"study-planner/internal/auth"

	"golang.org/x/oauth2"
)

type UserInfoSupplier interface {
	GetUserInfo(token *oauth2.Token) (*auth.UserInfo, error)
}

type CodeRequest struct {
	Code string `json:"code"`
}

type AuthenticationConfig struct {
	AuthenticationUrl string `json:"authenticationUrl"`
}

type OAuthPlatform struct {
	config   *oauth2.Config
	supplier UserInfoSupplier
}

func NewOAuthPlatform(config *oauth2.Config, supplier UserInfoSupplier) *OAuthPlatform {
	return &OAuthPlatform{config: config, supplier: supplier}
}

func (p *OAuthPlatform) GetAuthenticationConfig(_ context.Context) (*AuthenticationConfig, error) {
	return &AuthenticationConfig{
		AuthenticationUrl: p.config.AuthCodeURL(""),
	}, nil
}

func (p *OAuthPlatform) Authenticate(ctx context.Context, request *CodeRequest) (*auth.UserInfo, error) {
	token, err := p.config.Exchange(ctx, request.Code)
	if err != nil {
		return nil, fmt.Errorf("code exchange: %s", err)
	}

	return p.supplier.GetUserInfo(token)
}
