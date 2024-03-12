package types

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/barweiss/go-tuple"
	"github.com/tbellembois/gochimitheque-wasm/models"
	"github.com/tbellembois/gochimitheque-wasm/select2"
)

type Select2Names struct {
	Rows  []*Name `json:"rows"`
	Total int     `json:"total"`
}

type Name struct {
	*models.Name
}

func (elems Select2Names) GetRowConcreteTypeName() string {

	return "Name"

}

func (elems Select2Names) IsExactMatch() bool {

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

func (Select2Names) FromJsJSONValue(jsvalue js.Value) select2.Select2ResultAble {
	var (
		namesAjaxResponse tuple.T2[[]struct {
			MatchExactSearch bool   `json:"match_exact_search"`
			NameID           int64  `json:"name_id"`
			NameLabel        string `json:"name_label"`
		}, int]
		select2Names Select2Names
		err          error
	)

	jsNamesString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsNamesString), &namesAjaxResponse); err != nil {
		fmt.Println("(Select2Names) FromJsJSONValue:" + err.Error())
	}

	for _, name := range namesAjaxResponse.V1 {
		select2Names.Rows = append(select2Names.Rows, &Name{
			&models.Name{
				MatchExactSearch: name.MatchExactSearch,
				NameID:           int(name.NameID),
				NameLabel:        name.NameLabel,
			},
		})
	}

	select2Names.Total = namesAjaxResponse.V2

	return select2Names
}

func (n Select2Names) GetRows() []select2.Select2ItemAble {

	var select2ItemAble []select2.Select2ItemAble = make([]select2.Select2ItemAble, len(n.Rows))

	for i, row := range n.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (n Select2Names) GetTotal() int {

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
