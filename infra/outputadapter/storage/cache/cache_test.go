package cache

import (
	"context"
	"testing"

	"github.com/FluffyKebab/todo/domain/todo"
	"github.com/FluffyKebab/todo/infra/outputadapter/storage/mock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	c := New(
		nil,
		mock.TodoService{
			CreateTodoFunc: func(_ context.Context, todo todo.Todo) (string, error) {
				return uuid.NewString(), nil
			},
			GetTodoFunc: func(_ context.Context, id string) (todo.Todo, error) {
				return todo.Todo{}, nil
			},
			UpdateTodoFunc: func(_ context.Context, _ todo.UpdateTodoRequest) error {
				return nil
			},
		},
	)

	id, err := c.CreateTodo(context.Background(), todo.Todo{Body: "read book", Done: false})
	require.NoError(t, err)

	todoGotten, err := c.GetTodo(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, todo.Todo{ID: id, Body: "read book", Done: false}, todoGotten)

	err = c.UpdateTodo(context.Background(), todo.UpdateTodoRequest{
		ShouldUpdateBody: true,
		NewTodo:          todo.Todo{ID: id, Body: "do dishes"},
	})
	require.NoError(t, err)

	todoGotten, err = c.GetTodo(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, todo.Todo{ID: id, Body: "do dishes", Done: false}, todoGotten)
}
