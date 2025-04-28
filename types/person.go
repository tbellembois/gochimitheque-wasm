//go:build go1.24 && js && wasm

package types

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/select2"
	"github.com/tbellembois/gochimitheque/models"
)

type People struct {
	Rows  []*Person `json:"rows"`
	Total int       `json:"total"`
}
type Person struct {
	*models.Person
}

func (elems People) GetRowConcreteTypeName() string {

	return "Person"

}

func (elems People) IsExactMatch() bool {

	return false

}

func (p *Person) MarshalJSON() ([]byte, error) {
	type Copy Person
	return json.Marshal(&struct {
		Id   int    `json:"id"`
		Text string `json:"text"`

		*Copy
	}{
		Id:   p.GetSelect2Id(),
		Text: p.GetSelect2Text(),
		Copy: (*Copy)(p),
	})
}

func (People) FromJsJSONValue(jsvalue js.Value) select2.Select2ResultAble {

	var (
		people People
		err    error
	)

	jsPeopleString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsPeopleString), &people); err != nil {
		fmt.Println(err)
	}

	return people

}

func (p Person) FromJsJSONValue(jsvalue js.Value) select2.Select2ItemAble {

	var (
		person Person
		err    error
	)

	jsPersonString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsPersonString), &person); err != nil {
		fmt.Println(err)
	}

	return person

}

func (p People) GetRows() []select2.Select2ItemAble {

	var select2ItemAble []select2.Select2ItemAble = make([]select2.Select2ItemAble, len(p.Rows))

	for i, row := range p.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (p People) GetTotal() int {

	return p.Total

}

func (p Person) GetSelect2Id() int {

	return p.PersonID

}

func (p Person) GetSelect2Text() string {

	if p.Person != nil {
		return p.PersonEmail
	}

	return ""

}
