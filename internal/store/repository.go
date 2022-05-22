package store

import "github.com/tmrrwnxtsn/todo-lists-api/internal/model"

// AuthRepository ...
type AuthRepository interface {
	CreateUser(user model.User) (uint64, error)
	GetUser(username, passwordHash string) (model.User, error)
}

type TodoList interface {
}

type TodoItem interface {
}
