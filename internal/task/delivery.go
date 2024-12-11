package task

import (
	"github.com/ruscalworld/study-planner/internal/user"

	"github.com/gofiber/fiber/v2"
)

type Controller interface {
	GetTaskGroups(ctx *fiber.Ctx) (*[]Group, error)
	GetTaskGroup(ctx *fiber.Ctx) (*Group, error)
	CreateTaskGroup(ctx *fiber.Ctx, params *GroupParams) (*Group, error)
	UpdateTaskGroup(ctx *fiber.Ctx, params *GroupParams) (*Group, error)
	DeleteTaskGroup(ctx *fiber.Ctx) (*any, error)

	GetTaskGroupGoal(ctx *fiber.Ctx) (*user.Goal, error)
	UpdateTaskGroupGoal(ctx *fiber.Ctx, params *UpdateGoalParams) (*user.Goal, error)

	GetTasks(ctx *fiber.Ctx) (*[]Task, error)
	GetTask(ctx *fiber.Ctx) (*Task, error)
	CreateTask(ctx *fiber.Ctx, params *CreateParams) (*Task, error)
	UpdateTask(ctx *fiber.Ctx, params *Params) (*Task, error)
	DeleteTask(ctx *fiber.Ctx) (*any, error)

	GetTaskLinks(ctx *fiber.Ctx) (*[]Link, error)
	CreateTaskLink(ctx *fiber.Ctx, params *LinkParams) (*Link, error)
	UpdateTaskLink(ctx *fiber.Ctx, params *LinkParams) (*Link, error)
	DeleteTaskLink(ctx *fiber.Ctx) (*any, error)

	GetTaskProgress(ctx *fiber.Ctx) (*user.TaskProgress, error)
	UpdateTaskProgress(ctx *fiber.Ctx, params *UpdateProgressParams) (*user.TaskProgress, error)
}
