package uow

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/dzyanis/go-service-example/internal/transactions"
	transRepo "github.com/dzyanis/go-service-example/internal/transactions/repositories"
	"github.com/dzyanis/go-service-example/internal/users"
	usersRepo "github.com/dzyanis/go-service-example/internal/users/repositories"
	"github.com/dzyanis/go-service-example/internal/wallets"
	walletsRepo "github.com/dzyanis/go-service-example/internal/wallets/repositories"
	"github.com/dzyanis/go-service-example/pkg/database"
)

type StartFunc func() (UnitOfWork, error)

type Transaction interface {
	Commit() error
	Rollback() error
}

type UnitOfWork interface {
	Transaction

	Trans() transactions.Repository
	Users() users.Repository
	Wallets() wallets.Repository
}

type Repository struct {
	tx      *database.Transaction
	trans   transactions.Repository
	users   users.Repository
	wallets wallets.Repository
}

func NewRepository(dbc *sqlx.DB) (*Repository, error) {
	tx, err := dbc.Beginx()
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}

	dbtx := database.NewTx(tx)
	return &Repository{
		tx:      dbtx,
		trans:   transRepo.NewRepository(dbtx),
		users:   usersRepo.NewRepository(dbtx),
		wallets: walletsRepo.NewRepository(dbtx),
	}, nil
}

func (r *Repository) Trans() transactions.Repository { return r.trans }
func (r *Repository) Users() users.Repository        { return r.users }
func (r *Repository) Wallets() wallets.Repository    { return r.wallets }

func (r *Repository) Commit() error   { return r.tx.Commit() }
func (r *Repository) Rollback() error { return r.tx.Rollback() }
