package types

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/barweiss/go-tuple"
	"github.com/tbellembois/gochimitheque-wasm/select2"
	"github.com/tbellembois/gochimitheque/models"
)

type Select2PrecautionaryStatements struct {
	Rows  []*PrecautionaryStatement `json:"rows"`
	Total int                       `json:"total"`
}

type PrecautionaryStatement struct {
	*models.PrecautionaryStatement
}

func (elems Select2PrecautionaryStatements) GetRowConcreteTypeName() string {

	return "PrecautionaryStatement"

}

func (elems Select2PrecautionaryStatements) IsExactMatch() bool {

	return false

}

func (s *PrecautionaryStatement) MarshalJSON() ([]byte, error) {
	type Copy PrecautionaryStatement
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

func (Select2PrecautionaryStatements) FromJsJSONValue(jsvalue js.Value) select2.Select2ResultAble {
	var (
		precautionarystatementsAjaxResponse tuple.T2[[]struct {
			// MatchExactSearch     bool   `json:"match_exact_search"`
			PrecautionaryStatementID    int64  `json:"precautionarystatement_id"`
			PrecautionaryStatementLabel string `json:"precautionarystatement_label"`
		}, int]
		select2PrecautionaryStatements Select2PrecautionaryStatements
		err                            error
	)

	jsPrecautionaryStatementsString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsPrecautionaryStatementsString), &precautionarystatementsAjaxResponse); err != nil {
		fmt.Println("(Select2PrecautionaryStatements) FromJsJSONValue:" + err.Error())
	}

	for _, precautionarystatement := range precautionarystatementsAjaxResponse.V1 {
		select2PrecautionaryStatements.Rows = append(select2PrecautionaryStatements.Rows, &PrecautionaryStatement{
			&models.PrecautionaryStatement{
				// MatchExactSearch:     precautionarystatement.MatchExactSearch,
				PrecautionaryStatementID:    int(precautionarystatement.PrecautionaryStatementID),
				PrecautionaryStatementLabel: precautionarystatement.PrecautionaryStatementLabel,
			},
		})
	}

	select2PrecautionaryStatements.Total = precautionarystatementsAjaxResponse.V2

	return select2PrecautionaryStatements
}

func (s Select2PrecautionaryStatements) GetRows() []select2.Select2ItemAble {

	var select2ItemAble []select2.Select2ItemAble = make([]select2.Select2ItemAble, len(s.Rows))

	for i, row := range s.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (s Select2PrecautionaryStatements) GetTotal() int {

	return s.Total

}

func (s PrecautionaryStatement) FromJsJSONValue(jsvalue js.Value) select2.Select2ItemAble {

	var (
		precautionaryStatement PrecautionaryStatement
		err                    error
	)

	jsPrecautionaryStatementString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsPrecautionaryStatementString), &precautionaryStatement); err != nil {
		fmt.Println(err)
	}

	return precautionaryStatement

}

func (s PrecautionaryStatement) GetSelect2Id() int {

	return s.PrecautionaryStatementID

}

func (s PrecautionaryStatement) GetSelect2Text() string {

	if s.PrecautionaryStatement != nil {
		return s.PrecautionaryStatementReference
	}

	return ""

}
