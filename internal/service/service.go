package service

import (
	"github.com/tmrrwnxtsn/todo-lists-api/internal/store"
)

type Service struct {
	AuthService         Authorization
	TodoListService     TodoList
	TodoListItemService TodoListItem
}

func NewService(store store.Store) *Service {
	return &Service{
		AuthService:         NewAuthService(store.User()),
		TodoListService:     NewTodoListService(store.TodoList()),
		TodoListItemService: NewTodoListItemService(store.TodoListItem(), store.TodoList()),
	}
}
