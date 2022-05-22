package service

import (
	"github.com/tmrrwnxtsn/todo-lists-api/internal/store"
)

type TodoList interface {
}

type TodoItem interface {
}

type Service struct {
	AuthService Authorization
	TodoList
	TodoItem
}

func NewService(store store.Store) *Service {
	return &Service{
		AuthService: NewAuthService(store.Auth()),
	}
}
