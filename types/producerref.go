package types

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/select2"
	"github.com/tbellembois/gochimitheque/models"
)

type ProducerRefs struct {
	Rows  []*ProducerRef `json:"rows"`
	Total int            `json:"total"`
}

type ProducerRef struct {
	*models.ProducerRef
}

func (elems ProducerRefs) GetRowConcreteTypeName() string {

	return "ProducerRef"

}

func (p ProducerRef) ToJsValue() js.Value {

	var (
		marshalJson []byte
		err         error
	)

	if marshalJson, err = json.Marshal(p); err != nil {
		fmt.Println(err)
		return js.Null()
	}

	return js.Global().Get("JSON").Call("parse", string(marshalJson))

}

func (elems ProducerRefs) IsExactMatch() bool {

	for _, elem := range elems.Rows {
		if elem.MatchExactSearch {
			return true
		}
	}

	return false

}

func (r ProducerRef) MarshalJSON() ([]byte, error) {

	type Copy ProducerRef
	return json.Marshal(&struct {
		Id   int    `json:"id"`
		Text string `json:"text"`
		Copy
	}{
		Id:   r.GetSelect2Id(),
		Text: r.GetSelect2Text(),
		Copy: (Copy)(r),
	})

}

func (ProducerRefs) FromJsJSONValue(jsvalue js.Value) select2.Select2ResultAble {

	var (
		producerRefs ProducerRefs
		err          error
	)

	jsProducerRefsString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsProducerRefsString), &producerRefs); err != nil {
		fmt.Println(err)
	}

	return producerRefs

}

func (r ProducerRefs) GetRows() []select2.Select2ItemAble {

	var select2ItemAble []select2.Select2ItemAble = make([]select2.Select2ItemAble, len(r.Rows))

	for i, row := range r.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (r ProducerRefs) GetTotal() int {

	return r.Total

}

func (r ProducerRef) FromJsJSONValue(jsvalue js.Value) select2.Select2ItemAble {

	var (
		producerRef ProducerRef
		err         error
	)

	jsProducerRefString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsProducerRefString), &producerRef); err != nil {
		fmt.Println(err)
	}

	return producerRef

}

func (r ProducerRef) GetSelect2Id() int {

	return int(*r.ProducerRefID)

}

func (r ProducerRef) GetSelect2Text() string {

	if r.ProducerRef != nil {
		return *r.ProducerRefLabel
	}

	return ""

}
