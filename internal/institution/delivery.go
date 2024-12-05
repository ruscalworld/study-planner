package institution

import (
	"github.com/ruscalworld/study-planner/internal/curriculum"

	"github.com/gofiber/fiber/v2"
)

type Controller interface {
	GetInstitutions(ctx *fiber.Ctx) (*[]Institution, error)
	GetInstitution(ctx *fiber.Ctx) (*Institution, error)
	GetCurriculums(ctx *fiber.Ctx) (*[]curriculum.Curriculum, error)
	CreateCurriculum(ctx *fiber.Ctx, request *CreateCurriculumParams) (*curriculum.Curriculum, error)
}
