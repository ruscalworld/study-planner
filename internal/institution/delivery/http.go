package delivery

import (
	"study-planner/internal/curriculum"
	"study-planner/internal/institution"
	"study-planner/pkg/httputil"

	"github.com/gofiber/fiber/v2"
)

type InstitutionController struct {
	institutionRepository institution.Repository
	curriculumRepository  curriculum.Repository
}

func NewInstitutionController(institutionRepository institution.Repository, curriculumRepository curriculum.Repository) *InstitutionController {
	return &InstitutionController{institutionRepository: institutionRepository, curriculumRepository: curriculumRepository}
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
