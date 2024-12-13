package repository

import (
	"github.com/ruscalworld/study-planner/internal/stats"

	"github.com/jmoiron/sqlx"
)

type MySqlRepository struct {
	db *sqlx.DB
}

func NewMySqlRepository(db *sqlx.DB) *MySqlRepository {
	return &MySqlRepository{db: db}
}

const undoneUpcomingTasksQuery = `
select
    t.id as 'id',
    t.name as 'name',
    t.external_name as 'external_name',
    t.description as 'description',
    t.task_group_id as 'task_group_id',
    t.status as 'status',
    t.difficulty as 'difficulty',
    t.deadline as 'deadline',
    d.id as 'discipline.id',
    d.name as 'discipline.name'
from tasks t
         left join user_task_progress utp on t.id = utp.task_id and user_id = ?
         join task_groups tg on tg.id = t.task_group_id
         join disciplines d on d.id = tg.discipline_id
         join curriculums c on c.id = d.curriculum_id
where
    c.id = ? and
    t.deadline is not null and
    t.status = 'Available' and
    (utp.status != 'Completed' or utp.status is null)
order by t.deadline
limit ?
`

func (m *MySqlRepository) GetUndoneUpcomingTasks(curriculumId int64, userId int64, count int) (*[]stats.DisciplineTask, error) {
	tasks := make([]stats.DisciplineTask, 0)
	err := m.db.Select(&tasks, undoneUpcomingTasksQuery, userId, curriculumId, count)
	if err != nil {
		return nil, err
	}

	return &tasks, nil
}
