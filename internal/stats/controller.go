package stats

import (
	"github.com/gofiber/fiber/v2"
)

type Controller interface {
	GetUndoneUpcomingTasks(ctx *fiber.Ctx) (*[]DisciplineTask, error)
}
