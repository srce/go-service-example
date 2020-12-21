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

// type DB interface {
// 	PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)
// 	NamedQueryContext(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error)
// 	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
// 	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
// 	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
// 	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
// 	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
// 	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
// 	MustBeginTx(ctx context.Context, opts *sql.TxOptions) *sqlx.Tx
// 	MustExecContext(ctx context.Context, query string, args ...interface{}) sql.Result
// 	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
// }

// type Tx interface {
// 	StmtxContext(ctx context.Context, stmt interface{}) *sqlx.Stmt
// 	NamedStmtContext(ctx context.Context, stmt *sqlx.NamedStmt) *sqlx.NamedStmt
// 	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
// 	PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)
// 	MustExecContext(ctx context.Context, query string, args ...interface{}) sql.Result
// 	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
// 	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
// 	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
// 	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
// 	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
// }

// type SMTP interface {
// 	SelectContext(ctx context.Context, dest interface{}, args ...interface{}) error
// 	GetContext(ctx context.Context, dest interface{}, args ...interface{}) error
// 	MustExecContext(ctx context.Context, args ...interface{}) sql.Result
// 	QueryRowxContext(ctx context.Context, args ...interface{}) *sqlx.Row
// 	QueryxContext(ctx context.Context, args ...interface{}) (*sqlx.Rows, error)
// }

// type qStmt interface {
// 	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
// 	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
// 	QueryRowxContext(ctx context.Context, query string, args ...interface{})
// 	ExecContext(ctx context.Context, query string, args ...interface{})
// }
