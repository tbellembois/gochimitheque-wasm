package types

import (
	"encoding/json"
	"fmt"
	"syscall/js"
)

type Symbols struct {
	Rows  []*Symbol `json:"rows"`
	Total int       `json:"total"`
}

type Symbol struct {
	SymbolID    int    `json:"symbol_id"`
	SymbolLabel string `json:"symbol_label"`
	SymbolImage string `json:"symbol_image"`
}

func (elems Symbols) IsExactMatch() bool {

	return false

}

func (s *Symbol) MarshalJSON() ([]byte, error) {
	type Copy Symbol
	return json.Marshal(&struct {
		Id   int    `json:"id"`
		Text string `json:"text"`

		*Copy
	}{
		Id:   s.GetSelect2Id(),
		Text: s.GetSelect2Text(),
		Copy: (*Copy)(s),
	})
}

func (Symbols) FromJsJSONValue(jsvalue js.Value) Select2ResultAble {

	var (
		symbols Symbols
		err     error
	)

	jsSymbolsString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsSymbolsString), &symbols); err != nil {
		fmt.Println(err)
	}

	return symbols

}

func (s Symbols) GetRows() []Select2ItemAble {

	var select2ItemAble []Select2ItemAble = make([]Select2ItemAble, len(s.Rows))

	for i, row := range s.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (s Symbols) GetTotal() int {

	return s.Total

}

func (s Symbol) FromJsJSONValue(jsvalue js.Value) Select2ItemAble {

	var (
		symbol Symbol
		err    error
	)

	jsSymbolString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsSymbolString), &symbol); err != nil {
		fmt.Println(err)
	}

	return symbol

}

func (s Symbol) GetSelect2Id() int {

	return s.SymbolID

}

func (s Symbol) GetSelect2Text() string {

	return s.SymbolLabel

}
