package curriculum

import (
	"study-planner/pkg/stderrors"
)

var (
	ErrUnknownCurriculum = stderrors.NotFound("unknown curriculum")
)

type Curriculum struct {
	ID       int64  `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Semester int    `json:"semester" db:"semester"`
}
