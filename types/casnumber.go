package types

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/select2"
)

type CasNumbers struct {
	Rows  []*CasNumber `json:"rows"`
	Total int          `json:"total"`
}

type CasNumber struct {
	C              int            `json:"c"` // not stored in db but db:"c" set for sqlx
	CasNumberID    sql.NullInt64  `json:"casnumber_id"`
	CasNumberLabel sql.NullString `json:"casnumber_label"`
	CasNumberCMR   sql.NullString `json:"casnumber_cmr"`
}

func (elems CasNumbers) IsExactMatch() bool {

	for _, elem := range elems.Rows {
		if elem.C == 1 {
			return true
		}
	}

	return false

}

func (c *CasNumber) MarshalJSON() ([]byte, error) {
	type Copy CasNumber
	return json.Marshal(&struct {
		Id   int    `json:"id"`
		Text string `json:"text"`

		*Copy
	}{
		Id:   c.GetSelect2Id(),
		Text: c.GetSelect2Text(),

		Copy: (*Copy)(c),
	})
}

func (CasNumbers) FromJsJSONValue(jsvalue js.Value) select2.Select2ResultAble {

	var (
		casNumbers CasNumbers
		err        error
	)

	jsCasNumbersString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsCasNumbersString), &casNumbers); err != nil {
		fmt.Println(err)
	}

	return casNumbers

}

func (c CasNumber) FromJsJSONValue(jsvalue js.Value) select2.Select2ItemAble {

	var (
		casnumber CasNumber
		err       error
	)

	jsCasNumberString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsCasNumberString), &casnumber); err != nil {
		fmt.Println(err)
	}

	return casnumber

}

func (c CasNumbers) GetRows() []select2.Select2ItemAble {

	var select2ItemAble []select2.Select2ItemAble = make([]select2.Select2ItemAble, len(c.Rows))

	for i, row := range c.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (c CasNumbers) GetTotal() int {

	return c.Total

}

func (c CasNumber) GetSelect2Id() int {

	return int(c.CasNumberID.Int64)

}

func (c CasNumber) GetSelect2Text() string {

	return c.CasNumberLabel.String

}
