package postgres

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tmrrwnxtsn/todo-lists-api/internal/model"
	"github.com/tmrrwnxtsn/todo-lists-api/internal/store"
	"strings"
)

var _ store.TodoListRepository = (*TodoListRepository)(nil)

type TodoListRepository struct {
	store *Store
}

func NewTodoListRepository(store *Store) *TodoListRepository {
	return &TodoListRepository{store: store}
}

func (r *TodoListRepository) Create(userId uint64, list model.TodoList) (uint64, error) {
	tx, err := r.store.db.Begin()
	if err != nil {
		return 0, err
	}

	var todoListId uint64
	createTodoListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoListsTable)
	if err = tx.QueryRow(createTodoListQuery, list.Title, list.Description).Scan(&todoListId); err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	createUsersListsQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2)", usersListsTable)
	_, err = tx.Exec(createUsersListsQuery, userId, todoListId)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	return todoListId, tx.Commit()
}

func (r *TodoListRepository) GetAll(userId uint64) ([]model.TodoList, error) {
	query := fmt.Sprintf(`SELECT tl.id, tl.title, tl.description 
								FROM %s tl 
								INNER JOIN %s utl 
								ON tl.id = utl.list_id 
								WHERE utl.user_id = $1`,
		todoListsTable, usersListsTable)

	var todoLists []model.TodoList
	err := r.store.db.Select(&todoLists, query, userId)
	return todoLists, err
}

func (r *TodoListRepository) GetById(userId, todoListId uint64) (model.TodoList, error) {
	query := fmt.Sprintf(`SELECT tl.id, tl.title, tl.description 
								FROM %s tl 
								INNER JOIN %s utl 
								ON tl.id = utl.list_id 
								WHERE utl.user_id = $1 AND utl.list_id = $2`,
		todoListsTable, usersListsTable)

	var todoList model.TodoList
	err := r.store.db.Get(&todoList, query, userId, todoListId)
	return todoList, err
}

func (r *TodoListRepository) Update(userId, todoListId uint64, data model.UpdateTodoListData) error {
	setValues := make([]string, 0, 2)
	args := make([]interface{}, 0, 2)
	argId := 1

	if data.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *data.Title)
		argId++
	}

	if data.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *data.Description)
		argId++
	}

	// 1) title=$1
	// 2) description=$1
	// 3) title=$1, description=$2
	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE %s tl
								SET %s
								FROM %s utl 
								WHERE tl.id = utl.list_id AND utl.user_id = $%d AND utl.list_id = $%d`,
		todoListsTable, setQuery, usersListsTable, argId, argId+1)

	args = append(args, userId, todoListId)

	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %s", args)

	_, err := r.store.db.Exec(query, args...)
	return err
}

func (r *TodoListRepository) Delete(userId, todoListId uint64) error {
	query := fmt.Sprintf(`DELETE FROM %s tl
								USING %s utl 
								WHERE tl.id = utl.list_id AND utl.user_id = $1 AND utl.list_id = $2`,
		todoListsTable, usersListsTable)

	_, err := r.store.db.Exec(query, userId, todoListId)
	return err
}
