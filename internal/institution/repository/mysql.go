package repository

import (
	"database/sql"
	"errors"

	"study-planner/internal/institution"

	"github.com/jmoiron/sqlx"
)

type MySqlRepository struct {
	db *sqlx.DB
}

func NewMySqlRepository(db *sqlx.DB) *MySqlRepository {
	return &MySqlRepository{db: db}
}

func (m *MySqlRepository) GetInstitutions() (*[]institution.Institution, error) {
	i := make([]institution.Institution, 0)
	err := m.db.Select(&i, "select id, name from institutions")
	if err != nil {
		return nil, err
	}

	return &i, nil
}

func (m *MySqlRepository) GetInstitution(id int64) (*institution.Institution, error) {
	var i institution.Institution
	err := m.db.Get(&i, "select id, name from institutions where id = ?", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, institution.ErrUnknownInstitution
		}

		return nil, err
	}

	return &i, nil
}
