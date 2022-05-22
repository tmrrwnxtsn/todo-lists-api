package postgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/tmrrwnxtsn/todo-lists-api/internal/store"
)

var _ store.Store = (*Store)(nil)

type Store struct {
	db             *sqlx.DB
	authRepository *AuthRepository
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Auth() store.AuthRepository {
	if s.authRepository != nil {
		return s.authRepository
	}

	s.authRepository = NewAuthRepository(s)

	return s.authRepository
}
