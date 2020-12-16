package search

import (
	"syscall/js"

	. "github.com/tbellembois/gochimitheque-wasm/globals"
	. "github.com/tbellembois/gochimitheque-wasm/types"
)

func Search(this js.Value, args []js.Value) interface{} {

	if CurrentView == "storage" {
		Jq("#Storage_table").Bootstraptable(nil).Refresh(nil)
	} else {
		Jq("#Product_table").Bootstraptable(nil).Refresh(nil)
	}

	return nil

}

func ClearSearch(this js.Value, args []js.Value) interface{} {

	if Jq("select#s_storelocation").Select2IsInitialized() {
		Jq("select#s_storelocation").Select2Clear()
	}
	if Jq("select#s_name").Select2IsInitialized() {
		Jq("select#s_name").Select2Clear()
	}
	if Jq("select#s_empiricalformula").Select2IsInitialized() {
		Jq("select#s_empiricalformula").Select2Clear()
	}
	if Jq("select#s_casnumber").Select2IsInitialized() {
		Jq("select#s_casnumber").Select2Clear()
	}
	if Jq("select#s_signalword").Select2IsInitialized() {
		Jq("select#s_signalword").Select2Clear()
	}
	if Jq("select#s_symbols").Select2IsInitialized() {
		Jq("select#s_symbols").Select2Clear()
	}
	if Jq("select#s_hazardstatements").Select2IsInitialized() {
		Jq("select#s_hazardstatements").Select2Clear()
	}
	if Jq("select#s_precautionarystatements").Select2IsInitialized() {
		Jq("select#s_precautionarystatements").Select2Clear()
	}

	if Jq("#s_casnumber_cmr:checked").Object.Length() > 0 {
		Jq("#s_casnumber_cmr:checked").SetProp("checked", false)
	}
	if Jq("#s_borrowing:checked").Object.Length() > 0 {
		Jq("#s_borrowing:checked").SetProp("checked", false)
	}
	if Jq("#s_storage_to_destroy:checked").Object.Length() > 0 {
		Jq("#s_storage_to_destroy:checked").SetProp("checked", false)
	}

	Jq("#s_storage_barecode").SetVal("")
	Jq("#s_custom_name_part_of").SetVal("")

	BSTableQueryFilter.Clean()

	Search(js.Null(), nil)

	return nil

}
