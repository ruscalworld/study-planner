package draft

import (
	"github.com/ruscalworld/study-planner/internal/task"

	"github.com/gofiber/fiber/v2"
)

type Controller interface {
	GetDrafts(ctx *fiber.Ctx) (*[]Draft, error)
	CreateDraft(ctx *fiber.Ctx, params *Params) (*Draft, error)
	UpdateDraft(ctx *fiber.Ctx, params *Params) (*Draft, error)
	DeleteDraft(ctx *fiber.Ctx) (*any, error)

	MoveDraft(ctx *fiber.Ctx, params *MoveParams) (*task.Task, error)
}
