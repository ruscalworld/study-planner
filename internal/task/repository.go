package task

import (
	"github.com/ruscalworld/study-planner/internal/access"
)

type Repository interface {
	GetGroups(disciplineId int64) (*[]Group, error)
	GetGroup(disciplineId int64, groupId int64) (*Group, error)
	CreateGroup(disciplineId int64, group *Group) error
	UpdateGroup(group *Group) error
	DeleteGroup(taskGroupId int64) error

	GetTasks(disciplineId int64) (*[]Task, error)
	GetTask(disciplineId int64, taskId int64) (*Task, error)
	CreateTask(taskGroupId int64, t *Task) error
	UpdateTask(task *Task) error
	DeleteTask(taskId int64) error

	GetTaskLinks(disciplineId int64, taskId int64) (*[]Link, error)
	CreateTaskLink(taskId int64, t *Link) error
	UpdateTaskLink(link *Link) error
	DeleteTaskLink(taskLinkId int64) error

	GetCurriculumPrivileges(taskId int64, userId int64) (*access.CurriculumPrivileges, error)
}
