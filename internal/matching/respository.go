package matching

import (
	"fmt"

	"github.com/dzyanis/go-service-example/pkg/database"
)

type Repository struct {
	db *database.Database
}

func NewRepository(db *database.Database) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Save(m *Matching) error {
	query := `INSERT INTO matching
				(offer_id, status)
			VALUES
				(:offer_id, :status);`
	_, err := r.db.Write().NamedExec(query, m)
	if err != nil {
		return fmt.Errorf("repository: %w", err)
	}
	return nil
}
