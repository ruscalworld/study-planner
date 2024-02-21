package task

import (
	"study-planner/internal/user"

	"github.com/gofiber/fiber/v2"
)

type Controller interface {
	GetTaskGroups(ctx *fiber.Ctx) (*[]Group, error)
	GetTaskGroup(ctx *fiber.Ctx) (*Group, error)

	GetTaskGroupGoal(ctx *fiber.Ctx) (*user.Goal, error)
	UpdateTaskGroupGoal(ctx *fiber.Ctx, params *UpdateGoalParams) (*user.Goal, error)

	GetTasks(ctx *fiber.Ctx) (*[]Task, error)
	GetTask(ctx *fiber.Ctx) (*Task, error)
	GetTaskLinks(ctx *fiber.Ctx) (*[]Link, error)

	GetTaskProgress(ctx *fiber.Ctx) (*user.TaskProgress, error)
	UpdateTaskProgress(ctx *fiber.Ctx, params *UpdateProgressParams) (*user.TaskProgress, error)
}
