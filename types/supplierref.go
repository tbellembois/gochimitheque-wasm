package types

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/barweiss/go-tuple"
	"github.com/tbellembois/gochimitheque-wasm/models"
	"github.com/tbellembois/gochimitheque-wasm/select2"
)

type Select2SupplierRefs struct {
	Rows  []*SupplierRef `json:"rows"`
	Total int            `json:"total"`
}

type SupplierRef struct {
	*models.SupplierRef
}

func (elems Select2SupplierRefs) GetRowConcreteTypeName() string {

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

func (elems Select2SupplierRefs) IsExactMatch() bool {

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
		Id   int    `json:"id"`
		Text string `json:"text"`
		Copy
	}{
		Id:   s.GetSelect2Id(),
		Text: s.GetSelect2Text(),
		Copy: (Copy)(s),
	})

}

func (Select2SupplierRefs) FromJsJSONValue(jsvalue js.Value) select2.Select2ResultAble {

	var (
		suppliersAjaxResponse tuple.T2[[]struct {
			MatchExactSearch bool   `json:"match_exact_search"`
			SupplierRefID    int64  `json:"supplierref_id"`
			SupplierRefLabel string `json:"supplierref_label"`
		}, int]
		select2SupplierRefs Select2SupplierRefs
		err                 error
	)

	jsSuppliersString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsSuppliersString), &suppliersAjaxResponse); err != nil {
		fmt.Println("(Select2SupplierRefs) FromJsJSONValue:" + err.Error())
	}

	for _, supplierref := range suppliersAjaxResponse.V1 {
		select2SupplierRefs.Rows = append(select2SupplierRefs.Rows, &SupplierRef{
			&models.SupplierRef{
				MatchExactSearch: supplierref.MatchExactSearch,
				SupplierRefID:    int(supplierref.SupplierRefID),
				SupplierRefLabel: supplierref.SupplierRefLabel,
			},
		})
	}

	select2SupplierRefs.Total = suppliersAjaxResponse.V2

	return select2SupplierRefs

}

func (s Select2SupplierRefs) GetRows() []select2.Select2ItemAble {

	var select2ItemAble []select2.Select2ItemAble = make([]select2.Select2ItemAble, len(s.Rows))

	for i, row := range s.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (s Select2SupplierRefs) GetTotal() int {

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

func (s SupplierRef) GetSelect2Id() int {

	return s.SupplierRefID

}

func (s SupplierRef) GetSelect2Text() string {

	if s.SupplierRef != nil {
		return s.SupplierRefLabel
	}

	return ""

}
