package delivery

import (
	"time"

	"github.com/ruscalworld/study-planner/internal/auth"
	"github.com/ruscalworld/study-planner/internal/discipline"
	"github.com/ruscalworld/study-planner/internal/task"
	"github.com/ruscalworld/study-planner/internal/user"

	"github.com/ruscalworld/study-planner/pkg/httputil"

	"github.com/gofiber/fiber/v2"
)

type TaskController struct {
	disciplineRepository discipline.Repository
	taskRepository       task.Repository
	userRepository       user.Repository
}

func NewTaskController(
	disciplineRepository discipline.Repository,
	taskRepository task.Repository,
	userRepository user.Repository,
) *TaskController {
	return &TaskController{
		disciplineRepository: disciplineRepository,
		taskRepository:       taskRepository,
		userRepository:       userRepository,
	}
}

func (c *TaskController) GetTaskGroups(ctx *fiber.Ctx) (*[]task.Group, error) {
	disciplineId, err := httputil.ExtractId(ctx, "discipline_id")
	if err != nil {
		return nil, err
	}

	err = auth.Authorize(ctx, c.disciplineRepository.GetCurriculumPrivileges, auth.LevelPublic, auth.ActionRead, disciplineId)
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

	err = auth.Authorize(ctx, c.disciplineRepository.GetCurriculumPrivileges, auth.LevelPublic, auth.ActionRead, disciplineId)
	if err != nil {
		return nil, err
	}

	groupId, err := httputil.ExtractId(ctx, "group_id")
	if err != nil {
		return nil, err
	}

	return c.taskRepository.GetGroup(disciplineId, groupId)
}

func (c *TaskController) CreateTaskGroup(ctx *fiber.Ctx, params *task.GroupParams) (*task.Group, error) {
	disciplineId, err := httputil.ExtractId(ctx, "discipline_id")
	if err != nil {
		return nil, err
	}

	err = auth.Authorize(ctx, c.disciplineRepository.GetCurriculumPrivileges, auth.LevelPublic, auth.ActionUpdate, disciplineId)
	if err != nil {
		return nil, err
	}

	err = params.Validate()
	if err != nil {
		return nil, err
	}

	g := &task.Group{
		Name: params.Name,
	}

	err = c.taskRepository.CreateGroup(disciplineId, g)
	if err != nil {
		return nil, err
	}

	ctx.Status(fiber.StatusCreated)
	return g, nil
}

func (c *TaskController) UpdateTaskGroup(ctx *fiber.Ctx, params *task.GroupParams) (*task.Group, error) {
	disciplineId, err := httputil.ExtractId(ctx, "discipline_id")
	if err != nil {
		return nil, err
	}

	taskGroupId, err := httputil.ExtractId(ctx, "group_id")
	if err != nil {
		return nil, err
	}

	err = auth.Authorize(ctx, c.disciplineRepository.GetCurriculumPrivileges, auth.LevelPublic, auth.ActionUpdate, disciplineId)
	if err != nil {
		return nil, err
	}

	err = params.Validate()
	if err != nil {
		return nil, err
	}

	g := &task.Group{
		ID:   taskGroupId,
		Name: params.Name,
	}

	err = c.taskRepository.UpdateGroup(g)
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (c *TaskController) DeleteTaskGroup(ctx *fiber.Ctx) (*any, error) {
	disciplineId, err := httputil.ExtractId(ctx, "discipline_id")
	if err != nil {
		return nil, err
	}

	taskGroupId, err := httputil.ExtractId(ctx, "group_id")
	if err != nil {
		return nil, err
	}

	err = auth.Authorize(ctx, c.disciplineRepository.GetCurriculumPrivileges, auth.LevelPublic, auth.ActionUpdate, disciplineId)
	if err != nil {
		return nil, err
	}

	err = c.taskRepository.DeleteGroup(taskGroupId)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (c *TaskController) GetTaskGroupGoal(ctx *fiber.Ctx) (*user.Goal, error) {
	disciplineId, err := httputil.ExtractId(ctx, "discipline_id")
	if err != nil {
		return nil, err
	}

	err = auth.Authorize(ctx, c.disciplineRepository.GetCurriculumPrivileges, auth.LevelPublic, auth.ActionRead, disciplineId)
	if err != nil {
		return nil, err
	}

	userId := ctx.Locals("userid").(int64)
	groupId, err := httputil.ExtractId(ctx, "group_id")
	if err != nil {
		return nil, err
	}

	return c.userRepository.GetGoal(userId, groupId)
}

func (c *TaskController) UpdateTaskGroupGoal(ctx *fiber.Ctx, params *task.UpdateGoalParams) (*user.Goal, error) {
	disciplineId, err := httputil.ExtractId(ctx, "discipline_id")
	if err != nil {
		return nil, err
	}

	err = auth.Authorize(ctx, c.disciplineRepository.GetCurriculumPrivileges, auth.LevelPublic, auth.ActionRead, disciplineId)
	if err != nil {
		return nil, err
	}

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

	err = auth.Authorize(ctx, c.disciplineRepository.GetCurriculumPrivileges, auth.LevelPublic, auth.ActionRead, disciplineId)
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

	t, err := c.taskRepository.GetTask(disciplineId, taskId)
	if err != nil {
		return nil, err
	}

	err = auth.Authorize(ctx, c.taskRepository.GetCurriculumPrivileges, auth.LevelPublic, auth.ActionRead, t.GroupID)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (c *TaskController) CreateTask(ctx *fiber.Ctx, params *task.CreateParams) (*task.Task, error) {
	disciplineId, err := httputil.ExtractId(ctx, "discipline_id")
	if err != nil {
		return nil, err
	}

	taskGroupId, err := httputil.ExtractId(ctx, "group_id")
	if err != nil {
		return nil, err
	}

	err = auth.Authorize(ctx, c.disciplineRepository.GetCurriculumPrivileges, auth.LevelPublic, auth.ActionUpdate, disciplineId)
	if err != nil {
		return nil, err
	}

	err = params.Validate()
	if err != nil {
		return nil, err
	}

	var (
		externalName *string = nil
		description  *string = nil
	)

	if params.ExternalName != "" {
		externalName = &params.ExternalName
	}

	if params.Description != "" {
		description = &params.Description
	}

	t := &task.Task{
		Name:         params.Name,
		ExternalName: externalName,
		Description:  description,
		GroupID:      taskGroupId,
		Status:       task.StatusNotPublished,
		Difficulty:   1,
	}

	err = c.taskRepository.CreateTask(taskGroupId, t)
	if err != nil {
		return nil, err
	}

	ctx.Status(fiber.StatusCreated)
	return t, nil
}

func (c *TaskController) UpdateTask(ctx *fiber.Ctx, params *task.Params) (*task.Task, error) {
	disciplineId, err := httputil.ExtractId(ctx, "discipline_id")
	if err != nil {
		return nil, err
	}

	taskId, err := httputil.ExtractId(ctx, "task_id")
	if err != nil {
		return nil, err
	}

	err = auth.Authorize(ctx, c.disciplineRepository.GetCurriculumPrivileges, auth.LevelPublic, auth.ActionUpdate, disciplineId)
	if err != nil {
		return nil, err
	}

	err = params.Validate()
	if err != nil {
		return nil, err
	}

	t, err := c.taskRepository.GetTask(disciplineId, taskId)
	if err != nil {
		return nil, err
	}

	t.Name = params.Name
	t.Difficulty = params.Difficulty
	t.Deadline = params.Deadline
	t.Status = params.Status

	if params.ExternalName == "" {
		t.ExternalName = nil
	} else {
		t.ExternalName = &params.ExternalName
	}

	if params.Description == "" {
		t.Description = nil
	} else {
		t.Description = &params.Description
	}

	err = c.taskRepository.UpdateTask(t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (c *TaskController) DeleteTask(ctx *fiber.Ctx) (*any, error) {
	taskId, err := httputil.ExtractId(ctx, "task_id")
	if err != nil {
		return nil, err
	}

	err = auth.Authorize(ctx, c.taskRepository.GetCurriculumPrivileges, auth.LevelPublic, auth.ActionUpdate, taskId)
	if err != nil {
		return nil, err
	}

	err = c.taskRepository.DeleteTask(taskId)
	if err != nil {
		return nil, err
	}

	return nil, nil
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

	t, err := c.taskRepository.GetTask(disciplineId, taskId)
	if err != nil {
		return nil, err
	}

	err = auth.Authorize(ctx, c.taskRepository.GetCurriculumPrivileges, auth.LevelPublic, auth.ActionRead, t.GroupID)
	if err != nil {
		return nil, err
	}

	return c.taskRepository.GetTaskLinks(disciplineId, taskId)
}

func (c *TaskController) CreateTaskLink(ctx *fiber.Ctx, params *task.LinkParams) (*task.Link, error) {
	disciplineId, err := httputil.ExtractId(ctx, "discipline_id")
	if err != nil {
		return nil, err
	}

	taskId, err := httputil.ExtractId(ctx, "task_id")
	if err != nil {
		return nil, err
	}

	err = auth.Authorize(ctx, c.disciplineRepository.GetCurriculumPrivileges, auth.LevelPublic, auth.ActionUpdate, disciplineId)
	if err != nil {
		return nil, err
	}

	err = params.Validate()
	if err != nil {
		return nil, err
	}

	l := &task.Link{
		Name: params.Name,
		URL:  params.URL,
	}

	err = c.taskRepository.CreateTaskLink(taskId, l)
	if err != nil {
		return nil, err
	}

	ctx.Status(fiber.StatusCreated)
	return l, nil
}

func (c *TaskController) UpdateTaskLink(ctx *fiber.Ctx, params *task.LinkParams) (*task.Link, error) {
	disciplineId, err := httputil.ExtractId(ctx, "discipline_id")
	if err != nil {
		return nil, err
	}

	err = auth.Authorize(ctx, c.disciplineRepository.GetCurriculumPrivileges, auth.LevelPublic, auth.ActionUpdate, disciplineId)
	if err != nil {
		return nil, err
	}

	taskLinkId, err := httputil.ExtractId(ctx, "task_link_id")
	if err != nil {
		return nil, err
	}

	err = params.Validate()
	if err != nil {
		return nil, err
	}

	l := &task.Link{
		ID:   taskLinkId,
		Name: params.Name,
		URL:  params.URL,
	}

	err = c.taskRepository.UpdateTaskLink(l)
	if err != nil {
		return nil, err
	}

	return l, nil
}

func (c *TaskController) DeleteTaskLink(ctx *fiber.Ctx) (*any, error) {
	disciplineId, err := httputil.ExtractId(ctx, "discipline_id")
	if err != nil {
		return nil, err
	}

	err = auth.Authorize(ctx, c.disciplineRepository.GetCurriculumPrivileges, auth.LevelPublic, auth.ActionUpdate, disciplineId)
	if err != nil {
		return nil, err
	}

	taskLinkId, err := httputil.ExtractId(ctx, "task_link_id")
	if err != nil {
		return nil, err
	}

	err = c.taskRepository.DeleteTaskLink(taskLinkId)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (c *TaskController) GetTaskProgress(ctx *fiber.Ctx) (*user.TaskProgress, error) {
	disciplineId, err := httputil.ExtractId(ctx, "discipline_id")
	if err != nil {
		return nil, err
	}

	userId := ctx.Locals("userid").(int64)
	taskId, err := httputil.ExtractId(ctx, "task_id")
	if err != nil {
		return nil, err
	}

	t, err := c.taskRepository.GetTask(disciplineId, taskId)
	if err != nil {
		return nil, err
	}

	err = auth.Authorize(ctx, c.taskRepository.GetCurriculumPrivileges, auth.LevelPublic, auth.ActionRead, t.GroupID)
	if err != nil {
		return nil, err
	}

	return c.userRepository.GetProgress(userId, taskId)
}

func (c *TaskController) UpdateTaskProgress(ctx *fiber.Ctx, params *task.UpdateProgressParams) (*user.TaskProgress, error) {
	disciplineId, err := httputil.ExtractId(ctx, "discipline_id")
	if err != nil {
		return nil, err
	}

	userId := ctx.Locals("userid").(int64)
	taskId, err := httputil.ExtractId(ctx, "task_id")
	if err != nil {
		return nil, err
	}

	t, err := c.taskRepository.GetTask(disciplineId, taskId)
	if err != nil {
		return nil, err
	}

	err = auth.Authorize(ctx, c.taskRepository.GetCurriculumPrivileges, auth.LevelPublic, auth.ActionRead, t.GroupID)
	if err != nil {
		return nil, err
	}

	oldProgress, err := c.userRepository.GetProgress(userId, taskId)
	if err != nil {
		return nil, err
	}

	p := &user.TaskProgress{
		GenericTaskProgress: user.GenericTaskProgress{
			Status:      params.Status,
			Grade:       params.Grade,
			StartedAt:   oldProgress.StartedAt,
			CompletedAt: oldProgress.CompletedAt,
		},
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
