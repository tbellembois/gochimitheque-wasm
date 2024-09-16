package types

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/models"
	"github.com/tbellembois/gochimitheque-wasm/select2"
	"github.com/tbellembois/gochimitheque/models"
)

type Names struct {
	Rows  []*Name `json:"rows"`
	Total int     `json:"total"`
}

type Name struct {
	*models.Name
}

func (elems Names) GetRowConcreteTypeName() string {

	return "Name"

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

func (Names) FromJsJSONValue(jsvalue js.Value) select2.Select2ResultAble {

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

func (n Names) GetRows() []select2.Select2ItemAble {

	var select2ItemAble []select2.Select2ItemAble = make([]select2.Select2ItemAble, len(n.Rows))

	for i, row := range n.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (n Names) GetTotal() int {

	return n.Total

}

func (n Name) FromJsJSONValue(jsvalue js.Value) select2.Select2ItemAble {

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

	if n.Name != nil {
		return n.NameLabel
	}

	return ""

}
