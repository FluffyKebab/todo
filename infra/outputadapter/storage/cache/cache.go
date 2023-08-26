package cache

import (
	"context"

	"github.com/FluffyKebab/todo/domain/todo"
)

type Cache struct {
	userService todo.UserService
	todoService todo.TodoService

	todos map[string]todo.Todo
}

var _ todo.UserService = Cache{}
var _ todo.TodoService = Cache{}

func New(userService todo.UserService, todoService todo.TodoService) Cache {
	return Cache{
		userService: userService,
		todoService: todoService,
		todos:       make(map[string]todo.Todo),
	}
}

func (c Cache) CreateUser(ctx context.Context, u todo.User) (string, error) {
	return c.userService.CreateUser(ctx, u)
}

func (c Cache) CreateTodo(ctx context.Context, t todo.Todo) (string, error) {
	id, err := c.todoService.CreateTodo(ctx, t)
	if err != nil {
		return "", err
	}

	t.ID = id
	c.todos[id] = t
	return id, nil
}

func (c Cache) GetTodo(ctx context.Context, id string) (todo.Todo, error) {
	if t, ok := c.todos[id]; ok {
		return t, nil
	}

	t, err := c.todoService.GetTodo(ctx, id)
	if err != nil {
		return todo.Todo{}, err
	}

	c.todos[id] = t
	return t, nil
}

func (c Cache) UpdateTodo(ctx context.Context, req todo.UpdateTodoRequest) error {
	err := c.todoService.UpdateTodo(ctx, req)
	if err != nil {
		return err
	}

	if cached, ok := c.todos[req.NewTodo.ID]; ok {
		if req.ShouldUpdateBody {
			cached.Body = req.NewTodo.Body
		}
		if req.ShouldUpdateDone {
			cached.Done = req.NewTodo.Done
		}
		c.todos[req.NewTodo.ID] = cached
	}

	return nil
}

func (c Cache) DeleteTodo(ctx context.Context, id string) error {
	return c.DeleteTodo(ctx, id)
}

func (c Cache) GetUserTodos(ctx context.Context, userId string) ([]todo.Todo, error) {
	return c.todoService.GetUserTodos(ctx, userId)
}
