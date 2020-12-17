package wallets

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dzyanis/go-service-example/pkg/database"
)

type Wallet struct {
	ID        int64     `db:"id"`
	UserID    int64     `db:"user_id"`
	Value     int64     `db:"value"`
	Currency  string    `db:"currency"`
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

func (r *Repository) Create(ctx context.Context, w *Wallet) (int64, error) {
	query := `
		INSERT INTO wallets
			(user_id, value, currency, deleted, created_at, updated_at)
		VALUES
			($1, $2, $3, $4, $5, $6)
		RETURNING
			id;
	`

	row := r.db.Write().QueryRowContext(ctx, query,
		w.UserID, w.Value, w.Currency, w.Deleted, w.CreatedAt, w.UpdatedAt)
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

func (r *Repository) Delete(ctx context.Context, walletID int64) error {
	query := `
		UPDATE wallets
		SET
			deleted = TRUE
		WHERE
			id = $1;
	`
	res, err := r.db.Write().ExecContext(ctx, query, walletID)
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

func (r *Repository) Get(ctx context.Context, walletID int64) (*Wallet, error) {
	wallet := Wallet{}

	query := `SELECT * FROM wallets WHERE id = $1 LIMIT 1;`
	err := r.db.Write().GetContext(ctx, &wallet, query, walletID)
	if err != nil {
		return nil, fmt.Errorf("query row: %w", err)
	}

	return &wallet, nil
}

func (r *Repository) GetByUserID(ctx context.Context, userID int64) (*Wallet, error) {
	wallet := Wallet{}

	query := `SELECT * FROM wallets WHERE user_id = $1 LIMIT 1;`
	err := r.db.Write().GetContext(ctx, &wallet, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query row: %w", err)
	}

	return &wallet, nil
}
