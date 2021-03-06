package types

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"syscall/js"
)

type LinearFormulas struct {
	Rows  []*LinearFormula `json:"rows"`
	Total int              `json:"total"`
}

type LinearFormula struct {
	C                  int            `json:"c"` // not stored in db but db:"c" set for sqlx
	LinearFormulaID    sql.NullInt64  `json:"linearformula_id"`
	LinearFormulaLabel sql.NullString `json:"linearformula_label"`
}

func (elems LinearFormulas) IsExactMatch() bool {

	for _, elem := range elems.Rows {
		if elem.C == 1 {
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

func (LinearFormulas) FromJsJSONValue(jsvalue js.Value) Select2ResultAble {

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

func (e LinearFormula) FromJsJSONValue(jsvalue js.Value) Select2ItemAble {

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

func (e LinearFormulas) GetRows() []Select2ItemAble {

	var select2ItemAble []Select2ItemAble = make([]Select2ItemAble, len(e.Rows))

	for i, row := range e.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (e LinearFormulas) GetTotal() int {

	return e.Total

}

func (e LinearFormula) GetSelect2Id() int {

	return int(e.LinearFormulaID.Int64)

}

func (e LinearFormula) GetSelect2Text() string {

	return e.LinearFormulaLabel.String

}
