package auth

import (
	"github.com/ruscalworld/study-planner/internal/access"
	"github.com/ruscalworld/study-planner/pkg/stderrors"

	"github.com/gofiber/fiber/v2"
)

var (
	ErrNoPrivileges = stderrors.Forbidden("not enough privileges to perform the requested operation")
)

type Action string

const (
	ActionCreate Action = "Create"
	ActionRead   Action = "Read"
	ActionUpdate Action = "Update"
	ActionDelete Action = "Delete"
)

type ResourceLevel byte

const (
	LevelPublic ResourceLevel = iota
	LevelSecure ResourceLevel = iota
)

type PrivilegeFetcher[E any] func(entityId E, userId int64) (*access.CurriculumPrivileges, error)

func Authorize[E any](ctx *fiber.Ctx, fetcher PrivilegeFetcher[E], level ResourceLevel, action Action, entityId E) error {
	uc, err := fetcher(entityId, ctx.Locals("userid").(int64))
	if err != nil {
		return err
	}

	if !can(uc.Role, level, action) {
		return ErrNoPrivileges
	}

	return nil
}

func can(role access.Role, level ResourceLevel, action Action) bool {
	switch level {
	case LevelPublic:
		switch action {
		case ActionCreate, ActionUpdate, ActionDelete:
			return role == access.RoleOwner || role == access.RoleEditor
		case ActionRead:
			return role == access.RoleOwner || role == access.RoleEditor || role == access.RoleViewer
		}
	case LevelSecure:
		return role == access.RoleOwner
	}

	return false
}
