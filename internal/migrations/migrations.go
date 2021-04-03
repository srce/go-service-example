package migrations

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"

	// driver for reading from filesystem.
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/dzyanis/go-service-example/pkg/logger"
)

type Boot struct {
	log    *logger.Logger
	driver database.Driver
}

func NewBoot(log *logger.Logger, driver database.Driver) Boot {
	return Boot{
		log:    log,
		driver: driver,
	}
}

func (b Boot) Name() string {
	return "migrations"
}

func (b Boot) Init(ctx context.Context) error {
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", b.driver)
	if err != nil {
		return fmt.Errorf("db instance: %w", err)
	}
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			b.log.Infof("migration: %v", err)
		} else {
			return fmt.Errorf("up: %w", err)
		}
	}
	return nil
}

func (b Boot) Stop() error {
	return nil
}
