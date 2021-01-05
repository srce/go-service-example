package transactions

import (
	"context"
	"net/http"
	"time"

	"github.com/dzyanis/go-service-example/pkg/money"
)

const (
	CompanyBeneficiaryEmail = "company@beneficiary"
	CompanyFeePercent       = 1.5
)

type Transaction struct {
	ID            int64     `db:"id"`
	SenderID      int64     `db:"sender_id"`
	BeneficiaryID int64     `db:"beneficiary_id"`
	Amount        int64     `db:"amount"`
	Currency      string    `db:"currency"`
	CreatedAt     time.Time `db:"created_at"`
}

type Repository interface {
	Create(ctx context.Context, t *Transaction) (int64, error)
	Get(ctx context.Context, transactionID int64) (*Transaction, error)
}

type Service interface {
	Transfer(ctx context.Context, senderID, beneficiaryID int64, amount money.Money) error
}

type Controller interface {
	Transfer(w http.ResponseWriter, r *http.Request)
}
