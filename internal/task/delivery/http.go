package delivery

import (
	"study-planner/internal/task"
	"study-planner/pkg/httputil"

	"github.com/gofiber/fiber/v2"
)

type TaskController struct {
	taskRepository task.Repository
}

func NewTaskController(taskRepository task.Repository) *TaskController {
	return &TaskController{taskRepository: taskRepository}
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
