package repository

import (
	"database/sql"
	"errors"

	"github.com/ruscalworld/study-planner/internal/access"
	"github.com/ruscalworld/study-planner/internal/auth"
	"github.com/ruscalworld/study-planner/internal/curriculum"

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

func (m *MySqlRepository) CreateCurriculum(institutionId int64, c *curriculum.Curriculum) error {
	result, err := m.db.Exec(
		"insert into curriculums (name, semester, institution_id) values (?, ?, ?)",
		c.Name, c.Semester, institutionId,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	c.ID = id
	return nil
}

func (m *MySqlRepository) GetInstitutionCurriculums(institutionId int64) (*[]curriculum.Curriculum, error) {
	c := make([]curriculum.Curriculum, 0)
	err := m.db.Select(&c, "select id, name, semester from curriculums where institution_id = ?", institutionId)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (m *MySqlRepository) GetCurriculumCodes(curriculumId int64) (*[]curriculum.Code, error) {
	c := make([]curriculum.Code, 0)
	err := m.db.Select(&c, "select id, code, role, curriculum_id from curriculum_codes where curriculum_id = ?", curriculumId)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (m *MySqlRepository) CreateCurriculumCode(code *curriculum.Code) error {
	result, err := m.db.Exec(
		"insert into curriculum_codes (code, role, curriculum_id, user_id, expires_at, created_at) values (?, ?, ?, ?, ?, ?)",
		code.Code, code.Role, code.CurriculumId, code.UserId, code.ExpiresAt, code.CreatedAt,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	code.ID = id
	return nil
}

func (m *MySqlRepository) DeleteCurriculumCode(id int64) error {
	_, err := m.db.Exec("delete from curriculum_codes where id = ?", id)
	return err
}

func (m *MySqlRepository) GetCurriculumByCode(code string) (*curriculum.Privileges, error) {
	var cp curriculum.Privileges

	err := m.db.Get(&cp, "select c.id as id, c.name as name, c.semester as semester, cc.role as role from curriculum_codes cc join curriculums c on c.id = cc.curriculum_id where cc.code = ? and cc.expires_at > current_timestamp", code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, curriculum.ErrUnknownCurriculum
		}

		return nil, err
	}

	return &cp, nil
}

func (m *MySqlRepository) GetCurriculumPrivileges(curriculumId int64, userId int64) (*access.CurriculumPrivileges, error) {
	var cp access.CurriculumPrivileges
	err := m.db.Get(&cp, "select c.id as curriculum_id, uc.role as role from curriculums c join user_curriculums uc where c.id = ? and uc.user_id = ?", curriculumId, userId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, auth.ErrNoPrivileges
		}

		return nil, err
	}

	return &cp, nil
}
