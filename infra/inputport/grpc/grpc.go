package grpc

import (
	"context"
	"fmt"

	"github.com/FluffyKebab/todo/app/auth"
	"github.com/FluffyKebab/todo/app/log"
	"github.com/FluffyKebab/todo/domain/todo"
	"github.com/FluffyKebab/todo/infra/inputport/grpc/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.TodoServiceServer
	todoService todo.TodoService
	auth        auth.Authenticator
	logger      log.Logger
}

func NewServer(authenticator auth.Authenticator, todoService todo.TodoService) *Server {
	return &Server{
		todoService: todoService,
		auth:        authenticator,
	}
}

func (s *Server) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{
		Status: pb.HealthCheckResponse_SERVING,
	}, nil
}

func (s *Server) CreateTodo(ctx context.Context, req *pb.CreateTodoRequest) (*pb.CreateTodoResponse, error) {
	hasAccess, err := s.auth.HasAccess(req.Token, req.UserID)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	if !hasAccess {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	todoId, err := s.todoService.CreateTodo(todo.Todo{
		UserID: req.UserID,
		Body:   req.Body,
		Done:   req.Done,
	})
	if err != nil {
		s.logger.Error(fmt.Errorf("creating todo: %w", err))
		return nil, status.Error(codes.Internal, "error inserting todo in database")
	}

	return &pb.CreateTodoResponse{Id: todoId}, nil
}

/*  CreateTodo(context.Context, *CreateTodoRequest) (*CreateTodoResponse, error)
    UpdateTodoDone(context.Context, *UpdateTodoDoneRequest) (*Error, error)
    UpdateTodoBody(context.Context, *UpdateTodoBodyRequest) (*Error, error)
    DeleteTodo(context.Context, *DeleteTodoRequest) (*Error, error)
    ListUserTodos(context.Context, *ListUserTodosRequest) (*ListUserTodosResponse, error)
*/
