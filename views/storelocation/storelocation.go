//go:build go1.24 && js && wasm

package storelocation

import (
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/bstable"
	. "github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/jquery"
	"github.com/tbellembois/gochimitheque-wasm/locales"
	"github.com/tbellembois/gochimitheque-wasm/select2"
	. "github.com/tbellembois/gochimitheque-wasm/types"
	"github.com/tbellembois/gochimitheque-wasm/validate"
)

func storelocation_common() {

	// validate
	validate.NewValidate(jquery.Jq("#store_location"), &validate.ValidateConfig{
		Ignore:     "", // required to validate select2
		ErrorClass: "alert alert-danger",
		Rules: map[string]validate.ValidateRule{
			"store_location_name": {
				Required: js.FuncOf(func(this js.Value, args []js.Value) interface{} { return true }),
				Remote: validate.ValidateRemote{
					BeforeSend: js.FuncOf(func(this js.Value, args []js.Value) interface{} { return false }),
				},
			},
			"entity": {
				Required: js.FuncOf(func(this js.Value, args []js.Value) interface{} { return true }),
				Remote: validate.ValidateRemote{
					BeforeSend: js.FuncOf(func(this js.Value, args []js.Value) interface{} { return false }),
				},
			},
		},
		Messages: map[string]validate.ValidateMessage{
			"store_location_name": {
				Required: locales.Translate("required_input", HTTPHeaderAcceptLanguage),
			},
			"entity": {
				Required: locales.Translate("required_input", HTTPHeaderAcceptLanguage),
			},
		},
	}).Validate()

	// select2
	select2.NewSelect2(jquery.Jq("select#entity"), &select2.Select2Config{
		Placeholder:    locales.Translate("store_location_entity_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(Entity{})),
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "entities",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Entities{})),
		},
	}).Select2ify()
	jquery.Jq("select#entity").On("select2:select", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		select2StoreLocation := select2.NewSelect2(jquery.Jq("select#store_location"), nil)
		select2StoreLocation.Select2Clear()
		return nil
	}))

	select2.NewSelect2(jquery.Jq("select#store_location"), &select2.Select2Config{
		Placeholder:    locales.Translate("store_location_store_location_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(StoreLocation{})),
		AllowClear:     true,
		Tags:           true,
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "store_locations",
			DataType:       "json",
			Data:           js.FuncOf(Select2StoreLocationAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(StoreLocations{})),
		},
	}).Select2ify()

	jquery.Jq("#store_location_color").Object.Call("colorpicker")

	jquery.Jq("#searchbar").Hide()
	jquery.Jq("#actions").Hide()

}

func StoreLocation_createCallBack(this js.Value, args []js.Value) interface{} {

	storelocation_common()

	return nil

}

func StoreLocation_listCallback(this js.Value, args []js.Value) interface{} {

	storelocation_common()

	bstable.NewBootstraptable(jquery.Jq("#StoreLocation_table"), &bstable.BootstraptableParams{Ajax: "StoreLocation_getTableData"})
	jquery.Jq("#StoreLocation_table").On("load-success.bs.table", js.FuncOf(ShowIfAuthorizedActionButtons))

	return nil

}

func StoreLocation_SaveCallback(args ...interface{}) {

	search := args[0].(string)

	bstable.NewBootstraptable(jquery.Jq("#StoreLocation_table"), nil).ResetSearch(search)
	jquery.Jq("#StoreLocation_table").On("load-success.bs.table", js.FuncOf(ShowIfAuthorizedActionButtons))

	storelocation_common()

}
