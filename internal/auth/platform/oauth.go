package platform

import (
	"context"
	"fmt"

	"github.com/ruscalworld/study-planner/internal/auth"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type UserInfoSupplier interface {
	GetUserInfo(token *oauth2.Token) (*auth.UserInfo, error)
}

type CodeRequest struct {
	Code    string `json:"code"`
	IdToken string `json:"idToken"`
}

type AuthenticationConfig struct {
	AuthenticationUrl string `json:"authenticationUrl"`
}

type OAuthPlatform struct {
	config   *oauth2.Config
	provider *oidc.Provider
	verifier *oidc.IDTokenVerifier
}

func NewOAuthPlatform(config *oauth2.Config, oidIssuer string) (*OAuthPlatform, error) {
	provider, err := oidc.NewProvider(context.Background(), oidIssuer)
	if err != nil {
		return nil, err
	}

	return &OAuthPlatform{
		config:   config,
		provider: provider,
		verifier: provider.Verifier(&oidc.Config{ClientID: config.ClientID}),
	}, nil
}

func (p *OAuthPlatform) GetAuthenticationConfig(_ context.Context) (*AuthenticationConfig, error) {
	return &AuthenticationConfig{
		AuthenticationUrl: p.config.AuthCodeURL(""),
	}, nil
}

func (p *OAuthPlatform) Authenticate(ctx context.Context, request *CodeRequest) (*auth.UserInfo, error) {
	if request.IdToken == "" {
		token, err := p.config.Exchange(ctx, request.Code)
		if err != nil {
			return nil, fmt.Errorf("code exchange: %s", err)
		}

		rawIdToken, ok := token.Extra("id_token").(string)
		if !ok {
			return nil, fmt.Errorf("id_token is missing or is not a string")
		}

		request.IdToken = rawIdToken
	}

	idToken, err := p.verifier.Verify(ctx, request.IdToken)
	if err != nil {
		return nil, fmt.Errorf("id token verification: %s", err)
	}

	var u GoogleUser
	if err := idToken.Claims(&u); err != nil {
		return nil, fmt.Errorf("parsing claims: %s", err)
	}

	return &auth.UserInfo{
		ExternalID: u.ID,
		Platform:   "google",
		Name:       u.Name,
		AvatarURL:  u.Picture,
	}, nil
}
