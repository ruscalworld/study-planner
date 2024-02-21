package discipline

import (
	"study-planner/internal/user"

	"github.com/gofiber/fiber/v2"
)

type Controller interface {
	GetDisciplines(ctx *fiber.Ctx) (*[]Discipline, error)
	GetDiscipline(ctx *fiber.Ctx) (*Discipline, error)
	GetDisciplineLinks(ctx *fiber.Ctx) (*[]Link, error)
	GetDisciplineProgress(ctx *fiber.Ctx) (*[]user.ScopedTaskProgress, error)
}
