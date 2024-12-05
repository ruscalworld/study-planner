package delivery

import (
	"github.com/ruscalworld/study-planner/internal/auth"
	"github.com/ruscalworld/study-planner/internal/curriculum"
	"github.com/ruscalworld/study-planner/internal/discipline"
	"github.com/ruscalworld/study-planner/internal/user"

	"github.com/ruscalworld/study-planner/pkg/httputil"

	"github.com/gofiber/fiber/v2"
)

type DisciplineController struct {
	curriculumRepository curriculum.Repository
	disciplineRepository discipline.Repository
	userRepository       user.Repository
}

func NewDisciplineController(
	curriculumRepository curriculum.Repository,
	disciplineRepository discipline.Repository,
	userRepository user.Repository,
) *DisciplineController {
	return &DisciplineController{
		curriculumRepository: curriculumRepository,
		disciplineRepository: disciplineRepository,
		userRepository:       userRepository,
	}
}

func (c *DisciplineController) GetDisciplines(ctx *fiber.Ctx) (*[]discipline.Discipline, error) {
	curriculumId, err := httputil.ExtractId(ctx, "curriculum_id")
	if err != nil {
		return nil, err
	}

	err = auth.Authorize(ctx, c.curriculumRepository.GetCurriculumPrivileges, auth.LevelPublic, auth.ActionRead, curriculumId)
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

	err = auth.Authorize(ctx, c.curriculumRepository.GetCurriculumPrivileges, auth.LevelPublic, auth.ActionRead, curriculumId)
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

	err = auth.Authorize(ctx, c.disciplineRepository.GetCurriculumPrivileges, auth.LevelPublic, auth.ActionRead, disciplineId)
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

	err = auth.Authorize(ctx, c.disciplineRepository.GetCurriculumPrivileges, auth.LevelPublic, auth.ActionRead, disciplineId)
	if err != nil {
		return nil, err
	}

	return c.userRepository.GetDisciplineProgress(userId, disciplineId)
}

func (c *DisciplineController) GetDisciplineStats(ctx *fiber.Ctx) (*user.GenericStats, error) {
	userId := ctx.Locals("userid").(int64)

	disciplineId, err := httputil.ExtractId(ctx, "discipline_id")
	if err != nil {
		return nil, err
	}

	err = auth.Authorize(ctx, c.disciplineRepository.GetCurriculumPrivileges, auth.LevelPublic, auth.ActionRead, disciplineId)
	if err != nil {
		return nil, err
	}

	return c.userRepository.GetDisciplineStats(userId, disciplineId)
}
