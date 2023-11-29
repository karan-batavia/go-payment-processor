package testcontainers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type PostgresContainer struct {
	ctx        context.Context
	container  *postgres.PostgresContainer
	migrations string
	DSN        string
}

func NewPostgresContainer(ctx context.Context, migrationsPath string) (*PostgresContainer, error) {
	username := "postgres"
	password := "postgres"
	database := "test_db"

	container, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16.0-alpine"),
		postgres.WithUsername(username),
		postgres.WithPassword(password),
		postgres.WithDatabase(database),
		testcontainers.WithWaitStrategy(
			wait.
				ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second),
		),
	)
	if err != nil {
		return nil, err
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	port, err := container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("%s:%s@%s:%s/%s?sslmode=disable",
		username, password, host, port.Port(), database,
	)

	return &PostgresContainer{
		ctx:        ctx,
		container:  container,
		migrations: migrationsPath,
		DSN:        dsn,
	}, nil
}

func (c *PostgresContainer) TerminateContainer() error {
	return c.container.Terminate(c.ctx)
}

func (c *PostgresContainer) ClearDB() error {
	m, err := migrate.New("file://"+c.migrations, "postgres://"+c.DSN)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	if err := m.Up(); err != nil {
		return err
	}

	return nil
}
