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
