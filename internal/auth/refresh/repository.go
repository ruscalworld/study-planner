package refresh

type Repository interface {
	GetValidToken(prefix []byte) (*Token, error)
	CreateToken(token *Token) error
	DeleteToken(id int64) error
}
