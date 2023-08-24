package grpc

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/FluffyKebab/todo/domain/todo"
	"github.com/FluffyKebab/todo/infra/inputport/grpc/pb"
	"github.com/FluffyKebab/todo/infra/outputport/auth/testauth"
	"github.com/FluffyKebab/todo/infra/outputport/log/testlog"
	"github.com/FluffyKebab/todo/infra/outputport/storage/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

func bufDialer(lis *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(ctx context.Context, s string) (net.Conn, error) {
		return lis.Dial()
	}
}

func runTestServer(t *testing.T, s *Server) *bufconn.Listener {
	lis := bufconn.Listen(bufSize)
	grpcServer := grpc.NewServer()
	pb.RegisterTodoServiceServer(grpcServer, s.todoServer)
	pb.RegisterUserServiceServer(grpcServer, s.userServer)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("unexpected error: %s", err.Error())
		}
	}()

	return lis
}

func createUserClient(t *testing.T, lis *bufconn.Listener) pb.UserServiceClient {
	t.Helper()
	cc, err := grpc.Dial("bufnet",
		grpc.WithContextDialer(bufDialer(lis)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	return pb.NewUserServiceClient(cc)
}

func createTodoClient(t *testing.T, lis *bufconn.Listener) pb.TodoServiceClient {
	t.Helper()
	cc, err := grpc.Dial("bufnet",
		grpc.WithContextDialer(bufDialer(lis)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	return pb.NewTodoServiceClient(cc)
}

func TestCreateTodo(t *testing.T) {
	t.Parallel()
	s := NewServer(
		testlog.Logger{ErrorFunc: func(err error) {
			t.Fatalf("unexpected error: %s", err.Error())
		}},
		testauth.Authenticator{HasAccessReturn: true},
		mock.TodoService{
			CreateTodoFunc: func(_ context.Context, todo todo.Todo) (string, error) {
				require.Equal(t, "body", todo.Body)
				return "todo_id", nil
			},
		},
		nil,
	)

	lis := runTestServer(t, s)
	c := createTodoClient(t, lis)

	res, err := c.CreateTodo(context.Background(), &pb.CreateTodoRequest{Body: "body"})
	require.NoError(t, err)
	require.Equal(t, res.Id, "todo_id")
}

func TestUpdateTodoDone(t *testing.T) {
	t.Parallel()

	s := NewServer(
		testlog.Logger{ErrorFunc: func(err error) {
			t.Fatalf("unexpected error: %s", err.Error())
		}},
		testauth.Authenticator{HasAccessReturn: true},
		mock.TodoService{
			UpdateTodoFunc: func(_ context.Context, req todo.UpdateTodoRequest) error {
				require.Equal(t, "id", req.NewTodo.ID)
				require.True(t, req.NewTodo.Done)
				require.True(t, req.ShouldUpdateDone)
				require.False(t, req.ShouldUpdateBody)
				return nil
			},
		},
		nil,
	)

	lis := runTestServer(t, s)
	c := createTodoClient(t, lis)

	_, err := c.UpdateTodoDone(context.Background(), &pb.UpdateTodoDoneRequest{
		Done: true,
		Id:   "id",
	})
	require.NoError(t, err)
}

func TestUpdateTodoBody(t *testing.T) {
	t.Parallel()
	s := NewServer(
		testlog.Logger{ErrorFunc: func(err error) {
			t.Fatalf("unexpected error: %s", err.Error())
		}},
		testauth.Authenticator{HasAccessReturn: true},
		mock.TodoService{
			UpdateTodoFunc: func(_ context.Context, req todo.UpdateTodoRequest) error {
				require.Equal(t, "id", req.NewTodo.ID)
				require.Equal(t, "body", req.NewTodo.Body)
				require.False(t, req.ShouldUpdateDone)
				require.True(t, req.ShouldUpdateBody)
				return nil
			},
		},
		nil,
	)

	lis := runTestServer(t, s)
	c := createTodoClient(t, lis)

	_, err := c.UpdateTodoBody(context.Background(), &pb.UpdateTodoBodyRequest{
		Body: "body",
		Id:   "id",
	})
	require.NoError(t, err)
}

func TestDeleteTodo(t *testing.T) {
	t.Parallel()
	s := NewServer(
		testlog.Logger{ErrorFunc: func(err error) {
			t.Fatalf("unexpected error: %s", err.Error())
		}},
		testauth.Authenticator{HasAccessReturn: true},
		mock.TodoService{
			DeleteTodoFunc: func(_ context.Context, id string) error {
				require.Equal(t, "id", id)
				return nil
			},
		},
		nil,
	)

	lis := runTestServer(t, s)
	c := createTodoClient(t, lis)

	_, err := c.DeleteTodo(context.Background(), &pb.DeleteTodoRequest{
		Id: "id",
	})
	require.NoError(t, err)
}

func TestListUserTodos(t *testing.T) {
	t.Parallel()
	s := NewServer(
		testlog.Logger{ErrorFunc: func(err error) {
			t.Fatalf("unexpected error: %s", err.Error())
		}},
		testauth.Authenticator{HasAccessReturn: true},
		mock.TodoService{
			GetUserTodosFunc: func(_ context.Context, userId string) ([]todo.Todo, error) {
				require.Equal(t, "id", userId)
				return []todo.Todo{
					{Body: "foo"}, {Body: "baz"},
				}, nil
			},
		},
		nil,
	)

	lis := runTestServer(t, s)
	c := createTodoClient(t, lis)

	res, err := c.ListUserTodos(context.Background(), &pb.ListUserTodosRequest{
		UserID: "id",
	})
	require.NoError(t, err)
	require.Equal(t, []*pb.Todo{{Body: "foo"}, {Body: "baz"}}, res.Todos)
}
