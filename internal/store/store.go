package store

// Store ...
type Store interface {
	// Auth ...
	Auth() AuthRepository
	TodoList
	TodoItem
}
