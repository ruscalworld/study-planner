package repository

import (
	"database/sql"
	"errors"

	"study-planner/internal/task"

	"github.com/jmoiron/sqlx"
)

type MySqlRepository struct {
	db *sqlx.DB
}

func NewMySqlRepository(db *sqlx.DB) *MySqlRepository {
	return &MySqlRepository{db: db}
}

func (m *MySqlRepository) GetGroups(disciplineId int64) (*[]task.Group, error) {
	g := make([]task.Group, 0)
	err := m.db.Select(&g, "select id, name from task_groups where discipline_id = ?", disciplineId)
	if err != nil {
		return nil, err
	}

	return &g, nil
}

func (m *MySqlRepository) GetGroup(disciplineId int64, groupId int64) (*task.Group, error) {
	var g task.Group
	err := m.db.Get(&g, "select id, name from task_groups where discipline_id = ? and id = ?", disciplineId, groupId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, task.ErrUnknownGroup
		}

		return nil, err
	}

	return &g, nil
}

func (m *MySqlRepository) GetTasks(disciplineId int64) (*[]task.Task, error) {
	t := make([]task.Task, 0)
	err := m.db.Select(&t, "select t.id, t.name, t.external_name, t.description, t.task_group_id, t.status, t.difficulty, t.deadline from tasks t join task_groups g on t.task_group_id = g.id where g.discipline_id = ?", disciplineId)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (m *MySqlRepository) GetTask(disciplineId int64, taskId int64) (*task.Task, error) {
	var t task.Task
	err := m.db.Get(&t, "select t.id, t.name, t.external_name, t.description, t.task_group_id, t.status, t.difficulty, t.deadline from tasks t join task_groups g on t.task_group_id = g.id where g.discipline_id = ? and t.id = ?", disciplineId, taskId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, task.ErrUnknownTask
		}

		return nil, err
	}

	return &t, nil
}

func (m *MySqlRepository) GetTaskLinks(disciplineId int64, taskId int64) (*[]task.Link, error) {
	l := make([]task.Link, 0)
	err := m.db.Select(&l, "select l.id, l.name, l.url from task_links l join tasks t on l.task_id = t.id join task_groups g on t.task_group_id = g.id where g.discipline_id = ? and t.id = ?", disciplineId, taskId)
	if err != nil {
		return nil, err
	}

	return &l, nil
}
