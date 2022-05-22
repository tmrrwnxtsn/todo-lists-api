package postgres

import (
	"fmt"
	"github.com/tmrrwnxtsn/todo-lists-api/internal/model"
	"github.com/tmrrwnxtsn/todo-lists-api/internal/store"
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
