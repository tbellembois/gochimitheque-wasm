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

type Select2Units struct {
	Rows  []*Unit `json:"rows"`
	Total int     `json:"total"`
}

type Unit struct {
	*models.Unit
}

func (elems Select2Units) GetRowConcreteTypeName() string {

	return "Unit"

}

func (elems Select2Units) IsExactMatch() bool {

	return false

}

func (u *Unit) MarshalJSON() ([]byte, error) {
	type Copy Unit
	return json.Marshal(&struct {
		Id   int    `json:"id"`
		Text string `json:"text"`
		Copy
	}{
		Id:   u.GetSelect2Id(),
		Text: u.GetSelect2Text(),
		Copy: (Copy)(*u),
	})
}

func (Select2Units) FromJsJSONValue(jsvalue js.Value) select2.Select2ResultAble {
	var (
		unitsAjaxResponse tuple.T2[[]struct {
			// MatchExactSearch bool   `json:"match_exact_search"`
			UnitID    int64  `json:"unit_id"`
			UnitLabel string `json:"unit_label"`
		}, int]
		select2Units Select2Units
		err          error
	)

	jsUnitsString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsUnitsString), &unitsAjaxResponse); err != nil {
		fmt.Println("(Select2Units) FromJsJSONValue:" + err.Error())
	}

	for _, unit := range unitsAjaxResponse.V1 {
		select2Units.Rows = append(select2Units.Rows, &Unit{
			&models.Unit{
				// MatchExactSearch: unit.MatchExactSearch,
				UnitID:    sql.NullInt64{Int64: unit.UnitID, Valid: true},
				UnitLabel: sql.NullString{String: unit.UnitLabel, Valid: true},
			},
		})
	}

	select2Units.Total = unitsAjaxResponse.V2

	return select2Units
}

func (u Unit) FromJsJSONValue(jsvalue js.Value) select2.Select2ItemAble {

	var (
		unit Unit
		err  error
	)

	jsUnitString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsUnitString), &unit); err != nil {
		fmt.Println(err)
	}

	return unit

}

func (u Select2Units) GetRows() []select2.Select2ItemAble {

	var select2ItemAble []select2.Select2ItemAble = make([]select2.Select2ItemAble, len(u.Rows))

	for i, row := range u.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (u Select2Units) GetTotal() int {

	return u.Total

}

func (u Unit) GetSelect2Id() int {

	return int(u.UnitID.Int64)

}

func (u Unit) GetSelect2Text() string {

	if u.Unit != nil {
		return u.UnitLabel.String
	}

	return ""

}
