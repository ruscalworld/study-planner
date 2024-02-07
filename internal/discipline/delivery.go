package discipline

import "github.com/gofiber/fiber/v2"

type Controller interface {
	GetDisciplines(ctx *fiber.Ctx) (*[]Discipline, error)
	GetDiscipline(ctx *fiber.Ctx) (*Discipline, error)
	GetDisciplineLinks(ctx *fiber.Ctx) (*[]Link, error)
}
