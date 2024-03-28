package types

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/barweiss/go-tuple"
	"github.com/tbellembois/gochimitheque-wasm/select2"
	"github.com/tbellembois/gochimitheque/models"
)

type Select2Symbols struct {
	Rows  []*Symbol `json:"rows"`
	Total int       `json:"total"`
}

type Symbol struct {
	*models.Symbol
}

func (elems Select2Symbols) GetRowConcreteTypeName() string {

	return "Symbol"

}

func (elems Select2Symbols) IsExactMatch() bool {

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

func (Select2Symbols) FromJsJSONValue(jsvalue js.Value) select2.Select2ResultAble {
	var (
		symbolsAjaxResponse tuple.T2[[]struct {
			MatchExactSearch bool   `json:"match_exact_search"`
			SymbolID         int64  `json:"symbol_id"`
			SymbolLabel      string `json:"symbol_label"`
		}, int]
		select2Symbols Select2Symbols
		err            error
	)

	jsSymbolsString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsSymbolsString), &symbolsAjaxResponse); err != nil {
		fmt.Println("(Select2Symbols) FromJsJSONValue:" + err.Error())
	}

	for _, symbol := range symbolsAjaxResponse.V1 {
		select2Symbols.Rows = append(select2Symbols.Rows, &Symbol{
			&models.Symbol{
				// MatchExactSearch: symbol.MatchExactSearch,
				SymbolID:    int(symbol.SymbolID),
				SymbolLabel: symbol.SymbolLabel,
			},
		})
	}

	select2Symbols.Total = symbolsAjaxResponse.V2

	return select2Symbols
}

func (s Select2Symbols) GetRows() []select2.Select2ItemAble {

	var select2ItemAble []select2.Select2ItemAble = make([]select2.Select2ItemAble, len(s.Rows))

	for i, row := range s.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (s Select2Symbols) GetTotal() int {

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
