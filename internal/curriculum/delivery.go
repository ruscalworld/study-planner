package curriculum

import "github.com/gofiber/fiber/v2"

type Controller interface {
	GetCurriculums(ctx *fiber.Ctx) (*[]Curriculum, error)
	GetCurriculum(ctx *fiber.Ctx) (*Curriculum, error)
}
