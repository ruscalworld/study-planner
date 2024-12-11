package discipline

import (
	"net/url"

	"github.com/ruscalworld/study-planner/pkg/stderrors"
)

var (
	ErrUnknownDiscipline     = stderrors.NotFound("unknown discipline")
	ErrUnknownDisciplineLink = stderrors.NotFound("unknown discipline link")
)

type Discipline struct {
	ID   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type CreateDisciplineParams struct {
	Name string `json:"name"`
}

func (p *CreateDisciplineParams) Validate() error {
	if len([]rune(p.Name)) == 0 {
		return stderrors.UnprocessableEntity("name must not be empty")
	}

	return nil
}

type UpdateDisciplineParams struct {
	Name string `json:"name"`
}

func (p *UpdateDisciplineParams) Validate() error {
	if len([]rune(p.Name)) == 0 {
		return stderrors.UnprocessableEntity("name must not be empty")
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
