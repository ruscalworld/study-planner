package delivery

import (
	"study-planner/internal/task"
	"study-planner/internal/user"
	"study-planner/pkg/httputil"
	"time"

	"github.com/gofiber/fiber/v2"
)

type TaskController struct {
	taskRepository task.Repository
	userRepository user.Repository
}

func NewTaskController(taskRepository task.Repository, userRepository user.Repository) *TaskController {
	return &TaskController{taskRepository: taskRepository, userRepository: userRepository}
}

func (c *TaskController) GetTaskGroups(ctx *fiber.Ctx) (*[]task.Group, error) {
	disciplineId, err := httputil.ExtractId(ctx, "discipline_id")
	if err != nil {
		return nil, err
	}

	return c.taskRepository.GetGroups(disciplineId)
}

func (c *TaskController) GetTaskGroup(ctx *fiber.Ctx) (*task.Group, error) {
	disciplineId, err := httputil.ExtractId(ctx, "discipline_id")
	if err != nil {
		return nil, err
	}

	groupId, err := httputil.ExtractId(ctx, "group_id")
	if err != nil {
		return nil, err
	}

	return c.taskRepository.GetGroup(disciplineId, groupId)
}

func (c *TaskController) GetTaskGroupGoal(ctx *fiber.Ctx) (*user.Goal, error) {
	userId := ctx.Locals("userid").(int64)
	groupId, err := httputil.ExtractId(ctx, "group_id")
	if err != nil {
		return nil, err
	}

	return c.userRepository.GetGoal(userId, groupId)
}

func (c *TaskController) UpdateTaskGroupGoal(ctx *fiber.Ctx, params *task.UpdateGoalParams) (*user.Goal, error) {
	userId := ctx.Locals("userid").(int64)
	groupId, err := httputil.ExtractId(ctx, "group_id")
	if err != nil {
		return nil, err
	}

	g := &user.Goal{
		MinCompleted: params.MinCompleted,
	}

	err = c.userRepository.StoreGoal(userId, groupId, g)
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (c *TaskController) GetTasks(ctx *fiber.Ctx) (*[]task.Task, error) {
	disciplineId, err := httputil.ExtractId(ctx, "discipline_id")
	if err != nil {
		return nil, err
	}

	return c.taskRepository.GetTasks(disciplineId)
}

func (c *TaskController) GetTask(ctx *fiber.Ctx) (*task.Task, error) {
	disciplineId, err := httputil.ExtractId(ctx, "discipline_id")
	if err != nil {
		return nil, err
	}

	taskId, err := httputil.ExtractId(ctx, "task_id")
	if err != nil {
		return nil, err
	}

	return c.taskRepository.GetTask(disciplineId, taskId)
}

func (c *TaskController) GetTaskLinks(ctx *fiber.Ctx) (*[]task.Link, error) {
	disciplineId, err := httputil.ExtractId(ctx, "discipline_id")
	if err != nil {
		return nil, err
	}

	taskId, err := httputil.ExtractId(ctx, "task_id")
	if err != nil {
		return nil, err
	}

	return c.taskRepository.GetTaskLinks(disciplineId, taskId)
}

func (c *TaskController) GetTaskProgress(ctx *fiber.Ctx) (*user.TaskProgress, error) {
	userId := ctx.Locals("userid").(int64)
	taskId, err := httputil.ExtractId(ctx, "task_id")
	if err != nil {
		return nil, err
	}

	return c.userRepository.GetProgress(userId, taskId)
}

func (c *TaskController) UpdateTaskProgress(ctx *fiber.Ctx, params *task.UpdateProgressParams) (*user.TaskProgress, error) {
	userId := ctx.Locals("userid").(int64)
	taskId, err := httputil.ExtractId(ctx, "task_id")
	if err != nil {
		return nil, err
	}

	oldProgress, err := c.userRepository.GetProgress(userId, taskId)
	if err != nil {
		return nil, err
	}

	p := &user.TaskProgress{
		Status:      params.Status,
		Grade:       params.Grade,
		StartedAt:   oldProgress.StartedAt,
		CompletedAt: oldProgress.CompletedAt,
	}

	now := time.Now()

	// Update StartedAt if task was not started before
	if oldProgress.Status == user.TaskStatusNotStarted && params.Status != user.TaskStatusNotStarted {
		p.StartedAt = &now
	}

	// If task is being marked as NotStarted, then clear StartedAt
	if params.Status == user.TaskStatusNotStarted {
		p.StartedAt = nil
	}

	// If task is being marked not as Completed, then clear CompletedAt and Grade
	if params.Status != user.TaskStatusCompleted {
		p.Grade = nil
		p.CompletedAt = nil
	}

	// If task is being marked as Completed, then update CompletedAt
	if params.Status == user.TaskStatusCompleted {
		p.CompletedAt = &now
	}

	err = c.userRepository.StoreProgress(userId, taskId, p)
	if err != nil {
		return nil, err
	}

	return p, nil
}
