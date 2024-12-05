package delivery

import (
	"regexp"

	"github.com/ruscalworld/study-planner/internal/access"
	"github.com/ruscalworld/study-planner/internal/curriculum"
	"github.com/ruscalworld/study-planner/internal/institution"
	"github.com/ruscalworld/study-planner/internal/user"

	"github.com/ruscalworld/study-planner/pkg/httputil"
	"github.com/ruscalworld/study-planner/pkg/stderrors"

	"github.com/gofiber/fiber/v2"
)

type InstitutionController struct {
	institutionRepository institution.Repository
	curriculumRepository  curriculum.Repository
	userRepository        user.Repository
}

func NewInstitutionController(
	institutionRepository institution.Repository,
	curriculumRepository curriculum.Repository,
	userRepository user.Repository,
) *InstitutionController {
	return &InstitutionController{
		institutionRepository: institutionRepository,
		curriculumRepository:  curriculumRepository,
		userRepository:        userRepository,
	}
}

func (c *InstitutionController) GetInstitutions(_ *fiber.Ctx) (*[]institution.Institution, error) {
	return c.institutionRepository.GetInstitutions()
}

func (c *InstitutionController) GetInstitution(ctx *fiber.Ctx) (*institution.Institution, error) {
	id, err := httputil.ExtractId(ctx, "institution_id")
	if err != nil {
		return nil, err
	}

	return c.institutionRepository.GetInstitution(id)
}

func (c *InstitutionController) GetCurriculums(ctx *fiber.Ctx) (*[]curriculum.Curriculum, error) {
	id, err := httputil.ExtractId(ctx, "institution_id")
	if err != nil {
		return nil, err
	}

	return c.curriculumRepository.GetInstitutionCurriculums(id)
}

func (c *InstitutionController) CreateCurriculum(ctx *fiber.Ctx, request *institution.CreateCurriculumParams) (*curriculum.Curriculum, error) {
	institutionId, err := httputil.ExtractId(ctx, "institution_id")
	if err != nil {
		return nil, err
	}

	nameValid, err := regexp.MatchString("^[a-zA-Zа-яА-Яё0-9.-]{3,}$", request.Name)
	if err != nil {
		return nil, err
	}

	if !nameValid {
		return nil, stderrors.BadRequest("name is not valid")
	}

	cu := &curriculum.Curriculum{
		Name:     request.Name,
		Semester: request.Semester,
	}

	err = c.curriculumRepository.CreateCurriculum(institutionId, cu)
	if err != nil {
		return nil, err
	}

	userId := ctx.Locals("userid").(int64)
	err = c.userRepository.CreateUserCurriculum(userId, cu.ID, access.RoleOwner)
	if err != nil {
		return nil, err
	}

	return cu, nil
}
