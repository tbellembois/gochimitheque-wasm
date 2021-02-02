package types

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/select2"
)

type PhysicalStates struct {
	Rows  []*PhysicalState `json:"rows"`
	Total int              `json:"total"`
}

type PhysicalState struct {
	C                  int            `json:"c"` // not stored in db but db:"c" set for sqlx
	PhysicalStateID    sql.NullInt64  `json:"physicalstate_id"`
	PhysicalStateLabel sql.NullString `json:"physicalstate_label"`
}

func (elems PhysicalStates) IsExactMatch() bool {

	for _, elem := range elems.Rows {
		if elem.C == 1 {
			return true
		}
	}

	return false

}

func (p *PhysicalState) MarshalJSON() ([]byte, error) {
	type Copy PhysicalState
	return json.Marshal(&struct {
		Id   int    `json:"id"`
		Text string `json:"text"`

		*Copy
	}{
		Id:   p.GetSelect2Id(),
		Text: p.GetSelect2Text(),
		Copy: (*Copy)(p),
	})
}

func (PhysicalStates) FromJsJSONValue(jsvalue js.Value) select2.Select2ResultAble {

	var (
		physicalStates PhysicalStates
		err            error
	)

	jsPhysicalStatesString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsPhysicalStatesString), &physicalStates); err != nil {
		fmt.Println(err)
	}

	return physicalStates

}

func (p PhysicalState) FromJsJSONValue(jsvalue js.Value) select2.Select2ItemAble {

	var (
		physicalState PhysicalState
		err           error
	)

	jsPhysicalStateString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsPhysicalStateString), &physicalState); err != nil {
		fmt.Println(err)
	}

	return physicalState

}

func (p PhysicalStates) GetRows() []select2.Select2ItemAble {

	var select2ItemAble []select2.Select2ItemAble = make([]select2.Select2ItemAble, len(p.Rows))

	for i, row := range p.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (p PhysicalStates) GetTotal() int {

	return p.Total

}

func (p PhysicalState) GetSelect2Id() int {

	return int(p.PhysicalStateID.Int64)

}

func (p PhysicalState) GetSelect2Text() string {

	return p.PhysicalStateLabel.String

}
