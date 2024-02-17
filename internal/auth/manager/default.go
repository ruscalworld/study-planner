package manager

import (
	"errors"
	"fmt"
	"time"

	"study-planner/internal/auth"
	"study-planner/internal/user"
)

type Default struct {
	userRepository user.Repository
	tokenProvider  auth.TokenProvider
}

func NewAuthManager(userRepository user.Repository, tokenProvider auth.TokenProvider) *Default {
	return &Default{userRepository: userRepository, tokenProvider: tokenProvider}
}

func (d *Default) Authenticate(userInfo *auth.UserInfo) (*auth.Token, error) {
	u, err := d.userRepository.GetUserByExternalId(userInfo.ExternalID)
	if err != nil && !errors.Is(err, user.ErrUnknownUser) {
		return nil, err
	}

	if u != nil {
		if u.Platform != userInfo.Platform {
			return nil, fmt.Errorf("user was registered with another platform: %s", u.Platform)
		}

		return d.tokenProvider.MakeToken(u)
	}

	u = &user.User{
		Name:       userInfo.Name,
		AvatarURL:  userInfo.AvatarURL,
		Platform:   userInfo.Platform,
		ExternalID: userInfo.ExternalID,
		CreatedAt:  time.Now(),
	}

	err = d.userRepository.RegisterUser(u)
	if err != nil {
		return nil, fmt.Errorf("registering user: %s", err)
	}

	return d.tokenProvider.MakeToken(u)
}

func (d *Default) Authorize(token *auth.Token) (*auth.TokenInfo, error) {
	tokenInfo, err := d.tokenProvider.Verify(token)
	if err != nil {
		return nil, err
	}

	if tokenInfo.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("token is expired")
	}

	return tokenInfo, nil
}

func (d *Default) GetTokenProvider() auth.TokenProvider {
	return d.tokenProvider
}
