package store

type Authorization interface {
}

type TodoList interface {
}

type TodoItem interface {
}

type Store struct {
	Authorization
	TodoList
	TodoItem
}

func NewStore() *Store {
	return &Store{}
}
