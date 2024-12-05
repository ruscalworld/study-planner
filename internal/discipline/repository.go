package discipline

import "github.com/ruscalworld/study-planner/internal/access"

type Repository interface {
	GetDisciplines(curriculumId int64) (*[]Discipline, error)
	GetDiscipline(curriculumId int64, id int64) (*Discipline, error)
	GetDisciplineLinks(id int64) (*[]Link, error)
	GetCurriculumPrivileges(disciplineId int64, userId int64) (*access.CurriculumPrivileges, error)
}
