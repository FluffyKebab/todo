package mongo

import (
	"context"

	"github.com/FluffyKebab/todo/domain/todo"
	"github.com/google/uuid"
)

func (s *Storer) CreateUser(ctx context.Context, u todo.User) (string, error) {
	coll := s.c.Database("todo").Collection("users")
	u.ID = uuid.NewString()
	_, err := coll.InsertOne(ctx, u)
	return u.ID, err
}
