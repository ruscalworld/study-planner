package task

import (
	"net/url"
	"time"

	"github.com/ruscalworld/study-planner/internal/user"

	"github.com/ruscalworld/study-planner/pkg/stderrors"
)

var (
	ErrUnknownGroup    = stderrors.NotFound("unknown task group")
	ErrUnknownTask     = stderrors.NotFound("unknown task")
	ErrUnknownTaskLink = stderrors.NotFound("unknown task link")
)

type Group struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type GroupParams struct {
	Name string `json:"name"`
}

func (p *GroupParams) Validate() error {
	if p.Name == "" {
		return stderrors.UnprocessableEntity("name must not be empty")
	}

	return nil
}

type Status string

const (
	StatusNotPublished Status = "NotPublished"
	StatusAvailable    Status = "Available"
)

func (s Status) IsValid() bool {
	return s == StatusNotPublished || s == StatusAvailable
}

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

type CreateParams struct {
	Name         string `json:"name"`
	ExternalName string `json:"externalName"`
	Description  string `json:"description"`
}

func (p *CreateParams) Validate() error {
	if p.Name == "" {
		return stderrors.UnprocessableEntity("name must not be empty")
	}

	return nil
}

type Params struct {
	CreateParams
	Status     Status     `json:"status,omitempty"`
	Deadline   *time.Time `json:"deadline,omitempty"`
	Difficulty int        `json:"difficulty,omitempty"`
}

func (p *Params) Validate() error {
	if err := p.CreateParams.Validate(); err != nil {
		return err
	}

	if p.Difficulty < 0 {
		return stderrors.UnprocessableEntity("difficulty must not be negative")
	}

	if !p.Status.IsValid() {
		return stderrors.UnprocessableEntity("invalid task status")
	}

	return nil
}

type Link struct {
	ID   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	URL  string `json:"url" db:"url"`
}

type LinkParams struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func (p *LinkParams) Validate() error {
	if p.Name == "" {
		return stderrors.UnprocessableEntity("name must not be empty")
	}

	if p.URL == "" {
		return stderrors.UnprocessableEntity("url must not be empty")
	}

	_, err := url.Parse(p.URL)
	if err != nil {
		return stderrors.UnprocessableEntity("url must be a valid URL")
	}

	return nil
}

type UpdateGoalParams struct {
	MinCompleted int `json:"minCompleted"`
}

type UpdateProgressParams struct {
	Status user.TaskStatus `json:"status"`
	Grade  *user.Grade     `json:"grade"`
}
