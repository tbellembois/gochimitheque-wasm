package entity

import (
	"syscall/js"

	. "github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/locales"
	. "github.com/tbellembois/gochimitheque-wasm/types"
)

func entity_common() {

	// validate
	Jq("#entity").Validate(ValidateConfig{
		ErrorClass: "alert alert-danger",
		Rules: map[string]ValidateRule{
			"entity_name": {
				Required: js.FuncOf(func(this js.Value, args []js.Value) interface{} { return true }),
				Remote: ValidateRemote{
					URL:        "",
					Type:       "post",
					BeforeSend: js.FuncOf(ValidateEntityNameBeforeSend),
				},
			},
		},
		Messages: map[string]ValidateMessage{
			"entity_name": {
				Required: locales.Translate("required_input", HTTPHeaderAcceptLanguage),
			},
		},
	})

	// select2
	Jq("select#managers").Select2(Select2Config{
		Placeholder:    locales.Translate("entity_manager_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2GenericTemplateResults(Person{})),
		Ajax: Select2Ajax{
			URL:            ApplicationProxyPath + "people",
			DataType:       "json",
			Data:           js.FuncOf(Select2ManagerAjaxData),
			ProcessResults: js.FuncOf(Select2GenericAjaxProcessResults(People{})),
		},
	})

	Jq("#search").Hide()
	Jq("#actions").Hide()

}

func Entity_createCallBack(this js.Value, args []js.Value) interface{} {

	entity_common()

	return nil

}

func Entity_listCallback(this js.Value, args []js.Value) interface{} {

	entity_common()

	Jq("#Entity_table").Bootstraptable(&BootstraptableParams{Ajax: "Entity_getTableData"})
	Jq("#Entity_table").On("load-success.bs.table", js.FuncOf(ShowIfAuthorizedActionButtons))

	return nil

}

func Entity_SaveCallback(args ...interface{}) {

	search := args[0].(string)

	Jq("#Entity_table").Bootstraptable(nil).ResetSearch(search)
	Jq("#Entity_table").On("load-success.bs.table", js.FuncOf(ShowIfAuthorizedActionButtons))

	entity_common()

}
