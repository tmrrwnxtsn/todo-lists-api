package service

import (
	"github.com/tmrrwnxtsn/todo-lists-api/internal/model"
	"github.com/tmrrwnxtsn/todo-lists-api/internal/store"
)

type TodoList interface {
	Create(userId uint64, list model.TodoList) (uint64, error)
	GetAll(userId uint64) ([]model.TodoList, error)
	GetById(userId, listId uint64) (model.TodoList, error)
	Update(userId, listId uint64, data model.UpdateTodoListData) error
	Delete(userId, listId uint64) error
}

type TodoListService struct {
	listRepository store.TodoListRepository
}

func NewTodoListService(repo store.TodoListRepository) *TodoListService {
	return &TodoListService{listRepository: repo}
}

func (s *TodoListService) Create(userId uint64, list model.TodoList) (uint64, error) {
	return s.listRepository.Create(userId, list)
}

func (s *TodoListService) GetAll(userId uint64) ([]model.TodoList, error) {
	return s.listRepository.GetAll(userId)
}

func (s *TodoListService) GetById(userId, listId uint64) (model.TodoList, error) {
	return s.listRepository.GetById(userId, listId)
}

func (s *TodoListService) Update(userId, listId uint64, data model.UpdateTodoListData) error {
	if err := data.Validate(); err != nil {
		return err
	}
	return s.listRepository.Update(userId, listId, data)
}

func (s *TodoListService) Delete(userId, listId uint64) error {
	return s.listRepository.Delete(userId, listId)
}
