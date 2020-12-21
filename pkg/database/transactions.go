package database

import (
	"github.com/jmoiron/sqlx"
)

type Transaction struct {
	tx *sqlx.Tx
}

func NewTx(tx *sqlx.Tx) *Transaction {
	return &Transaction{tx: tx}
}

func (t *Transaction) Write() Ext {
	return t.tx
}

func (t *Transaction) Read() Ext {
	return t.tx
}

func (t *Transaction) Open() error {
	return nil
}

func (t *Transaction) Close() error {
	if err := t.Commit(); err != nil {
		return t.Rollback()
	}
	return nil
}

func (t *Transaction) Rollback() error {
	return t.tx.Rollback()
}

func (t *Transaction) Commit() error {
	return t.tx.Commit()
}
