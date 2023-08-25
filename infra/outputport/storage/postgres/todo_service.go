package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/FluffyKebab/todo/domain/todo"
	"github.com/google/uuid"
)

func (s *Storer) CreateTodo(ctx context.Context, t todo.Todo) (string, error) {
	todoId := uuid.NewString()
	_, err := s.db.ExecContext(ctx, "INSERT INTO todos (id, userId, body, done) VALUES ($1, $2, $3, $4)",
		todoId,
		t.UserID,
		t.Body,
		t.Done,
	)
	return todoId, err
}

func (s *Storer) GetTodo(ctx context.Context, id string) (todo.Todo, error) {
	row := s.db.QueryRowContext(ctx, "SELECT * FROM todos WHERE id = $1", id)

	var t todo.Todo
	err := row.Scan(&t.ID, &t.UserID, &t.Body, &t.Done)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return todo.Todo{}, todo.ErrTodoNotFound
		}
		return todo.Todo{}, err
	}

	return t, err
}

func (s *Storer) UpdateTodo(ctx context.Context, t todo.UpdateTodoRequest) error {
	if t.ShouldUpdateBody {
		res, err := s.db.ExecContext(ctx, "UPDATE todos SET body = $1 WHERE id = $2", t.NewTodo.Body, t.NewTodo.ID)
		if err != nil {
			return err
		}
		numEffected, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if numEffected == 0 {
			return todo.ErrTodoNotFound
		}
	}

	if t.ShouldUpdateDone {
		res, err := s.db.ExecContext(ctx, "UPDATE todos SET done = $1 WHERE id = $2", t.NewTodo.Done, t.NewTodo.ID)
		if err != nil {
			return err
		}
		numEffected, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if numEffected == 0 {
			return todo.ErrTodoNotFound
		}
	}

	return nil
}

func (s *Storer) DeleteTodo(ctx context.Context, id string) error {
	res, err := s.db.ExecContext(ctx, "DELETE FROM todos WHERE id = $1", id)
	if err != nil {
		return err
	}

	numEffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if numEffected == 0 {
		return todo.ErrTodoNotFound
	}

	return err
}

func (s *Storer) GetUserTodos(ctx context.Context, userId string) ([]todo.Todo, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT * FROM todos WHERE userId = $1", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := make([]todo.Todo, 0)
	for rows.Next() {
		var todo todo.Todo
		err := rows.Scan(&todo.ID, &todo.UserID, &todo.Body, &todo.Done)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}
