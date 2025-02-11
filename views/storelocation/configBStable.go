package storelocation

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
	"github.com/tbellembois/gochimitheque-wasm/widgets"
	"github.com/tbellembois/gochimitheque-wasm/widgets/themes"
	"github.com/tbellembois/gochimitheque/models"
)

func OperateEventsDelete(this js.Value, args []js.Value) interface{} {

	row := args[2]
	storeLocation := StoreLocation{StoreLocation: &models.StoreLocation{}}.FromJsJSONValue(row).(StoreLocation)

	jquery.Jq(fmt.Sprintf("button#delete%d", storeLocation.StoreLocationID.Int64)).On("click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		url := fmt.Sprintf("%sstorelocations/%d", ApplicationProxyPath, storeLocation.StoreLocationID.Int64)
		method := "delete"

		done := func(data js.Value) {

			jsutils.DisplaySuccessMessage(locales.Translate("store_location_deleted_message", HTTPHeaderAcceptLanguage))
			bstable.NewBootstraptable(jquery.Jq("#StoreLocation_table"), nil).ResetSearch("")
			bstable.NewBootstraptable(jquery.Jq("#StoreLocation_table"), nil).Refresh(nil)

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
	jquery.Jq(fmt.Sprintf("button#delete%d", storeLocation.StoreLocationID.Int64)).SetHtml("")
	jquery.Jq(fmt.Sprintf("button#delete%d", storeLocation.StoreLocationID.Int64)).Append(buttonTitle.OuterHTML())

	return nil

}

func OperateEventsEdit(this js.Value, args []js.Value) interface{} {

	row := args[2]
	index := args[3].Int()
	storeLocation := StoreLocation{StoreLocation: &models.StoreLocation{}}.FromJsJSONValue(row).(StoreLocation)

	url := fmt.Sprintf("%sstorelocations?store_location=%d", ApplicationProxyPath, storeLocation.StoreLocationID.Int64)
	method := "get"

	done := func(data js.Value) {

		var (
			storeLocation StoreLocation
			err           error
		)

		if err = json.Unmarshal([]byte(data.String()), &storeLocation); err != nil {
			fmt.Println(err)
		}

		FillInStoreLocationForm(storeLocation, "edit-collapse")

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
	storelocation := StoreLocation{StoreLocation: &models.StoreLocation{}}.FromJsJSONValue(row).(StoreLocation)

	buttonEdit := widgets.NewBSButtonWithIcon(
		widgets.ButtonAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Id:         "edit" + strconv.Itoa(int(storelocation.StoreLocationID.Int64)),
				Classes:    []string{"edit"},
				Visible:    false,
				Attributes: map[string]string{"slid": strconv.Itoa(int(storelocation.StoreLocationID.Int64))},
			},
			Title: locales.Translate("edit", HTTPHeaderAcceptLanguage),
		},
		widgets.IconAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible:    true,
				Attributes: map[string]string{"slid": strconv.Itoa(int(storelocation.StoreLocationID.Int64))},
			},
			Text: "",
			Icon: themes.NewMdiIcon(themes.MDI_EDIT, ""),
		},
		[]themes.BSClass{themes.BS_BTN, themes.BS_BNT_LINK},
	).OuterHTML()

	buttonDelete := widgets.NewBSButtonWithIcon(
		widgets.ButtonAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Id:         "delete" + strconv.Itoa(int(storelocation.StoreLocationID.Int64)),
				Classes:    []string{"delete"},
				Visible:    false,
				Attributes: map[string]string{"slid": strconv.Itoa(int(storelocation.StoreLocationID.Int64))},
			},
			Title: locales.Translate("delete", HTTPHeaderAcceptLanguage),
		},
		widgets.IconAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible:    true,
				Attributes: map[string]string{"slid": strconv.Itoa(int(storelocation.StoreLocationID.Int64))},
			},
			Text: "",
			Icon: themes.NewMdiIcon(themes.MDI_DELETE, ""),
		},
		[]themes.BSClass{themes.BS_BTN, themes.BS_BNT_LINK},
	).OuterHTML()

	return buttonEdit + buttonDelete

}

func ColorFormatter(this js.Value, args []js.Value) interface{} {

	row := args[1]
	storelocation := StoreLocation{StoreLocation: &models.StoreLocation{}}.FromJsJSONValue(row).(StoreLocation)

	span := widgets.NewSpan(widgets.SpanAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
		},
		Text: "&nbsp;",
	})
	div := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Attributes: map[string]string{
				"style": "background-color:" + storelocation.StoreLocationColor.String,
			},
		}})

	div.AppendChild(span)

	return div.OuterHTML()

}

func CanStoreFormatter(this js.Value, args []js.Value) interface{} {

	row := args[1]
	storelocation := StoreLocation{StoreLocation: &models.StoreLocation{}}.FromJsJSONValue(row).(StoreLocation)

	var icon themes.IconFace

	if storelocation.StoreLocationCanStore.Bool {
		icon = themes.MDI_CHECK
	} else {
		icon = themes.MDI_NO_CHECK
	}

	i := widgets.NewIcon(widgets.IconAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
		},
		Icon: themes.NewMdiIcon(icon, themes.MDI_24PX),
	})

	return i.OuterHTML()

}

func StoreLocationFullPathFormatter(this js.Value, args []js.Value) interface{} {

	row := args[1]
	storelocation := StoreLocation{StoreLocation: &models.StoreLocation{}}.FromJsJSONValue(row).(StoreLocation)

	label := storelocation.StoreLocationFullPath + "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;[<i>"

	if storelocation.StoreLocationNbStorage != nil {
		label += " st:" + strconv.Itoa(int(*storelocation.StoreLocationNbStorage))
	}
	if storelocation.StoreLocationNbChildren != nil {
		label += " chil.: " + strconv.Itoa(int(*storelocation.StoreLocationNbChildren))
	}
	label += " </i>]"

	span := widgets.NewSpan(widgets.SpanAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
		},
		Text: label,
	})

	return span.OuterHTML()

}

func StoreLocationFormatter(this js.Value, args []js.Value) interface{} {

	row := args[1]
	storelocation := StoreLocation{StoreLocation: &models.StoreLocation{}}.FromJsJSONValue(row).(StoreLocation)

	var storeLocationName string

	if storelocation.StoreLocation.StoreLocation != nil {
		storeLocationName = storelocation.StoreLocation.StoreLocation.StoreLocationName.String
	}

	span := widgets.NewSpan(widgets.SpanAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
		},
		Text: storeLocationName,
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

		u := url.URL{Path: ApplicationProxyPath + "storelocations"}
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
			storeLocationId := button.GetAttribute("slid")

			jsutils.HasPermission("storelocations", storeLocationId, "put", func() {
				jquery.Jq("#edit" + storeLocationId).FadeIn()
			}, func() {
			})
			jsutils.HasPermission("storelocations", storeLocationId, "delete", func() {
				jquery.Jq("#delete" + storeLocationId).FadeIn()
			}, func() {
			})
		}
	}

	return nil

}
