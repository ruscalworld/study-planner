package delivery

import (
	"strings"
	"study-planner/internal/auth"
	"study-planner/pkg/stderrors"

	"github.com/gofiber/fiber/v2"
)

var (
	ErrInvalidHeader = stderrors.Unauthorized("invalid Authorization header")
	ErrInvalidToken  = stderrors.Unauthorized("invalid token")
)

func NewMiddleware(manager auth.Manager) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		token, err := extractToken(ctx)
		if err != nil {
			return err
		}

		tokenInfo, err := manager.Authorize(token)
		if err != nil {
			return stderrors.Populate(ErrInvalidToken, err)
		}

		ctx.Locals("userid", tokenInfo.UserId)
		return ctx.Next()
	}
}

func extractToken(ctx *fiber.Ctx) (*auth.Token, error) {
	header := ctx.Get("Authorization")
	parts := strings.Split(header, " ")

	if len(parts) != 2 {
		return nil, ErrInvalidHeader
	}

	return &auth.Token{
		TokenType:   auth.TokenType(parts[0]),
		AccessToken: parts[1],
	}, nil
}
