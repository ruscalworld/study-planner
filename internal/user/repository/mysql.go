package repository

import (
	"database/sql"
	"errors"

	"study-planner/internal/user"

	"github.com/jmoiron/sqlx"
)

type MySqlRepository struct {
	db *sqlx.DB
}

func NewMySqlRepository(db *sqlx.DB) *MySqlRepository {
	return &MySqlRepository{db: db}
}

func (m *MySqlRepository) GetUserById(userId int64) (*user.User, error) {
	var u user.User
	err := m.db.Get(&u, "select id, name, avatar_url, external_id, platform, created_at from users where id = ?", userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrUnknownUser
		}

		return nil, err
	}

	return &u, nil
}

func (m *MySqlRepository) GetUserByExternalId(externalId string) (*user.User, error) {
	var u user.User
	err := m.db.Get(&u, "select id, name, avatar_url, external_id, platform, created_at from users where external_id = ?", externalId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrUnknownUser
		}

		return nil, err
	}

	return &u, nil
}

func (m *MySqlRepository) RegisterUser(user *user.User) error {
	result, err := m.db.Exec(
		"insert into users (name, avatar_url, platform, external_id) values (?, ?, ?, ?)",
		user.Name, user.AvatarURL, user.Platform, user.ExternalID,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = id
	return nil
}

func (m *MySqlRepository) GetGoal(userId int64, taskGroupId int64) (*user.Goal, error) {
	var g user.Goal
	err := m.db.Get(&g, "select id, min_completed from user_goals where user_id = ? and task_group_id = ?", userId, taskGroupId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &user.Goal{}, nil
		}

		return nil, err
	}

	return &g, nil
}

func (m *MySqlRepository) StoreGoal(userId int64, taskGroupId int64, goal *user.Goal) error {
	result, err := m.db.Exec(
		"insert into user_goals (user_id, task_group_id, min_completed) values (?, ?, ?) on duplicate key update min_completed = ?",
		userId, taskGroupId, goal.MinCompleted, goal.MinCompleted,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	goal.ID = id
	return nil
}

func (m *MySqlRepository) DeleteGoal(userId int64, taskGroupId int64, goalId int64) error {
	_, err := m.db.Exec("delete from user_goals where user_id = ? and task_group_id = ? and id = ?", userId, taskGroupId, goalId)
	return err
}

func (m *MySqlRepository) GetProgress(userId int64, taskId int64) (*user.TaskProgress, error) {
	var p user.TaskProgress
	err := m.db.Get(&p,
		"select id, status, grade, started_at, completed_at from user_task_progress where user_id = ? and task_id = ?",
		userId, taskId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &user.TaskProgress{
				GenericTaskProgress: user.GenericTaskProgress{Status: user.TaskStatusNotStarted},
			}, nil
		}

		return nil, err
	}

	return &p, nil
}

func (m *MySqlRepository) StoreProgress(userId int64, taskId int64, progress *user.TaskProgress) error {
	result, err := m.db.Exec(
		"insert into user_task_progress (user_id, task_id, status, grade, started_at, completed_at) values (?, ?, ?, ?, ?, ?) on duplicate key update status = ?, grade = ?, started_at = ?, completed_at = ?",
		userId, taskId, progress.Status, progress.Grade, progress.StartedAt, progress.CompletedAt, progress.Status, progress.Grade, progress.StartedAt, progress.CompletedAt,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	progress.ID = id
	return nil
}

func (m *MySqlRepository) GetDisciplineProgress(userId int64, disciplineId int64) (*[]user.ScopedTaskProgress, error) {
	p := make([]user.ScopedTaskProgress, 0)
	err := m.db.Select(&p,
		"select p.id, p.task_id, t.task_group_id, p.status, p.grade, p.started_at, p.completed_at from user_task_progress p join tasks t on t.id = p.task_id join task_groups g on g.id = t.task_group_id where p.user_id = ? and g.discipline_id = ?",
		userId, disciplineId,
	)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

const taskGroupStatsQuery = `
with t_task_statuses as (
	select 
	    p.status as status,
	    count(*) as count
	from user_task_progress p 
	    join tasks t on t.id = p.task_id
	    join task_groups g on g.id = t.task_group_id
	where p.user_id = ? and g.discipline_id = ?
	group by p.status
), t_all_tasks as (
    select 
        t.status as status,
        count(*) as count 
    from tasks t 
        join task_groups g on t.task_group_id = g.id
    where g.discipline_id = ?
    group by t.status
), t_goal_tasks as (
    select ug.min_completed
    from user_goals ug
        join task_groups g on g.id = ug.task_group_id
    where ug.user_id = ? and g.discipline_id = ?
)
select 
    coalesce((select count from t_task_statuses where status = 'Completed'), 0) as completed_tasks,
    coalesce((select sum(count) from t_task_statuses where status in ('InProgress', 'NeedsProtection')), 0) as in_progress_tasks,
    coalesce((select count from t_all_tasks where status = 'Available'), 0) as available_tasks,
    coalesce((select sum(count) from t_all_tasks), 0) as total_tasks,
    coalesce((select * from t_goal_tasks), 0) as goal_tasks
`

func (m *MySqlRepository) GetDisciplineStats(userId int64, disciplineId int64) (*user.GenericStats, error) {
	var s user.GenericStats
	err := m.db.Get(&s, taskGroupStatsQuery, userId, disciplineId, disciplineId, userId, disciplineId)
	if err != nil {
		return nil, err
	}

	return &s, nil
}
