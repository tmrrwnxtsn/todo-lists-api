package service

import (
	"github.com/tmrrwnxtsn/todo-lists-api/internal/model"
	"github.com/tmrrwnxtsn/todo-lists-api/internal/store"
)

type TodoListItem interface {
	Create(userId, listId uint64, item model.TodoListItem) (uint64, error)
	GetAll(userId, listId uint64) ([]model.TodoListItem, error)
	GetById(userId, itemId uint64) (model.TodoListItem, error)
	Update(userId, itemId uint64, data model.UpdateTodoListItemData) error
	Delete(userId, itemId uint64) error
}

type TodoListItemService struct {
	itemRepository store.TodoListItemRepository
	listRepository store.TodoListRepository
}

func NewTodoListItemService(itemRepo store.TodoListItemRepository, listRepo store.TodoListRepository) *TodoListItemService {
	return &TodoListItemService{itemRepository: itemRepo, listRepository: listRepo}
}

func (s *TodoListItemService) Create(userId, listId uint64, item model.TodoListItem) (uint64, error) {
	_, err := s.listRepository.GetById(userId, listId)
	if err != nil {
		// list does not exists or not belongs to the user
		return 0, err
	}
	return s.itemRepository.Create(listId, item)
}

func (s *TodoListItemService) GetAll(userId, listId uint64) ([]model.TodoListItem, error) {
	return s.itemRepository.GetAll(userId, listId)
}

func (s *TodoListItemService) GetById(userId, itemId uint64) (model.TodoListItem, error) {
	return s.itemRepository.GetById(userId, itemId)
}

func (s *TodoListItemService) Update(userId, itemId uint64, data model.UpdateTodoListItemData) error {
	if err := data.Validate(); err != nil {
		return err
	}
	return s.itemRepository.Update(userId, itemId, data)
}

func (s *TodoListItemService) Delete(userId, itemId uint64) error {
	return s.itemRepository.Delete(userId, itemId)
}
