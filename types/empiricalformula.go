package types

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/select2"
)

type EmpiricalFormulas struct {
	Rows  []*EmpiricalFormula `json:"rows"`
	Total int                 `json:"total"`
}

type EmpiricalFormula struct {
	C                     int            `json:"c"` // not stored in db but db:"c" set for sqlx
	EmpiricalFormulaID    sql.NullInt64  `json:"empiricalformula_id"`
	EmpiricalFormulaLabel sql.NullString `json:"empiricalformula_label"`
}

func (elems EmpiricalFormulas) GetRowConcreteTypeName() string {

	return "EmpiricalFormula"

}

func (elems EmpiricalFormulas) IsExactMatch() bool {

	for _, elem := range elems.Rows {
		if elem.C == 1 {
			return true
		}
	}

	return false

}

func (e *EmpiricalFormula) MarshalJSON() ([]byte, error) {
	type Copy EmpiricalFormula
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

func (EmpiricalFormulas) FromJsJSONValue(jsvalue js.Value) select2.Select2ResultAble {

	var (
		empiricalFormulas EmpiricalFormulas
		err               error
	)

	jsEmpiricalFormulasString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsEmpiricalFormulasString), &empiricalFormulas); err != nil {
		fmt.Println(err)
	}

	return empiricalFormulas

}

func (e EmpiricalFormula) FromJsJSONValue(jsvalue js.Value) select2.Select2ItemAble {

	var (
		empiricalFormula EmpiricalFormula
		err              error
	)

	jsEmpiricalFormulaString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsEmpiricalFormulaString), &empiricalFormula); err != nil {
		fmt.Println(err)
	}

	return empiricalFormula

}

func (e EmpiricalFormulas) GetRows() []select2.Select2ItemAble {

	var select2ItemAble []select2.Select2ItemAble = make([]select2.Select2ItemAble, len(e.Rows))

	for i, row := range e.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (e EmpiricalFormulas) GetTotal() int {

	return e.Total

}

func (e EmpiricalFormula) GetSelect2Id() int {

	return int(e.EmpiricalFormulaID.Int64)

}

func (e EmpiricalFormula) GetSelect2Text() string {

	return e.EmpiricalFormulaLabel.String

}
