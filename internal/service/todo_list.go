package service

import (
	"github.com/tmrrwnxtsn/todo-lists-api/internal/model"
	"github.com/tmrrwnxtsn/todo-lists-api/internal/store"
)

type TodoList interface {
	Create(userId uint64, list model.TodoList) (uint64, error)
	GetAll(userId uint64) ([]model.TodoList, error)
	GetById(userId, todoListId uint64) (model.TodoList, error)
	Update(userId, todoListId uint64, data model.UpdateTodoListData) error
	Delete(userId, todoListId uint64) error
}

type TodoListService struct {
	repo store.TodoListRepository
}

func NewTodoListService(repo store.TodoListRepository) *TodoListService {
	return &TodoListService{repo: repo}
}

func (s *TodoListService) Create(userId uint64, list model.TodoList) (uint64, error) {
	return s.repo.Create(userId, list)
}

func (s *TodoListService) GetAll(userId uint64) ([]model.TodoList, error) {
	return s.repo.GetAll(userId)
}

func (s *TodoListService) GetById(userId, todoListId uint64) (model.TodoList, error) {
	return s.repo.GetById(userId, todoListId)
}

func (s *TodoListService) Update(userId, todoListId uint64, data model.UpdateTodoListData) error {
	if err := data.Validate(); err != nil {
		return err
	}
	return s.repo.Update(userId, todoListId, data)
}

func (s *TodoListService) Delete(userId, todoListId uint64) error {
	return s.repo.Delete(userId, todoListId)
}
