package models

type Searchable interface {
	SetMatchExactSearch(bool) Searchable
	GetTableName() string
	GetIDFieldName() string
	GetTextFieldName() string
	GetID() int64
}
