package matching

import "time"

type Status string

const (
	StatusNew Status = "new"
)

type Matching struct {
	ID               int64      `db:"id"`
	OfferID          int64      `db:"offer_id"`
	Status           Status     `db:"status"`
	MatchedWineID    *int64     `db:"matched_wine_id"`
	SuggestedWineIDs []byte     `db:"suggested_wine_ids"`
	MatchedBy        *string    `db:"matched_by"`
	MatchedAt        *time.Time `db:"matched_at"`
	CreatedAt        time.Time  `db:"created_at"`
	UpdatedAt        *time.Time `db:"updated_at"`
}
