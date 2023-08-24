package grpc

import (
	"github.com/FluffyKebab/todo/app/auth"
	"github.com/FluffyKebab/todo/app/log"
	"github.com/FluffyKebab/todo/domain/todo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	*todoServer
	*userServer
}

func NewServer(
	logger log.Logger,
	auth auth.Authenticator,
	todoService todo.TodoService,
	userService todo.UserService,
) *Server {
	return &Server{
		todoServer: &todoServer{
			todoService: todoService,
			auth:        auth,
			logger:      logger,
		},
		userServer: &userServer{
			userService: userService,
			auth:        auth,
			logger:      logger,
		},
	}
}

func authenticate(auth auth.Authenticator, token string, userId string) error {
	hasAccess, err := auth.HasAccess(token, userId)
	if err != nil {
		return status.Error(codes.Unauthenticated, err.Error())
	}
	if !hasAccess {
		return status.Error(codes.PermissionDenied, "permission denied")
	}

	return nil
}
