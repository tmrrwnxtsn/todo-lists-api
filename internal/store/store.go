package store

type Store interface {
	Authorization
	TodoList
	TodoItem
}
