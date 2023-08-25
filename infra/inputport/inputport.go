package inputport

import (
	"github.com/FluffyKebab/todo/app"
	"github.com/FluffyKebab/todo/app/log"
	"github.com/FluffyKebab/todo/infra/inputport/grpc"
)

type Inputport struct {
	grpcServer *grpc.Server
	l          log.Logger
}

func New(app *app.App) *Inputport {
	return &Inputport{
		grpcServer: grpc.NewServer(app.Logger, app.Auth, app.TodoService, app.UserService),
		l:          app.Logger,
	}
}

func (i *Inputport) Run() error {
	i.l.Info("running server...")
	return i.grpcServer.ListenAndServe(":9090")
}
