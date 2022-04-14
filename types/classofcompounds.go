package types

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/select2"
	"github.com/tbellembois/gochimitheque/models"
)

type ClassesOfCompound struct {
	Rows  []*ClassOfCompound `json:"rows"`
	Total int                `json:"total"`
}

type ClassOfCompound struct {
	*models.ClassOfCompound
}

func (elems ClassesOfCompound) GetRowConcreteTypeName() string {

	return "ClassOfCompound"

}

func (elems ClassesOfCompound) IsExactMatch() bool {

	for _, elem := range elems.Rows {
		if elem.C == 1 {
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

func (ClassesOfCompound) FromJsJSONValue(jsvalue js.Value) select2.Select2ResultAble {

	var (
		classesOfCompound ClassesOfCompound
		err               error
	)

	jsClassesOfCompoundString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsClassesOfCompoundString), &classesOfCompound); err != nil {
		fmt.Println(err)
	}

	return classesOfCompound

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

func (c ClassesOfCompound) GetRows() []select2.Select2ItemAble {

	var select2ItemAble []select2.Select2ItemAble = make([]select2.Select2ItemAble, len(c.Rows))

	for i, row := range c.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (c ClassesOfCompound) GetTotal() int {

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
