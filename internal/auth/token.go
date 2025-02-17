package auth

import (
	"time"

	"github.com/ruscalworld/study-planner/internal/auth/refresh"
	"github.com/ruscalworld/study-planner/internal/user"
)

type TokenType string

const (
	TokenBearer TokenType = "Bearer"
)

type Token struct {
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	TokenType    TokenType `json:"tokenType"`
}

type TokenProvider interface {
	MakeToken(u *user.User) (*Token, error)
	Verify(token *Token) (*TokenInfo, error)
	UseRefreshToken(rawToken []byte) (*refresh.Token, error)
}

type TokenInfo struct {
	UserId    int64
	ExpiresAt time.Time
}
