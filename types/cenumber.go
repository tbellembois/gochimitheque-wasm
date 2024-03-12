package types

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/barweiss/go-tuple"
	"github.com/tbellembois/gochimitheque-wasm/models"
	"github.com/tbellembois/gochimitheque-wasm/select2"
)

type Select2CeNumbers struct {
	Rows  []*CeNumber `json:"rows"`
	Total int         `json:"total"`
}

type CeNumber struct {
	*models.CeNumber
}

func (elems Select2CeNumbers) GetRowConcreteTypeName() string {

	return "CeNumber"

}

func (elems Select2CeNumbers) IsExactMatch() bool {

	for _, elem := range elems.Rows {
		if elem.MatchExactSearch {
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

func (Select2CeNumbers) FromJsJSONValue(jsvalue js.Value) select2.Select2ResultAble {
	var (
		cesNumbersAjaxResponse tuple.T2[[]struct {
			MatchExactSearch bool   `json:"match_exact_search"`
			CeNumberID       int64  `json:"cenumber_id"`
			CeNumberLabel    string `json:"cenumber_label"`
		}, int]
		select2CeNumbers Select2CeNumbers
		err              error
	)

	jsCasNumbersString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsCasNumbersString), &cesNumbersAjaxResponse); err != nil {
		fmt.Println("(Select2CeNumbers) FromJsJSONValue:" + err.Error())
	}

	for _, cenumber := range cesNumbersAjaxResponse.V1 {
		select2CeNumbers.Rows = append(select2CeNumbers.Rows, &CeNumber{
			&models.CeNumber{
				MatchExactSearch: cenumber.MatchExactSearch,
				CeNumberID:       sql.NullInt64{Int64: cenumber.CeNumberID, Valid: true},
				CeNumberLabel:    sql.NullString{String: cenumber.CeNumberLabel, Valid: true},
			},
		})
	}

	select2CeNumbers.Total = cesNumbersAjaxResponse.V2

	return select2CeNumbers
}

func (c CeNumber) FromJsJSONValue(jsvalue js.Value) select2.Select2ItemAble {

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

func (c Select2CeNumbers) GetRows() []select2.Select2ItemAble {

	var select2ItemAble []select2.Select2ItemAble = make([]select2.Select2ItemAble, len(c.Rows))

	for i, row := range c.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (c Select2CeNumbers) GetTotal() int {

	return c.Total

}

func (c CeNumber) GetSelect2Id() int {

	return int(c.CeNumberID.Int64)

}

func (c CeNumber) GetSelect2Text() string {

	if c.CeNumber != nil {
		return c.CeNumberLabel.String
	}

	return ""

}
