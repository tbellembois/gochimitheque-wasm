package types

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"syscall/js"
)

type CeNumbers struct {
	Rows  []*CeNumber `json:"rows"`
	Total int         `json:"total"`
}

type CeNumber struct {
	C             int            `json:"c"` // not stored in db but db:"c" set for sqlx
	CeNumberID    sql.NullInt64  `json:"cenumber_id"`
	CeNumberLabel sql.NullString `json:"cenumber_label"`
}

func (elems CeNumbers) IsExactMatch() bool {

	for _, elem := range elems.Rows {
		if elem.C == 1 {
			return true
		}
	}

	return false

}

func (c *CeNumber) MarshalJSON() ([]byte, error) {
	type Copy CeNumber
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

func (CeNumbers) FromJsJSONValue(jsvalue js.Value) Select2ResultAble {

	var (
		ceNumbers CeNumbers
		err       error
	)

	jsCeNumbersString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsCeNumbersString), &ceNumbers); err != nil {
		fmt.Println(err)
	}

	return ceNumbers

}

func (c CeNumber) FromJsJSONValue(jsvalue js.Value) Select2ItemAble {

	var (
		ceNumber CeNumber
		err      error
	)

	jsCeNumberString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsCeNumberString), &ceNumber); err != nil {
		fmt.Println(err)
	}

	return ceNumber

}

func (c CeNumbers) GetRows() []Select2ItemAble {

	var select2ItemAble []Select2ItemAble = make([]Select2ItemAble, len(c.Rows))

	for i, row := range c.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (c CeNumbers) GetTotal() int {

	return c.Total

}

func (c CeNumber) GetSelect2Id() int {

	return int(c.CeNumberID.Int64)

}

func (c CeNumber) GetSelect2Text() string {

	return c.CeNumberLabel.String

}
