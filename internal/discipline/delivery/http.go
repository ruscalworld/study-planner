package delivery

import (
	"study-planner/internal/discipline"
	"study-planner/internal/user"
	"study-planner/pkg/httputil"

	"github.com/gofiber/fiber/v2"
)

type DisciplineController struct {
	disciplineRepository discipline.Repository
	userRepository       user.Repository
}

func NewDisciplineController(disciplineRepository discipline.Repository, userRepository user.Repository) *DisciplineController {
	return &DisciplineController{disciplineRepository: disciplineRepository, userRepository: userRepository}
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

func (c *DisciplineController) GetDisciplineProgress(ctx *fiber.Ctx) (*[]user.ScopedTaskProgress, error) {
	userId := ctx.Locals("userid").(int64)

	disciplineId, err := httputil.ExtractId(ctx, "discipline_id")
	if err != nil {
		return nil, err
	}

	return c.userRepository.GetDisciplineProgress(userId, disciplineId)
}
