package config

import (
	"context"
	"errors"
	"fmt"
	"path"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const DotEnv = ".env"

var (
	ErrKillContainer = errors.New("failed to terminate container")
)

func GetDotEnvPath() (string, error) {
	projectPath, err := GetProjectRoot()

	if err != nil {
		return "", err
	}

	return path.Join(projectPath, DotEnv), nil
}

func KillContainer(container testcontainers.Container) error {
	if err := testcontainers.TerminateContainer(container); err != nil {
		return ErrKillContainer
	}

	return nil
}

func SpawnPostgresqlContainer() (*postgres.PostgresContainer, string, error) {
	ctx := context.Background()

	dbName := "beat-identity"
	dbUser := "ecoprove"
	dbPassword := "ecoprove"

	postgresContainer, err := postgres.Run(ctx, "postgres:16.4-alpine3.20",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)

	if err != nil {
		return nil, "", err
	}

	connectionString, err := postgresContainer.ConnectionString(ctx)

	if err != nil {
		return nil, "", err
	}

	return postgresContainer, fmt.Sprintf("%ssslmode=disable", connectionString), nil
}
