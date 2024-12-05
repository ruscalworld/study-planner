package delivery

import (
	"errors"

	"github.com/ruscalworld/study-planner/internal/auth"
	"github.com/ruscalworld/study-planner/internal/curriculum"
	"github.com/ruscalworld/study-planner/internal/user"
	"github.com/ruscalworld/study-planner/pkg/httputil"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	userRepository       user.Repository
	curriculumRepository curriculum.Repository
}

func NewUserController(userRepository user.Repository, curriculumRepository curriculum.Repository) *UserController {
	return &UserController{userRepository: userRepository, curriculumRepository: curriculumRepository}
}

func (c *UserController) GetUserCurriculums(ctx *fiber.Ctx) (*[]user.Curriculum, error) {
	userId := ctx.Locals("userid").(int64)
	curriculums, err := c.userRepository.GetUserCurriculums(userId)
	if err != nil {
		return nil, err
	}

	return curriculums, nil
}

func (c *UserController) CreateUserCurriculum(ctx *fiber.Ctx, request *user.CreateUserCurriculumParams) (*any, error) {
	cp, err := c.curriculumRepository.GetCurriculumByCode(request.Code)
	if err != nil {
		return nil, err
	}

	userId := ctx.Locals("userid").(int64)
	privileges, err := c.curriculumRepository.GetCurriculumPrivileges(cp.Curriculum.ID, userId)
	if err != nil && !errors.Is(err, auth.ErrNoPrivileges) {
		return nil, err
	}

	if privileges == nil || cp.Role.GreaterThan(privileges.Role) {
		err = c.userRepository.CreateUserCurriculum(userId, cp.ID, cp.Role)
		if err != nil {
			return nil, err
		}

		ctx.Status(fiber.StatusCreated)
	} else {
		ctx.Status(fiber.StatusNoContent)
	}

	return nil, nil
}

func (c *UserController) DeleteUserCurriculum(ctx *fiber.Ctx) (*any, error) {
	userId := ctx.Locals("userid").(int64)
	curriculumIdd, err := httputil.ExtractId(ctx, "curriculum_id")
	if err != nil {
		return nil, err
	}

	err = c.userRepository.DeleteUserCurriculum(userId, curriculumIdd)
	if err != nil {
		return nil, err
	}

	ctx.Status(fiber.StatusNoContent)
	return nil, nil
}
