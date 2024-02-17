package task

import (
	"time"

	"study-planner/pkg/stderrors"
)

var (
	ErrUnknownGroup = stderrors.NotFound("unknown task group")
	ErrUnknownTask  = stderrors.NotFound("unknown task")
)

type Group struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Status string

const (
	StatusNotPublished Status = "NotPublished"
	StatusAvailable    Status = "Available"
)

type Task struct {
	ID           int64      `json:"id" db:"id"`
	Name         string     `json:"name" db:"name"`
	ExternalName *string    `json:"externalName" db:"external_name"`
	Description  *string    `json:"description" db:"description"`
	GroupID      int64      `json:"groupId" db:"task_group_id"`
	Status       Status     `json:"status" db:"status"`
	Difficulty   int        `json:"difficulty" db:"difficulty"`
	Deadline     *time.Time `json:"deadline" db:"deadline"`
}

type Link struct {
	ID   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	URL  string `json:"url" db:"url"`
}

type UpdateGoalParams struct {
	MinCompleted int `json:"minCompleted"`
}
