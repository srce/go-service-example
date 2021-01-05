package wallets

import (
	"context"
	"net/http"
	"time"

	"github.com/dzyanis/go-service-example/pkg/money"
)

type Wallet struct {
	ID        int64     `db:"id"`
	UserID    int64     `db:"user_id"`
	Amount    int64     `db:"amount"`
	Currency  string    `db:"currency"`
	Deleted   bool      `db:"deleted"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Request struct {
	ID       int64  `json:"id"`
	UserID   int64  `json:"user_id"`
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

type Response struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Amount    int64     `json:"amount"`
	Currency  string    `json:"currency"`
	Deleted   bool      `json:"deleted"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Repository interface {
	Create(ctx context.Context, w *Wallet) (int64, error)
	Update(ctx context.Context, wallet *Wallet) error
	Delete(ctx context.Context, walletID int64) error
	Get(ctx context.Context, walletID int64) (*Wallet, error)
	GetByUserIDAndCurrency(ctx context.Context, userID int64, currency string) (*Wallet, error)
}

type Service interface {
	Create(ctx context.Context, userID int64, amount money.Money) (*Response, error)
	Update(ctx context.Context, r Request) error
	Delete(ctx context.Context, walletID int64) error
	Get(ctx context.Context, walletID int64) (*Response, error)
}

type Controller interface {
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
}
