package grpc

import (
	"github.com/FluffyKebab/todo/app/auth"
	"github.com/FluffyKebab/todo/app/log"
	"github.com/FluffyKebab/todo/domain/todo"
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
