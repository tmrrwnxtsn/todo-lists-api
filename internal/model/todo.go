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

type TodoItem struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
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

func (r UpdateTodoListData) Validate() error {
	if r.Title == nil && r.Description == nil {
		return errors.New("update todo list request has no values")
	}
	return nil
}
