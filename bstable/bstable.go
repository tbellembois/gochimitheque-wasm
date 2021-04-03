package bstable

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/ajax"
	"github.com/tbellembois/gochimitheque-wasm/jquery"
)

type Bootstraptable struct {
	jquery.Jquery
	params *BootstraptableParams
}

type BootstraptableParams struct {
	Ajax                 string      `json:"ajax"` // name of the Ajax function to call.
	FormatLoadingMessage interface{} `json:"formatLoadingMessage,omitempty"`
}

type BootstraptableRefreshQuery struct {
	// Query map[string]string `json:"query"`
	Query ajax.QueryFilter `json:"query"`
}

func structToMap(item interface{}) map[string]interface{} {

	res := map[string]interface{}{}
	if item == nil {
		return res
	}
	v := reflect.TypeOf(item)
	reflectValue := reflect.ValueOf(item)
	reflectValue = reflect.Indirect(reflectValue)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	var (
		omitempty bool
	)

	for i := 0; i < v.NumField(); i++ {

		tag := v.Field(i).Tag.Get("json")
		if strings.Contains(tag, "omitempty") {
			omitempty = true
			tag = strings.Replace(tag, ",omitempty", "", 1)
		} else {
			omitempty = false
		}

		field := reflectValue.Field(i).Interface()

		if tag != "" && tag != "-" {
			if v.Field(i).Type.Kind() == reflect.String {
				if !omitempty && len(field.(string)) > 0 {
					res[tag] = field
				}
			} else if v.Field(i).Type.Kind() == reflect.Bool {
				res[tag] = field
			} else if v.Field(i).Type.Kind() == reflect.Struct {
				res[tag] = structToMap(field)
			} else if v.Field(i).Type.Kind() == reflect.Map {

				m := field.(map[string]interface{})

				res2 := map[string]interface{}{}
				for k, v := range m {
					res2[k] = v
				}
				res[tag] = res2
			} else if v.Field(i).Type.Kind() == reflect.Interface {
				if !reflectValue.Field(i).IsNil() {
					res[tag] = field
				}
			} else {
				res[tag] = field
			}
		}
	}
	return res
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

	if params == nil {
		params = &BootstraptableParams{}
	}
	params.FormatLoadingMessage = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return "<span>&nbsp;</span>"
	})

	if params != nil {
		// jq.Object = jq.Object.Call("bootstrapTable", params.ToJsValue())
		jq.Object = jq.Object.Call("bootstrapTable", structToMap(params))
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

func (bt Bootstraptable) ResetSearch(search string) {

	bt.Jquery.Object.Call("bootstrapTable", "resetSearch", search)

}
