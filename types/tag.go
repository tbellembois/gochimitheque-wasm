package types

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/models"
	"github.com/tbellembois/gochimitheque-wasm/select2"
)

type Tags struct {
	Rows  []*Tag `json:"rows"`
	Total int    `json:"total"`
}

type Tag struct {
	*models.Tag
}

func (elems Tags) GetRowConcreteTypeName() string {

	return "Tag"

}

func (elems Tags) IsExactMatch() bool {

	for _, elem := range elems.Rows {
		if elem.C == 1 {
			return true
		}
	}

	return false

}

func (t *Tag) MarshalJSON() ([]byte, error) {
	type Copy Tag
	return json.Marshal(&struct {
		Id   int    `json:"id"`
		Text string `json:"text"`

		*Copy
	}{
		Id:   t.GetSelect2Id(),
		Text: t.GetSelect2Text(),
		Copy: (*Copy)(t),
	})
}

func (Tags) FromJsJSONValue(jsvalue js.Value) select2.Select2ResultAble {

	var (
		tags Tags
		err  error
	)

	jsTagsString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsTagsString), &tags); err != nil {
		fmt.Println(err)
	}

	return tags

}

func (t Tag) FromJsJSONValue(jsvalue js.Value) select2.Select2ItemAble {

	var (
		tag Tag
		err error
	)

	jsTagString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsTagString), &tag); err != nil {
		fmt.Println(err)
	}

	return tag

}

func (t Tags) GetRows() []select2.Select2ItemAble {

	var select2ItemAble []select2.Select2ItemAble = make([]select2.Select2ItemAble, len(t.Rows))

	for i, row := range t.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (t Tags) GetTotal() int {

	return t.Total

}

func (t Tag) GetSelect2Id() int {

	return t.TagID

}

func (t Tag) GetSelect2Text() string {

	if t.Tag != nil {
		return t.TagLabel
	}

	return ""

}
