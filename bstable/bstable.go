package bstable

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/ajax"
	"github.com/tbellembois/gochimitheque-wasm/jquery"
)

type Bootstraptable struct {
	jquery.Jquery
	params *BootstraptableParams
}

type BootstraptableParams struct {
	Ajax string `json:"ajax"` // name of the Ajax function to call.
}

type BootstraptableRefreshQuery struct {
	// Query map[string]string `json:"query"`
	Query ajax.QueryFilter `json:"query"`
}

// QueryParamsFromJsJSONValue converts a JS JSON into a
// Go queryParams.
func QueryParamsFromJsJSONValue(jsvalue js.Value) ajax.QueryParams {

	var (
		queryParams ajax.QueryParams
		err         error
	)

	jsQueryParamsString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsQueryParamsString), &queryParams); err != nil {
		fmt.Println(err)
	}

	return queryParams

}

// ToJsValue converts a Go BootstraptableRefreshQuery
// into a JS JSON.
func (p BootstraptableRefreshQuery) ToJsValue() js.Value {

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

// ToJsValue converts a Go BootstraptableParams
// into a JS JSON.
func (p BootstraptableParams) ToJsValue() js.Value {

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

func NewBootstraptable(jq jquery.Jquery, params *BootstraptableParams) Bootstraptable {

	if params != nil {
		jq.Object = jq.Object.Call("bootstrapTable", params.ToJsValue())
	} else {
		jq.Object = jq.Object.Call("bootstrapTable")
	}

	return Bootstraptable{
		Jquery: jq,
		params: params,
	}

}

func (bt Bootstraptable) Refresh(params *BootstraptableRefreshQuery) {

	if params != nil {
		bt.Jquery.Object.Call("bootstrapTable", "refresh", params.ToJsValue())
	} else {
		bt.Jquery.Object.Call("bootstrapTable", "refresh")
	}

}

func (bt Bootstraptable) TotalRows() int {

	return bt.Jquery.Object.Call("bootstrapTable", "getOptions").Get("totalRows").Int()

}

type nullData []interface{}

func (n nullData) ToJsValue() js.Value {

	var (
		marshalJson []byte
		err         error
	)

	if marshalJson, err = json.Marshal(n); err != nil {
		fmt.Println(err)
		return js.Null()
	}

	return js.Global().Get("JSON").Call("parse", string(marshalJson))

}

func (bt Bootstraptable) RemoveAll() {

	//bt.Jquery.Object.Call("bootstrapTable", "removeAll")

	jqSave := bt.Jquery
	var nullData nullData
	nullData = make([]interface{}, 1)

	jqSave.Object.Call("bootstrapTable", "destroy")
	jqSave.Object.Call("bootstrapTable", nullData.ToJsValue())

}

func (bt Bootstraptable) ResetSearch(search string) {

	bt.Jquery.Object.Call("bootstrapTable", "resetSearch", search)

}

func (bt Bootstraptable) HideLoading() {

	bt.Jquery.Object.Call("bootstrapTable", "hideLoading")

}
