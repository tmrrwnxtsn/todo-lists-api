package service

import (
	"github.com/tmrrwnxtsn/todo-lists-api/internal/store"
)

type TodoItem interface {
}

type Service struct {
	AuthService     Authorization
	TodoListService TodoList
	TodoItem
}

func NewService(store store.Store) *Service {
	return &Service{
		AuthService:     NewAuthService(store.User()),
		TodoListService: NewTodoListService(store.TodoList()),
	}
}
