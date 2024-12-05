package user

import "github.com/gofiber/fiber/v2"

type Controller interface {
	GetUserCurriculums(ctx *fiber.Ctx) (*[]Curriculum, error)
	CreateUserCurriculum(ctx *fiber.Ctx, request *CreateUserCurriculumParams) (*any, error)
}
