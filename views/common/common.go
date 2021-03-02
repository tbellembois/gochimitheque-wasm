package common

import (
	"fmt"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/bstable"
	"github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/jquery"
	"github.com/tbellembois/gochimitheque-wasm/jsutils"
	"github.com/tbellembois/gochimitheque-wasm/locales"
	"github.com/tbellembois/gochimitheque-wasm/views/product"
	"github.com/tbellembois/gochimitheque-wasm/views/storage"
	"github.com/tbellembois/gochimitheque-wasm/widgets"
)

func Export(this js.Value, args []js.Value) interface{} {

	globals.BSTableQueryFilter.Lock()
	globals.BSTableQueryFilter.QueryFilter.Export = true

	jquery.Jq("#export-body").SetHtml(widgets.NewSpan(widgets.SpanAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"mdi", "mdi-loading", "mdi-spin"},
		},
		Text: locales.Translate("export_progress", globals.HTTPHeaderAcceptLanguage),
	}).OuterHTML())

	jquery.Jq("#export").Show()
	jquery.Jq("button#export").SetProp("disabled", true)

	if globals.CurrentView != "storage" {
		bstable.NewBootstraptable(jquery.Jq("#Product_table"), nil).Refresh(nil)
	} else {
		bstable.NewBootstraptable(jquery.Jq("#Storage_table"), nil).Refresh(nil)
	}

	return nil

}

func SwitchProductStorageWrapper(this js.Value, args []js.Value) interface{} {

	storageCallbackWrapper := func(args ...interface{}) {
		storage.Storage_listCallback(js.Null(), nil)
	}
	productCallbackWrapper := func(args ...interface{}) {
		product.Product_listCallback(js.Null(), nil)
	}

	if globals.CurrentView != "storage" {
		jsutils.LoadContent("div#content", "storage", fmt.Sprintf("%sv/storages", globals.ApplicationProxyPath), storageCallbackWrapper, nil)
	} else {
		jsutils.LoadContent("div#content", "product", fmt.Sprintf("%sv/products", globals.ApplicationProxyPath), productCallbackWrapper, nil)
	}

	return nil

}
