package storelocation

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"syscall/js"

	"honnef.co/go/js/dom/v2"

	. "github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/locales"
	. "github.com/tbellembois/gochimitheque-wasm/types"
	"github.com/tbellembois/gochimitheque-wasm/utils"
	"github.com/tbellembois/gochimitheque-wasm/widgets"
	"github.com/tbellembois/gochimitheque-wasm/widgets/themes"
)

func OperateEventsDelete(this js.Value, args []js.Value) interface{} {

	row := args[2]
	storeLocation := StoreLocation{}.FromJsJSONValue(row).(StoreLocation)

	url := fmt.Sprintf("%sstorelocations/%d", ApplicationProxyPath, storeLocation.StoreLocationID.Int64)
	method := "delete"

	done := func(data js.Value) {

		utils.DisplaySuccessMessage(locales.Translate("storelocation_deleted_message", HTTPHeaderAcceptLanguage))
		Jq("#StoreLocation_table").Bootstraptable(nil).ResetSearch("")

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
	storeLocation := StoreLocation{}.FromJsJSONValue(row).(StoreLocation)

	url := fmt.Sprintf("%sstorelocations/%d", ApplicationProxyPath, storeLocation.StoreLocationID.Int64)
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
	storelocation := StoreLocation{}.FromJsJSONValue(row).(StoreLocation)

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
	storelocation := StoreLocation{}.FromJsJSONValue(row).(StoreLocation)

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
	storelocation := StoreLocation{}.FromJsJSONValue(row).(StoreLocation)

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

func StoreLocationFormatter(this js.Value, args []js.Value) interface{} {

	row := args[1]
	storelocation := StoreLocation{}.FromJsJSONValue(row).(StoreLocation)

	var storeLocationName string

	if storelocation.StoreLocation != nil {
		storeLocationName = storelocation.StoreLocation.StoreLocationName.String
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

		u := url.URL{Path: ApplicationProxyPath + "storelocations"}
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
			storeLocationId := button.GetAttribute("slid")

			utils.HasPermission("storelocations", storeLocationId, "put", func() {
				Jq("#edit" + storeLocationId).FadeIn()
			}, func() {
			})
			utils.HasPermission("storelocations", storeLocationId, "delete", func() {
				Jq("#delete" + storeLocationId).FadeIn()
			}, func() {
			})
		}
	}

	return nil

}
