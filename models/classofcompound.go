package models

// ClassOfCompound is a product class of compound.
type ClassOfCompound struct {
	// nullable values to handle optional Product foreign key (gorilla shema nil values)
	MatchExactSearch     bool   `db:"match_exact_search" json:"match_exact_search"` // not stored in db but db:"c" set for sqlx	CeNumberID    sql.NullInt64  `db:"cenumber_id" json:"cenumber_id" schema:"cenumber_id" `
	ClassOfCompoundID    int    `db:"classofcompound_id" json:"classofcompound_id" schema:"classofcompound_id" `
	ClassOfCompoundLabel string `db:"classofcompound_label" json:"classofcompound_label" schema:"classofcompound_label" `
}

func (coc ClassOfCompound) SetMatchExactSearch(MatchExactSearch bool) Searchable {
	coc.MatchExactSearch = MatchExactSearch

	return coc
}

func (coc ClassOfCompound) GetTableName() string {
	return ("classofcompound")
}

func (coc ClassOfCompound) GetIDFieldName() string {
	return ("classofcompound_id")
}

func (coc ClassOfCompound) GetTextFieldName() string {
	return ("classofcompound_label")
}

func (coc ClassOfCompound) GetID() int64 {
	return int64(coc.ClassOfCompoundID)
}
