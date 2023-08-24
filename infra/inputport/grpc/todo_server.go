package grpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/FluffyKebab/todo/app/auth"
	"github.com/FluffyKebab/todo/app/log"
	"github.com/FluffyKebab/todo/domain/todo"
	"github.com/FluffyKebab/todo/infra/inputport/grpc/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type todoServer struct {
	pb.UnimplementedTodoServiceServer
	todoService todo.TodoService
	auth        auth.Authenticator
	logger      log.Logger
}

func (s *todoServer) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{
		Status: pb.HealthCheckResponse_SERVING,
	}, nil
}

func (s *todoServer) CreateTodo(ctx context.Context, req *pb.CreateTodoRequest) (*pb.CreateTodoResponse, error) {
	err := authenticate(s.auth, req.Token, req.UserID)
	if err != nil {
		return nil, err
	}

	todoId, err := s.todoService.CreateTodo(ctx, todo.Todo{
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

func (s *todoServer) UpdateTodoDone(ctx context.Context, req *pb.UpdateTodoDoneRequest) (*pb.EmptyResponse, error) {
	err := authenticate(s.auth, req.Token, req.UserID)
	if err != nil {
		return nil, err
	}

	err = s.todoService.UpdateTodo(ctx, todo.UpdateTodoRequest{
		NewTodo:          todo.Todo{Done: req.Done, ID: req.Id},
		ShouldUpdateDone: true,
	})
	if err != nil {
		s.logger.Error(fmt.Errorf("updating todo: %w", err))
		return nil, status.Error(codes.Internal, "error updating todo")
	}

	return &pb.EmptyResponse{}, nil
}

func (s *todoServer) UpdateTodoBody(ctx context.Context, req *pb.UpdateTodoBodyRequest) (*pb.EmptyResponse, error) {
	err := authenticate(s.auth, req.Token, req.UserID)
	if err != nil {
		return nil, err
	}

	err = s.todoService.UpdateTodo(ctx, todo.UpdateTodoRequest{
		NewTodo:          todo.Todo{Body: req.Body, ID: req.Id},
		ShouldUpdateBody: true,
	})
	if err != nil {
		s.logger.Error(fmt.Errorf("updating todo: %w", err))
		return nil, status.Error(codes.Internal, "error updating todo")
	}

	return &pb.EmptyResponse{}, nil
}

func (s *todoServer) DeleteTodo(ctx context.Context, req *pb.DeleteTodoRequest) (*pb.EmptyResponse, error) {
	err := authenticate(s.auth, req.Token, req.UserID)
	if err != nil {
		return nil, err
	}

	err = s.todoService.DeleteTodo(ctx, req.Id)
	if errors.Is(err, todo.ErrTodoNotFound) {
		s.logger.Warning("delete on error not found")
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if err != nil {
		s.logger.Error(fmt.Errorf("updating todo: %w", err))
		return nil, status.Error(codes.Internal, "error updating todo")
	}

	return &pb.EmptyResponse{}, nil
}

func (s *todoServer) ListUserTodos(ctx context.Context, req *pb.ListUserTodosRequest) (*pb.ListUserTodosResponse, error) {
	err := authenticate(s.auth, req.Token, req.UserID)
	if err != nil {
		return nil, err
	}

	todos, err := s.todoService.GetUserTodos(ctx, req.UserID)
	if err != nil {
		s.logger.Error(err)
		return nil, status.Error(codes.Internal, "error getting todos")
	}

	pbTodos := make([]*pb.Todo, 0, len(todos))
	for _, t := range todos {
		pbTodos = append(pbTodos, &pb.Todo{
			Id:     t.ID,
			UserID: t.UserID,
			Body:   t.Body,
			Done:   t.Done,
		})
	}

	return &pb.ListUserTodosResponse{Todos: pbTodos}, nil
}
