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

	var listId uint64
	createListQuery := fmt.Sprintf(`INSERT INTO %s 
										  (title, description) 
										  VALUES ($1, $2) 
										  RETURNING id`,
		todoListsTable)
	if err = tx.QueryRow(createListQuery, list.Title, list.Description).Scan(&listId); err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	createUsersListsQuery := fmt.Sprintf(`INSERT INTO %s 
												(user_id, list_id) 
												VALUES ($1, $2)`,
		usersListsTable)
	_, err = tx.Exec(createUsersListsQuery, userId, listId)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	return listId, tx.Commit()
}

func (r *TodoListRepository) GetAll(userId uint64) ([]model.TodoList, error) {
	getAllListsQuery := fmt.Sprintf(`SELECT tl.id, tl.title, tl.description
										   FROM %s tl 
										   INNER JOIN %s utl 
										   ON tl.id = utl.list_id 
										   WHERE utl.user_id = $1`,
		todoListsTable, usersListsTable)

	var lists []model.TodoList
	err := r.store.db.Select(&lists, getAllListsQuery, userId)
	if err != nil {
		return nil, err
	}
	return lists, nil
}

func (r *TodoListRepository) GetById(userId, listId uint64) (model.TodoList, error) {
	getListByIdQuery := fmt.Sprintf(`SELECT tl.id, tl.title, tl.description 
										   FROM %s tl 
										   INNER JOIN %s utl 
										   ON tl.id = utl.list_id 
										   WHERE utl.user_id = $1 AND utl.list_id = $2`,
		todoListsTable, usersListsTable)

	var list model.TodoList
	err := r.store.db.Get(&list, getListByIdQuery, userId, listId)
	if err != nil {
		return model.TodoList{}, err
	}
	return list, nil
}

func (r *TodoListRepository) Update(userId, listId uint64, data model.UpdateTodoListData) error {
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

	setQuery := strings.Join(setValues, ", ")

	updateListQuery := fmt.Sprintf(`UPDATE %s tl
										  SET %s
										  FROM %s utl 
                                          WHERE tl.id = utl.list_id AND utl.user_id = $%d AND utl.list_id = $%d`,
		todoListsTable, setQuery, usersListsTable, argId, argId+1)

	args = append(args, userId, listId)

	logrus.Debugf("updateQuery: %s", updateListQuery)
	logrus.Debugf("args: %s", args)

	_, err := r.store.db.Exec(updateListQuery, args...)
	return err
}

func (r *TodoListRepository) Delete(userId, listId uint64) error {
	deleteListQuery := fmt.Sprintf(`DELETE FROM %s tl
										  USING %s utl 
										  WHERE tl.id = utl.list_id AND utl.user_id = $1 AND utl.list_id = $2`,
		todoListsTable, usersListsTable)

	_, err := r.store.db.Exec(deleteListQuery, userId, listId)
	return err
}
