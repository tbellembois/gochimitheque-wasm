package storelocation

import (
	"syscall/js"

	. "github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/locales"
	. "github.com/tbellembois/gochimitheque-wasm/types"
)

func storelocation_common() {

	// validate
	Jq("#storelocation").Validate(ValidateConfig{
		Ignore:     "", // required to validate select2
		ErrorClass: "alert alert-danger",
		Rules: map[string]ValidateRule{
			"storelocation_name": {
				Required: js.FuncOf(func(this js.Value, args []js.Value) interface{} { return true }),
				Remote: ValidateRemote{
					BeforeSend: js.FuncOf(func(this js.Value, args []js.Value) interface{} { return false }),
				},
			},
			"entity": {
				Required: js.FuncOf(func(this js.Value, args []js.Value) interface{} { return true }),
				Remote: ValidateRemote{
					BeforeSend: js.FuncOf(func(this js.Value, args []js.Value) interface{} { return false }),
				},
			},
		},
		Messages: map[string]ValidateMessage{
			"storelocation_name": {
				Required: locales.Translate("required_input", HTTPHeaderAcceptLanguage),
			},
			"entity": {
				Required: locales.Translate("required_input", HTTPHeaderAcceptLanguage),
			},
		},
	})

	// select2
	Jq("select#entity").Select2(Select2Config{
		Placeholder:    locales.Translate("storelocation_entity_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2GenericTemplateResults(Entity{})),
		Ajax: Select2Ajax{
			URL:            ApplicationProxyPath + "entities",
			DataType:       "json",
			Data:           js.FuncOf(Select2GenericAjaxData),
			ProcessResults: js.FuncOf(Select2GenericAjaxProcessResults(Entities{})),
		},
	})
	Jq("select#entity").On("select2:select", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		Jq("select#storelocation").Select2Clear()
		return nil
	}))

	Jq("select#storelocation").Select2(Select2Config{
		Placeholder:    locales.Translate("storelocation_storelocation_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2GenericTemplateResults(StoreLocation{})),
		Ajax: Select2Ajax{
			URL:            ApplicationProxyPath + "storelocations",
			DataType:       "json",
			Data:           js.FuncOf(Select2StoreLocationAjaxData),
			ProcessResults: js.FuncOf(Select2GenericAjaxProcessResults(StoreLocations{})),
		},
	})

	Jq("#storelocation_color").Object.Call("colorpicker")

	Jq("#search").Hide()
	Jq("#actions").Hide()

}

func StoreLocation_createCallBack(this js.Value, args []js.Value) interface{} {

	storelocation_common()

	return nil

}

func StoreLocation_listCallback(this js.Value, args []js.Value) interface{} {

	storelocation_common()

	Jq("#StoreLocation_table").Bootstraptable(&BootstraptableParams{Ajax: "StoreLocation_getTableData"})
	Jq("#StoreLocation_table").On("load-success.bs.table", js.FuncOf(ShowIfAuthorizedActionButtons))

	return nil

}

func StoreLocation_SaveCallback(args ...interface{}) {

	search := args[0].(string)

	Jq("#StoreLocation_table").Bootstraptable(nil).ResetSearch(search)
	Jq("#StoreLocation_table").On("load-success.bs.table", js.FuncOf(ShowIfAuthorizedActionButtons))

	storelocation_common()

}
