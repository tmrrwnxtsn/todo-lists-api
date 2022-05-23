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
	GetById(userId, listId uint64) (model.TodoList, error)
	// Update ...
	Update(userId, listId uint64, data model.UpdateTodoListData) error
	// Delete ...
	Delete(userId, listId uint64) error
}

// TodoListItemRepository ...
type TodoListItemRepository interface {
	// Create ...
	Create(listId uint64, item model.TodoListItem) (uint64, error)
	// GetAll ...
	GetAll(userId, listId uint64) ([]model.TodoListItem, error)
	// GetById ...
	GetById(userId, itemId uint64) (model.TodoListItem, error)
	// Update ...
	Update(userId, itemId uint64, data model.UpdateTodoListItemData) error
	// Delete ...
	Delete(userId, itemId uint64) error
}
