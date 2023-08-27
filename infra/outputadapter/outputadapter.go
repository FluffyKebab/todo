package outputadapter

import (
	"fmt"

	"github.com/FluffyKebab/todo/app"
	"github.com/FluffyKebab/todo/app/log"
	"github.com/FluffyKebab/todo/domain/todo"
	"github.com/FluffyKebab/todo/infra/outputadapter/auth/jwt"
	"github.com/FluffyKebab/todo/infra/outputadapter/log/console"
	"github.com/FluffyKebab/todo/infra/outputadapter/storage/cache"
	"github.com/FluffyKebab/todo/infra/outputadapter/storage/mongo"
	"github.com/FluffyKebab/todo/infra/outputadapter/storage/postgres"
)

type DatabaseType string

const (
	DatabaseTypePostgres = "postgres"
	DatabaseTypeMongoDB  = "mongo"
)

type LogLevel string

const (
	// LogLevelNone is a log level for logging nothing.
	LogLevelNone = "NONE"
	// LogLevelLow is a log level for only logging internal errors.
	LogLevelLow = "LOW"
	// LogLevelMedium is a log level for logging only warnings and errors.
	LogLevelMedium = "MEDIUM"
	// LogLevelAll is a log level for logging all messages.
	LogLevelAll = "ALL"
)

type Config struct {
	LogLevel
	DatabaseType
	AuthTokenSecretKey string
	DatabaseName       string
	DatabaseUser       string
	DatabasePassword   string
	DatabasePort       string
	DatabaseHost       string
	ServerPort         string
}

func NewServices(c Config) (*app.App, error) {
	l := createLogger(c.LogLevel)
	userService, todoService, err := getDBServices(c, l)
	if err != nil {
		return nil, err
	}

	memoryCache := cache.New(userService, todoService)

	return &app.App{
		Logger:      l,
		Auth:        jwt.NewAuthenticator(c.AuthTokenSecretKey),
		UserService: memoryCache,
		TodoService: memoryCache,
	}, nil
}

func getDBServices(c Config, l log.Logger) (todo.UserService, todo.TodoService, error) {
	if c.DatabaseType == DatabaseTypeMongoDB {
		l.Info("connecting to database...")
		s, err := mongo.New(getMongoConnString(c))
		return s, s, err
	}

	if c.DatabaseType == DatabaseTypePostgres {
		l.Info("connecting to database...")
		s, err := postgres.New(getPostgresConnString(c))
		if err != nil {
			return nil, nil, err
		}

		l.Info("running database migrations...")
		err = s.Migrate()
		if err != nil {
			return nil, nil, fmt.Errorf("running migrations: %w", err)
		}

		return s, s, err
	}

	return nil, nil, fmt.Errorf("invalid database type %s", c.DatabaseType)
}

func getMongoConnString(c Config) string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%s", c.DatabaseUser, c.DatabasePassword, c.DatabaseHost, c.DatabasePort)
}

func getPostgresConnString(c Config) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DatabaseUser,
		c.DatabasePassword,
		c.DatabaseHost,
		c.DatabasePort,
		c.DatabaseName,
	)
}

func createLogger(level LogLevel) console.Logger {
	switch level {
	case LogLevelNone:
		return console.Logger{}
	case LogLevelLow:
		return console.Logger{LogErrors: true}
	case LogLevelMedium:
		return console.Logger{LogErrors: true, LogWarnings: true}
	case LogLevelAll:
		return console.Logger{LogErrors: true, LogWarnings: true, LogInfo: true}
	default:
		return console.Logger{LogErrors: true, LogWarnings: true, LogInfo: true}
	}
}
