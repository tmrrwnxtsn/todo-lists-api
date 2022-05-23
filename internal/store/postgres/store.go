package postgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/tmrrwnxtsn/todo-lists-api/internal/store"
)

var _ store.Store = (*Store)(nil)

type Store struct {
	db             *sqlx.DB
	userRepository *UserRepository
	listRepository *TodoListRepository
	itemRepository *TodoListItemRepository
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
	if s.listRepository != nil {
		return s.listRepository
	}

	s.listRepository = NewTodoListRepository(s)

	return s.listRepository
}

func (s *Store) TodoListItem() store.TodoListItemRepository {
	if s.itemRepository != nil {
		return s.itemRepository
	}

	s.itemRepository = NewTodoListItemRepository(s)

	return s.itemRepository
}
