package storelocation

import (
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/ajax"
	"github.com/tbellembois/gochimitheque-wasm/jquery"
	"github.com/tbellembois/gochimitheque-wasm/select2"
)

func Select2StoreLocationAjaxData(this js.Value, args []js.Value) interface{} {

	params := args[0]

	var entityId string

	select2Entity := select2.NewSelect2(jquery.Jq("select#entity"), nil)
	if len(select2Entity.Select2Data()) > 0 {
		select2ItemEntity := select2Entity.Select2Data()[0]
		if !select2ItemEntity.IsEmpty() {
			entityId = select2ItemEntity.Id
		}
	}

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
		Entity: entityId,
		Search: search,
		Offset: offset,
		Page:   page,
		Limit:  limit,
	}.ToJsValue()

}
