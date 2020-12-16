package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dzyanis/go-service-example/pkg/logger"
)

type Migrator interface {
	Apply(tx *sql.Tx) error
	Rollback(tx *sql.Tx) error

	Name() string
}

type Migration struct {
	NameString   string
	ApplyFunc    func(tx *sql.Tx) error
	RollbackFunc func(tx *sql.Tx) error
}

func (s Migration) Apply(tx *sql.Tx) error {
	return s.ApplyFunc(tx)
}

func (s Migration) Rollback(tx *sql.Tx) error {
	return s.RollbackFunc(tx)
}

func (s Migration) Name() string {
	return s.NameString
}

// Enumerates all migrations
var list = []Migrator{}

type Boot struct {
	log *logger.Logger
	sch *Schema
}

func NewBoot(log *logger.Logger, sch *Schema) Boot {
	return Boot{
		log: log,
		sch: sch,
	}
}

func (b Boot) Name() string {
	return "migrations"
}

func (b Boot) Init(ctx context.Context) error {
	if err := b.sch.Init(); err != nil {
		return fmt.Errorf("initialization: %w", err)
	}

	unapplied, err := b.sch.FindUnapplied(list)
	if err != nil {
		return fmt.Errorf("finding unapplied: %w", err)
	}

	n, err := b.sch.Apply(unapplied)
	if err != nil {
		return fmt.Errorf("applying: %w", err)
	}
	b.log.Printf("migrations applied: %d", n)
	return nil
}

func (b Boot) Stop() error {
	return nil
}
