package types

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/barweiss/go-tuple"
	"github.com/tbellembois/gochimitheque-wasm/models"
	"github.com/tbellembois/gochimitheque-wasm/select2"
)

type Select2ClassesOfCompound struct {
	Rows  []*ClassOfCompound `json:"rows"`
	Total int                `json:"total"`
}

type ClassOfCompound struct {
	*models.ClassOfCompound
}

func (elems Select2ClassesOfCompound) GetRowConcreteTypeName() string {

	return "ClassOfCompound"

}

func (elems Select2ClassesOfCompound) IsExactMatch() bool {

	for _, elem := range elems.Rows {
		if elem.MatchExactSearch {
			return true
		}
	}

	return false

}

func (c *ClassOfCompound) MarshalJSON() ([]byte, error) {
	type Copy ClassOfCompound
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

func (Select2ClassesOfCompound) FromJsJSONValue(jsvalue js.Value) select2.Select2ResultAble {
	var (
		classesOfCompoundNumbersAjaxResponse tuple.T2[[]struct {
			MatchExactSearch     bool   `json:"match_exact_search"`
			ClassOfCompoundID    int64  `json:"classofcompound_id"`
			ClassOfCompoundLabel string `json:"classofcompound_label"`
		}, int]
		select2ClassesOfCompound Select2ClassesOfCompound
		err                      error
	)

	jsCasNumbersString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsCasNumbersString), &classesOfCompoundNumbersAjaxResponse); err != nil {
		fmt.Println("(Select2ClassesOfCompound) FromJsJSONValue:" + err.Error())
	}

	for _, classofcompound := range classesOfCompoundNumbersAjaxResponse.V1 {
		select2ClassesOfCompound.Rows = append(select2ClassesOfCompound.Rows, &ClassOfCompound{
			&models.ClassOfCompound{
				MatchExactSearch:     classofcompound.MatchExactSearch,
				ClassOfCompoundID:    int(classofcompound.ClassOfCompoundID),
				ClassOfCompoundLabel: classofcompound.ClassOfCompoundLabel,
			},
		})
	}

	select2ClassesOfCompound.Total = classesOfCompoundNumbersAjaxResponse.V2

	return select2ClassesOfCompound
}

func (c ClassOfCompound) FromJsJSONValue(jsvalue js.Value) select2.Select2ItemAble {

	var (
		classOfCompound ClassOfCompound
		err             error
	)

	jsClassOfCompoundString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsClassOfCompoundString), &classOfCompound); err != nil {
		fmt.Println(err)
	}

	return classOfCompound

}

func (c Select2ClassesOfCompound) GetRows() []select2.Select2ItemAble {

	var select2ItemAble []select2.Select2ItemAble = make([]select2.Select2ItemAble, len(c.Rows))

	for i, row := range c.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (c Select2ClassesOfCompound) GetTotal() int {

	return c.Total

}

func (c ClassOfCompound) GetSelect2Id() int {

	return c.ClassOfCompoundID

}

func (c ClassOfCompound) GetSelect2Text() string {

	if c.ClassOfCompound != nil {
		return c.ClassOfCompoundLabel
	}

	return ""

}
