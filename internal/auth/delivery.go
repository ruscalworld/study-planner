package auth

import (
	"study-planner/internal/user"

	"github.com/gofiber/fiber/v2"
)

type Controller[C, T comparable] interface {
	GetCurrentUser(ctx *fiber.Ctx) (*user.User, error)
	Authenticate(ctx *fiber.Ctx, credentials *T) (*Token, error)
	GetConfig(ctx *fiber.Ctx) (*C, error)
}
