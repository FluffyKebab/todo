package postgres

import (
	"context"

	"github.com/FluffyKebab/todo/domain/todo"
	"github.com/google/uuid"
)

func (s *Storer) CreateUser(ctx context.Context, u todo.User) (string, error) {
	id := uuid.NewString()
	_, err := s.db.ExecContext(ctx, "INSERT INTO users (id, name) VALUES ($1, $2)", id, u.Name)
	return id, err
}
