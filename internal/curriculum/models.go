package curriculum

import (
	"time"

	"github.com/ruscalworld/study-planner/internal/access"
	"github.com/ruscalworld/study-planner/pkg/stderrors"
)

var (
	ErrUnknownCurriculum     = stderrors.NotFound("unknown curriculum")
	ErrUnknownCurriculumUser = stderrors.NotFound("unknown curriculum user")
)

type Curriculum struct {
	ID       int64  `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Semester int    `json:"semester" db:"semester"`
}

type Code struct {
	ID           int64       `json:"id" db:"id"`
	Code         string      `json:"code" db:"code"`
	Role         access.Role `json:"role" db:"role"`
	UserId       int64       `json:"user_id" db:"user_id"`
	CurriculumId int64       `json:"-" db:"curriculum_id"`
	ExpiresAt    time.Time   `json:"expiresAt" db:"expires_at"`
	CreatedAt    time.Time   `json:"createdAt" db:"created_at"`
}

type CreateCodeParams struct {
	Role      access.Role `json:"role"`
	ExpiresAt time.Time   `json:"expiresAt"`
}

type Privileges struct {
	Curriculum
	access.CurriculumPrivileges
}

type User struct {
	UserId    int64       `json:"userId" db:"user_id"`
	Name      string      `json:"name" db:"name"`
	AvatarUrl string      `json:"avatarUrl" db:"avatar_url"`
	Role      access.Role `json:"role" db:"role"`
}

type Params struct {
	Name     string `json:"name"`
	Semester int    `json:"semester"`
}

func (p *Params) Validate() error {
	if p.Name == "" {
		return stderrors.UnprocessableEntity("name must not be empty")
	}

	if p.Semester < 0 {
		return stderrors.UnprocessableEntity("semester must not be negative")
	}

	if p.Semester > 32 {
		return stderrors.UnprocessableEntity("semester must not be greater than 32")
	}

	return nil
}

type UpdateCurriculumUserParams struct {
	Role access.Role `json:"role"`
}
