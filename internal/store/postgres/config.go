package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

const (
	usersTable      = "users"
	todoListsTable  = "todo_lists"
	usersListsTable = "users_lists"
	todoItemsTable  = "todo_items"
	listsItemsTable = "lists_items"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewDB(cfg Config) (*sqlx.DB, error) {
	return sqlx.Connect("postgres",
		fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
}
