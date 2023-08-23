package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/FluffyKebab/todo/app/auth"
	"github.com/FluffyKebab/todo/app/log"
	"github.com/FluffyKebab/todo/domain/todo"
	"github.com/FluffyKebab/todo/infra/inputport/grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userServer struct {
	pb.UnimplementedUserServiceServer
	userService todo.UserService
	logger      log.Logger
	auth        auth.Authenticator
}

func (s userServer) run(port string) error {
	lister, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, s)
	return server.Serve(lister)
}

func (s userServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	userId, err := s.userService.CreateUser(ctx, todo.User{Name: req.Name})
	if err != nil {
		s.logger.Error(fmt.Errorf("creating user: %w", err))
		return nil, status.Error(codes.Internal, "error creating user")
	}

	token, err := s.auth.CreateToken(userId)
	if err != nil {
		s.logger.Error(err)
		return nil, status.Error(codes.Internal, "error creating token")
	}

	return &pb.CreateUserResponse{
		UserID: userId,
		Token:  token,
	}, nil
}
