package types

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/barweiss/go-tuple"
	"github.com/tbellembois/gochimitheque-wasm/select2"
	"github.com/tbellembois/gochimitheque/models"
)

type Select2EmpiricalFormulas struct {
	Rows  []*EmpiricalFormula `json:"rows"`
	Total int                 `json:"total"`
}

type EmpiricalFormula struct {
	*models.EmpiricalFormula
}

func (elems Select2EmpiricalFormulas) GetRowConcreteTypeName() string {

	return "EmpiricalFormula"

}

func (elems Select2EmpiricalFormulas) IsExactMatch() bool {

	for _, elem := range elems.Rows {
		if elem.MatchExactSearch {
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

func (Select2EmpiricalFormulas) FromJsJSONValue(jsvalue js.Value) select2.Select2ResultAble {
	var (
		categoriesAjaxResponse tuple.T2[[]struct {
			MatchExactSearch      bool   `json:"match_exact_search"`
			EmpiricalFormulaID    int64  `json:"empiricalformula_id"`
			EmpiricalFormulaLabel string `json:"empiricalformula_label"`
		}, int]
		select2EmpiricalFormulas Select2EmpiricalFormulas
		err                      error
	)

	jsCategoriesString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsCategoriesString), &categoriesAjaxResponse); err != nil {
		fmt.Println("(Select2EmpiricalFormulas) FromJsJSONValue:" + err.Error())
	}

	for _, empiricalformula := range categoriesAjaxResponse.V1 {
		select2EmpiricalFormulas.Rows = append(select2EmpiricalFormulas.Rows, &EmpiricalFormula{
			&models.EmpiricalFormula{
				MatchExactSearch:      empiricalformula.MatchExactSearch,
				EmpiricalFormulaID:    sql.NullInt64{Int64: empiricalformula.EmpiricalFormulaID, Valid: true},
				EmpiricalFormulaLabel: sql.NullString{String: empiricalformula.EmpiricalFormulaLabel, Valid: true},
			},
		})
	}

	select2EmpiricalFormulas.Total = categoriesAjaxResponse.V2

	return select2EmpiricalFormulas
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

func (e Select2EmpiricalFormulas) GetRows() []select2.Select2ItemAble {

	var select2ItemAble []select2.Select2ItemAble = make([]select2.Select2ItemAble, len(e.Rows))

	for i, row := range e.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (e Select2EmpiricalFormulas) GetTotal() int {

	return e.Total

}

func (e EmpiricalFormula) GetSelect2Id() int {

	return int(e.EmpiricalFormulaID.Int64)

}

func (e EmpiricalFormula) GetSelect2Text() string {

	if e.EmpiricalFormula != nil {
		return e.EmpiricalFormulaLabel.String
	}

	return ""

}
