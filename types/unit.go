//go:build go1.24 && js && wasm

package types

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/select2"
	"github.com/tbellembois/gochimitheque/models"
)

type Units struct {
	Rows  []*Unit `json:"rows"`
	Total int     `json:"total"`
}

type Unit struct {
	*models.Unit
}

func (elems Units) GetRowConcreteTypeName() string {

	return "Unit"

}

func (elems Units) IsExactMatch() bool {

	return false

}

func (u *Unit) MarshalJSON() ([]byte, error) {
	type Copy Unit
	return json.Marshal(&struct {
		Id   int64  `json:"id"`
		Text string `json:"text"`
		Copy
	}{
		Id:   u.GetSelect2Id(),
		Text: u.GetSelect2Text(),
		Copy: (Copy)(*u),
	})
}

func (Units) FromJsJSONValue(jsvalue js.Value) select2.Select2ResultAble {

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

func (u Units) GetRows() []select2.Select2ItemAble {

	var select2ItemAble []select2.Select2ItemAble = make([]select2.Select2ItemAble, len(u.Rows))

	for i, row := range u.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (u Units) GetTotal() int {

	return u.Total

}

func (u Unit) GetSelect2Id() int64 {

	return *u.UnitID

}

func (u Unit) GetSelect2Text() string {

	if u.Unit != nil {
		return *u.UnitLabel
	}

	return ""

}
