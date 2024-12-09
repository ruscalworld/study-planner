package institution

import (
	"github.com/ruscalworld/study-planner/pkg/stderrors"
)

var (
	ErrUnknownInstitution = stderrors.NotFound("unknown institution")
)

type Institution struct {
	ID   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type CreateCurriculumParams struct {
	Name     string `json:"name"`
	Semester int    `json:"semester"`
}

func (p *CreateCurriculumParams) Validate() error {
	if p.Name == "" {
		return stderrors.UnprocessableEntity("name must not be empty")
	}

	if p.Semester < 0 {
		return stderrors.UnprocessableEntity("semester must not be negative")
	}

	if p.Semester > 32 {
		return stderrors.UnprocessableEntity("semester must not be greater than 32")
	}

	return nil
}
