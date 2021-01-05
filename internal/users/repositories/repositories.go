package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/dzyanis/go-service-example/internal/users"
	"github.com/dzyanis/go-service-example/pkg/database"
)

type Repository struct {
	db database.Database
}

func NewRepository(db database.Database) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, u *users.User) (int64, error) {
	query := `
		INSERT INTO users
			(name, email, deleted, created_at, updated_at)
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING
			id;
	`

	res := struct {
		LastInsertID int64 `db:"id"`
	}{}
	err := r.db.Write().GetContext(ctx, &res, query,
		u.Name, u.Email, u.Deleted, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		return 0, fmt.Errorf("query: %w", err)
	}

	return res.LastInsertID, nil
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

func (r *Repository) Get(ctx context.Context, userID int64) (*users.User, error) {
	var (
		user  = users.User{}
		query = `SELECT * FROM users WHERE id = $1 LIMIT 1;`
	)

	err := r.db.Write().GetContext(ctx, &user, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query row: %w", err)
	}

	return &user, nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*users.User, error) {
	var (
		user  = users.User{}
		query = `SELECT * FROM users WHERE email = $1 LIMIT 1;`
	)
	err := r.db.Write().GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, fmt.Errorf("query row: %w", err)
	}

	return &user, nil
}
