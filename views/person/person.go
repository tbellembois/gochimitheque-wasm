package person

import (
	"fmt"
	"strconv"
	"syscall/js"

	. "github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/locales"
	. "github.com/tbellembois/gochimitheque-wasm/types"
	"github.com/tbellembois/gochimitheque-wasm/utils"
	"github.com/tbellembois/gochimitheque-wasm/widgets"
)

func person_common() {

	// validate
	Jq("#person").Validate(ValidateConfig{
		ErrorClass: "alert alert-danger",
		Rules: map[string]ValidateRule{
			"person_email": {
				Required: js.FuncOf(func(this js.Value, args []js.Value) interface{} { return true }),
				Email:    true,
				Remote: ValidateRemote{
					URL:        "",
					Type:       "post",
					BeforeSend: js.FuncOf(ValidatePersonEmailBeforeSend),
				},
			},
		},
		Messages: map[string]ValidateMessage{
			"person_email": {
				Required: locales.Translate("required_input", HTTPHeaderAcceptLanguage),
			},
		},
	})

	// select2
	Jq("select#entities").Select2(Select2Config{
		Placeholder:    locales.Translate("person_entity_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2GenericTemplateResults(Entity{})),
		Ajax: Select2Ajax{
			URL:            ApplicationProxyPath + "entities",
			DataType:       "json",
			Data:           js.FuncOf(Select2GenericAjaxData),
			ProcessResults: js.FuncOf(Select2GenericAjaxProcessResults(Entities{})),
		},
	})
	Jq("select#entities").On("select2:unselecting", js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		select2EntitySelected := args[0].Get("params").Get("args").Get("data")
		entitySelected := Select2ItemFromJsJSONValue(select2EntitySelected)

		managedEntities := Jq("option.manageentities").Object
		for i := 0; i < managedEntities.Length(); i++ {
			entityId := managedEntities.Index(i).Call("getAttribute", "value").String()

			if entityId == entitySelected.Id {
				utils.DisplayErrorMessage(locales.Translate("person_can_not_remove_entity_manager", HTTPHeaderAcceptLanguage))
				args[0].Call("preventDefault")

				return nil
			}
		}

		Jq(fmt.Sprintf("#perm%s", entitySelected.Id)).Remove()

		return nil

	}))
	Jq("select#entities").On("select2:select", js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		var (
			entitySelectedId int
			err              error
		)

		select2EntitySelected := args[0].Get("params").Get("data")
		entitySelected := Select2ItemFromJsJSONValue(select2EntitySelected)

		if entitySelectedId, err = strconv.Atoi(entitySelected.Id); err != nil {
			fmt.Println(err)
			return nil
		}

		Jq("#permissions").Append(widgets.Permission(entitySelectedId, entitySelected.Text, false))

		return nil

	}))

	Jq("#search").Hide()
	Jq("#actions").Hide()

}

func Person_createCallBack(this js.Value, args []js.Value) interface{} {

	person_common()

	return nil

}

func Person_listCallback(this js.Value, args []js.Value) interface{} {

	person_common()

	Jq("#Person_table").Bootstraptable(&BootstraptableParams{Ajax: "Person_getTableData"})
	Jq("#Person_table").On("load-success.bs.table", js.FuncOf(ShowIfAuthorizedActionButtons))

	return nil

}

func Person_SaveCallback(args ...interface{}) {

	search := args[0].(string)

	Jq("#Person_table").Bootstraptable(nil).ResetSearch(search)
	Jq("#Person_table").On("load-success.bs.table", js.FuncOf(ShowIfAuthorizedActionButtons))

	person_common()

}
