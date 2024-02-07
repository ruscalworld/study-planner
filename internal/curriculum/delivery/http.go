package delivery

import (
	"study-planner/internal/curriculum"
	"study-planner/pkg/httputil"

	"github.com/gofiber/fiber/v2"
)

type CurriculumController struct {
	curriculumRepository curriculum.Repository
}

func NewCurriculumController(curriculumRepository curriculum.Repository) *CurriculumController {
	return &CurriculumController{curriculumRepository: curriculumRepository}
}

func (c *CurriculumController) GetCurriculums(_ *fiber.Ctx) (*[]curriculum.Curriculum, error) {
	// TODO: user curriculums
	panic("not implemented")
}

func (c *CurriculumController) GetCurriculum(ctx *fiber.Ctx) (*curriculum.Curriculum, error) {
	id, err := httputil.ExtractId(ctx, "curriculum_id")
	if err != nil {
		return nil, err
	}

	return c.curriculumRepository.GetCurriculum(id)
}
