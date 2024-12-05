package jsutils

import (
	"fmt"
	"syscall/js"

	"honnef.co/go/js/dom/v2"

	"github.com/tbellembois/gochimitheque-wasm/ajax"
	"github.com/tbellembois/gochimitheque-wasm/bstable"
	"github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/jquery"
	"github.com/tbellembois/gochimitheque-wasm/select2"
	"github.com/tbellembois/gochimitheque-wasm/widgets"
	"github.com/tbellembois/gochimitheque-wasm/widgets/themes"
)

func ConsoleLog(s string) {
	js.Global().Get("console").Call("log", s)
}

func CloseQR(this js.Value, args []js.Value) interface{} {

	js.Global().Get("window").Get("qrScanner").Call("destroy")
	js.Global().Get("window").Set("qrScanner", nil)

	jquery.Jq("#video").AddClass("invisible")

	return nil

}

func RedirectTo(href string) {

	js.Global().Get("window").Get("location").Set("href", href)

}

func CreateJsHTMLElementFromString(s string) js.Value {

	t := js.Global().Get("document").Call("createElement", "template")
	t.Set("innerHTML", s)
	return t.Get("content").Get("firstChild")

}

func LoadContent(containerId string, viewName string, url string, callback func(args ...interface{}), args ...interface{}) {

	if viewName != "" {
		globals.CurrentView = viewName
	}

	done := func(data js.Value) {
		jquery.Jq(containerId).SetHtml(data.String())
		if callback != nil {
			callback(args...)
		}
	}
	fail := func(data js.Value) {
		DisplayGenericErrorMessage()
	}

	ajax.Ajax{
		URL:    url,
		Method: "get",
		Done:   done,
		Fail:   fail,
	}.Send()

}

func Search(this js.Value, args []js.Value) interface{} {

	if globals.CurrentView == "storage" {
		bstable.NewBootstraptable(jquery.Jq("#Storage_table"), nil).Refresh(nil)
	} else {
		bstable.NewBootstraptable(jquery.Jq("#Product_table"), nil).Refresh(nil)
	}

	return nil

}

func clearSearchForm() {

	if select2.NewSelect2(jquery.Jq("select#s_tags"), nil).Select2IsInitialized() {
		select2.NewSelect2(jquery.Jq("select#s_tags"), nil).Select2Clear()
	}
	if select2.NewSelect2(jquery.Jq("select#s_category"), nil).Select2IsInitialized() {
		select2.NewSelect2(jquery.Jq("select#s_category"), nil).Select2Clear()
	}
	if select2.NewSelect2(jquery.Jq("select#s_store_location"), nil).Select2IsInitialized() {
		select2.NewSelect2(jquery.Jq("select#s_store_location"), nil).Select2Clear()
	}
	if select2.NewSelect2(jquery.Jq("select#s_name"), nil).Select2IsInitialized() {
		select2.NewSelect2(jquery.Jq("select#s_name"), nil).Select2Clear()
	}
	if select2.NewSelect2(jquery.Jq("select#s_empirical_formula"), nil).Select2IsInitialized() {
		select2.NewSelect2(jquery.Jq("select#s_empirical_formula"), nil).Select2Clear()
	}
	if select2.NewSelect2(jquery.Jq("select#s_producer_ref"), nil).Select2IsInitialized() {
		select2.NewSelect2(jquery.Jq("select#s_producer_ref"), nil).Select2Clear()
	}
	if select2.NewSelect2(jquery.Jq("select#s_cas_number"), nil).Select2IsInitialized() {
		select2.NewSelect2(jquery.Jq("select#s_cas_number"), nil).Select2Clear()
	}
	if select2.NewSelect2(jquery.Jq("select#s_signal_word"), nil).Select2IsInitialized() {
		select2.NewSelect2(jquery.Jq("select#s_signal_word"), nil).Select2Clear()
	}
	if select2.NewSelect2(jquery.Jq("select#s_symbols"), nil).Select2IsInitialized() {
		select2.NewSelect2(jquery.Jq("select#s_symbols"), nil).Select2Clear()
	}
	if select2.NewSelect2(jquery.Jq("select#s_hazard_statements"), nil).Select2IsInitialized() {
		select2.NewSelect2(jquery.Jq("select#s_hazard_statements"), nil).Select2Clear()
	}
	if select2.NewSelect2(jquery.Jq("select#s_precautionary_statements"), nil).Select2IsInitialized() {
		select2.NewSelect2(jquery.Jq("select#s_precautionary_statements"), nil).Select2Clear()
	}

	jquery.Jq("#s_cas_number_cmr:checked").SetProp("checked", false)
	jquery.Jq("#s_borrowing:checked").SetProp("checked", false)
	jquery.Jq("#s_storage_to_destroy:checked").SetProp("checked", false)
	jquery.Jq("#searchshowchem").SetProp("checked", true)
	jquery.Jq("#searchshowbio").SetProp("checked", true)
	jquery.Jq("#searchshowconsu").SetProp("checked", true)

	jquery.Jq("#s_storage_batch_number").SetVal("")
	jquery.Jq("#s_storage_barecode").SetVal("")
	jquery.Jq("#s_custom_name_part_of").SetVal("")

	jquery.Jq("#s_storage_stock_button").SetInvisible()
	jquery.Jq("#stock").SetHtml("")

}

func ClearSearch(this js.Value, args []js.Value) interface{} {

	clearSearchForm()
	globals.BSTableQueryFilter.Clean()
	Search(js.Null(), nil)

	return nil

}

func ClearSearchExceptProduct(this js.Value, args []js.Value) interface{} {

	clearSearchForm()
	globals.BSTableQueryFilter.CleanExceptProduct()
	Search(js.Null(), nil)

	return nil

}

func DisplayGenericErrorMessage() {

	div := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Classes: []string{"animated", "fadeOutUp", "delay-2s", "fixed-top", "w-100", "p-3", "text-center", "alert", "alert-danger"},
			Visible: true,
			Attributes: map[string]string{
				"role":  "alert",
				"style": "z-index:2",
			},
		}})
	icon := widgets.NewIcon(widgets.IconAttributes{Icon: themes.NewMdiIcon(themes.MDI_BUG, themes.MDI_48PX)})
	span := widgets.NewSpan(widgets.SpanAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Classes: []string{"pl-sm-2"},
			Visible: true,
		},
	})

	div.AppendChild(icon)
	div.AppendChild(span)

	Win := dom.GetWindow()
	Doc := Win.Document()

	Doc.GetElementByID("message").SetInnerHTML("")
	Doc.GetElementByID("message").AppendChild(div)

	Win.SetTimeout(func() {
		Doc.GetElementByID("message").SetInnerHTML("")
	}, 5000)

}

func DisplaySuccessMessage(message string) {

	div := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Classes: []string{"animated", "fadeOutUp", "delay-2s", "fixed-top", "w-100", "p-3", "text-center", "alert", "alert-success"},
			Visible: true,
			Attributes: map[string]string{
				"role":  "alert",
				"style": "z-index:2",
			},
		}})
	icon := widgets.NewIcon(widgets.IconAttributes{Icon: themes.NewMdiIcon(themes.MDI_INFO, themes.MDI_24PX)})
	span := widgets.NewSpan(widgets.SpanAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Classes: []string{"pl-sm-2"},
			Visible: true,
		},
		Text: message,
	})

	div.AppendChild(icon)
	div.AppendChild(span)

	Win := dom.GetWindow()
	Doc := Win.Document()

	Doc.GetElementByID("message").SetInnerHTML("")
	Doc.GetElementByID("message").AppendChild(div)

	Win.SetTimeout(func() {
		Doc.GetElementByID("message").SetInnerHTML("")
	}, 5000)

}

func DisplayErrorMessage(message string) {

	div := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Classes: []string{"animated", "fadeOutUp", "delay-2s", "fixed-top", "w-100", "p-3", "text-center", "alert", "alert-danger"},
			Visible: true,
			Attributes: map[string]string{
				"role":  "alert",
				"style": "z-index:2",
			},
		}})
	icon := widgets.NewIcon(widgets.IconAttributes{Icon: themes.NewMdiIcon(themes.MDI_ERROR, themes.MDI_24PX)})
	span := widgets.NewSpan(widgets.SpanAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Classes: []string{"pl-sm-2"},
			Visible: true,
		},
		Text: message,
	})

	div.AppendChild(icon)
	div.AppendChild(span)

	Win := dom.GetWindow()
	Doc := Win.Document()

	Doc.GetElementByID("message").SetInnerHTML("")
	Doc.GetElementByID("message").AppendChild(div)

	Win.SetTimeout(func() {
		Doc.GetElementByID("message").SetInnerHTML("")
	}, 5000)

}

func DisplayMessageWrapper(this js.Value, args []js.Value) interface{} {

	DisplayMessage(args[0].String(), args[1].String())
	return nil

}

// DisplayMessage display fading messages at the
// top of the screen
func DisplayMessage(msgText string, msgType string) {

	Win := dom.GetWindow()
	Doc := Win.Document()

	d := Doc.CreateElement("div").(*dom.HTMLDivElement)
	s := Doc.CreateElement("span").(*dom.HTMLSpanElement)
	d.SetAttribute("role", "alert")
	d.SetAttribute("style", "z-index:2;")
	d.Class().SetString("animated fadeOutUp delay-2s fixed-top w-100 p-3 text-center alert alert-" + msgType)
	s.SetTextContent(msgText)
	d.AppendChild(s)

	Doc.GetElementByID("message").SetInnerHTML("")
	Doc.GetElementByID("message").AppendChild(d)

	Win.SetTimeout(func() {
		Doc.GetElementByID("message").SetInnerHTML("")
	}, 5000)

}

func CloseEdit(this js.Value, args []js.Value) interface{} {

	jquery.Jq("#list-collapse").Show()
	jquery.Jq("#edit-collapse").Hide()

	return nil

}

func DumpJsObject(object js.Value) {

	fmt.Println(js.Global().Get("JSON").Call("stringify", object).String())

}

func DisplayFilter(q ajax.QueryFilter) {

	//var isFilter bool

	jquery.Jq("#filter-item").SetHtml("")

	if q.Product != "" {
		// isFilter = true
		jquery.Jq("#filter-item").Append(widgets.FilterItem("storage_product_table_header", q.ProductFilterLabel))
		jquery.Jq("#removefilterstorage_product_table_header").On("click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			globals.BSTableQueryFilter.Lock()
			globals.BSTableQueryFilter.Product = ""
			globals.BSTableQueryFilter.Unlock()

			jquery.Jq("#s_storage_stock_button").SetInvisible()
			jquery.Jq("#stock").SetHtml("")

			Search(js.Null(), nil)
			return nil
		}))
	}
	if q.Storage != "" {
		// isFilter = true
		jquery.Jq("#filter-item").Append(widgets.FilterItem("storage", q.StorageFilterLabel))
		jquery.Jq("#removefilterstorage").On("click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			globals.BSTableQueryFilter.Lock()
			globals.BSTableQueryFilter.Storage = ""
			globals.BSTableQueryFilter.Unlock()
			Search(js.Null(), nil)
			return nil
		}))
	}
	if len(q.Storages) > 0 {
		// isFilter = true
		jquery.Jq("#filter-item").Append(widgets.FilterItem("storages", q.StoragesFilterLabel))
		jquery.Jq("#removefilterstorages").On("click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			globals.BSTableQueryFilter.Lock()
			globals.BSTableQueryFilter.Storages = nil
			globals.BSTableQueryFilter.Unlock()
			Search(js.Null(), nil)
			return nil
		}))
	}

	// if q.CustomNamePartOf != "" {
	// 	// isFilter = true
	// 	jquery.Jq("#filter-item").Append(widgets.FilterItem("s_custom_name_part_of", q.CustomNamePartOf))
	// 	jquery.Jq("#removefilters_custom_name_part_of").On("click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	// 		globals.BSTableQueryFilter.Lock()
	// 		globals.BSTableQueryFilter.CustomNamePartOf = ""
	// 		globals.BSTableQueryFilter.Unlock()
	// 		jquery.Jq("#s_custom_name_part_of").SetVal("")
	// 		Search(js.Null(), nil)
	// 		return nil
	// 	}))
	// }
	// if q.CasNumber != "" {
	// 	// isFilter = true
	// 	jquery.Jq("#filter-item").Append(widgets.FilterItem("s_cas_number", q.CasNumberFilterLabel))
	// 	jquery.Jq("#removefilters_cas_number").On("click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	// 		globals.BSTableQueryFilter.Lock()
	// 		globals.BSTableQueryFilter.CasNumber = ""
	// 		globals.BSTableQueryFilter.Unlock()
	// 		if select2.NewSelect2(jquery.Jq("select#s_cas_number"), nil).Select2IsInitialized() {
	// 			select2.NewSelect2(jquery.Jq("select#s_cas_number"), nil).Select2Clear()
	// 		}
	// 		Search(js.Null(), nil)
	// 		return nil
	// 	}))
	// }
	// if q.EmpiricalFormula != "" {
	// 	// isFilter = true
	// 	jquery.Jq("#filter-item").Append(widgets.FilterItem("s_empirical_formula", q.EmpiricalFormulaFilterLabel))
	// 	jquery.Jq("#removefilters_empirical_formula").On("click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	// 		globals.BSTableQueryFilter.Lock()
	// 		globals.BSTableQueryFilter.EmpiricalFormula = ""
	// 		globals.BSTableQueryFilter.Unlock()
	// 		if select2.NewSelect2(jquery.Jq("select#s_empirical_formula"), nil).Select2IsInitialized() {
	// 			select2.NewSelect2(jquery.Jq("select#s_empirical_formula"), nil).Select2Clear()
	// 		}
	// 		Search(js.Null(), nil)
	// 		return nil
	// 	}))
	// }
	// if q.StorageBarecode != "" {
	// 	// isFilter = true
	// 	jquery.Jq("#filter-item").Append(widgets.FilterItem("s_storage_barecode", q.StorageBarecode))
	// 	jquery.Jq("#removefilters_storage_barecode").On("click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	// 		globals.BSTableQueryFilter.Lock()
	// 		globals.BSTableQueryFilter.StorageBarecode = ""
	// 		globals.BSTableQueryFilter.Unlock()
	// 		jquery.Jq("#s_storage_barecode").SetVal("")
	// 		Search(js.Null(), nil)
	// 		return nil
	// 	}))
	// }
	// if q.StorageBatchNumber != "" {
	// 	// isFilter = true
	// 	jquery.Jq("#filter-item").Append(widgets.FilterItem("storage_batch_number_title", q.StorageBatchNumberFilterLabel))
	// 	jquery.Jq("#removefilterstorage_batch_number_title").On("click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	// 		globals.BSTableQueryFilter.Lock()
	// 		globals.BSTableQueryFilter.StorageBatchNumber = ""
	// 		globals.BSTableQueryFilter.Unlock()
	// 		jquery.Jq("#s_storage_batch_number").SetVal("")
	// 		Search(js.Null(), nil)
	// 		return nil
	// 	}))
	// }
	// if q.StoreLocation != "" {
	// 	// isFilter = true
	// 	jquery.Jq("#filter-item").Append(widgets.FilterItem("s_store_location", q.StoreLocationFilterLabel))
	// 	jquery.Jq("#removefilters_store_location").On("click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	// 		globals.BSTableQueryFilter.Lock()
	// 		globals.BSTableQueryFilter.StoreLocation = ""
	// 		globals.BSTableQueryFilter.Unlock()
	// 		if select2.NewSelect2(jquery.Jq("select#s_store_location"), nil).Select2IsInitialized() {
	// 			select2.NewSelect2(jquery.Jq("select#s_store_location"), nil).Select2Clear()
	// 		}
	// 		Search(js.Null(), nil)
	// 		return nil
	// 	}))
	// }
	// if q.Name != "" {
	// 	// isFilter = true
	// 	jquery.Jq("#filter-item").Append(widgets.FilterItem("s_name", q.NameFilterLabel))
	// 	jquery.Jq("#removefilters_name").On("click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	// 		globals.BSTableQueryFilter.Lock()
	// 		globals.BSTableQueryFilter.Name = ""
	// 		globals.BSTableQueryFilter.Unlock()
	// 		if select2.NewSelect2(jquery.Jq("select#s_name"), nil).Select2IsInitialized() {
	// 			select2.NewSelect2(jquery.Jq("select#s_name"), nil).Select2Clear()
	// 		}
	// 		Search(js.Null(), nil)
	// 		return nil
	// 	}))
	// }
	// if q.ProducerRef != "" {
	// 	// isFilter = true
	// 	jquery.Jq("#filter-item").Append(widgets.FilterItem("s_producer_ref", q.ProducerRefFilterLabel))
	// 	jquery.Jq("#removefilters_producer_ref").On("click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	// 		globals.BSTableQueryFilter.Lock()
	// 		globals.BSTableQueryFilter.ProducerRef = ""
	// 		globals.BSTableQueryFilter.Unlock()
	// 		if select2.NewSelect2(jquery.Jq("select#s_producer_ref"), nil).Select2IsInitialized() {
	// 			select2.NewSelect2(jquery.Jq("select#s_producer_ref"), nil).Select2Clear()
	// 		}
	// 		Search(js.Null(), nil)
	// 		return nil
	// 	}))
	// }

	if q.ProductBookmark {
		// isFilter = true
		//jquery.Jq("#filter-item").Append(locales.Translate("menu_bookmark", globals.HTTPHeaderAcceptLanguage))
		jquery.Jq("#filter-item").Append(widgets.FilterItem("menu_bookmark", ""))
		jquery.Jq("#removefiltermenu_bookmark").On("click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			globals.BSTableQueryFilter.Lock()
			globals.BSTableQueryFilter.ProductBookmark = false
			globals.BSTableQueryFilter.Unlock()
			Search(js.Null(), nil)
			return nil
		}))
	}
	// if q.ProductBookmark {
	// 	// isFilter = true
	// 	jquery.Jq("#filter-item").Append(widgets.FilterItem("", locales.Translate("menu_bookmark", globals.HTTPHeaderAcceptLanguage)))
	// }
	// if q.StorageArchive {
	// 	// isFilter = true
	// 	jquery.Jq("#filter-item").Append(widgets.FilterItem("", locales.Translate("archives", globals.HTTPHeaderAcceptLanguage)))
	// }
	// if q.StorageHistory {
	// 	// isFilter = true
	// 	jquery.Jq("#filter-item").Append(widgets.FilterItem("", locales.Translate("storage_history", globals.HTTPHeaderAcceptLanguage)))
	// }

	// if q.CasNumberCMR {
	// 	// isFilter = true
	// 	jquery.Jq("#filter-item").Append(widgets.FilterItem("s_cas_number_cmr", q.CasNumberCMRFilterLabel))
	// 	jquery.Jq("#removefilters_cas_number_cmr").On("click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	// 		globals.BSTableQueryFilter.Lock()
	// 		globals.BSTableQueryFilter.CasNumberCMR = false
	// 		globals.BSTableQueryFilter.Unlock()
	// 		jquery.Jq("#s_cas_number_cmr").SetProp("checked", false)
	// 		Search(js.Null(), nil)
	// 		return nil
	// 	}))
	// }
	// if q.Borrowing {
	// 	// isFilter = true
	// 	jquery.Jq("#filter-item").Append(widgets.FilterItem("s_borrowing", q.BorrowingFilterLabel))
	// 	jquery.Jq("#removefilters_borrowing").On("click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	// 		globals.BSTableQueryFilter.Lock()
	// 		globals.BSTableQueryFilter.Borrowing = false
	// 		globals.BSTableQueryFilter.Unlock()
	// 		jquery.Jq("#s_borrowing").SetProp("checked", false)
	// 		Search(js.Null(), nil)
	// 		return nil
	// 	}))
	// }
	// if q.StorageToDestroy {
	// 	// isFilter = true
	// 	jquery.Jq("#filter-item").Append(widgets.FilterItem("s_storage_to_destroy", q.StorageToDestroyFilterLabel))
	// 	jquery.Jq("#removefilters_storage_to_destroy").On("click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	// 		globals.BSTableQueryFilter.Lock()
	// 		globals.BSTableQueryFilter.StorageToDestroy = false
	// 		globals.BSTableQueryFilter.Unlock()
	// 		jquery.Jq("#s_storage_to_destroy").SetProp("checked", false)
	// 		Search(js.Null(), nil)
	// 		return nil
	// 	}))
	// }

	// if !isFilter {
	// 	jquery.Jq("#filter-item").Append(widgets.NewSpan(widgets.SpanAttributes{
	// 		BaseAttributes: widgets.BaseAttributes{
	// 			Visible: true,
	// 			Classes: []string{"iconlabel"},
	// 		},
	// 		Text: locales.Translate("no_filter", globals.HTTPHeaderAcceptLanguage),
	// 	}).OuterHTML())
	// }

}
