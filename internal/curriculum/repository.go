package curriculum

type Repository interface {
	GetCurriculum(id int64) (*Curriculum, error)
	GetInstitutionCurriculums(institutionId int64) (*[]Curriculum, error)
	GetUserCurriculums(userId int64) (*[]Curriculum, error)
}
