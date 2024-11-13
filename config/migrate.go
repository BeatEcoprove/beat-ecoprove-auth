package config

import (
	"database/sql"
	"errors"
	"path"
	"path/filepath"
	"runtime"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

const MigrationDIr = "migrations"

var (
	ErrFileNotRunning = errors.New("file is not running")
	ErrConnectToDb    = errors.New("failed to connect to the database")
	ErrFailToMigrate  = errors.New("fail to apply migrations")
)

func GetProjectRoot() (string, error) {
	_, filename, _, ok := runtime.Caller(0)

	if !ok {
		return "", ErrFileNotRunning
	}

	absPath, err := filepath.Abs(filename)

	if err != nil {
		return "", err
	}

	return path.Dir(path.Dir(absPath)), nil
}

func getMigrationDir() (string, error) {
	projectPath, err := GetProjectRoot()

	if err != nil {
		return "", err
	}

	return path.Join(projectPath, MigrationDIr), nil
}

func Migrate(connectionString string, verbose bool) error {
	migrationDIr, err := getMigrationDir()

	if err != nil {
		return err
	}

	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		return ErrConnectToDb
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		return ErrConnectToDb
	}

	goose.SetVerbose(verbose)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db, migrationDIr); err != nil {
		return ErrFailToMigrate
	}

	return nil
}
