package institution

type Repository interface {
	GetInstitutions() (*[]Institution, error)
	GetInstitution(id int64) (*Institution, error)
}
