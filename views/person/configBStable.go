package person

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"syscall/js"

	. "github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/locales"
	. "github.com/tbellembois/gochimitheque-wasm/types"
	"github.com/tbellembois/gochimitheque-wasm/utils"
	"github.com/tbellembois/gochimitheque-wasm/widgets"
	"github.com/tbellembois/gochimitheque-wasm/widgets/themes"
	"honnef.co/go/js/dom/v2"
)

func OperateEventsEdit(this js.Value, args []js.Value) interface{} {

	row := args[2]
	index := args[3].Int()
	person := Person{}.FromJsJSONValue(row).(Person)

	url := fmt.Sprintf("%speople/%d", ApplicationProxyPath, person.PersonId)
	method := "get"

	done := func(data js.Value) {

		var (
			person Person
			err    error
		)

		if err = json.Unmarshal([]byte(data.String()), &person); err != nil {
			fmt.Println(err)
		}

		FillInPersonForm(person, "edit-collapse")

		Jq("input#index").SetVal(index)
		Jq("#edit-collapse").Show()

	}
	fail := func(data js.Value) {

		utils.DisplayGenericErrorMessage()

	}

	Ajax{
		Method: method,
		URL:    url,
		Done:   done,
		Fail:   fail,
	}.Send()

	return nil

}

func OperateFormatter(this js.Value, args []js.Value) interface{} {

	row := args[1]
	person := Person{}.FromJsJSONValue(row).(Person)

	buttonEdit := widgets.NewBSButtonWithIcon(
		widgets.ButtonAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Id:         "edit" + strconv.Itoa(person.PersonId),
				Classes:    []string{"edit"},
				Visible:    false,
				Attributes: map[string]string{"pid": strconv.Itoa(person.PersonId)},
			},
			Title: locales.Translate("edit", HTTPHeaderAcceptLanguage),
		},
		widgets.IconAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible:    true,
				Attributes: map[string]string{"pid": strconv.Itoa(person.PersonId)},
			},
			Text: "",
			Icon: themes.NewMdiIcon(themes.MDI_EDIT, ""),
		},
		[]themes.BSClass{themes.BS_BTN, themes.BS_BNT_LINK},
	).OuterHTML()

	buttonDelete := widgets.NewBSButtonWithIcon(
		widgets.ButtonAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Id:         "delete" + strconv.Itoa(person.PersonId),
				Classes:    []string{"delete"},
				Visible:    false,
				Attributes: map[string]string{"pid": strconv.Itoa(person.PersonId)},
			},
			Title: locales.Translate("delete", HTTPHeaderAcceptLanguage),
		},
		widgets.IconAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible:    true,
				Attributes: map[string]string{"pid": strconv.Itoa(person.PersonId)},
			},
			Text: "",
			Icon: themes.NewMdiIcon(themes.MDI_DELETE, ""),
		},
		[]themes.BSClass{themes.BS_BTN, themes.BS_BNT_LINK},
	).OuterHTML()

	return buttonEdit + buttonDelete

}

func DataQueryParams(this js.Value, args []js.Value) interface{} {

	params := args[0]

	entity := URLParameters.Get("entity")
	if entity != "" {
		params.Set("entity", entity)
	}

	return params

}

// Get remote bootstrap table data - params object defined here:
// https://bootstrap-table.com/docs/api/table-options/#queryparams
func GetTableData(this js.Value, args []js.Value) interface{} {

	row := args[0]
	params := QueryParamsFromJsJSONValue(row)

	go func() {

		u := url.URL{Path: ApplicationProxyPath + "people"}
		u.RawQuery = params.Data.ToRawQuery()

		ajax := Ajax{
			URL:    u.String(),
			Method: "get",
			Done: func(data js.Value) {

				row.Call("success", js.ValueOf(js.Global().Get("JSON").Call("parse", data)))

			},
			Fail: func(jqXHR js.Value) {

				utils.DisplayGenericErrorMessage()

			},
		}

		ajax.Send()

	}()

	return nil

}

func ShowIfAuthorizedActionButtons(this js.Value, args []js.Value) interface{} {

	// Iterating other the button with the class "edit"
	// (we could choose "delete")
	// to retrieve once the entity id.
	buttons := dom.GetWindow().Document().GetElementsByTagName("button")
	for _, button := range buttons {
		if button.Class().Contains("edit") {
			personId := button.GetAttribute("pid")

			utils.HasPermission("people", personId, "put", func() {
				Jq("#edit" + personId).FadeIn()
			}, func() {
			})
			utils.HasPermission("people", personId, "delete", func() {
				Jq("#delete" + personId).FadeIn()
			}, func() {
			})
		}
	}

	return nil

}
