package delivery

import (
	"time"

	"github.com/ruscalworld/study-planner/internal/auth"
	"github.com/ruscalworld/study-planner/internal/curriculum"

	"github.com/ruscalworld/study-planner/pkg/code"
	"github.com/ruscalworld/study-planner/pkg/httputil"
	"github.com/ruscalworld/study-planner/pkg/stderrors"

	"github.com/gofiber/fiber/v2"
)

type CurriculumController struct {
	curriculumRepository curriculum.Repository
}

func NewCurriculumController(curriculumRepository curriculum.Repository) *CurriculumController {
	return &CurriculumController{curriculumRepository: curriculumRepository}
}

func (c *CurriculumController) GetCurriculum(ctx *fiber.Ctx) (*curriculum.Curriculum, error) {
	id, err := httputil.ExtractId(ctx, "curriculum_id")
	if err != nil {
		return nil, err
	}

	err = auth.Authorize(ctx, c.curriculumRepository.GetCurriculumPrivileges, auth.LevelPublic, auth.ActionRead, id)
	if err != nil {
		return nil, err
	}

	return c.curriculumRepository.GetCurriculum(id)
}

func (c *CurriculumController) GetCurriculumCodes(ctx *fiber.Ctx) (*[]curriculum.Code, error) {
	id, err := httputil.ExtractId(ctx, "curriculum_id")
	if err != nil {
		return nil, err
	}

	err = auth.Authorize(ctx, c.curriculumRepository.GetCurriculumPrivileges, auth.LevelSecure, auth.ActionRead, id)
	if err != nil {
		return nil, err
	}

	return c.curriculumRepository.GetCurriculumCodes(id)
}

func (c *CurriculumController) CreateCurriculumCode(ctx *fiber.Ctx, request *curriculum.CreateCodeParams) (*curriculum.Code, error) {
	curriculumId, err := httputil.ExtractId(ctx, "curriculum_id")
	if err != nil {
		return nil, err
	}

	err = auth.Authorize(ctx, c.curriculumRepository.GetCurriculumPrivileges, auth.LevelSecure, auth.ActionCreate, curriculumId)
	if err != nil {
		return nil, err
	}

	if !request.Role.IsValid() {
		return nil, stderrors.UnprocessableEntity("invalid role")
	}

	if request.ExpiresAt.Before(time.Now()) {
		return nil, stderrors.UnprocessableEntity("expiry date must be in future")
	}

	cd := &curriculum.Code{
		Code:         code.MakeCode(7),
		Role:         request.Role,
		UserId:       ctx.Locals("userid").(int64),
		CurriculumId: curriculumId,
		ExpiresAt:    request.ExpiresAt,
		CreatedAt:    time.Now(),
	}

	err = c.curriculumRepository.CreateCurriculumCode(cd)
	if err != nil {
		return nil, err
	}

	return cd, nil
}

func (c *CurriculumController) DeleteCurriculumCode(ctx *fiber.Ctx) (*any, error) {
	curriculumId, err := httputil.ExtractId(ctx, "curriculum_id")
	if err != nil {
		return nil, err
	}

	codeId, err := httputil.ExtractId(ctx, "code_id")
	if err != nil {
		return nil, err
	}

	err = auth.Authorize(ctx, c.curriculumRepository.GetCurriculumPrivileges, auth.LevelSecure, auth.ActionCreate, curriculumId)
	if err != nil {
		return nil, err
	}

	err = c.curriculumRepository.DeleteCurriculumCode(codeId)
	if err != nil {
		return nil, err
	}

	ctx.Status(fiber.StatusNoContent)
	return nil, nil
}

func (c *CurriculumController) GetCurriculumUsers(ctx *fiber.Ctx) (*[]curriculum.User, error) {
	id, err := httputil.ExtractId(ctx, "curriculum_id")
	if err != nil {
		return nil, err
	}

	err = auth.Authorize(ctx, c.curriculumRepository.GetCurriculumPrivileges, auth.LevelSecure, auth.ActionRead, id)
	if err != nil {
		return nil, err
	}

	return c.curriculumRepository.GetCurriculumUsers(id)
}

func (c *CurriculumController) UpdateCurriculumUser(ctx *fiber.Ctx, params *curriculum.UpdateCurriculumUserParams) (*any, error) {
	curriculumId, err := httputil.ExtractId(ctx, "curriculum_id")
	if err != nil {
		return nil, err
	}

	userId, err := httputil.ExtractId(ctx, "user_id")
	if err != nil {
		return nil, err
	}

	currentUserId := ctx.Locals("userid").(int64)
	if currentUserId == userId {
		return nil, stderrors.Conflict("you cannot alter yourself")
	}

	if !params.Role.IsValid() {
		return nil, stderrors.UnprocessableEntity("invalid role")
	}

	err = auth.Authorize(ctx, c.curriculumRepository.GetCurriculumPrivileges, auth.LevelSecure, auth.ActionUpdate, curriculumId)
	if err != nil {
		return nil, err
	}

	err = c.curriculumRepository.UpdateCurriculumUser(curriculumId, userId, params.Role)
	if err != nil {
		return nil, err
	}

	ctx.Status(fiber.StatusNoContent)
	return nil, nil
}

func (c *CurriculumController) DeleteCurriculumUser(ctx *fiber.Ctx) (*any, error) {
	curriculumId, err := httputil.ExtractId(ctx, "curriculum_id")
	if err != nil {
		return nil, err
	}

	userId, err := httputil.ExtractId(ctx, "user_id")
	if err != nil {
		return nil, err
	}

	currentUserId := ctx.Locals("userid").(int64)
	if currentUserId == userId {
		return nil, stderrors.Conflict("you cannot delete yourself from curriculum user list")
	}

	err = auth.Authorize(ctx, c.curriculumRepository.GetCurriculumPrivileges, auth.LevelSecure, auth.ActionUpdate, curriculumId)
	if err != nil {
		return nil, err
	}

	err = c.curriculumRepository.DeleteCurriculumUser(curriculumId, userId)
	if err != nil {
		return nil, err
	}

	ctx.Status(fiber.StatusNoContent)
	return nil, nil
}
