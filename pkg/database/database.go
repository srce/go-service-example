package database

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	// Register postgres driver.
	_ "github.com/lib/pq"
)

const (
	DefaultMigrationTableName string = "schema_migrations"
	DefaultSchemaName         string = "public"
)

type Ext interface {
	sqlx.ExtContext
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type Database interface {
	Write() Ext
	Read() Ext
	Open() error
	Close() error
}

type Boot struct {
	db Database
}

func NewBoot(db Database) Boot {
	return Boot{db: db}
}

func (b Boot) Name() string {
	return "database"
}

func (b Boot) Init(ctx context.Context) error {
	if err := b.db.Open(); err != nil {
		return fmt.Errorf("database init: %w", err)
	}

	return nil
}

func (b Boot) Stop() error {
	if err := b.db.Close(); err != nil {
		return fmt.Errorf("database stopping: %w", err)
	}

	return nil
}
