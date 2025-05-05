//go:build go1.24 && js && wasm

package person

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/ajax"
	"github.com/tbellembois/gochimitheque-wasm/bstable"
	. "github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/jquery"
	"github.com/tbellembois/gochimitheque-wasm/jsutils"
	"github.com/tbellembois/gochimitheque-wasm/locales"
	. "github.com/tbellembois/gochimitheque-wasm/types"
	"github.com/tbellembois/gochimitheque-wasm/widgets"
	"github.com/tbellembois/gochimitheque-wasm/widgets/themes"
	"github.com/tbellembois/gochimitheque/models"
	"honnef.co/go/js/dom/v2"
)

func OperateEventsDelete(this js.Value, args []js.Value) interface{} {

	row := args[2]
	person := Person{Person: &models.Person{}}.FromJsJSONValue(row).(Person)

	jquery.Jq(fmt.Sprintf("button#delete%d", person.PersonID)).On("click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		url := fmt.Sprintf("%speople/%d", ApplicationProxyPath, person.PersonID)
		method := "delete"

		done := func(data js.Value) {

			jsutils.DisplaySuccessMessage(locales.Translate("person_deleted_message", HTTPHeaderAcceptLanguage))
			bstable.NewBootstraptable(jquery.Jq("#Person_table"), nil).ResetSearch("")
			bstable.NewBootstraptable(jquery.Jq("#Person_table"), nil).Refresh(nil)

		}
		fail := func(data js.Value) {

			jsutils.DisplayGenericErrorMessage()

		}

		ajax.Ajax{
			Method: method,
			URL:    url,
			Done:   done,
			Fail:   fail,
		}.Send()

		return nil

	}))

	buttonTitle := widgets.NewIcon(widgets.IconAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"iconlabel"},
		},
		Icon: themes.NewMdiIcon(themes.MDI_CONFIRM, ""),
		Text: locales.Translate("confirm", HTTPHeaderAcceptLanguage),
	})
	jquery.Jq(fmt.Sprintf("button#delete%d", person.PersonID)).SetHtml("")
	jquery.Jq(fmt.Sprintf("button#delete%d", person.PersonID)).Append(buttonTitle.OuterHTML())

	return nil

}

func OperateEventsEdit(this js.Value, args []js.Value) interface{} {

	row := args[2]
	index := args[3].Int()
	person := Person{Person: &models.Person{}}.FromJsJSONValue(row).(Person)

	url := fmt.Sprintf("%speople/%d", ApplicationProxyPath, person.PersonID)
	method := "get"

	done := func(data js.Value) {

		var (
			person Person
			err    error
		)

		if err = json.Unmarshal([]byte(data.String()), &person); err != nil {
			fmt.Println(err)
		}

		// js.Global().Get("console").Call("log", fmt.Sprintf("%#v", person))

		FillInPersonForm(person, "edit-collapse")

		jquery.Jq("input#index").SetVal(index)
		jquery.Jq("#edit-collapse").Show()

	}
	fail := func(data js.Value) {

		jsutils.DisplayGenericErrorMessage()

	}

	ajax.Ajax{
		Method: method,
		URL:    url,
		Done:   done,
		Fail:   fail,
	}.Send()

	return nil

}

func OperateFormatter(this js.Value, args []js.Value) interface{} {

	row := args[1]
	person := Person{Person: &models.Person{}}.FromJsJSONValue(row).(Person)

	buttonEdit := widgets.NewBSButtonWithIcon(
		widgets.ButtonAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Id:         "edit" + strconv.Itoa(person.PersonID),
				Classes:    []string{"edit"},
				Visible:    false,
				Attributes: map[string]string{"pid": strconv.Itoa(person.PersonID)},
			},
			Title: locales.Translate("edit", HTTPHeaderAcceptLanguage),
		},
		widgets.IconAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible:    true,
				Attributes: map[string]string{"pid": strconv.Itoa(person.PersonID)},
			},
			Text: "",
			Icon: themes.NewMdiIcon(themes.MDI_EDIT, ""),
		},
		[]themes.BSClass{themes.BS_BTN, themes.BS_BNT_LINK},
	).OuterHTML()

	buttonDelete := widgets.NewBSButtonWithIcon(
		widgets.ButtonAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Id:         "delete" + strconv.Itoa(person.PersonID),
				Classes:    []string{"delete"},
				Visible:    false,
				Attributes: map[string]string{"pid": strconv.Itoa(person.PersonID)},
			},
			Title: locales.Translate("delete", HTTPHeaderAcceptLanguage),
		},
		widgets.IconAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible:    true,
				Attributes: map[string]string{"pid": strconv.Itoa(person.PersonID)},
			},
			Text: "",
			Icon: themes.NewMdiIcon(themes.MDI_DELETE, ""),
		},
		[]themes.BSClass{themes.BS_BTN, themes.BS_BNT_LINK},
	).OuterHTML()

	return buttonEdit + buttonDelete

}

func MemberOfFormatter(this js.Value, args []js.Value) interface{} {

	row := args[1]
	person := Person{Person: &models.Person{}}.FromJsJSONValue(row).(Person)

	label := "<ul>"

	for _, entity := range person.Entities {
		label += "<li>" + entity.EntityName + "</li>"
	}
	label += "<ul>"

	span := widgets.NewSpan(widgets.SpanAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
		},
		Text: label,
	})

	return span.OuterHTML()
}

func ManagerOfFormatter(this js.Value, args []js.Value) interface{} {

	row := args[1]
	person := Person{Person: &models.Person{}}.FromJsJSONValue(row).(Person)

	label := "<ul>"

	for _, entity := range person.ManagedEntities {
		label += "<li>" + entity.EntityName + "</li>"
	}
	label += "<ul>"

	span := widgets.NewSpan(widgets.SpanAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
		},
		Text: label,
	})

	return span.OuterHTML()
}

func DataQueryParams(this js.Value, args []js.Value) interface{} {

	params := args[0]

	queryFilter := ajax.QueryFilterFromJsJSONValue(params)

	queryFilter.Entity = BSTableQueryFilter.Entity
	BSTableQueryFilter.Unlock()

	return queryFilter.ToJsValue()

}

// Get remote bootstrap table data - params object defined here:
// https://bootstrap-table.com/docs/api/table-options/#queryparams
func GetTableData(this js.Value, args []js.Value) interface{} {

	row := args[0]
	params := bstable.QueryParamsFromJsJSONValue(row)

	go func() {

		u := url.URL{Path: ApplicationProxyPath + "people"}
		u.RawQuery = params.Data.ToRawQuery()

		ajax := ajax.Ajax{
			URL:    u.String(),
			Method: "get",
			Done: func(data js.Value) {

				row.Call("success", js.ValueOf(js.Global().Get("JSON").Call("parse", data)))

			},
			Fail: func(jqXHR js.Value) {

				jsutils.DisplayGenericErrorMessage()

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

			jsutils.HasPermission("people", personId, "put", func() {
				jquery.Jq("#edit" + personId).FadeIn()
			}, func() {
			})
			jsutils.HasPermission("people", personId, "delete", func() {
				jquery.Jq("#delete" + personId).FadeIn()
			}, func() {
			})
		}
	}

	return nil

}
