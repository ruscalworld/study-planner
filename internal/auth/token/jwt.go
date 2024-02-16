package token

import (
	"fmt"
	"strconv"
	"time"

	"study-planner/internal/auth"
	"study-planner/internal/user"

	"github.com/golang-jwt/jwt/v5"
)

const (
	jwtIssuer = "StudyPlanner"
	tokenType = auth.TokenBearer
)

type JwtTokenProvider struct {
	signingKey    []byte
	audience      string
	tokenLifetime time.Duration
}

func NewJwtTokenProvider(signingKey []byte, audience string, tokenLifetime time.Duration) *JwtTokenProvider {
	return &JwtTokenProvider{signingKey: signingKey, audience: audience, tokenLifetime: tokenLifetime}
}

type Claims struct {
	jwt.RegisteredClaims
}

func (j *JwtTokenProvider) MakeToken(u *user.User) (*auth.Token, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, j.makeClaims(u))

	signedToken, err := token.SignedString(j.signingKey)
	if err != nil {
		return nil, err
	}

	return &auth.Token{
		AccessToken: signedToken,
		TokenType:   tokenType,
	}, nil
}

func (j *JwtTokenProvider) makeClaims(u *user.User) *Claims {
	return &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    jwtIssuer,
			Audience:  jwt.ClaimStrings{j.audience},
			Subject:   fmt.Sprintf("%d", u.ID),
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(j.tokenLifetime)},
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
		},
	}
}

func (j *JwtTokenProvider) Verify(token *auth.Token) (*auth.TokenInfo, error) {
	if token.TokenType != tokenType {
		return nil, fmt.Errorf("unsupported token type \"%s\"", token.TokenType)
	}

	parsedToken, err := jwt.ParseWithClaims(token.AccessToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.signingKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %s", err)
	}

	tokenInfo, err := j.verifyClaims(parsedToken.Claims)
	if err != nil {
		return nil, err
	}

	return tokenInfo, nil
}

func (j *JwtTokenProvider) verifyClaims(claims jwt.Claims) (*auth.TokenInfo, error) {
	err := j.verifyIssuer(claims)
	if err != nil {
		return nil, fmt.Errorf("invalid issuer: %s", err)
	}

	err = j.verifyAudience(claims)
	if err != nil {
		return nil, fmt.Errorf("invalid audience: %s", err)
	}

	userId, err := j.verifySubject(claims)
	if err != nil {
		return nil, fmt.Errorf("invalid subject: %s", err)
	}

	expiresAt, err := claims.GetExpirationTime()
	if err != nil {
		return nil, fmt.Errorf("invalid expiration time: %s", err)
	}

	return &auth.TokenInfo{
		UserId:    userId,
		ExpiresAt: expiresAt.Time,
	}, nil
}

func (j *JwtTokenProvider) verifyIssuer(claims jwt.Claims) error {
	issuer, err := claims.GetIssuer()
	if err != nil {
		return err
	}

	if issuer != jwtIssuer {
		return fmt.Errorf("unsupported value")
	}

	return nil
}

func (j *JwtTokenProvider) verifyAudience(claims jwt.Claims) error {
	audience, err := claims.GetAudience()
	if err != nil {
		return err
	}

	if len(audience) != 1 || audience[0] != j.audience {
		return fmt.Errorf("unsupported value")
	}

	return nil
}

func (j *JwtTokenProvider) verifySubject(claims jwt.Claims) (int64, error) {
	rawUserId, err := claims.GetSubject()
	if err != nil {
		return 0, err
	}

	userId, err := strconv.ParseInt(rawUserId, 10, 64)
	if err != nil {
		return 0, err
	}

	return userId, nil
}
