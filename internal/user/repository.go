package user

type Repository interface {
	GetUserById(userId int64) (*User, error)
	GetUserByExternalId(externalId string) (*User, error)
	RegisterUser(user *User) error

	GetGoal(userId int64, taskGroupId int64) (*Goal, error)
	StoreGoal(userId int64, taskGroupId int64, goal *Goal) error
	DeleteGoal(userId int64, taskGroupId int64, goalId int64) error

	GetProgress(userId int64, taskId int64) (*TaskProgress, error)
	StoreProgress(userId int64, taskId int64, progress *TaskProgress) error
	GetDisciplineProgress(userId int64, disciplineId int64) (*[]ScopedTaskProgress, error)
	GetDisciplineStats(userId int64, disciplineId int64) (*GenericStats, error)
}
