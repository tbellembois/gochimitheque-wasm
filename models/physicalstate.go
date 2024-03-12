package models

import "database/sql"

// PhysicalState is a product physical state.
type PhysicalState struct {
	// nullable values to handle optional Product foreign key (gorilla shema nil values)
	MatchExactSearch   bool           `db:"match_exact_search" json:"match_exact_search"` // not stored in db but db:"c" set for sqlx
	PhysicalStateID    sql.NullInt64  `db:"physicalstate_id" json:"physicalstate_id" schema:"physicalstate_id" `
	PhysicalStateLabel sql.NullString `db:"physicalstate_label" json:"physicalstate_label" schema:"physicalstate_label" `
}

func (physicalstate PhysicalState) SetMatchExactSearch(MatchExactSearch bool) Searchable {
	physicalstate.MatchExactSearch = MatchExactSearch

	return physicalstate
}

func (physicalstate PhysicalState) GetTableName() string {
	return ("physicalstate")
}

func (physicalstate PhysicalState) GetIDFieldName() string {
	return ("physicalstate_id")
}

func (physicalstate PhysicalState) GetTextFieldName() string {
	return ("physicalstate_label")
}

func (physicalstate PhysicalState) GetID() int64 {
	return physicalstate.PhysicalStateID.Int64
}
