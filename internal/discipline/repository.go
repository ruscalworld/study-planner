package discipline

type Repository interface {
	GetDisciplines(curriculumId int64) (*[]Discipline, error)
	GetDiscipline(curriculumId int64, id int64) (*Discipline, error)
	GetDisciplineLinks(id int64) (*[]Link, error)
}
