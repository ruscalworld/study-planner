package discipline

import "github.com/ruscalworld/study-planner/internal/access"

type Repository interface {
	GetDisciplines(curriculumId int64) (*[]Discipline, error)
	GetDiscipline(curriculumId int64, id int64) (*Discipline, error)
	CreateDiscipline(curriculumId int64, d *Discipline) error
	UpdateDiscipline(d *Discipline) error
	DeleteDiscipline(curriculumId int64, id int64) error

	GetDisciplineLinks(id int64) (*[]Link, error)
	CreateDisciplineLink(disciplineId int64, l *Link) error
	UpdateDisciplineLink(d *Link) error
	DeleteDisciplineLink(disciplineLinkId int64) error

	GetCurriculumPrivileges(disciplineId int64, userId int64) (*access.CurriculumPrivileges, error)
}
