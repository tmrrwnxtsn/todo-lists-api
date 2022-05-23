package model

import "errors"

type TodoList struct {
	Id          int    `json:"id"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

type UsersList struct {
	Id     int
	UserId int
	ListId int
}

type TodoListItem struct {
	Id          int    `json:"id"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

type ListsItem struct {
	Id     int
	ListId int
	ItemId int
}

type UpdateTodoListData struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

func (d UpdateTodoListData) Validate() error {
	if d.Title == nil && d.Description == nil {
		return errors.New("update todo list request has no values")
	}
	return nil
}

type UpdateTodoListItemData struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Done        *bool   `json:"done"`
}

func (d UpdateTodoListItemData) Validate() error {
	if d.Title == nil && d.Description == nil && d.Done == nil {
		return errors.New("update todo list item request has no values")
	}
	return nil
}
