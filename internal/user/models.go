package user

import (
	"time"

	"study-planner/pkg/stderrors"
)

var (
	ErrUnknownUser = stderrors.NotFound("unknown user")
)

type User struct {
	ID         int64     `json:"id" db:"id"`
	Name       string    `json:"name" db:"name"`
	AvatarURL  string    `json:"avatarUrl" db:"avatar_url"`
	Platform   string    `json:"platform" db:"platform"`
	ExternalID string    `json:"externalId" db:"external_id"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
}

type TaskStatus string

const (
	TaskStatusNotStarted      TaskStatus = "NotStarted"
	TaskStatusInProgress      TaskStatus = "InProgress"
	TaskStatusNeedsProtection TaskStatus = "NeedsProtection"
	TaskStatusCompleted       TaskStatus = "Completed"
)

type Grade string

const (
	GradeExcellent    Grade = "Excellent"
	GradeGood         Grade = "Good"
	GradeSatisfactory Grade = "Satisfactory"
	GradeCredited     Grade = "Credited"
)

type TaskProgress struct {
	ID int64 `json:"-" db:"id"`
	GenericTaskProgress
}

type GenericTaskProgress struct {
	Status      TaskStatus `json:"status" db:"status"`
	Grade       *Grade     `json:"grade" db:"grade"`
	StartedAt   *time.Time `json:"startedAt" db:"started_at"`
	CompletedAt *time.Time `json:"completedAt" db:"completed_at"`
}

type ScopedTaskProgress struct {
	ID          int64 `json:"-" db:"id"`
	TaskId      int64 `json:"taskId" db:"task_id"`
	TaskGroupId int64 `json:"taskGroupId" db:"task_group_id"`
	GenericTaskProgress
}

type Goal struct {
	ID           int64 `json:"-" db:"id"`
	MinCompleted int   `json:"minCompleted" db:"min_completed"`
}

type TaskGroupStats struct {
	TaskGroupId int64 `json:"taskGroupId" db:"task_group_id"`
	GenericStats
}

type GenericStats struct {
	CompletedTasks  int `json:"completedTasks" db:"completed_tasks"`
	InProgressTasks int `json:"inProgressTasks" db:"in_progress_tasks"`
	GoalTasks       int `json:"goalTasks" db:"goal_tasks"`
	AvailableTasks  int `json:"availableTasks" db:"available_tasks"`
	TotalTasks      int `json:"totalTasks" db:"total_tasks"`
}
