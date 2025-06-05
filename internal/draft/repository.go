package draft

import (
	"github.com/ruscalworld/study-planner/internal/task"
)

type Repository interface {
	GetDraft(id int64, userId int64) (*Draft, error)
	GetUserDrafts(userId int64) (*[]Draft, error)
	CreateDraft(draft *Draft) error
	UpdateDraft(draft *Draft) error
	DeleteDraft(id int64) error
	MoveDraft(draft *Draft, taskGroup *task.Group, taskName string) (*task.Task, error)
}
