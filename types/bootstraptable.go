package types

import (
	"encoding/json"
	"fmt"
	"syscall/js"
)

type Bootstraptable struct {
	Jquery
}

type BootstraptableParams struct {
	Ajax string `json:"ajax"` // name of the Ajax function to call.
}

type BootstraptableRefreshQuery struct {
	// Query map[string]string `json:"query"`
	Query QueryFilter `json:"query"`
}

// QueryParamsFromJsJSONValue converts a JS JSON into a
// Go queryParams.
func QueryParamsFromJsJSONValue(jsvalue js.Value) QueryParams {

	var (
		queryParams QueryParams
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

func (jq Jquery) Bootstraptable(params *BootstraptableParams) Bootstraptable {

	if params != nil {
		jq.Object = jq.Object.Call("bootstrapTable", params.ToJsValue())
	} else {
		jq.Object = jq.Object.Call("bootstrapTable")
	}

	return Bootstraptable{Jquery: jq}

}

func (bt Bootstraptable) Refresh(params *BootstraptableRefreshQuery) {

	if params != nil {
		bt.Jquery.Object.Call("bootstrapTable", "refresh", params.ToJsValue())
	} else {
		bt.Jquery.Object.Call("bootstrapTable", "refresh")
	}

}

func (bt Bootstraptable) RemoveAll() {

	bt.Jquery.Object.Call("bootstrapTable", "removeAll")

}

func (bt Bootstraptable) ResetSearch(search string) {

	bt.Jquery.Object.Call("bootstrapTable", "resetSearch", search)

}
