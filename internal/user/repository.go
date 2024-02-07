package user

type Repository interface {
	GetUserByExternalId(externalId string) (*User, error)
	RegisterUser(user *User) error

	GetGoals(userId int64, taskGroupId int64) (*[]Goal, error)
	StoreGoal(userId int64, taskGroupId int64, goal *Goal) error
	DeleteGoal(userId int64, taskGroupId int64, goalId int64) error
}
