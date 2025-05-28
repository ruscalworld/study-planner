package draft

import (
	"time"

	"github.com/ruscalworld/study-planner/pkg/stderrors"
)

var (
	ErrUnknownDraft = stderrors.NotFound("unknown draft")
)

type Draft struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"-" db:"user_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

type Params struct {
	Text string `json:"text"`
}

func (p *Params) Validate() error {
	if p.Text == "" {
		return stderrors.UnprocessableEntity("text must not be empty")
	}

	return nil
}

type MoveParams struct {
	TaskName     string `json:"taskName"`
	DisciplineID int64  `json:"disciplineId"`
	TaskGroupID  int64  `json:"taskGroupId"`
}

func (p *MoveParams) Validate() error {
	if p.TaskName == "" {
		return stderrors.UnprocessableEntity("task name must not be empty")
	}

	if p.DisciplineID == 0 {
		return stderrors.UnprocessableEntity("discipline id must be provided")
	}

	if p.TaskGroupID == 0 {
		return stderrors.UnprocessableEntity("task group id must not be provided")
	}

	return nil
}
