package transactions

import (
	"fmt"

	"github.com/dzyanis/go-service-example/internal/users"
	"github.com/dzyanis/go-service-example/internal/wallets"
	"github.com/dzyanis/go-service-example/pkg/database"
	"github.com/jmoiron/sqlx"
)

type UOWStartFunc func() (*UOW, error)

type UOW struct {
	tx      *database.Transaction
	trans   *Repository
	users   *users.Repository
	wallets *wallets.Repository
}

func NewUOW(dbc *sqlx.DB) (*UOW, error) {
	tx, err := dbc.Beginx()
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}

	db := database.NewTx(tx)
	return &UOW{
		tx:      db,
		trans:   NewRepository(db),
		users:   users.NewRepository(db),
		wallets: wallets.NewRepository(db),
	}, nil
}

func (u *UOW) Trans() *Repository           { return u.trans }
func (u *UOW) Users() *users.Repository     { return u.users }
func (u *UOW) Wallets() *wallets.Repository { return u.wallets }

func (u *UOW) Commit() error   { return u.tx.Commit() }
func (u *UOW) Rollback() error { return u.tx.Rollback() }
