package delivery

import (
	"study-planner/internal/discipline"
	"study-planner/pkg/httputil"

	"github.com/gofiber/fiber/v2"
)

type DisciplineController struct {
	disciplineRepository discipline.Repository
}

func NewDisciplineController(disciplineRepository discipline.Repository) *DisciplineController {
	return &DisciplineController{disciplineRepository: disciplineRepository}
}

func (c *DisciplineController) GetDisciplines(ctx *fiber.Ctx) (*[]discipline.Discipline, error) {
	curriculumId, err := httputil.ExtractId(ctx, "curriculum_id")
	if err != nil {
		return nil, err
	}

	return c.disciplineRepository.GetDisciplines(curriculumId)
}

func (c *DisciplineController) GetDiscipline(ctx *fiber.Ctx) (*discipline.Discipline, error) {
	curriculumId, err := httputil.ExtractId(ctx, "curriculum_id")
	if err != nil {
		return nil, err
	}

	disciplineId, err := httputil.ExtractId(ctx, "discipline_id")
	if err != nil {
		return nil, err
	}

	return c.disciplineRepository.GetDiscipline(curriculumId, disciplineId)
}

func (c *DisciplineController) GetDisciplineLinks(ctx *fiber.Ctx) (*[]discipline.Link, error) {
	disciplineId, err := httputil.ExtractId(ctx, "discipline_id")
	if err != nil {
		return nil, err
	}

	return c.disciplineRepository.GetDisciplineLinks(disciplineId)
}
