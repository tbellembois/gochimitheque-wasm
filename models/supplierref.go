package models

// SupplierRef is a product supplier reference.
type SupplierRef struct {
	MatchExactSearch bool      `db:"match_exact_search" json:"match_exact_search"` // not stored in db but db:"c" set for sqlx
	SupplierRefID    int       `db:"supplierref_id" json:"supplierref_id" schema:"supplierref_id"`
	SupplierRefLabel string    `db:"supplierref_label" json:"supplierref_label" schema:"supplierref_label"`
	Supplier         *Supplier `db:"supplier" json:"supplier" schema:"supplier"`
}
