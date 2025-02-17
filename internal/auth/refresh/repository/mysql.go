package repository

import (
	"github.com/ruscalworld/study-planner/internal/auth/refresh"

	"github.com/jmoiron/sqlx"
)

type MySqlRepository struct {
	db *sqlx.DB
}

func NewMySqlRepository(db *sqlx.DB) *MySqlRepository {
	return &MySqlRepository{db: db}
}

func (r *MySqlRepository) GetValidToken(prefix []byte) (*refresh.Token, error) {
	var token refresh.Token

	err := r.db.Get(
		&token,
		"select id, user_id, prefix, token, created_at, expires_at from refresh_tokens where prefix = ? and expires_at > current_timestamp",
		prefix,
	)

	if err != nil {
		return nil, refresh.ErrInvalidToken
	}

	return &token, nil
}

func (r *MySqlRepository) CreateToken(token *refresh.Token) error {
	result, err := r.db.Exec(
		"insert into refresh_tokens (user_id, prefix, token, created_at, expires_at) values (?, ?, ?, ?, ?)",
		token.UserId, token.Prefix, token.Token, token.CreatedAt, token.ExpiresAt,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	token.ID = id
	return nil
}

func (r *MySqlRepository) DeleteToken(id int64) error {
	_, err := r.db.Exec("delete from refresh_tokens where id = ?", id)
	return err
}
