//go:build go1.24 && js && wasm

package types

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/select2"
	"github.com/tbellembois/gochimitheque/models"
)

type PrecautionaryStatements struct {
	Rows  []*PrecautionaryStatement `json:"rows"`
	Total int                       `json:"total"`
}

type PrecautionaryStatement struct {
	*models.PrecautionaryStatement
}

func (elems PrecautionaryStatements) GetRowConcreteTypeName() string {

	return "PrecautionaryStatement"

}

func (elems PrecautionaryStatements) IsExactMatch() bool {

	return false

}

func (s *PrecautionaryStatement) MarshalJSON() ([]byte, error) {
	type Copy PrecautionaryStatement
	return json.Marshal(&struct {
		Id   int64  `json:"id"`
		Text string `json:"text"`

		*Copy
	}{
		Id:   s.GetSelect2Id(),
		Text: s.GetSelect2Text(),
		Copy: (*Copy)(s),
	})
}

func (PrecautionaryStatements) FromJsJSONValue(jsvalue js.Value) select2.Select2ResultAble {

	var (
		PrecautionaryStatements PrecautionaryStatements
		err                     error
	)

	jsPrecautionaryStatementsString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsPrecautionaryStatementsString), &PrecautionaryStatements); err != nil {
		fmt.Println(err)
	}

	return PrecautionaryStatements

}

func (s PrecautionaryStatements) GetRows() []select2.Select2ItemAble {

	var select2ItemAble []select2.Select2ItemAble = make([]select2.Select2ItemAble, len(s.Rows))

	for i, row := range s.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (s PrecautionaryStatements) GetTotal() int {

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

func (s PrecautionaryStatement) GetSelect2Id() int64 {

	return *s.PrecautionaryStatementID

}

func (s PrecautionaryStatement) GetSelect2Text() string {

	if s.PrecautionaryStatement != nil {
		return s.PrecautionaryStatementReference
	}

	return ""

}
