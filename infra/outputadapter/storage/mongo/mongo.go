package mongo

import (
	"context"

	"github.com/FluffyKebab/todo/domain/todo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storer struct {
	c *mongo.Client
}

var _ todo.UserService = &Storer{}
var _ todo.TodoService = &Storer{}

func New(uri string) (*Storer, error) {
	client, err := mongo.Connect(
		context.TODO(),
		options.Client().ApplyURI(uri),
	)
	if err != nil {
		return nil, err
	}

	return &Storer{c: client}, nil
}
