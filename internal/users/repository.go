package users

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dzyanis/go-service-example/pkg/database"
)

type User struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Deleted   bool      `db:"deleted"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Repository struct {
	db *database.Database
}

func NewRepository(db *database.Database) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, u *User) (int64, error) {
	query := `
		INSERT INTO users
			(name, email, deleted, created_at, updated_at)
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING
			id;
	`

	row := r.db.Write().QueryRowContext(ctx, query,
		u.Name, u.Email, u.Deleted, u.CreatedAt, u.UpdatedAt)
	if err := row.Err(); err != nil {
		return 0, fmt.Errorf("query: %w", err)
	}

	var lastInsertId int64
	if err := row.Scan(&lastInsertId); err != nil {
		return 0, fmt.Errorf("scanning: %w", err)
	}
	return lastInsertId, nil
}

func (r *Repository) Update(ctx context.Context, userID int64) error {
	// TODO: implement
	return errors.New("not implemented")
}

func (r *Repository) Delete(ctx context.Context, userID int64) error {
	query := `
		UPDATE users
		SET
			deleted = TRUE
		WHERE
			id = $1;
	`
	res, err := r.db.Write().ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("named exec: %w", err)
	}

	if n, err := res.RowsAffected(); err != nil {
		return fmt.Errorf("rows affected: %w", err)
	} else if n < 0 {
		return errors.New("no row was affected")
	}
	return nil
}

func (r *Repository) Get(ctx context.Context, userID int64) (*User, error) {
	user := User{}

	query := `SELECT * FROM users WHERE id = $1 LIMIT 1;`
	err := r.db.Write().GetContext(ctx, &user, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query row: %w", err)
	}

	return &user, nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	user := User{}

	query := `SELECT * FROM users WHERE email = $1 LIMIT 1;`
	err := r.db.Write().GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, fmt.Errorf("query row: %w", err)
	}

	return &user, nil
}
