//go:build go1.24 && js && wasm

package types

type DictionnaryTerms struct {
	Compound []string `json:"compound"`
}
type Autocomplete struct {
	Total           uint64           `json:"total"`
	DictionaryTerms DictionnaryTerms `json:"dictionary_terms"`
}
