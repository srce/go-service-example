package uow

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/dzyanis/go-service-example/internal/transactions"
	transactionsRepository "github.com/dzyanis/go-service-example/internal/transactions/repositories"
	"github.com/dzyanis/go-service-example/internal/users"
	usersRepositories "github.com/dzyanis/go-service-example/internal/users/repositories"
	"github.com/dzyanis/go-service-example/internal/wallets"
	walletsRepositories "github.com/dzyanis/go-service-example/internal/wallets/repositories"
	"github.com/dzyanis/go-service-example/pkg/database"
)

type StartFunc func() (*UOW, error)

type UOW struct {
	tx      *database.Transaction
	trans   transactions.Repository
	users   users.Repository
	wallets wallets.Repository
}

func NewUOW(dbc *sqlx.DB) (*UOW, error) {
	tx, err := dbc.Beginx()
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}

	db := database.NewTx(tx)
	return &UOW{
		tx:      db,
		trans:   transactionsRepository.NewRepository(db),
		users:   usersRepositories.NewRepository(db),
		wallets: walletsRepositories.NewRepository(db),
	}, nil
}

func (u *UOW) Trans() transactions.Repository { return u.trans }
func (u *UOW) Users() users.Repository        { return u.users }
func (u *UOW) Wallets() wallets.Repository    { return u.wallets }

func (u *UOW) Commit() error   { return u.tx.Commit() }
func (u *UOW) Rollback() error { return u.tx.Rollback() }
