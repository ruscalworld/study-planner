package discipline

import (
	"github.com/ruscalworld/study-planner/internal/user"

	"github.com/gofiber/fiber/v2"
)

type Controller interface {
	GetDisciplines(ctx *fiber.Ctx) (*[]Discipline, error)
	GetCurriculumDiscipline(ctx *fiber.Ctx) (*Discipline, error)
	CreateDiscipline(ctx *fiber.Ctx, request *CreateDisciplineParams) (*Discipline, error)
	UpdateDiscipline(ctx *fiber.Ctx, request *UpdateDisciplineParams) (*Discipline, error)
	DeleteDiscipline(ctx *fiber.Ctx) (*any, error)
	GetDisciplineLinks(ctx *fiber.Ctx) (*[]Link, error)
	CreateDisciplineLink(ctx *fiber.Ctx, params *LinkParams) (*Link, error)
	UpdateDisciplineLink(ctx *fiber.Ctx, params *LinkParams) (*Link, error)
	DeleteDisciplineLink(ctx *fiber.Ctx) (*any, error)
	GetDisciplineProgress(ctx *fiber.Ctx) (*[]user.ScopedTaskProgress, error)
	GetDisciplineStats(ctx *fiber.Ctx) (*user.GenericStats, error)
}
