package postgres

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tmrrwnxtsn/todo-lists-api/internal/model"
	"github.com/tmrrwnxtsn/todo-lists-api/internal/store"
	"strings"
)

var _ store.TodoListItemRepository = (*TodoListItemRepository)(nil)

type TodoListItemRepository struct {
	store *Store
}

func NewTodoListItemRepository(store *Store) *TodoListItemRepository {
	return &TodoListItemRepository{store: store}
}

func (r *TodoListItemRepository) Create(listId uint64, item model.TodoListItem) (uint64, error) {
	tx, err := r.store.db.Begin()
	if err != nil {
		return 0, err
	}

	var itemId uint64
	createItemQuery := fmt.Sprintf(`INSERT INTO %s 
										  (title, description) 
										  VALUES ($1, $2) 
										  RETURNING id`,
		todoItemsTable)
	if err = tx.QueryRow(createItemQuery, item.Title, item.Description).Scan(&itemId); err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	createItemsListsQuery := fmt.Sprintf(`INSERT INTO %s 
												(list_id, item_id) 
												VALUES ($1, $2)`,
		listsItemsTable)
	_, err = tx.Exec(createItemsListsQuery, listId, itemId)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	return itemId, tx.Commit()
}

func (r *TodoListItemRepository) GetAll(userId, listId uint64) ([]model.TodoListItem, error) {
	getAllItemsQuery := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done 
										   FROM %s ti 
										   INNER JOIN %s tlti 
										   ON ti.id = tlti.item_id
										   INNER JOIN %s utl 
										   ON tlti.list_id = utl.list_id 
										   WHERE tlti.list_id = $1 AND utl.user_id = $2`,
		todoItemsTable, listsItemsTable, usersListsTable)

	var items []model.TodoListItem
	err := r.store.db.Select(&items, getAllItemsQuery, listId, userId)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *TodoListItemRepository) GetById(userId, itemId uint64) (model.TodoListItem, error) {
	getItemByIdQuery := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done 
										   FROM %s ti 
										   INNER JOIN %s tlti 
										   ON ti.id = tlti.item_id
										   INNER JOIN %s utl 
										   ON tlti.list_id = utl.list_id 
										   WHERE utl.user_id = $1 AND ti.id = $2`,
		todoItemsTable, listsItemsTable, usersListsTable)

	var item model.TodoListItem
	err := r.store.db.Get(&item, getItemByIdQuery, userId, itemId)
	if err != nil {
		return model.TodoListItem{}, err
	}
	return item, nil
}

func (r *TodoListItemRepository) Update(userId, itemId uint64, data model.UpdateTodoListItemData) error {
	setValues := make([]string, 0, 3)
	args := make([]interface{}, 0, 3)
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

	if data.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done=$%d", argId))
		args = append(args, *data.Done)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	updateItemQuery := fmt.Sprintf(`UPDATE %s ti
										  SET %s
										  FROM %s tlti, %s utl
										  WHERE ti.id = tlti.item_id AND tlti.list_id = utl.list_id 
										  AND utl.user_id = $%d AND ti.id = $%d`,
		todoItemsTable, setQuery, listsItemsTable, usersListsTable, argId, argId+1)

	args = append(args, userId, itemId)

	logrus.Debugf("updateQuery: %s", updateItemQuery)
	logrus.Debugf("args: %s", args)

	_, err := r.store.db.Exec(updateItemQuery, args...)
	return err
}

func (r *TodoListItemRepository) Delete(userId, itemId uint64) error {
	deleteItemQuery := fmt.Sprintf(`DELETE FROM %s ti
										  USING %s tlti, %s utl
										  WHERE ti.id = tlti.item_id AND tlti.list_id = utl.list_id 
										  AND utl.user_id = $1 AND ti.id = $2`,
		todoItemsTable, listsItemsTable, usersListsTable)

	_, err := r.store.db.Exec(deleteItemQuery, userId, itemId)
	return err
}
