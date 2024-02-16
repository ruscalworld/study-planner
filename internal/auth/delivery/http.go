package delivery

import (
	"study-planner/internal/auth"
	"study-planner/internal/user"
	"study-planner/pkg/stderrors"

	"github.com/gofiber/fiber/v2"
)

var (
	ErrInvalidCredentials = stderrors.Unauthorized("invalid credentials")
)

type AuthController[C, T comparable] struct {
	userRepository user.Repository
	platform       auth.Platform[C, T]
	authManager    auth.Manager
}

func NewAuthController[C, T comparable](userRepository user.Repository, platform auth.Platform[C, T], authManager auth.Manager) auth.Controller[C, T] {
	return &AuthController[C, T]{userRepository: userRepository, platform: platform, authManager: authManager}
}

func (a *AuthController[C, T]) GetCurrentUser(ctx *fiber.Ctx) (*user.User, error) {
	userId := ctx.Locals("userid").(int64)
	return a.userRepository.GetUserById(userId)
}

func (a *AuthController[C, T]) GetConfig(ctx *fiber.Ctx) (*C, error) {
	return a.platform.GetAuthenticationConfig(ctx.UserContext())
}

func (a *AuthController[C, T]) Authenticate(ctx *fiber.Ctx, credentials *T) (*auth.Token, error) {
	userInfo, err := a.platform.Authenticate(ctx.UserContext(), credentials)
	if err != nil {
		return nil, stderrors.Populate(ErrInvalidCredentials, err)
	}

	return a.authManager.Authenticate(userInfo)
}
