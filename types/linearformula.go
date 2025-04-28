//go:build go1.24 && js && wasm

package types

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/select2"
	"github.com/tbellembois/gochimitheque/models"
)

type LinearFormulas struct {
	Rows  []*LinearFormula `json:"rows"`
	Total int              `json:"total"`
}

type LinearFormula struct {
	*models.LinearFormula
}

func (elems LinearFormulas) GetRowConcreteTypeName() string {

	return "LinearFormula"

}

func (elems LinearFormulas) IsExactMatch() bool {

	for _, elem := range elems.Rows {
		if elem.MatchExactSearch {
			return true
		}
	}

	return false

}

func (e *LinearFormula) MarshalJSON() ([]byte, error) {
	type Copy LinearFormula
	return json.Marshal(&struct {
		Id   int    `json:"id"`
		Text string `json:"text"`

		*Copy
	}{
		Id:   e.GetSelect2Id(),
		Text: e.GetSelect2Text(),
		Copy: (*Copy)(e),
	})
}

func (LinearFormulas) FromJsJSONValue(jsvalue js.Value) select2.Select2ResultAble {

	var (
		linearFormulas LinearFormulas
		err            error
	)

	jsLinearFormulasString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsLinearFormulasString), &linearFormulas); err != nil {
		fmt.Println(err)
	}

	return linearFormulas

}

func (e LinearFormula) FromJsJSONValue(jsvalue js.Value) select2.Select2ItemAble {

	var (
		linearFormula LinearFormula
		err           error
	)

	jsLinearFormulaString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsLinearFormulaString), &linearFormula); err != nil {
		fmt.Println(err)
	}

	return linearFormula

}

func (e LinearFormulas) GetRows() []select2.Select2ItemAble {

	var select2ItemAble []select2.Select2ItemAble = make([]select2.Select2ItemAble, len(e.Rows))

	for i, row := range e.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (e LinearFormulas) GetTotal() int {

	return e.Total

}

func (e LinearFormula) GetSelect2Id() int {

	return int(*e.LinearFormulaID)

}

func (e LinearFormula) GetSelect2Text() string {

	if e.LinearFormula != nil {
		return *e.LinearFormulaLabel
	}

	return ""

}
