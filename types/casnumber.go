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

type Select2CasNumbers struct {
	Rows  []*CasNumber `json:"rows"`
	Total int          `json:"total"`
}

type CasNumber struct {
	*models.CasNumber
}

func (elems Select2CasNumbers) GetRowConcreteTypeName() string {
	return "CasNumber"
}

func (elems Select2CasNumbers) IsExactMatch() bool {
	for _, elem := range elems.Rows {
		if elem.MatchExactSearch {
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

func (Select2CasNumbers) FromJsJSONValue(jsvalue js.Value) select2.Select2ResultAble {
	var (
		casNumbersAjaxResponse tuple.T2[[]struct {
			MatchExactSearch bool   `json:"match_exact_search"`
			CasNumberID      int64  `json:"casnumber_id"`
			CasNumberLabel   string `json:"casnumber_label"`
		}, int]
		select2CasNumbers Select2CasNumbers
		err               error
	)

	jsCasNumbersString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsCasNumbersString), &casNumbersAjaxResponse); err != nil {
		fmt.Println("(Select2CasNumbers) FromJsJSONValue:" + err.Error())
	}

	for _, casnumber := range casNumbersAjaxResponse.V1 {
		select2CasNumbers.Rows = append(select2CasNumbers.Rows, &CasNumber{
			&models.CasNumber{
				MatchExactSearch: casnumber.MatchExactSearch,
				CasNumberID:      sql.NullInt64{Int64: casnumber.CasNumberID, Valid: true},
				CasNumberLabel:   sql.NullString{String: casnumber.CasNumberLabel, Valid: true},
			},
		})
	}

	select2CasNumbers.Total = casNumbersAjaxResponse.V2

	return select2CasNumbers
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

func (c Select2CasNumbers) GetRows() []select2.Select2ItemAble {
	var select2ItemAble []select2.Select2ItemAble = make([]select2.Select2ItemAble, len(c.Rows))

	for i, row := range c.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble
}

func (c Select2CasNumbers) GetTotal() int {
	return c.Total
}

func (c CasNumber) GetSelect2Id() int {
	return int(c.CasNumberID.Int64)
}

func (c CasNumber) GetSelect2Text() string {
	if c.CasNumber != nil {
		return c.CasNumberLabel.String
	}

	return ""
}
