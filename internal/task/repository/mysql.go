package repository

import (
	"database/sql"
	"errors"

	"github.com/ruscalworld/study-planner/internal/access"
	"github.com/ruscalworld/study-planner/internal/auth"
	"github.com/ruscalworld/study-planner/internal/task"

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

func (m *MySqlRepository) CreateGroup(disciplineId int64, group *task.Group) error {
	result, err := m.db.Exec("insert into task_groups (name, discipline_id) values (?, ?)", group.Name, disciplineId)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	group.ID = id
	return nil
}

func (m *MySqlRepository) UpdateGroup(group *task.Group) error {
	_, err := m.db.Exec("update task_groups set name = ? where id = ?", group.Name, group.ID)
	return err
}

func (m *MySqlRepository) DeleteGroup(taskGroupId int64) error {
	result, err := m.db.Exec("delete from task_groups where id = ?", taskGroupId)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return task.ErrUnknownGroup
	}

	return nil
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

func (m *MySqlRepository) CreateTask(taskGroupId int64, t *task.Task) error {
	result, err := m.db.Exec(
		"insert into tasks (name, external_name, description, task_group_id) values (?, ?, ?, ?)",
		t.Name, t.ExternalName, t.Description, taskGroupId,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	t.ID = id
	return nil
}

func (m *MySqlRepository) UpdateTask(t *task.Task) error {
	_, err := m.db.Exec(
		"update tasks set name = ?, external_name = ?, description = ?, deadline = ?, difficulty = ? where id = ?",
		t.Name, t.ExternalName, t.Description, t.Deadline, t.Difficulty, t.ID,
	)

	return err
}

func (m *MySqlRepository) DeleteTask(taskId int64) error {
	result, err := m.db.Exec("delete from tasks where id = ?", taskId)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return task.ErrUnknownTask
	}

	return nil
}

func (m *MySqlRepository) GetTaskLinks(disciplineId int64, taskId int64) (*[]task.Link, error) {
	l := make([]task.Link, 0)
	err := m.db.Select(&l, "select l.id, l.name, l.url from task_links l join tasks t on l.task_id = t.id join task_groups g on t.task_group_id = g.id where g.discipline_id = ? and t.id = ?", disciplineId, taskId)
	if err != nil {
		return nil, err
	}

	return &l, nil
}

func (m *MySqlRepository) CreateTaskLink(taskId int64, link *task.Link) error {
	result, err := m.db.Exec("insert into task_links (task_id, name, url) values (?, ?, ?)", taskId, link.Name, link.URL)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	link.ID = id
	return nil
}

func (m *MySqlRepository) UpdateTaskLink(link *task.Link) error {
	_, err := m.db.Exec("update task_links set name = ?, url = ? where id = ?", link.Name, link.URL, link.ID)
	return err
}

func (m *MySqlRepository) DeleteTaskLink(taskLinkId int64) error {
	result, err := m.db.Exec("delete from task_links where id = ?", taskLinkId)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return task.ErrUnknownTaskLink
	}

	return nil
}

func (m *MySqlRepository) GetCurriculumPrivileges(taskGroupId int64, userId int64) (*access.CurriculumPrivileges, error) {
	var cp access.CurriculumPrivileges
	err := m.db.Get(&cp, "select c.id as curriculum_id, uc.role as role from task_groups tg join disciplines d on tg.discipline_id = d.id join curriculums c on d.curriculum_id = c.id join user_curriculums uc on c.id = uc.curriculum_id where tg.id = ? and uc.user_id = ?", taskGroupId, userId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, auth.ErrNoPrivileges
		}

		return nil, err
	}

	return &cp, nil
}
