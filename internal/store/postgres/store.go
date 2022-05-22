package postgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/tmrrwnxtsn/todo-lists-api/internal/store"
)

var _ store.Store = (*Store)(nil)

type Store struct {
	db                 *sqlx.DB
	userRepository     *UserRepository
	todoListRepository *TodoListRepository
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}

func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = NewUserRepository(s)

	return s.userRepository
}

func (s *Store) TodoList() store.TodoListRepository {
	if s.todoListRepository != nil {
		return s.todoListRepository
	}

	s.todoListRepository = NewTodoListRepository(s)

	return s.todoListRepository
}
