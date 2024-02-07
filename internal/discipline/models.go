package discipline

import (
	"study-planner/pkg/stderrors"
)

var (
	ErrUnknownDiscipline = stderrors.NotFound("unknown discipline")
)

type Discipline struct {
	ID   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type Link struct {
	ID   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	URL  string `json:"url" db:"url"`
}
