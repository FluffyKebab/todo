package mongo

import (
	"context"
	"errors"

	"github.com/FluffyKebab/todo/domain/todo"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Storer) CreateTodo(ctx context.Context, t todo.Todo) (string, error) {
	t.ID = uuid.NewString()
	_, err := s.c.Database("todo").Collection("todos").InsertOne(ctx, t)
	return t.ID, err
}

func (s *Storer) GetTodo(ctx context.Context, id string) (todo.Todo, error) {
	var t todo.Todo
	res := s.c.Database("todo").Collection("todos").FindOne(ctx, bson.M{"id": id})
	err := res.Decode(&t)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return todo.Todo{}, todo.ErrTodoNotFound
		}

		return todo.Todo{}, err
	}

	return t, err
}

func (s *Storer) UpdateTodo(ctx context.Context, t todo.UpdateTodoRequest) error {
	coll := s.c.Database("todo").Collection("todos")

	if t.ShouldUpdateBody {
		updateResult, err := coll.UpdateOne(
			ctx,
			bson.M{"id": t.NewTodo.ID},
			bson.D{{Key: "$set", Value: bson.D{{Key: "body", Value: t.NewTodo.Body}}}},
		)
		if err != nil {
			return err
		}

		if updateResult.MatchedCount != 1 {
			return todo.ErrTodoNotFound
		}
	}

	if t.ShouldUpdateDone {
		updateResult, err := coll.UpdateOne(
			ctx,
			bson.M{"id": t.NewTodo.ID},
			bson.D{{Key: "$set", Value: bson.D{{Key: "done", Value: t.NewTodo.Done}}}},
		)
		if err != nil {
			return err
		}

		if updateResult.MatchedCount != 1 {
			return todo.ErrTodoNotFound
		}
	}

	return nil
}

func (s *Storer) DeleteTodo(ctx context.Context, id string) error {
	res, err := s.c.Database("todo").Collection("todos").DeleteOne(
		ctx,
		bson.M{"id": id},
	)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return todo.ErrTodoNotFound
	}

	return err
}

func (s *Storer) GetUserTodos(ctx context.Context, userId string) ([]todo.Todo, error) {
	cur, err := s.c.Database("todo").Collection("todos").Find(ctx, bson.M{"userId": userId})
	if err != nil {
		return nil, err
	}

	todos := make([]todo.Todo, 0)
	for cur.Next(ctx) {
		var todo todo.Todo
		if err := cur.Decode(&todo); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}
