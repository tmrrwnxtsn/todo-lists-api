package service

import "github.com/tmrrwnxtsn/todo-lists-api/internal/store"

type Authorization interface {
}

type TodoList interface {
}

type TodoItem interface {
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(store *store.Store) *Service {
	return &Service{}
}
