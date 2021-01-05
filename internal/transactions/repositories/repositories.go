package transactions

import (
	"context"
	"fmt"

	"github.com/dzyanis/go-service-example/internal/transactions"
	"github.com/dzyanis/go-service-example/pkg/database"
)

type Repository struct {
	db database.Database
}

func NewRepository(db database.Database) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, t *transactions.Transaction) (int64, error) {
	query := `
		INSERT INTO transactions
			(sender_id, beneficiary_id, amount, currency)
		VALUES
			($1, $2, $3, $4)
		RETURNING
			id;
	`

	res := struct {
		LastInsertID int64 `db:"id"`
	}{}
	err := r.db.Write().GetContext(ctx, &res, query,
		t.SenderID, t.BeneficiaryID, t.Amount, t.Currency)
	if err != nil {
		return 0, fmt.Errorf("query: %w", err)
	}

	return res.LastInsertID, nil
}

func (r *Repository) Get(ctx context.Context, transactionID int64) (*transactions.Transaction, error) {
	var (
		transaction = transactions.Transaction{}
		query       = `SELECT * FROM transactions WHERE id = $1 LIMIT 1;`
	)
	err := r.db.Write().GetContext(ctx, &transaction, query, transactionID)
	if err != nil {
		return nil, fmt.Errorf("query row: %w", err)
	}

	return &transaction, nil
}
