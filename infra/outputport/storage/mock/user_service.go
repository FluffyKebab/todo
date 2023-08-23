package mock

import (
	"context"

	"github.com/FluffyKebab/todo/domain/todo"
)

type UserService struct {
	CreateUserFunc    func(ctx context.Context, u todo.User) (string, error)
	CreateUserInvoked bool
}

var _ todo.UserService = UserService{}

func (s UserService) CreateUser(ctx context.Context, u todo.User) (string, error) {
	s.CreateUserInvoked = true
	return s.CreateUserFunc(ctx, u)
}
