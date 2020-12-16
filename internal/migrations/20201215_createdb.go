package migrations

import "database/sql"

var migrationCreateMerchant = Migration{
	NameString: "20200902_create_merchant",
	ApplyFunc: func(tx *sql.Tx) error {
		query := `CREATE TABLE IF NOT EXISTS merchant (` +
			`  id BIGSERIAL PRIMARY KEY` +
			`, name TEXT NOT NULL` +
			`, priority INT` +
			`, company_name TEXT` +
			`, website TEXT` +
			`, address_id BIGINT` +
			`, languages TEXT` +
			`, default_language TEXT` +
			`, shipping_costs INT` +
			`, shipping_costs_free INT` +
			`, shipping_partner TEXT` +
			`, created_at TIMESTAMPTZ` +
			`, updated_at TIMESTAMPTZ` +
			`, channel_id BIGINT` +
			`);`

		if _, err := tx.Exec(query); err != nil {
			return err
		}
		return nil
	},
	RollbackFunc: func(tx *sql.Tx) error {
		query := `DROP TABLE IF EXISTS merchant;`
		if _, err := tx.Exec(query); err != nil {
			return err
		}
		return nil
	},
}
