package task

type Repository interface {
	GetGroups(disciplineId int64) (*[]Group, error)
	GetGroup(disciplineId int64, groupId int64) (*Group, error)
	GetTasks(disciplineId int64) (*[]Task, error)
	GetTask(disciplineId int64, taskId int64) (*Task, error)
	GetTaskLinks(disciplineId int64, taskId int64) (*[]Link, error)
}
