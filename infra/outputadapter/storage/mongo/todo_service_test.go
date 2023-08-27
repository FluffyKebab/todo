package mongo

import (
	"context"
	"testing"

	"github.com/FluffyKebab/todo/domain/todo"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestTodoService(t *testing.T) {
	s, err := connectToDatabase()
	require.NoError(t, err)
	userId := uuid.NewString()

	todosCreated := []todo.Todo{
		{Body: "foo", Done: false, UserID: userId},
		{Body: "baz", Done: false, UserID: userId},
		{Body: "bal", Done: true, UserID: userId},
	}

	todosCreated[0].ID, err = s.CreateTodo(context.Background(), todosCreated[0])
	require.NoError(t, err)
	todosCreated[1].ID, err = s.CreateTodo(context.Background(), todosCreated[1])
	require.NoError(t, err)
	todosCreated[2].ID, err = s.CreateTodo(context.Background(), todosCreated[2])
	require.NoError(t, err)

	todosInDatabase, err := s.GetUserTodos(context.Background(), userId)
	require.NoError(t, err)
	require.Equal(t, todosCreated, todosInDatabase, "todos created not the same as gotten from database")

	updatedTodo := todo.Todo{ID: todosCreated[0].ID, UserID: userId, Body: "zoo", Done: true}
	err = s.UpdateTodo(context.Background(), todo.UpdateTodoRequest{
		NewTodo:          updatedTodo,
		ShouldUpdateDone: true,
		ShouldUpdateBody: true,
	})
	require.NoError(t, err)

	updatedTodoFromDB, err := s.GetTodo(context.Background(), updatedTodo.ID)
	require.NoError(t, err)
	require.Equal(t, updatedTodo, updatedTodoFromDB, "update from database wrong")

	err = s.DeleteTodo(context.Background(), todosCreated[2].ID)
	require.NoError(t, err)

	todosInDatabase, err = s.GetUserTodos(context.Background(), userId)
	require.NoError(t, err)
	require.Contains(t, todosInDatabase, updatedTodo, "wrong todos in database")
	require.Contains(t, todosInDatabase, todosCreated[1], "wrong todos in database")
}
