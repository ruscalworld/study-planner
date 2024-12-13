package curriculum

import (
	"github.com/gofiber/fiber/v2"
)

type Controller interface {
	GetCurriculum(ctx *fiber.Ctx) (*Curriculum, error)
	CreateCurriculum(ctx *fiber.Ctx, params *Params) (*Curriculum, error)
	UpdateCurriculum(ctx *fiber.Ctx, params *Params) (*Curriculum, error)
	DeleteCurriculum(ctx *fiber.Ctx) (*Curriculum, error)

	GetCurriculumCodes(ctx *fiber.Ctx) (*[]Code, error)
	CreateCurriculumCode(ctx *fiber.Ctx, params *CreateCodeParams) (*Code, error)
	DeleteCurriculumCode(ctx *fiber.Ctx) (*any, error)

	GetCurriculumUsers(ctx *fiber.Ctx) (*[]User, error)
	UpdateCurriculumUser(ctx *fiber.Ctx, params *UpdateCurriculumUserParams) (*any, error)
	DeleteCurriculumUser(ctx *fiber.Ctx) (*any, error)
}
