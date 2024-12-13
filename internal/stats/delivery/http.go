package delivery

import (
	"github.com/ruscalworld/study-planner/internal/auth"
	"github.com/ruscalworld/study-planner/internal/curriculum"
	"github.com/ruscalworld/study-planner/internal/stats"
	"github.com/ruscalworld/study-planner/pkg/httputil"

	"github.com/gofiber/fiber/v2"
)

type StatsController struct {
	curriculumRepository curriculum.Repository
	statsRepository      stats.Repository
}

func NewStatsController(
	curriculumRepository curriculum.Repository,
	statsRepository stats.Repository,
) *StatsController {
	return &StatsController{
		curriculumRepository: curriculumRepository,
		statsRepository:      statsRepository,
	}
}

func (c *StatsController) GetUndoneUpcomingTasks(ctx *fiber.Ctx) (*[]stats.DisciplineTask, error) {
	curriculumId, err := httputil.ExtractId(ctx, "curriculum_id")
	if err != nil {
		return nil, err
	}

	err = auth.Authorize(ctx, c.curriculumRepository.GetCurriculumPrivileges, auth.LevelPublic, auth.ActionRead, curriculumId)
	if err != nil {
		return nil, err
	}

	userId := ctx.Locals("userid").(int64)
	tasks, err := c.statsRepository.GetUndoneUpcomingTasks(curriculumId, userId, 10)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
