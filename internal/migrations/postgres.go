package migrations

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
)

func NewPostgres(dbc *sqlx.DB) (database.Driver, error) {
	driver, err := postgres.WithInstance(dbc.DB, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("postgres instance: %w", err)
	}
	return driver, nil
}
