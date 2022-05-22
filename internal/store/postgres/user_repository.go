package postgres

import (
	"fmt"
	"github.com/tmrrwnxtsn/todo-lists-api/internal/model"
	"github.com/tmrrwnxtsn/todo-lists-api/internal/store"
)

var _ store.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	store *Store
}

func NewUserRepository(store *Store) *UserRepository {
	return &UserRepository{store: store}
}

func (r *UserRepository) Create(user model.User) (uint64, error) {
	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) VALUES ($1, $2, $3) RETURNING id", usersTable)

	var id uint64
	err := r.store.db.QueryRow(query, user.Name, user.Username, user.Password).Scan(&id)
	return id, err
}

func (r *UserRepository) Get(username, passwordHash string) (model.User, error) {
	query := fmt.Sprintf(`SELECT id, name, username, password_hash 
								FROM %s 
								WHERE username = $1 AND password_hash = $2`,
		usersTable)

	var user model.User
	err := r.store.db.Get(&user, query, username, passwordHash)
	return user, err
}
