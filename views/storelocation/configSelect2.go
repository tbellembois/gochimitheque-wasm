package storelocation

import (
	"syscall/js"

	. "github.com/tbellembois/gochimitheque-wasm/types"
)

// func Select2EntityAjaxProcessResults(this js.Value, args []js.Value) interface{} {

// 	data := args[0]
// 	params := args[1]
// 	page := 1
// 	if params.Get("page").Truthy() {
// 		page = params.Get("page").Int()
// 	}

// 	entities := EntitiesFromJsJSONValue(data)

// 	var select2Items []Select2Item
// 	for _, entity := range entities.Rows {
// 		select2Item := Select2Item{
// 			Id:   strconv.Itoa(entity.EntityID),
// 			Text: entity.EntityName,
// 		}
// 		select2Items = append(select2Items, select2Item)
// 	}

// 	select2Data := Select2Data{
// 		Results: select2Items,
// 		Pagination: Select2Pagination{
// 			More: (page * 10) < entities.Total,
// 		},
// 	}

// 	return select2Data.ToJsValue()

// }

// func Select2EntityAjaxData(this js.Value, args []js.Value) interface{} {

// 	params := args[0]

// 	search := ""
// 	if params.Get("term").Truthy() {
// 		search = params.Get("term").String()
// 	}
// 	page := 1
// 	if params.Get("page").Truthy() {
// 		page = params.Get("page").Int()
// 	}
// 	offset := (page - 1) * 10
// 	limit := 10

// 	if offset < 0 {
// 		offset = 0
// 	}

// 	return QueryFilter{
// 		Search: search,
// 		Offset: offset,
// 		Page:   page,
// 		Limit:  limit,
// 	}.ToJsValue()

// }

func Select2StoreLocationAjaxData(this js.Value, args []js.Value) interface{} {

	params := args[0]

	var entityId string

	if len(Jq("select#entity").Select2Data()) > 0 {
		select2ItemEntity := Jq("select#entity").Select2Data()[0]
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

	return QueryFilter{
		Entity: entityId,
		Search: search,
		Offset: offset,
		Page:   page,
		Limit:  limit,
	}.ToJsValue()

}
