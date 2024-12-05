package repository

import (
	"database/sql"
	"errors"

	"github.com/ruscalworld/study-planner/internal/access"
	"github.com/ruscalworld/study-planner/internal/auth"
	"github.com/ruscalworld/study-planner/internal/discipline"

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

func (m *MySqlRepository) GetCurriculumPrivileges(disciplineId int64, userId int64) (*access.CurriculumPrivileges, error) {
	var cp access.CurriculumPrivileges
	err := m.db.Get(&cp, "select c.id as curriculum_id, uc.role as role from  disciplines d join curriculums c on d.curriculum_id = c.id join user_curriculums uc on c.id = uc.curriculum_id where d.id = ? and uc.user_id = ?", disciplineId, userId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, auth.ErrNoPrivileges
		}

		return nil, err
	}

	return &cp, nil
}
