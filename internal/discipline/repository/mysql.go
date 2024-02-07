package repository

import (
	"database/sql"
	"errors"

	"study-planner/internal/discipline"

	"github.com/jmoiron/sqlx"
)

type MySqlRepository struct {
	db *sqlx.DB
}

func NewMySqlRepository(db *sqlx.DB) *MySqlRepository {
	return &MySqlRepository{db: db}
}

func (m *MySqlRepository) GetDisciplines(curriculumId int64) (*[]discipline.Discipline, error) {
	d := make([]discipline.Discipline, 0)
	err := m.db.Select(&d, "select id, name from disciplines where curriculum_id = ?", curriculumId)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

func (m *MySqlRepository) GetDiscipline(curriculumId int64, id int64) (*discipline.Discipline, error) {
	var d discipline.Discipline
	err := m.db.Get(&d, "select id, name from disciplines where curriculum_id = ? and id = ?", curriculumId, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, discipline.ErrUnknownDiscipline
		}

		return nil, err
	}

	return &d, nil
}

func (m *MySqlRepository) GetDisciplineLinks(id int64) (*[]discipline.Link, error) {
	l := make([]discipline.Link, 0)
	err := m.db.Select(&l, "select id, name, url from discipline_links where discipline_id = ?", id)
	if err != nil {
		return nil, err
	}

	return &l, nil
}
