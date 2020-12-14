package types

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"syscall/js"
)

type Units struct {
	Rows  []*Unit `json:"rows"`
	Total int     `json:"total"`
}

type Unit struct {
	UnitID         sql.NullInt64  `json:"unit_id"`
	UnitLabel      sql.NullString `json:"unit_label"`
	UnitType       sql.NullString `json:"unit_type"`
	Unit           *Unit          `json:"unit"` // reference
	UnitMultiplier int            `json:"-"`
}

func (elems Units) IsExactMatch() bool {

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

func (Units) FromJsJSONValue(jsvalue js.Value) Select2ResultAble {

	var (
		units Units
		err   error
	)

	jsUnitsString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsUnitsString), &units); err != nil {
		fmt.Println(err)
	}

	return units

}

func (u Unit) FromJsJSONValue(jsvalue js.Value) Select2ItemAble {

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

func (u Units) GetRows() []Select2ItemAble {

	var select2ItemAble []Select2ItemAble = make([]Select2ItemAble, len(u.Rows))

	for i, row := range u.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (u Units) GetTotal() int {

	return u.Total

}

func (u Unit) GetSelect2Id() int {

	return int(u.UnitID.Int64)

}

func (u Unit) GetSelect2Text() string {

	return u.UnitLabel.String

}
