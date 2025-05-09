//go:build go1.24 && js && wasm

package types

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/select2"
	"github.com/tbellembois/gochimitheque/models"
)

type HazardStatements struct {
	Rows  []*HazardStatement `json:"rows"`
	Total int                `json:"total"`
}

type HazardStatement struct {
	*models.HazardStatement
}

func (elems HazardStatements) GetRowConcreteTypeName() string {

	return "HazardStatement"

}

func (elems HazardStatements) IsExactMatch() bool {

	return false

}

func (s *HazardStatement) MarshalJSON() ([]byte, error) {
	type Copy HazardStatement
	return json.Marshal(&struct {
		Id   int    `json:"id"`
		Text string `json:"text"`

		*Copy
	}{
		Id:   s.GetSelect2Id(),
		Text: s.GetSelect2Text(),
		Copy: (*Copy)(s),
	})
}

func (HazardStatements) FromJsJSONValue(jsvalue js.Value) select2.Select2ResultAble {

	var (
		hazardStatements HazardStatements
		err              error
	)

	jsHazardStatementsString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsHazardStatementsString), &hazardStatements); err != nil {
		fmt.Println(err)
	}

	return hazardStatements

}

func (s HazardStatements) GetRows() []select2.Select2ItemAble {

	var select2ItemAble []select2.Select2ItemAble = make([]select2.Select2ItemAble, len(s.Rows))

	for i, row := range s.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (s HazardStatements) GetTotal() int {

	return s.Total

}

func (h HazardStatement) FromJsJSONValue(jsvalue js.Value) select2.Select2ItemAble {

	var (
		hazardStatement HazardStatement
		err             error
	)

	jsHazardStatementString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsHazardStatementString), &hazardStatement); err != nil {
		fmt.Println(err)
	}

	return hazardStatement

}

func (s HazardStatement) GetSelect2Id() int {

	return s.HazardStatementID

}

func (s HazardStatement) GetSelect2Text() string {

	if s.HazardStatement != nil {
		return s.HazardStatementReference
	}

	return ""

}
