package postgres

import (
	"database/sql"

	"github.com/FluffyKebab/todo/domain/todo"
)

type Storer struct {
	db *sql.DB
}

var _ todo.UserService = &Storer{}
var _ todo.TodoService = &Storer{}

func New(dsn string) (*Storer, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return &Storer{
		db: db,
	}, nil
}

func (s *Storer) Migrate() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT NOT NULL, 
			name TEXT NOT NULL,
		);`,
	)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS todos (
			id TEXT NOT NULL, 
			userID TEXT NOT NULL,
			body TEXT NOT NULL,
			done BOOLEAN NOT NULL,
		);`,
	)
	return err
}
