package users

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type Request struct {
	ID    int64   `json:"id"`
	Name  *string `json:"name"`
	Email *string `json:"email"`
}

type Response struct {
	ID        int64     `json:"id,"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
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

func (s *Service) Create(ctx context.Context, name string, email string) (*Response, error) {
	userID, err := s.repo.Create(ctx, &User{
		Name:      name,
		Email:     email,
		Deleted:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return nil, fmt.Errorf("creating: %w", err)
	}

	return s.Get(ctx, userID)
}

func (s *Service) Update(ctx context.Context, r Request) error {
	// TODO: implement
	return errors.New("not implemented")
}

func (s *Service) Delete(ctx context.Context, userID int64) error {
	if err := s.repo.Delete(ctx, userID); err != nil {
		return fmt.Errorf("deleting: %w", err)
	}
	return nil
}

func (s *Service) Get(ctx context.Context, userID int64) (*Response, error) {
	u, err := s.repo.Get(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("getting: %w", err)
	}
	return &Response{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Deleted:   u.Deleted,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}, nil
}
