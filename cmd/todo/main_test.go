package main

import (
	"context"
	"testing"

	"github.com/FluffyKebab/todo/domain/todo"
	"github.com/FluffyKebab/todo/infra/inputport/grpc/pb"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestServer(t *testing.T) {
	cc, err := grpc.Dial("localhost:9090",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	defer cc.Close()

	userClient := pb.NewUserServiceClient(cc)
	user, err := userClient.CreateUser(context.Background(), &pb.CreateUserRequest{Name: "bob"})
	require.NoError(t, err)

	todoClient := pb.NewTodoServiceClient(cc)
	createTodoRes1, err := todoClient.CreateTodo(context.Background(), &pb.CreateTodoRequest{
		UserID: user.UserID, Token: user.Token, Body: "do dishes", Done: false,
	})
	require.NoError(t, err)

	createTodoRes2, err := todoClient.CreateTodo(context.Background(), &pb.CreateTodoRequest{
		UserID: user.UserID, Token: user.Token, Body: "write text", Done: false,
	})
	require.NoError(t, err)

	getTodosResponse, err := todoClient.ListUserTodos(context.Background(), &pb.ListUserTodosRequest{
		UserID: user.UserID,
		Token:  user.Token,
	})

	expectedTodosResponse := []todo.Todo{
		{ID: createTodoRes1.Id, UserID: user.UserID, Body: "do dishes", Done: false},
		{ID: createTodoRes2.Id, UserID: user.UserID, Body: "write text", Done: false},
	}

	require.Len(t, getTodosResponse.Todos, 2, "wrong number of todos gotten")
	require.Contains(t, expectedTodosResponse, pbTodoToTodoTodo(getTodosResponse.Todos[0]))
	require.Contains(t, expectedTodosResponse, pbTodoToTodoTodo(getTodosResponse.Todos[1]))
}

func pbTodoToTodoTodo(t *pb.Todo) todo.Todo {
	return todo.Todo{
		ID:     t.Id,
		UserID: t.UserID,
		Body:   t.Body,
		Done:   t.Done,
	}
}
