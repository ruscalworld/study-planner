package repository

import (
	"database/sql"
	"errors"

	"study-planner/internal/curriculum"

	"github.com/jmoiron/sqlx"
)

type MySqlRepository struct {
	db *sqlx.DB
}

func NewMySqlRepository(db *sqlx.DB) *MySqlRepository {
	return &MySqlRepository{db: db}
}

func (m *MySqlRepository) GetCurriculum(id int64) (*curriculum.Curriculum, error) {
	var c curriculum.Curriculum
	err := m.db.Get(&c, "select id, name, semester from curriculums where id = ?", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, curriculum.ErrUnknownCurriculum
		}

		return nil, err
	}

	return &c, nil
}

func (m *MySqlRepository) GetUserCurriculums(userId int64) (*[]curriculum.Curriculum, error) {
	c := make([]curriculum.Curriculum, 0)
	err := m.db.Select(&c, "select c.id, c.name, c.semester from user_curriculums u join curriculums c on c.id = u.curriculum_id where u.user_id = ?", userId)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (m *MySqlRepository) GetInstitutionCurriculums(institutionId int64) (*[]curriculum.Curriculum, error) {
	c := make([]curriculum.Curriculum, 0)
	err := m.db.Select(&c, "select id, name, semester from curriculums where institution_id = ?", institutionId)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
