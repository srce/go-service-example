package users

import (
	"context"
	"net/http"
	"time"
)

type User struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Deleted   bool      `db:"deleted"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

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

type Repository interface {
	Create(ctx context.Context, u *User) (int64, error)
	Update(ctx context.Context, userID int64) error
	Delete(ctx context.Context, userID int64) error
	Get(ctx context.Context, userID int64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
}

type Service interface {
	Create(ctx context.Context, name string, email string) (*Response, error)
	Update(ctx context.Context, r Request) error
	Delete(ctx context.Context, userID int64) error
	Get(ctx context.Context, userID int64) (*Response, error)
}

type Controller interface {
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
}
