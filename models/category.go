package models

import "database/sql"

// Category is a product category.
type Category struct {
	MatchExactSearch bool           `db:"c" json:"c"` // not stored in db but db:"c" set for sqlx
	CategoryID       sql.NullInt64  `db:"category_id" json:"category_id" schema:"category_id" `
	CategoryLabel    sql.NullString `db:"category_label" json:"category_label" schema:"category_label" `
}

func (category Category) SetMatchExactSearch(MatchExactSearch bool) Searchable {
	category.MatchExactSearch = MatchExactSearch

	return category
}

func (category Category) GetTableName() string {
	return ("category")
}

func (category Category) GetIDFieldName() string {
	return ("category_id")
}

func (category Category) GetTextFieldName() string {
	return ("category_label")
}

func (category Category) GetID() int64 {
	return category.CategoryID.Int64
}
