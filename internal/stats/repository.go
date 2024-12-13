package stats

type Repository interface {
	GetUndoneUpcomingTasks(curriculumId int64, userId int64, count int) (*[]DisciplineTask, error)
}
