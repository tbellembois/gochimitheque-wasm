package types

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/barweiss/go-tuple"
	"github.com/tbellembois/gochimitheque-wasm/select2"
	"github.com/tbellembois/gochimitheque/models"
)

type Select2HazardStatements struct {
	Rows  []*HazardStatement `json:"rows"`
	Total int                `json:"total"`
}

type HazardStatement struct {
	*models.HazardStatement
}

func (elems Select2HazardStatements) GetRowConcreteTypeName() string {

	return "HazardStatement"

}

func (elems Select2HazardStatements) IsExactMatch() bool {

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

func (Select2HazardStatements) FromJsJSONValue(jsvalue js.Value) select2.Select2ResultAble {
	var (
		hazardstatementsAjaxResponse tuple.T2[[]struct {
			// MatchExactSearch     bool   `json:"match_exact_search"`
			HazardStatementID    int64  `json:"hazardstatement_id"`
			HazardStatementLabel string `json:"hazardstatement_label"`
		}, int]
		select2HazardStatements Select2HazardStatements
		err                     error
	)

	jsHazardStatementsString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsHazardStatementsString), &hazardstatementsAjaxResponse); err != nil {
		fmt.Println("(Select2HazardStatements) FromJsJSONValue:" + err.Error())
	}

	for _, hazardstatement := range hazardstatementsAjaxResponse.V1 {
		select2HazardStatements.Rows = append(select2HazardStatements.Rows, &HazardStatement{
			&models.HazardStatement{
				// MatchExactSearch:     hazardstatement.MatchExactSearch,
				HazardStatementID:    int(hazardstatement.HazardStatementID),
				HazardStatementLabel: hazardstatement.HazardStatementLabel,
			},
		})
	}

	select2HazardStatements.Total = hazardstatementsAjaxResponse.V2

	return select2HazardStatements
}

func (s Select2HazardStatements) GetRows() []select2.Select2ItemAble {

	var select2ItemAble []select2.Select2ItemAble = make([]select2.Select2ItemAble, len(s.Rows))

	for i, row := range s.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (s Select2HazardStatements) GetTotal() int {

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
