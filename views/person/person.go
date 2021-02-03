package person

import (
	"fmt"
	"strconv"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/bstable"
	. "github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/jquery"
	"github.com/tbellembois/gochimitheque-wasm/jsutils"
	"github.com/tbellembois/gochimitheque-wasm/locales"
	"github.com/tbellembois/gochimitheque-wasm/select2"
	. "github.com/tbellembois/gochimitheque-wasm/types"
	"github.com/tbellembois/gochimitheque-wasm/validate"
	"github.com/tbellembois/gochimitheque-wasm/widgets"
)

func person_common() {

	// validate
	validate.NewValidate(jquery.Jq("#person"), &validate.ValidateConfig{
		ErrorClass: "alert alert-danger",
		Rules: map[string]validate.ValidateRule{
			"person_email": {
				Required: js.FuncOf(func(this js.Value, args []js.Value) interface{} { return true }),
				Email:    true,
				Remote: validate.ValidateRemote{
					URL:        "",
					Type:       "post",
					BeforeSend: js.FuncOf(ValidatePersonEmailBeforeSend),
				},
			},
		},
		Messages: map[string]validate.ValidateMessage{
			"person_email": {
				Required: locales.Translate("required_input", HTTPHeaderAcceptLanguage),
			},
		},
	}).Validate()

	// select2
	select2.NewSelect2(jquery.Jq("select#entities"), &select2.Select2Config{
		Placeholder:    locales.Translate("person_entity_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(Entity{})),
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "entities",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Entities{})),
		},
	}).Select2ify()
	jquery.Jq("select#entities").On("select2:unselecting", js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		select2EntitySelected := args[0].Get("params").Get("args").Get("data")
		entitySelected := select2.Select2ItemFromJsJSONValue(select2EntitySelected)

		managedEntities := jquery.Jq("option.manageentities").Object
		for i := 0; i < managedEntities.Length(); i++ {
			entityId := managedEntities.Index(i).Call("getAttribute", "value").String()

			if entityId == entitySelected.Id {
				jsutils.DisplayErrorMessage(locales.Translate("person_can_not_remove_entity_manager", HTTPHeaderAcceptLanguage))
				args[0].Call("preventDefault")

				return nil
			}
		}

		jquery.Jq(fmt.Sprintf("#perm%s", entitySelected.Id)).Remove()

		return nil

	}))
	jquery.Jq("select#entities").On("select2:select", js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		var (
			entitySelectedId int
			err              error
		)

		select2EntitySelected := args[0].Get("params").Get("data")
		entitySelected := select2.Select2ItemFromJsJSONValue(select2EntitySelected)

		if entitySelectedId, err = strconv.Atoi(entitySelected.Id); err != nil {
			fmt.Println(err)
			return nil
		}

		jquery.Jq("#permissions").Append(widgets.Permission(entitySelectedId, entitySelected.Text, false))

		return nil

	}))

	jquery.Jq("#search").Hide()
	jquery.Jq("#actions").Hide()

}

func Person_createCallBack(this js.Value, args []js.Value) interface{} {

	person_common()

	return nil

}

func Person_listCallback(this js.Value, args []js.Value) interface{} {

	//person_common()

	bstable.NewBootstraptable(jquery.Jq("#Person_table"), &bstable.BootstraptableParams{Ajax: "Person_getTableData"})
	jquery.Jq("#Person_table").On("load-success.bs.table", js.FuncOf(ShowIfAuthorizedActionButtons))

	return nil

}

func Person_SaveCallback(args ...interface{}) {

	search := args[0].(string)

	bstable.NewBootstraptable(jquery.Jq("#Person_table"), nil).ResetSearch(search)
	jquery.Jq("#Person_table").On("load-success.bs.table", js.FuncOf(ShowIfAuthorizedActionButtons))

	//person_common()

}
