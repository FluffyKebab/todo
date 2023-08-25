package postgres

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/FluffyKebab/todo/domain/todo"
	"github.com/google/uuid"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"
)

var databaseUrl string

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "11",
		Env: []string{
			"POSTGRES_PASSWORD=password",
			"POSTGRES_USER=postgres",
			"POSTGRES_DB=dbname",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	resource.Expire(120)

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl = fmt.Sprintf("postgres://postgres:password@%s/dbname?sslmode=disable", hostAndPort)

	log.Println("Connecting to database...")
	if err := pool.Retry(func() error {
		s, err := connectToDatabase()
		if err != nil {
			return err
		}
		err = s.db.Ping()
		if err != nil {
			return err
		}

		err = s.Migrate()
		if err != nil {
			log.Fatalf("migrations failed: %s", err.Error())
		}
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err.Error())
	}
	log.Println("Success")

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err.Error())
	}

	os.Exit(code)
}

func TestTodoService(t *testing.T) {
	s, err := connectToDatabase()
	require.NoError(t, err)
	userId := uuid.NewString()

	todosCreated := []todo.Todo{
		{Body: "foo", Done: false, UserID: userId},
		{Body: "baz", Done: false, UserID: userId},
		{Body: "bal", Done: true, UserID: userId},
	}

	todosCreated[0].ID, err = s.CreateTodo(context.Background(), todosCreated[0])
	require.NoError(t, err)
	todosCreated[1].ID, err = s.CreateTodo(context.Background(), todosCreated[1])
	require.NoError(t, err)
	todosCreated[2].ID, err = s.CreateTodo(context.Background(), todosCreated[2])
	require.NoError(t, err)

	todosInDatabase, err := s.GetUserTodos(context.Background(), userId)
	require.NoError(t, err)
	require.Equal(t, todosCreated, todosInDatabase, "todos created not the same as gotten from database")

	updatedTodo := todo.Todo{ID: todosCreated[0].ID, UserID: userId, Body: "zoo", Done: true}
	err = s.UpdateTodo(context.Background(), todo.UpdateTodoRequest{
		NewTodo:          updatedTodo,
		ShouldUpdateDone: true,
		ShouldUpdateBody: true,
	})
	require.NoError(t, err)

	updatedTodoFromDB, err := s.GetTodo(context.Background(), updatedTodo.ID)
	require.NoError(t, err)
	require.Equal(t, updatedTodo, updatedTodoFromDB, "update from database wrong")

	err = s.DeleteTodo(context.Background(), todosCreated[2].ID)
	require.NoError(t, err)

	todosInDatabase, err = s.GetUserTodos(context.Background(), userId)
	require.NoError(t, err)
	require.Contains(t, todosInDatabase, updatedTodo, "wrong todos in database")
	require.Contains(t, todosInDatabase, todosCreated[1], "wrong todos in database")
}

func connectToDatabase() (*Storer, error) {
	return New(databaseUrl)
}
