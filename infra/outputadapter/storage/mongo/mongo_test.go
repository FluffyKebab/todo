package mongo

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
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
		Repository: "mongo",
		Tag:        "latest",
		Env: []string{
			"MONGO_INITDB_ROOT_PASSWORD=password",
			"MONGO_INITDB_ROOT_USERNAME=user",
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

	hostAndPort := resource.GetHostPort("27017/tcp")
	databaseUrl = fmt.Sprintf("mongodb://user:password@%s", hostAndPort)

	log.Println("Connecting to database...")
	if err := pool.Retry(func() error {
		s, err := connectToDatabase()
		if err != nil {
			return err
		}
		return s.c.Ping(context.Background(), nil)
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

func connectToDatabase() (*Storer, error) {
	return New(databaseUrl)
}
