package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/ruscalworld/study-planner/internal/draft"
	"github.com/ruscalworld/study-planner/internal/task"

	"github.com/jmoiron/sqlx"
)

type MySqlRepository struct {
	db *sqlx.DB
}

func NewMySqlRepository(db *sqlx.DB) *MySqlRepository {
	return &MySqlRepository{db: db}
}

func (r *MySqlRepository) GetDraft(id int64, userId int64) (*draft.Draft, error) {
	var d draft.Draft

	err := r.db.Get(&d, `select id, user_id, text, created_at from drafts where id = ? and user_id = ?`, id, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, draft.ErrUnknownDraft
		}

		return nil, err
	}

	return &d, nil
}

func (r *MySqlRepository) GetUserDrafts(userId int64) (*[]draft.Draft, error) {
	t := make([]draft.Draft, 0)

	err := r.db.Select(&t, "select id, user_id, text, created_at from drafts where user_id = ?", userId)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *MySqlRepository) CreateDraft(draft *draft.Draft) error {
	result, err := r.db.Exec("insert into drafts (user_id, text) values (?, ?)", draft.UserID, draft.Text)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	draft.ID = id
	return nil
}

func (r *MySqlRepository) UpdateDraft(draft *draft.Draft) error {
	_, err := r.db.Exec("update drafts set text = ? where id = ?", draft.Text, draft.ID)
	return err
}

func (r *MySqlRepository) DeleteDraft(id int64) error {
	_, err := r.db.Exec("delete from drafts where id = ?", id)
	return err
}

func (r *MySqlRepository) MoveDraft(draft *draft.Draft, taskGroup *task.Group, taskName string) (*task.Task, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}

	result, err := tx.Exec(
		"insert into tasks (name, description, task_group_id, status) values (?, ?, ?, 'Available')",
		taskName, draft.Text, taskGroup.ID,
	)
	if err != nil {
		terr := tx.Rollback()
		if terr != nil {
			return nil, fmt.Errorf("rolling back transaction due to error %q failed: %w", err, terr)
		}

		return nil, terr
	}

	taskId, err := result.LastInsertId()
	if err != nil {
		terr := tx.Rollback()
		if terr != nil {
			return nil, fmt.Errorf("rolling back transaction due to error %q failed: %w", err, terr)
		}

		return nil, err
	}

	_, err = tx.Exec("delete from drafts where id = ?", draft.ID)
	if err != nil {
		terr := tx.Rollback()
		if terr != nil {
			return nil, fmt.Errorf("rolling back transaction due to error %q failed: %w", err, terr)
		}

		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		terr := tx.Rollback()
		if terr != nil {
			return nil, fmt.Errorf("rolling back transaction due to error %q failed: %w", err, terr)
		}

		return nil, err
	}

	return &task.Task{
		ID:          taskId,
		Name:        taskName,
		Description: &draft.Text,
		GroupID:     taskGroup.ID,
		Status:      task.StatusAvailable,
		Difficulty:  1,
	}, nil
}
