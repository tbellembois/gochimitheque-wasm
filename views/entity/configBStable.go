package entity

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

func OperateEventsStorelocations(this js.Value, args []js.Value) interface{} {

	row := args[2]
	entity := Entity{}.FromJsJSONValue(row).(Entity)

	href := fmt.Sprintf("%sv/storelocations?entity=%d", ApplicationProxyPath, entity.EntityID)
	utils.RedirectTo(href)

	return nil

}

func OperateEventsMembers(this js.Value, args []js.Value) interface{} {

	row := args[2]
	entity := Entity{}.FromJsJSONValue(row).(Entity)

	href := fmt.Sprintf("%sv/people?entity=%d", ApplicationProxyPath, entity.EntityID)
	utils.RedirectTo(href)

	return nil

}

func OperateEventsDelete(this js.Value, args []js.Value) interface{} {

	row := args[2]
	entity := Entity{}.FromJsJSONValue(row).(Entity)

	url := fmt.Sprintf("%sentities/%d", ApplicationProxyPath, entity.EntityID)
	method := "delete"

	done := func(data js.Value) {

		utils.DisplaySuccessMessage(locales.Translate("entity_deleted_message", HTTPHeaderAcceptLanguage))
		Jq("#Entity_table").Bootstraptable(nil).ResetSearch("")

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
	params := QueryParamsFromJsJSONValue(row)

	go func() {

		u := url.URL{Path: ApplicationProxyPath + "entities"}
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

	// Iterating other the button with the class "storelocation"
	// (we could choose "members" or "delete")
	// to retrieve once the entity id.
	buttons := dom.GetWindow().Document().GetElementsByTagName("button")
	for _, button := range buttons {
		if button.Class().Contains("storelocations") {
			entityId := button.GetAttribute("eid")

			utils.HasPermission("storelocations", entityId, "get", func() {
				Jq("#storelocations" + entityId).FadeIn()
			}, func() {
			})
			utils.HasPermission("people", entityId, "get", func() {
				Jq("#members" + entityId).FadeIn()
			}, func() {
			})
			utils.HasPermission("entities", entityId, "put", func() {
				Jq("#edit" + entityId).FadeIn()
			}, func() {
			})
			utils.HasPermission("entities", entityId, "delete", func() {
				Jq("#delete" + entityId).FadeIn()
			}, func() {
			})
		}
	}

	return nil

}
