package types

import (
	"encoding/json"
	"fmt"
	"syscall/js"
)

type Names struct {
	Rows  []*Name `json:"rows"`
	Total int     `json:"total"`
}

type Name struct {
	C         int    `json:"c"` // not stored in db but db:"c" set for sqlx
	NameID    int    `json:"name_id"`
	NameLabel string `json:"name_label"`
}

func (elems Names) IsExactMatch() bool {

	return false

}

func (n *Name) MarshalJSON() ([]byte, error) {
	type Copy Name
	return json.Marshal(&struct {
		Id   int    `json:"id"`
		Text string `json:"text"`

		*Copy
	}{
		Id:   n.GetSelect2Id(),
		Text: n.GetSelect2Text(),
		Copy: (*Copy)(n),
	})
}

func (Names) FromJsJSONValue(jsvalue js.Value) Select2ResultAble {

	var (
		names Names
		err   error
	)

	jsNamesString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsNamesString), &names); err != nil {
		fmt.Println(err)
	}

	return names

}

func (n Names) GetRows() []Select2ItemAble {

	var select2ItemAble []Select2ItemAble = make([]Select2ItemAble, len(n.Rows))

	for i, row := range n.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (n Names) GetTotal() int {

	return n.Total

}

func (n Name) FromJsJSONValue(jsvalue js.Value) Select2ItemAble {

	var (
		name Name
		err  error
	)

	jsNameString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsNameString), &name); err != nil {
		fmt.Println(err)
	}

	return name

}

func (n Name) GetSelect2Id() int {

	return n.NameID

}

func (n Name) GetSelect2Text() string {

	return n.NameLabel

}
