package auth

import (
	"time"

	"study-planner/internal/user"
)

type TokenType string

const (
	TokenBearer TokenType = "Bearer"
)

type Token struct {
	AccessToken string    `json:"accessToken"`
	TokenType   TokenType `json:"tokenType"`
}

type TokenProvider interface {
	MakeToken(u *user.User) (*Token, error)
	Verify(token *Token) (*TokenInfo, error)
}

type TokenInfo struct {
	UserId    int64
	ExpiresAt time.Time
}
