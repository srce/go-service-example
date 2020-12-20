package transactions

import (
	"context"
	"fmt"
	"time"

	"github.com/dzyanis/go-service-example/pkg/database"
)

type Transaction struct {
	ID            int64     `db:"id"`
	SenderID      int64     `db:"sender_id"`
	BeneficiaryID int64     `db:"beneficiary_id"`
	Amount        int64     `db:"amount"`
	Currency      string    `db:"currency"`
	CreatedAt     time.Time `db:"created_at"`
}

type Repository struct {
	db *database.Database
}

func NewRepository(db *database.Database) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, t *Transaction) (int64, error) {
	query := `
		INSERT INTO transactions
			(sender_id, beneficiary_id, amount, currency)
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING
			id;
	`

	row := r.db.Write().QueryRowContext(ctx, query,
		t.SenderID, t.BeneficiaryID, t.Amount, t.Currency)
	if err := row.Err(); err != nil {
		return 0, fmt.Errorf("query: %w", err)
	}

	var lastInsertID int64
	if err := row.Scan(&lastInsertID); err != nil {
		return 0, fmt.Errorf("scanning: %w", err)
	}
	return lastInsertID, nil
}

func (r *Repository) Get(ctx context.Context, transactionID int64) (*Transaction, error) {
	transaction := Transaction{}

	query := `SELECT * FROM transactions WHERE id = $1 LIMIT 1;`
	err := r.db.Write().GetContext(ctx, &transaction, query, transactionID)
	if err != nil {
		return nil, fmt.Errorf("query row: %w", err)
	}

	return &transaction, nil
}
