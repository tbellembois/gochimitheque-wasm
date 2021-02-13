package entity

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"syscall/js"

	"honnef.co/go/js/dom/v2"

	"github.com/tbellembois/gochimitheque-wasm/ajax"
	"github.com/tbellembois/gochimitheque-wasm/bstable"
	. "github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/jquery"
	"github.com/tbellembois/gochimitheque-wasm/jsutils"
	"github.com/tbellembois/gochimitheque-wasm/locales"
	. "github.com/tbellembois/gochimitheque-wasm/types"
	"github.com/tbellembois/gochimitheque-wasm/views/person"
	"github.com/tbellembois/gochimitheque-wasm/views/storelocation"
	"github.com/tbellembois/gochimitheque-wasm/widgets"
	"github.com/tbellembois/gochimitheque-wasm/widgets/themes"
)

func OperateEventsStorelocations(this js.Value, args []js.Value) interface{} {

	storelocationCallbackWrapper := func(args ...interface{}) {
		storelocation.StoreLocation_listCallback(js.Null(), nil)
	}

	row := args[2]
	entity := Entity{}.FromJsJSONValue(row).(Entity)

	BSTableQueryFilter.Lock()
	BSTableQueryFilter.QueryFilter.Entity = strconv.Itoa(entity.EntityID)

	href := fmt.Sprintf("%sv/storelocations", ApplicationProxyPath)
	jsutils.LoadContent("div#content", "storelocation", href, storelocationCallbackWrapper)

	return nil

}

func OperateEventsMembers(this js.Value, args []js.Value) interface{} {

	personCallbackWrapper := func(args ...interface{}) {
		person.Person_listCallback(js.Null(), nil)
	}

	row := args[2]
	entity := Entity{}.FromJsJSONValue(row).(Entity)

	BSTableQueryFilter.Lock()
	BSTableQueryFilter.QueryFilter.Entity = strconv.Itoa(entity.EntityID)

	href := fmt.Sprintf("%sv/people", ApplicationProxyPath)
	jsutils.LoadContent("div#content", "person", href, personCallbackWrapper)

	return nil

}

func OperateEventsDelete(this js.Value, args []js.Value) interface{} {

	row := args[2]
	entity := Entity{}.FromJsJSONValue(row).(Entity)

	jquery.Jq(fmt.Sprintf("button#delete%d", entity.EntityID)).On("click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		url := fmt.Sprintf("%sentities/%d", ApplicationProxyPath, entity.EntityID)
		method := "delete"

		done := func(data js.Value) {

			jsutils.DisplaySuccessMessage(locales.Translate("entity_deleted_message", HTTPHeaderAcceptLanguage))
			bstable.NewBootstraptable(jquery.Jq("#Entity_table"), nil).ResetSearch("")
			bstable.NewBootstraptable(jquery.Jq("#Entity_table"), nil).Refresh(nil)

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
	jquery.Jq(fmt.Sprintf("button#delete%d", entity.EntityID)).SetHtml("")
	jquery.Jq(fmt.Sprintf("button#delete%d", entity.EntityID)).Append(buttonTitle.OuterHTML())

	return nil

}

func OperateEventsEdit(this js.Value, args []js.Value) interface{} {

	row := args[2]
	index := args[3].Int()
	entity := Entity{}.FromJsJSONValue(row).(Entity)

	url := fmt.Sprintf("%sentities/%d", ApplicationProxyPath, entity.EntityID)
	method := "get"

	done := func(data js.Value) {

		var (
			entity Entity
			err    error
		)

		if err = json.Unmarshal([]byte(data.String()), &entity); err != nil {
			fmt.Println(err)
		}

		FillInEntityForm(entity, "edit-collapse")

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
	entity := Entity{}.FromJsJSONValue(row).(Entity)

	buttonStorelocations := widgets.NewBSButtonWithIcon(
		widgets.ButtonAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Id:         "storelocations" + strconv.Itoa(entity.EntityID),
				Classes:    []string{"storelocations"},
				Visible:    false,
				Attributes: map[string]string{"eid": strconv.Itoa(entity.EntityID)},
			},
			Title: locales.Translate("storelocations", HTTPHeaderAcceptLanguage),
		},
		widgets.IconAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible:    true,
				Attributes: map[string]string{"eid": strconv.Itoa(entity.EntityID)},
			},
			Text: strconv.Itoa(entity.EntitySLC),
			Icon: themes.NewMdiIcon(themes.MDI_STORELOCATION, ""),
		},
		[]themes.BSClass{themes.BS_BTN, themes.BS_BNT_LINK},
	).OuterHTML()

	buttonMembers := widgets.NewBSButtonWithIcon(
		widgets.ButtonAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Id:         "members" + strconv.Itoa(entity.EntityID),
				Classes:    []string{"members"},
				Visible:    false,
				Attributes: map[string]string{"eid": strconv.Itoa(entity.EntityID)},
			},
			Title: locales.Translate("members", HTTPHeaderAcceptLanguage),
		},
		widgets.IconAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible:    true,
				Attributes: map[string]string{"eid": strconv.Itoa(entity.EntityID)},
			},
			Text: strconv.Itoa(entity.EntityPC),
			Icon: themes.NewMdiIcon(themes.MDI_PEOPLE, ""),
		},
		[]themes.BSClass{themes.BS_BTN, themes.BS_BNT_LINK},
	).OuterHTML()

	buttonEdit := widgets.NewBSButtonWithIcon(
		widgets.ButtonAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Id:         "edit" + strconv.Itoa(entity.EntityID),
				Classes:    []string{"edit"},
				Visible:    false,
				Attributes: map[string]string{"eid": strconv.Itoa(entity.EntityID)},
			},
			Title: locales.Translate("edit", HTTPHeaderAcceptLanguage),
		},
		widgets.IconAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible:    true,
				Attributes: map[string]string{"eid": strconv.Itoa(entity.EntityID)},
			},
			Text: "",
			Icon: themes.NewMdiIcon(themes.MDI_EDIT, ""),
		},
		[]themes.BSClass{themes.BS_BTN, themes.BS_BNT_LINK},
	).OuterHTML()

	buttonDelete := widgets.NewBSButtonWithIcon(
		widgets.ButtonAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Id:         "delete" + strconv.Itoa(entity.EntityID),
				Classes:    []string{"delete"},
				Visible:    false,
				Attributes: map[string]string{"eid": strconv.Itoa(entity.EntityID)},
			},
			Title: locales.Translate("delete", HTTPHeaderAcceptLanguage),
		},
		widgets.IconAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible:    true,
				Attributes: map[string]string{"eid": strconv.Itoa(entity.EntityID)},
			},
			Text: "",
			Icon: themes.NewMdiIcon(themes.MDI_DELETE, ""),
		},
		[]themes.BSClass{themes.BS_BTN, themes.BS_BNT_LINK},
	).OuterHTML()

	return buttonStorelocations + buttonMembers + buttonEdit + buttonDelete

}

func ManagersFormatter(this js.Value, args []js.Value) interface{} {

	row := args[1]
	entity := Entity{}.FromJsJSONValue(row).(Entity)

	ul := widgets.NewUl(widgets.UlAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
		}})

	if entity.Managers != nil {
		for _, manager := range entity.Managers {
			li := widgets.NewLi(widgets.LiAttributes{
				BaseAttributes: widgets.BaseAttributes{Visible: true},
				Text:           manager.PersonEmail,
			})
			ul.AppendChild(li)
		}
	}

	return ul.OuterHTML()

}

// Get remote bootstrap table data - params object defined here:
// https://bootstrap-table.com/docs/api/table-options/#queryparams
func GetTableData(this js.Value, args []js.Value) interface{} {

	row := args[0]
	params := bstable.QueryParamsFromJsJSONValue(row)

	go func() {

		u := url.URL{Path: ApplicationProxyPath + "entities"}
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

	// Iterating other the button with the class "storelocation"
	// (we could choose "members" or "delete")
	// to retrieve once the entity id.
	buttons := dom.GetWindow().Document().GetElementsByTagName("button")
	for _, button := range buttons {
		if button.Class().Contains("storelocations") {
			entityId := button.GetAttribute("eid")

			jsutils.HasPermission("storelocations", entityId, "get", func() {
				jquery.Jq("#storelocations" + entityId).FadeIn()
			}, func() {
			})
			jsutils.HasPermission("people", entityId, "get", func() {
				jquery.Jq("#members" + entityId).FadeIn()
			}, func() {
			})
			jsutils.HasPermission("entities", entityId, "put", func() {
				jquery.Jq("#edit" + entityId).FadeIn()
			}, func() {
			})
			jsutils.HasPermission("entities", entityId, "delete", func() {
				jquery.Jq("#delete" + entityId).FadeIn()
			}, func() {
			})
		}
	}

	return nil

}
