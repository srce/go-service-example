package matching

import (
	"context"
	"fmt"

	"github.com/dzyanis/go-service-example/pkg/database"
)

type Repository struct {
	db database.Database
}

func NewRepository(db database.Database) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Save(ctx context.Context, m *Matching) error {
	query := `
		INSERT INTO matching
			(offer_id, status)
		VALUES
			($1, $2);`
	_, err := r.db.Write().ExecContext(ctx, query, m.OfferID, m.Status)
	if err != nil {
		return fmt.Errorf("repository: %w", err)
	}
	return nil
}
