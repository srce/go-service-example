package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	// Register postgres driver.
	_ "github.com/lib/pq"
)

type Config struct {
	Host           string   `env:"POSTGRES_HOST"`
	Port           int      `env:"POSTGRES_PORT"`
	User           string   `env:"POSTGRES_USER"`
	Password       string   `env:"POSTGRES_PASSWORD"`
	DB             string   `env:"POSTGRES_DB"`
	Options        []string `env:"POSTGRES_OPTIONS"`
	MaxConnections int      `env:"POSTGRES_MAX_CONN"`
}

func (c Config) Addr() string {
	addr := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", c.User, c.Password, c.Host, c.Port, c.DB)
	if len(c.Options) > 0 {
		addr += "?" + strings.Join(c.Options, ",")
	}

	return addr
}

const (
	DefaultMigrationTableName string = "schema_migrations"
	DefaultSchemaName         string = "public"
)

type Database struct {
	cnf  Config
	conn *sqlx.DB
}

func NewDatabase(cnf Config) *Database {
	return &Database{
		cnf:  cnf,
		conn: &sqlx.DB{},
	}
}

// Write returns connection for writing.
func (db *Database) Write() *sqlx.DB {
	return db.conn
}

// Read returns connection for reading.
func (db *Database) Read() *sqlx.DB {
	return db.conn
}

// Open establishes connection to database or returns error.
func (db *Database) Open() error {
	conn, err := sqlx.Open("postgres", db.cnf.Addr())
	if err != nil {
		return fmt.Errorf("database opening connection: %w", err)
	}
	db.conn = conn

	if err := db.conn.Ping(); err != nil {
		return fmt.Errorf("database ping: %w", err)
	}

	db.conn.SetMaxIdleConns(db.cnf.MaxConnections)

	return nil
}

// Close closes current connection to database.
func (db *Database) Close() error {
	if err := db.conn.Close(); err != nil {
		return fmt.Errorf("database close: %w", err)
	}

	return nil
}

type Boot struct {
	db *Database
}

func NewBoot(db *Database) Boot {
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
