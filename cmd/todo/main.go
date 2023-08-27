package main

import (
	"fmt"
	"log"
	"os"

	"github.com/FluffyKebab/todo/infra/inputport"
	"github.com/FluffyKebab/todo/infra/outputadapter"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err.Error())
	}
}

func run() error {
	config, err := readEnvironment()
	if err != nil {
		return err
	}

	services, err := outputadapter.NewServices(config)
	if err != nil {
		return err
	}

	return inputport.New(services).Run()
}

func readEnvironment() (outputadapter.Config, error) {
	config := outputadapter.Config{
		LogLevel: outputadapter.LogLevelAll,
	}

	config.DatabasePassword = os.Getenv("DB_PASSWORD")
	config.DatabaseUser = os.Getenv("DB_USER")
	config.DatabasePort = os.Getenv("DB_PORT")
	config.DatabaseName = os.Getenv("DB_NAME")
	config.DatabaseHost = os.Getenv("DB_HOST")
	config.DatabaseType = outputadapter.DatabaseType(os.Getenv("DB_TYPE"))
	config.AuthTokenSecretKey = os.Getenv("AUTH_SECRET_KEY")
	config.ServerPort = os.Getenv("PORT")

	if config.DatabasePassword == "" || config.DatabaseUser == "" || config.DatabasePort == "" ||
		config.DatabaseName == "" || config.DatabaseHost == "" || config.AuthTokenSecretKey == "" ||
		config.ServerPort == "" {

		return outputadapter.Config{}, fmt.Errorf("missing env variable: %v", config)
	}

	return config, nil
}
