package store

// Store ...
type Store interface {
	// User ...
	User() UserRepository
	// TodoList ...
	TodoList() TodoListRepository
	TodoItem
}
