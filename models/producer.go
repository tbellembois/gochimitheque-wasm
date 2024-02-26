package models

import "database/sql"

// Producer is a product producer.
type Producer struct {
	MatchExactSearch bool           `db:"match_exact_search" json:"match_exact_search"` // not stored in db but db:"c" set for sqlx
	ProducerID       sql.NullInt64  `db:"producer_id" json:"producer_id" schema:"producer_id" `
	ProducerLabel    sql.NullString `db:"producer_label" json:"producer_label" schema:"producer_label" `
}
