package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/dzyanis/go-service-example/internal/funds"
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
var list = []Migrator{
	Migration{
		NameString: "202012151_create_users",
		ApplyFunc: func(tx *sql.Tx) error {
			query := `
				CREATE TABLE IF NOT EXISTS users (
					id BIGSERIAL PRIMARY KEY
					, name TEXT NOT NULL
					, email TEXT NOT NULL
					, deleted BOOLEAN  NOT NULL
					, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
					, updated_at TIMESTAMPTZ
				);`

			if _, err := tx.Exec(query); err != nil {
				return fmt.Errorf("apply: %w", err)
			}
			return nil
		},
		RollbackFunc: func(tx *sql.Tx) error {
			query := `DROP TABLE IF EXISTS users;`
			if _, err := tx.Exec(query); err != nil {
				return fmt.Errorf("rollback: %w", err)
			}
			return nil
		},
	},

	Migration{
		NameString: "202012152_create_company_beneficiary",
		ApplyFunc: func(tx *sql.Tx) error {
			query := `
				INSERT INTO users
					(name, email, deleted, created_at, updated_at)
				VALUES
					($1, $2, $3, $4, $5);
			`

			now := time.Now()
			_, err := tx.Exec(query, "Company Beneficiary", funds.CompanyBeneficiaryEmail, false, now, now)
			if err != nil {
				return fmt.Errorf("apply: %w", err)
			}
			return nil
		},
		RollbackFunc: func(tx *sql.Tx) error {
			query := `DELETE FROM users WHERE email = $1;`
			if _, err := tx.Exec(query, funds.CompanyBeneficiaryEmail); err != nil {
				return fmt.Errorf("rollback: %w", err)
			}
			return nil
		},
	},
}

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
