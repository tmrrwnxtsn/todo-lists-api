package store

import (
	"github.com/tmrrwnxtsn/todo-lists-api/internal/model"
)

// UserRepository ...
type UserRepository interface {
	// Create ...
	Create(user model.User) (uint64, error)
	// Get ...
	Get(username, passwordHash string) (model.User, error)
}

// TodoListRepository ...
type TodoListRepository interface {
	// Create ...
	Create(userId uint64, list model.TodoList) (uint64, error)
	// GetAll ...
	GetAll(userId uint64) ([]model.TodoList, error)
	// GetById ...
	GetById(userId, todoListId uint64) (model.TodoList, error)
	// Update ...
	Update(userId, todoListId uint64, data model.UpdateTodoListData) error
	// Delete ...
	Delete(userId, todoListId uint64) error
}

type TodoItem interface {
}
