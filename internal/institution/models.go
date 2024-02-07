package institution

import (
	"study-planner/pkg/stderrors"
)

var (
	ErrUnknownInstitution = stderrors.NotFound("unknown institution")
)

type Institution struct {
	ID   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}
