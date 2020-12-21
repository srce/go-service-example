package database

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
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

type Postgres struct {
	cnf  Config
	conn *sqlx.DB
}

func NewPostgres(cnf Config) *Postgres {
	return &Postgres{cnf: cnf}
}

func (p *Postgres) Connection() *sqlx.DB {
	return p.conn
}

// Write returns connection for writing.
func (p *Postgres) Write() Ext {
	return p.Connection()
}

// Read returns connection for reading.
func (p *Postgres) Read() Ext {
	return p.Connection()
}

// Open establishes connection to database or returns error.
func (p *Postgres) Open() error {
	conn, err := sqlx.Open("postgres", p.cnf.Addr())
	if err != nil {
		return fmt.Errorf("database opening connection: %w", err)
	}
	p.conn = conn

	if err := p.conn.Ping(); err != nil {
		return fmt.Errorf("database ping: %w", err)
	}

	p.conn.SetMaxIdleConns(p.cnf.MaxConnections)

	return nil
}

// Close closes current connection to database.
func (p *Postgres) Close() error {
	if err := p.conn.Close(); err != nil {
		return fmt.Errorf("database close: %w", err)
	}

	return nil
}
