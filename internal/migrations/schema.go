package migrations

import (
	"fmt"
	"sort"
	"time"

	"github.com/jmoiron/sqlx"
)

// Schema is the single database's schema representation.
type Schema struct {
	dbc          *sqlx.DB
	schemaName   string
	migTableName string
}

func NewSchema(dbc *sqlx.DB, schemaName, migTableName string) *Schema {
	return &Schema{
		dbc:          dbc,
		schemaName:   schemaName,
		migTableName: migTableName,
	}
}

// Init creates a migrations table in the database.
func (s *Schema) Init() error {
	var err error
	q := `CREATE SCHEMA IF NOT EXISTS "` + s.schemaName + `"`
	_, err = s.dbc.Exec(q)
	if err != nil {
		return fmt.Errorf("creating schema: %w", err)
	}

	q = `CREATE TABLE IF NOT EXISTS "` + s.schemaName + `"` + `."` + s.migTableName + `" ` +
		`(name TEXT UNIQUE, applied_at TIMESTAMP)`
	_, err = s.dbc.Exec(q)
	if err != nil {
		return fmt.Errorf("creating table: %w", err)
	}
	return nil
}

// ErrNameNotUnique is returned whenever a non-unique migration name is found.
type ErrNameNotUnique struct {
	Name string
}

// Error implements the error interface for ErrNameNotUnique.
func (err ErrNameNotUnique) Error() string {
	return fmt.Sprintf("migration name not unique: %q", err.Name)
}

// ErrorPair is a pair of errors.
type ErrorPair struct {
	Err1, Err2 error
}

// Error implements the error interface for ErrorPair.
func (err ErrorPair) Error() string {
	return fmt.Sprintf("err1: %q, err2: %q", err.Err1, err.Err2)
}

var _ error = ErrNameNotUnique{}

func (s *Schema) FindUnapplied(migrations []Migrator) (res []Migrator, err error) {
	if len(migrations) == 0 {
		return nil, nil
	}

	migByName := map[string]Migrator{}
	for _, m := range migrations {
		if migByName[m.Name()] != nil {
			return nil, ErrNameNotUnique{Name: m.Name()}
		}
		migByName[m.Name()] = m
	}

	q := `SELECT name FROM "` + s.schemaName + `"` + `."` + s.migTableName + `"` +
		`ORDER BY name`

	rows, err := s.dbc.Query(q)
	if err != nil {
		return nil, fmt.Errorf("executing query: %w", err)
	}

	defer func() {
		closeErr := rows.Close()
		if closeErr != nil {
			if err != nil {
				err = ErrorPair{Err1: err, Err2: closeErr}
			} else {
				err = closeErr
			}
		}
	}()

	var resNames []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}

		resNames = append(resNames, name)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}

	for _, name := range resNames {
		delete(migByName, name)
	}

	for _, m := range migByName {
		res = append(res, m)
	}

	sort.Sort(migrationsByName(res))

	return res, nil
}

// Apply applies all migrations in a single transaction. It returns the number
// of applied migrations and error if any.
func (s *Schema) Apply(migrations []Migrator) (n int, err error) {
	tx, err := s.dbc.Begin()
	if err != nil {
		return 0, fmt.Errorf("creating tx: %w", err)
	}

	defer func() {
		if err == nil {
			err = tx.Commit()
		} else {
			rbErr := tx.Rollback()
			if rbErr != nil {
				err = ErrorPair{
					Err1: err,
					Err2: rbErr,
				}
			}
		}
	}()

	now := time.Now()
	q := `INSERT INTO "` + s.schemaName + `"` + `."` + s.migTableName + `" (name, applied_at) ` +
		`VALUES ($1, $2)`
	for _, m := range migrations {
		err = m.Apply(tx)
		if err != nil {
			return 0, fmt.Errorf("applying: %w", err)
		}

		_, err = tx.Exec(q, m.Name(), now)
		if err != nil {
			return 0, fmt.Errorf("executing: %w", err)
		}

		n++
	}

	return n, nil
}

type migrationsByName []Migrator

func (ms migrationsByName) Len() int           { return len(ms) }
func (ms migrationsByName) Less(i, j int) bool { return ms[i].Name() < ms[j].Name() }
func (ms migrationsByName) Swap(i, j int)      { ms[i], ms[j] = ms[j], ms[i] }
