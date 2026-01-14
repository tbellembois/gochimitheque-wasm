//go:build go1.24 && js && wasm

package entity

import (
	"encoding/json"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/bstable"
	. "github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/jquery"
	"github.com/tbellembois/gochimitheque-wasm/locales"
	"github.com/tbellembois/gochimitheque-wasm/select2"
	. "github.com/tbellembois/gochimitheque-wasm/types"
	"github.com/tbellembois/gochimitheque-wasm/validate"
)

func entity_common() {

	// validate
	validate.NewValidate(jquery.Jq("#entity"), &validate.ValidateConfig{
		ErrorClass: "alert alert-danger",
		Rules: map[string]validate.ValidateRule{
			"entity_name": {
				Required: js.FuncOf(func(this js.Value, args []js.Value) interface{} { return true }),
				// Remote: validate.ValidateRemote{
				// 	URL:        "",
				// 	Type:       "post",
				// 	BeforeSend: js.FuncOf(ValidateEntityNameBeforeSend),
				// },
			},
		},
		Messages: map[string]validate.ValidateMessage{
			"entity_name": {
				Required: locales.Translate("required_input", HTTPHeaderAcceptLanguage),
			},
		},
	}).Validate()

	window := js.Global()
	var keycloak js.Value
	keycloak = window.Get("keycloak")
	token := keycloak.Get("token").String()
	marshalToken, _ := json.Marshal(map[string]string{"Authorization": "Bearer " + token})

	// select2
	select2.NewSelect2(jquery.Jq("select#managers"), &select2.Select2Config{
		Placeholder:    locales.Translate("entity_manager_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(Person{})),
		Ajax: select2.Select2Ajax{
			URL:            BackProxyPath + "people_old",
			DataType:       "json",
			Data:           js.FuncOf(Select2ManagerAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(People{})),
			Headers:        js.Global().Get("JSON").Call("parse", string(marshalToken)),
		},
	}).Select2ify()

	// select2
	// select2.NewSelect2(jquery.Jq("select#ldapgroups"), &select2.Select2Config{
	// 	Placeholder:    locales.Translate("entity_ldap_group_placeholder", HTTPHeaderAcceptLanguage),
	// 	TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(LDAPEntry{})),
	// 	Ajax: select2.Select2Ajax{
	// 		URL:            ApplicationProxyPath + "ldapgroup",
	// 		DataType:       "json",
	// 		Data:           js.FuncOf(Select2LDAPGroupAjaxData),
	// 		ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(LDAPSearchResults{})),
	// 	},
	// }).Select2ify()

	jquery.Jq("#searchbar").Hide()
	jquery.Jq("#actions").Hide()

}

func Entity_createCallBack(this js.Value, args []js.Value) interface{} {

	entity_common()

	return nil

}

func Entity_listCallback(this js.Value, args []js.Value) interface{} {

	entity_common()

	bstable.NewBootstraptable(jquery.Jq("#Entity_table"), &bstable.BootstraptableParams{Ajax: "Entity_getTableData"})
	jquery.Jq("#Entity_table").On("load-success.bs.table", js.FuncOf(ShowIfAuthorizedActionButtons))

	return nil

}

func Entity_SaveCallback(args ...interface{}) {

	search := args[0].(string)

	bstable.NewBootstraptable(jquery.Jq("#Entity_table"), nil).ResetSearch(search)
	jquery.Jq("#Entity_table").On("load-success.bs.table", js.FuncOf(ShowIfAuthorizedActionButtons))

	entity_common()

}
