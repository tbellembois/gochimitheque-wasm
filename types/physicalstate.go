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

type Select2PhysicalStates struct {
	Rows  []*PhysicalState `json:"rows"`
	Total int              `json:"total"`
}

type PhysicalState struct {
	*models.PhysicalState
}

func (elems Select2PhysicalStates) GetRowConcreteTypeName() string {

	return "PhysicalState"

}

func (elems Select2PhysicalStates) IsExactMatch() bool {

	for _, elem := range elems.Rows {
		if elem.MatchExactSearch {
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

func (Select2PhysicalStates) FromJsJSONValue(jsvalue js.Value) select2.Select2ResultAble {
	var (
		physicalstatesAjaxResponse tuple.T2[[]struct {
			MatchExactSearch   bool   `json:"match_exact_search"`
			PhysicalStateID    int64  `json:"physicalstate_id"`
			PhysicalStateLabel string `json:"physicalstate_label"`
		}, int]
		select2PhysicalStates Select2PhysicalStates
		err                   error
	)

	jsPhysicalStatesString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsPhysicalStatesString), &physicalstatesAjaxResponse); err != nil {
		fmt.Println("(Select2PhysicalStates) FromJsJSONValue:" + err.Error())
	}

	for _, physicalstate := range physicalstatesAjaxResponse.V1 {
		select2PhysicalStates.Rows = append(select2PhysicalStates.Rows, &PhysicalState{
			&models.PhysicalState{
				MatchExactSearch:   physicalstate.MatchExactSearch,
				PhysicalStateID:    sql.NullInt64{Int64: physicalstate.PhysicalStateID, Valid: true},
				PhysicalStateLabel: sql.NullString{String: physicalstate.PhysicalStateLabel, Valid: true},
			},
		})
	}

	select2PhysicalStates.Total = physicalstatesAjaxResponse.V2

	return select2PhysicalStates
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

func (p Select2PhysicalStates) GetRows() []select2.Select2ItemAble {

	var select2ItemAble []select2.Select2ItemAble = make([]select2.Select2ItemAble, len(p.Rows))

	for i, row := range p.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (p Select2PhysicalStates) GetTotal() int {

	return p.Total

}

func (p PhysicalState) GetSelect2Id() int {

	return int(p.PhysicalStateID.Int64)

}

func (p PhysicalState) GetSelect2Text() string {

	if p.PhysicalState != nil {
		return p.PhysicalStateLabel.String
	}

	return ""

}
