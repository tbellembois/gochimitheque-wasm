package types

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/select2"
	"github.com/tbellembois/gochimitheque/models"
)

type Symbols struct {
	Rows  []*Symbol `json:"rows"`
	Total int       `json:"total"`
}

type Symbol struct {
	*models.Symbol
}

func (elems Symbols) GetRowConcreteTypeName() string {

	return "Symbol"

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

func (Symbols) FromJsJSONValue(jsvalue js.Value) select2.Select2ResultAble {

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

func (s Symbols) GetRows() []select2.Select2ItemAble {

	var select2ItemAble []select2.Select2ItemAble = make([]select2.Select2ItemAble, len(s.Rows))

	for i, row := range s.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (s Symbols) GetTotal() int {

	return s.Total

}

func (s Symbol) FromJsJSONValue(jsvalue js.Value) select2.Select2ItemAble {

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

	if s.Symbol != nil {
		return s.SymbolLabel
	}

	return ""

}
