package entity

import (
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/ajax"
)

func Select2LDAPGroupAjaxData(this js.Value, args []js.Value) interface{} {

	params := args[0]

	search := ""
	if params.Get("term").Truthy() {
		search = params.Get("term").String()
	}

	return ajax.QueryFilter{
		Search: search,
	}.ToJsValue()

}

func Select2ManagerAjaxData(this js.Value, args []js.Value) interface{} {

	params := args[0]

	search := ""
	if params.Get("term").Truthy() {
		search = params.Get("term").String()
	}
	page := 1
	if params.Get("page").Truthy() {
		page = params.Get("page").Int()
	}
	offset := (page - 1) * 10
	limit := 10

	if offset < 0 {
		offset = 0
	}

	return ajax.QueryFilter{
		Search: search,
		Offset: offset,
		Page:   page,
		Limit:  limit,
	}.ToJsValue()

}
