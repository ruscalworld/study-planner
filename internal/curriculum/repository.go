package curriculum

import "github.com/ruscalworld/study-planner/internal/access"

type Repository interface {
	GetCurriculum(id int64) (*Curriculum, error)
	CreateCurriculum(institutionId int64, curriculum *Curriculum) error
	GetInstitutionCurriculums(institutionId int64) (*[]Curriculum, error)

	GetCurriculumCodes(curriculumId int64) (*[]Code, error)
	CreateCurriculumCode(code *Code) error
	DeleteCurriculumCode(id int64) error
	GetCurriculumByCode(code string) (*Privileges, error)

	GetCurriculumPrivileges(curriculumId int64, userId int64) (*access.CurriculumPrivileges, error)
	GetCurriculumUsers(curriculumId int64) (*[]User, error)
	UpdateCurriculumUser(curriculumId int64, userId int64, role access.Role) error
	DeleteCurriculumUser(curriculumId int64, userId int64) error
}
