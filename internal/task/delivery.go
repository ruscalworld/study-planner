package task

import "github.com/gofiber/fiber/v2"

type Controller interface {
	GetTaskGroups(ctx *fiber.Ctx) (*[]Group, error)
	GetTaskGroup(ctx *fiber.Ctx) (*Group, error)
	GetTasks(ctx *fiber.Ctx) (*[]Task, error)
	GetTask(ctx *fiber.Ctx) (*Task, error)
	GetTaskLinks(ctx *fiber.Ctx) (*[]Link, error)
}
