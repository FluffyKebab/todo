package app

import (
	"github.com/FluffyKebab/todo/app/auth"
	"github.com/FluffyKebab/todo/app/log"
	"github.com/FluffyKebab/todo/domain/todo"
)

type App struct {
	Logger      log.Logger
	Auth        auth.Authenticator
	UserService todo.UserService
	TodoService todo.TodoService
}

func New(l log.Logger, a auth.Authenticator, u todo.UserService, t todo.TodoService) *App {
	return &App{
		Logger:      l,
		Auth:        a,
		UserService: u,
		TodoService: t,
	}
}
