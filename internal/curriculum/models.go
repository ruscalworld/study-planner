package curriculum

import (
	"time"

	"github.com/ruscalworld/study-planner/internal/access"
	"github.com/ruscalworld/study-planner/pkg/stderrors"
)

var (
	ErrUnknownCurriculum = stderrors.NotFound("unknown curriculum")
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
