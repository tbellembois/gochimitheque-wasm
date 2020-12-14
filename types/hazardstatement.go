package types

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"syscall/js"
)

type HazardStatements struct {
	Rows  []*HazardStatement `json:"rows"`
	Total int                `json:"total"`
}

type HazardStatement struct {
	HazardStatementID        int            `json:"hazardstatement_id"`
	HazardStatementLabel     string         `json:"hazardstatement_label"`
	HazardStatementReference string         `json:"hazardstatement_reference"`
	HazardStatementCMR       sql.NullString `json:"hazardstatement_cmr"`
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

func (HazardStatements) FromJsJSONValue(jsvalue js.Value) Select2ResultAble {

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

func (s HazardStatements) GetRows() []Select2ItemAble {

	var select2ItemAble []Select2ItemAble = make([]Select2ItemAble, len(s.Rows))

	for i, row := range s.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (s HazardStatements) GetTotal() int {

	return s.Total

}

func (h HazardStatement) FromJsJSONValue(jsvalue js.Value) Select2ItemAble {

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

	return s.HazardStatementReference

}
