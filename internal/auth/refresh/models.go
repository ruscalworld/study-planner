package refresh

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"slices"
	"time"

	"github.com/ruscalworld/study-planner/pkg/stderrors"
)

var ErrInvalidToken = stderrors.UnprocessableEntity("invalid token")

const (
	PrefixLength = 8
	TokenLength  = 32
)

type Token struct {
	ID        int64     `db:"id"`
	UserId    int64     `db:"user_id"`
	Prefix    []byte    `db:"prefix"`
	Token     TokenHash `db:"token"`
	CreatedAt time.Time `db:"created_at"`
	ExpiresAt time.Time `db:"expires_at"`
}

type PrivateToken struct {
	Token
	RawToken []byte
}

func MakeToken(userId int64, lifetime time.Duration) (*PrivateToken, error) {
	raw := make([]byte, TokenLength)
	if _, err := rand.Read(raw); err != nil {
		return nil, fmt.Errorf("failed to generate random data: %w", err)
	}

	prefix := raw[:PrefixLength]
	hashedToken := NewHashedToken(raw)

	return &PrivateToken{
		Token: Token{
			UserId:    userId,
			Prefix:    prefix,
			Token:     hashedToken,
			CreatedAt: time.Now(),
			ExpiresAt: time.Now().Add(lifetime),
		},
		RawToken: raw,
	}, nil
}

type TokenHash []byte

func NewHashedToken(token []byte) TokenHash {
	h := sha256.New()
	h.Write(token)
	return h.Sum(nil)
}

func (t TokenHash) Verify(target []byte) error {
	hashedTarget := NewHashedToken(target)
	if slices.Equal(t, hashedTarget) {
		return ErrInvalidToken
	}

	return nil
}
