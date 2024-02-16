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

func (m *MySqlRepository) GetGoals(userId int64, taskGroupId int64) (*[]user.Goal, error) {
	g := make([]user.Goal, 0)
	err := m.db.Select(&g, "select id, min_completed from user_goals where user_id = ? and task_group_id = ?", userId, taskGroupId)
	if err != nil {
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
