package types

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/barweiss/go-tuple"
	"github.com/tbellembois/gochimitheque-wasm/models"
	"github.com/tbellembois/gochimitheque-wasm/select2"
)

type Select2LinearFormulas struct {
	Rows  []*LinearFormula `json:"rows"`
	Total int              `json:"total"`
}

type LinearFormula struct {
	*models.LinearFormula
}

func (elems Select2LinearFormulas) GetRowConcreteTypeName() string {

	return "LinearFormula"

}

func (elems Select2LinearFormulas) IsExactMatch() bool {

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

func (Select2LinearFormulas) FromJsJSONValue(jsvalue js.Value) select2.Select2ResultAble {
	var (
		categoriesAjaxResponse tuple.T2[[]struct {
			MatchExactSearch   bool   `json:"match_exact_search"`
			LinearFormulaID    int64  `json:"linearformula_id"`
			LinearFormulaLabel string `json:"linearformula_label"`
		}, int]
		select2LinearFormulas Select2LinearFormulas
		err                   error
	)

	jsCategoriesString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsCategoriesString), &categoriesAjaxResponse); err != nil {
		fmt.Println("(Select2LinearFormulas) FromJsJSONValue:" + err.Error())
	}

	for _, linearformula := range categoriesAjaxResponse.V1 {
		select2LinearFormulas.Rows = append(select2LinearFormulas.Rows, &LinearFormula{
			&models.LinearFormula{
				MatchExactSearch:   linearformula.MatchExactSearch,
				LinearFormulaID:    sql.NullInt64{Int64: linearformula.LinearFormulaID, Valid: true},
				LinearFormulaLabel: sql.NullString{String: linearformula.LinearFormulaLabel, Valid: true},
			},
		})
	}

	select2LinearFormulas.Total = categoriesAjaxResponse.V2

	return select2LinearFormulas
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

func (e Select2LinearFormulas) GetRows() []select2.Select2ItemAble {

	var select2ItemAble []select2.Select2ItemAble = make([]select2.Select2ItemAble, len(e.Rows))

	for i, row := range e.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (e Select2LinearFormulas) GetTotal() int {

	return e.Total

}

func (e LinearFormula) GetSelect2Id() int {

	return int(e.LinearFormulaID.Int64)

}

func (e LinearFormula) GetSelect2Text() string {

	if e.LinearFormula != nil {
		return e.LinearFormulaLabel.String
	}

	return ""

}
