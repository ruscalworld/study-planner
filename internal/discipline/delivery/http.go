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

func (c *DisciplineController) GetCurriculumDiscipline(ctx *fiber.Ctx) (*discipline.Discipline, error) {
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

func (c *DisciplineController) CreateDiscipline(ctx *fiber.Ctx, request *discipline.CreateDisciplineParams) (*discipline.Discipline, error) {
	curriculumId, err := httputil.ExtractId(ctx, "curriculum_id")
	if err != nil {
		return nil, err
	}

	err = auth.Authorize(ctx, c.curriculumRepository.GetCurriculumPrivileges, auth.LevelPublic, auth.ActionUpdate, curriculumId)
	if err != nil {
		return nil, err
	}

	err = request.Validate()
	if err != nil {
		return nil, err
	}

	d := &discipline.Discipline{
		Name: request.Name,
	}

	err = c.disciplineRepository.CreateDiscipline(curriculumId, d)
	if err != nil {
		return nil, err
	}

	ctx.Status(fiber.StatusCreated)
	return d, nil
}

func (c *DisciplineController) UpdateDiscipline(ctx *fiber.Ctx, request *discipline.UpdateDisciplineParams) (*discipline.Discipline, error) {
	curriculumId, err := httputil.ExtractId(ctx, "curriculum_id")
	if err != nil {
		return nil, err
	}

	err = auth.Authorize(ctx, c.curriculumRepository.GetCurriculumPrivileges, auth.LevelPublic, auth.ActionUpdate, curriculumId)
	if err != nil {
		return nil, err
	}

	err = request.Validate()
	if err != nil {
		return nil, err
	}

	disciplineId, err := httputil.ExtractId(ctx, "discipline_id")
	if err != nil {
		return nil, err
	}

	d, err := c.disciplineRepository.GetDiscipline(curriculumId, disciplineId)
	if err != nil {
		return nil, err
	}

	d.Name = request.Name
	err = c.disciplineRepository.UpdateDiscipline(d)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (c *DisciplineController) DeleteDiscipline(ctx *fiber.Ctx) (*any, error) {
	curriculumId, err := httputil.ExtractId(ctx, "curriculum_id")
	if err != nil {
		return nil, err
	}

	disciplineId, err := httputil.ExtractId(ctx, "discipline_id")
	if err != nil {
		return nil, err
	}

	err = auth.Authorize(ctx, c.curriculumRepository.GetCurriculumPrivileges, auth.LevelPublic, auth.ActionUpdate, curriculumId)
	if err != nil {
		return nil, err
	}

	err = c.disciplineRepository.DeleteDiscipline(curriculumId, disciplineId)
	if err != nil {
		return nil, err
	}

	ctx.Status(fiber.StatusNoContent)
	return nil, nil
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

func (c *DisciplineController) CreateDisciplineLink(ctx *fiber.Ctx, params *discipline.LinkParams) (*discipline.Link, error) {
	disciplineId, err := httputil.ExtractId(ctx, "discipline_id")
	if err != nil {
		return nil, err
	}

	err = auth.Authorize(ctx, c.disciplineRepository.GetCurriculumPrivileges, auth.LevelPublic, auth.ActionUpdate, disciplineId)
	if err != nil {
		return nil, err
	}

	err = params.Validate()
	if err != nil {
		return nil, err
	}

	l := &discipline.Link{
		Name: params.Name,
		URL:  params.URL,
	}

	err = c.disciplineRepository.CreateDisciplineLink(disciplineId, l)
	if err != nil {
		return nil, err
	}

	ctx.Status(fiber.StatusCreated)
	return l, nil
}

func (c *DisciplineController) UpdateDisciplineLink(ctx *fiber.Ctx, params *discipline.LinkParams) (*discipline.Link, error) {
	disciplineId, err := httputil.ExtractId(ctx, "discipline_id")
	if err != nil {
		return nil, err
	}

	err = auth.Authorize(ctx, c.disciplineRepository.GetCurriculumPrivileges, auth.LevelPublic, auth.ActionUpdate, disciplineId)
	if err != nil {
		return nil, err
	}

	disciplineLinkId, err := httputil.ExtractId(ctx, "discipline_link_id")
	if err != nil {
		return nil, err
	}

	err = params.Validate()
	if err != nil {
		return nil, err
	}

	l := &discipline.Link{
		ID:   disciplineLinkId,
		Name: params.Name,
		URL:  params.URL,
	}

	err = c.disciplineRepository.UpdateDisciplineLink(l)
	if err != nil {
		return nil, err
	}

	return l, nil
}

func (c *DisciplineController) DeleteDisciplineLink(ctx *fiber.Ctx) (*any, error) {
	disciplineId, err := httputil.ExtractId(ctx, "discipline_id")
	if err != nil {
		return nil, err
	}

	err = auth.Authorize(ctx, c.disciplineRepository.GetCurriculumPrivileges, auth.LevelPublic, auth.ActionUpdate, disciplineId)
	if err != nil {
		return nil, err
	}

	linkId, err := httputil.ExtractId(ctx, "discipline_link_id")
	if err != nil {
		return nil, err
	}

	err = c.disciplineRepository.DeleteDisciplineLink(linkId)
	if err != nil {
		return nil, err
	}

	ctx.Status(fiber.StatusNoContent)
	return nil, nil
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
