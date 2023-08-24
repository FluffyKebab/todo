package mock

import (
	"context"

	"github.com/FluffyKebab/todo/domain/todo"
)

type TodoService struct {
	CreateTodoFunc   func(ctx context.Context, todo todo.Todo) (string, error)
	UpdateTodoFunc   func(ctx context.Context, req todo.UpdateTodoRequest) error
	DeleteTodoFunc   func(ctx context.Context, id string) error
	GetUserTodosFunc func(ctx context.Context, userId string) ([]todo.Todo, error)

	CreateTodoInvoked   bool
	UpdateTodoInvoked   bool
	DeleteTodoInvoked   bool
	GetUserTodosInvoked bool
}

var _ todo.TodoService = TodoService{}

func (m TodoService) CreateTodo(ctx context.Context, todo todo.Todo) (string, error) {
	m.CreateTodoInvoked = true
	return m.CreateTodoFunc(ctx, todo)
}

func (m TodoService) UpdateTodo(ctx context.Context, req todo.UpdateTodoRequest) error {
	m.UpdateTodoInvoked = true
	return m.UpdateTodoFunc(ctx, req)
}

func (m TodoService) DeleteTodo(ctx context.Context, id string) error {
	m.DeleteTodoInvoked = true
	return m.DeleteTodoFunc(ctx, id)
}

func (m TodoService) GetUserTodos(ctx context.Context, userId string) ([]todo.Todo, error) {
	m.GetUserTodosInvoked = true
	return m.GetUserTodosFunc(ctx, userId)
}
