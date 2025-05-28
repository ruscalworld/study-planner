package delivery

import (
	"time"

	"github.com/ruscalworld/study-planner/internal/draft"
	"github.com/ruscalworld/study-planner/internal/task"
	"github.com/ruscalworld/study-planner/pkg/httputil"

	"github.com/gofiber/fiber/v2"
)

type DraftController struct {
	draftRepository draft.Repository
	taskRepository  task.Repository
}

func NewDraftController(draftRepository draft.Repository, taskRepository task.Repository) *DraftController {
	return &DraftController{
		draftRepository: draftRepository,
		taskRepository:  taskRepository,
	}
}

func (c *DraftController) GetDrafts(ctx *fiber.Ctx) (*[]draft.Draft, error) {
	userId := ctx.Locals("userid").(int64)
	return c.draftRepository.GetUserDrafts(userId)
}

func (c *DraftController) CreateDraft(ctx *fiber.Ctx, params *draft.Params) (*draft.Draft, error) {
	userId := ctx.Locals("userid").(int64)

	err := params.Validate()
	if err != nil {
		return nil, err
	}

	d := &draft.Draft{
		UserID:    userId,
		Text:      params.Text,
		CreatedAt: time.Now(),
	}

	err = c.draftRepository.CreateDraft(d)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (c *DraftController) UpdateDraft(ctx *fiber.Ctx, params *draft.Params) (*draft.Draft, error) {
	userId := ctx.Locals("userid").(int64)
	draftId, err := httputil.ExtractId(ctx, "draft_id")
	if err != nil {
		return nil, err
	}

	err = params.Validate()
	if err != nil {
		return nil, err
	}

	d, err := c.draftRepository.GetDraft(draftId, userId)
	if err != nil {
		return nil, err
	}

	d.Text = params.Text

	err = c.draftRepository.UpdateDraft(d)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (c *DraftController) DeleteDraft(ctx *fiber.Ctx) (*any, error) {
	userId := ctx.Locals("userid").(int64)
	draftId, err := httputil.ExtractId(ctx, "draft_id")
	if err != nil {
		return nil, err
	}

	d, err := c.draftRepository.GetDraft(draftId, userId)
	if err != nil {
		return nil, err
	}

	err = c.draftRepository.DeleteDraft(d.ID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (c *DraftController) MoveDraft(ctx *fiber.Ctx, params *draft.MoveParams) (*task.Task, error) {
	userId := ctx.Locals("userid").(int64)
	draftId, err := httputil.ExtractId(ctx, "draft_id")
	if err != nil {
		return nil, err
	}

	d, err := c.draftRepository.GetDraft(draftId, userId)
	if err != nil {
		return nil, err
	}

	taskGroup, err := c.taskRepository.GetGroup(params.DisciplineID, params.TaskGroupID)
	if err != nil {
		return nil, err
	}

	return c.draftRepository.MoveDraft(d, taskGroup, params.TaskName)
}
