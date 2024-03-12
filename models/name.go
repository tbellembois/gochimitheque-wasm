package models

// Name is a product name.
type Name struct {
	MatchExactSearch bool   `db:"match_exact_search" json:"match_exact_search"` // not stored in db but db:"c" set for sqlx
	NameID           int    `db:"name_id" json:"name_id" schema:"name_id"`
	NameLabel        string `db:"name_label" json:"name_label" schema:"name_label"`
}

func (name Name) SetMatchExactSearch(MatchExactSearch bool) Searchable {
	name.MatchExactSearch = MatchExactSearch

	return name
}

func (name Name) GetTableName() string {
	return ("name")
}

func (name Name) GetIDFieldName() string {
	return ("name_id")
}

func (name Name) GetTextFieldName() string {
	return ("name_label")
}

func (name Name) GetID() int64 {
	return int64(name.NameID)
}
