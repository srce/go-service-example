package wallets

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dzyanis/go-service-example/pkg/currencies"
)

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

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context,
	userID int64, amount int64, currency currencies.Currency) (*Response, error) {
	walletID, err := s.repo.Create(ctx, &Wallet{
		UserID:    userID,
		Amount:    amount,
		Currency:  string(currency),
		Deleted:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return nil, fmt.Errorf("creating: %w", err)
	}

	return s.Get(ctx, walletID)
}

func (s *Service) Update(ctx context.Context, r Request) error {
	// TODO: implement
	return errors.New("not implemented")
}

func (s *Service) Delete(ctx context.Context, walletID int64) error {
	if err := s.repo.Delete(ctx, walletID); err != nil {
		return fmt.Errorf("deleting: %w", err)
	}
	return nil
}

func (s *Service) Get(ctx context.Context, walletID int64) (*Response, error) {
	w, err := s.repo.Get(ctx, walletID)
	if err != nil {
		return nil, fmt.Errorf("getting: %w", err)
	}
	return &Response{
		ID:        w.ID,
		UserID:    w.UserID,
		Amount:    w.Amount,
		Currency:  w.Currency,
		Deleted:   w.Deleted,
		CreatedAt: w.CreatedAt,
		UpdatedAt: w.UpdatedAt,
	}, nil
}
