package types

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/models"
	"github.com/tbellembois/gochimitheque-wasm/select2"
)

type Entities struct {
	Rows  []*Entity `json:"rows"`
	Total int       `json:"total"`
}

type Entity struct {
	*models.Entity
}

func (elems Entities) GetRowConcreteTypeName() string {

	return "Entity"

}

func (elems Entities) IsExactMatch() bool {

	return false

}

func (e *Entity) MarshalJSON() ([]byte, error) {
	type Copy Entity
	return json.Marshal(&struct {
		Id   int    `json:"id"`
		Text string `json:"text"`

		*Copy
	}{
		Id:   e.GetSelect2Id(),
		Text: e.GetSelect2Text(),
		Copy: (*Copy)(e),
	})
}

func (Entities) FromJsJSONValue(jsvalue js.Value) select2.Select2ResultAble {

	var (
		entities Entities
		err      error
	)

	jsEntitiesString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsEntitiesString), &entities); err != nil {
		fmt.Println(err)
	}

	return entities

}

func (e Entities) GetRows() []select2.Select2ItemAble {

	var select2ItemAble []select2.Select2ItemAble = make([]select2.Select2ItemAble, len(e.Rows))

	for i, row := range e.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (e Entities) GetTotal() int {

	return e.Total

}

func (e Entity) FromJsJSONValue(jsvalue js.Value) select2.Select2ItemAble {

	var (
		entity Entity
		err    error
	)

	jsEntityString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsEntityString), &entity); err != nil {
		fmt.Println(err)
	}

	return entity

}

func (e Entity) GetSelect2Id() int {

	return e.EntityID

}

func (e Entity) GetSelect2Text() string {

	if e.Entity != nil {
		return e.EntityName
	}

	return ""

}
