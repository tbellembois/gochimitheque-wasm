//go:build go1.24 && js && wasm

package types

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/select2"
	"github.com/tbellembois/gochimitheque/models"
)

type SupplierRefs struct {
	Rows  []*SupplierRef `json:"rows"`
	Total int            `json:"total"`
}

type SupplierRef struct {
	*models.SupplierRef
}

func (elems SupplierRefs) GetRowConcreteTypeName() string {

	return "SupplierRef"

}

func (s SupplierRef) ToJsValue() js.Value {

	var (
		marshalJson []byte
		err         error
	)

	if marshalJson, err = json.Marshal(s); err != nil {
		fmt.Println(err)
		return js.Null()
	}

	return js.Global().Get("JSON").Call("parse", string(marshalJson))

}

func (elems SupplierRefs) IsExactMatch() bool {

	for _, elem := range elems.Rows {
		if elem.MatchExactSearch {
			return true
		}
	}

	return false

}

func (s SupplierRef) MarshalJSON() ([]byte, error) {

	type Copy SupplierRef
	return json.Marshal(&struct {
		Id   int64  `json:"id"`
		Text string `json:"text"`
		Copy
	}{
		Id:   s.GetSelect2Id(),
		Text: s.GetSelect2Text(),
		Copy: (Copy)(s),
	})

}

func (SupplierRefs) FromJsJSONValue(jsvalue js.Value) select2.Select2ResultAble {

	var (
		supplierRefs SupplierRefs
		err          error
	)

	jsSupplierRefsString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsSupplierRefsString), &supplierRefs); err != nil {
		fmt.Println(err)
	}

	return supplierRefs

}

func (s SupplierRefs) GetRows() []select2.Select2ItemAble {

	var select2ItemAble []select2.Select2ItemAble = make([]select2.Select2ItemAble, len(s.Rows))

	for i, row := range s.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (s SupplierRefs) GetTotal() int {

	return s.Total

}

func (s SupplierRef) FromJsJSONValue(jsvalue js.Value) select2.Select2ItemAble {

	var (
		supplierRef SupplierRef
		err         error
	)

	jsSupplierRefString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsSupplierRefString), &supplierRef); err != nil {
		fmt.Println(err)
	}

	return supplierRef

}

func (s SupplierRef) GetSelect2Id() int64 {

	return *s.SupplierRefID

}

func (s SupplierRef) GetSelect2Text() string {

	if s.SupplierRef != nil {
		return s.SupplierRefLabel
	}

	return ""

}
