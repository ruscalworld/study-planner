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
			return &user.TaskProgress{Status: user.TaskStatusNotStarted}, nil
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
