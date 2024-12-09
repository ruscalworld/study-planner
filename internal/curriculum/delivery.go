package curriculum

import "github.com/gofiber/fiber/v2"

type Controller interface {
	GetCurriculum(ctx *fiber.Ctx) (*Curriculum, error)
	CreateCurriculum(ctx *fiber.Ctx, request *CreateCurriculumParams) (*Curriculum, error)

	GetCurriculumCodes(ctx *fiber.Ctx) (*[]Code, error)
	CreateCurriculumCode(ctx *fiber.Ctx, request *CreateCodeParams) (*Code, error)
	DeleteCurriculumCode(ctx *fiber.Ctx) (*any, error)

	GetCurriculumUsers(ctx *fiber.Ctx) (*[]User, error)
	DeleteCurriculumUser(ctx *fiber.Ctx) (*any, error)
}
