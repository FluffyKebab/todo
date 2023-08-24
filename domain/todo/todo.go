package todo

import (
	"context"
	"errors"
)

type User struct {
	ID   string
	Name string
}

type Todo struct {
	ID     string
	UserID string
	Body   string
	Done   bool
}

type UserService interface {
	CreateUser(ctx context.Context, u User) (string, error)
}

type UpdateTodoRequest struct {
	NewTodo          Todo
	ShouldUpdateDone bool
	ShouldUpdateBody bool
}

type TodoService interface {
	CreateTodo(context.Context, Todo) (string, error)
	UpdateTodo(context.Context, UpdateTodoRequest) error
	DeleteTodo(ctx context.Context, id string) error
	GetUserTodos(ctx context.Context, userId string) ([]Todo, error)
}

var ErrTodoNotFound = errors.New("todo not found")
