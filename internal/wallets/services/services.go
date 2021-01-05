package wallets

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dzyanis/go-service-example/internal/wallets"
	"github.com/dzyanis/go-service-example/pkg/money"
)

type Service struct {
	repo wallets.Repository
}

func NewService(repo wallets.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context,
	userID int64, amount money.Money) (*wallets.Response, error) {
	walletID, err := s.repo.Create(ctx, &wallets.Wallet{
		UserID:    userID,
		Amount:    amount.Units(),
		Currency:  amount.Currency().String(),
		Deleted:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return nil, fmt.Errorf("creating: %w", err)
	}

	return s.Get(ctx, walletID)
}

func (s *Service) Update(ctx context.Context, r wallets.Request) error {
	// TODO: implement
	return errors.New("not implemented")
}

func (s *Service) Delete(ctx context.Context, walletID int64) error {
	if err := s.repo.Delete(ctx, walletID); err != nil {
		return fmt.Errorf("deleting: %w", err)
	}
	return nil
}

func (s *Service) Get(ctx context.Context, walletID int64) (*wallets.Response, error) {
	w, err := s.repo.Get(ctx, walletID)
	if err != nil {
		return nil, fmt.Errorf("getting: %w", err)
	}
	return &wallets.Response{
		ID:        w.ID,
		UserID:    w.UserID,
		Amount:    w.Amount,
		Currency:  w.Currency,
		Deleted:   w.Deleted,
		CreatedAt: w.CreatedAt,
		UpdatedAt: w.UpdatedAt,
	}, nil
}
