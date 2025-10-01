//go:build go1.24 && js && wasm

package search

import (
	"fmt"
	"syscall/js"

	"honnef.co/go/js/dom/v2"

	. "github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/jsutils"
	. "github.com/tbellembois/gochimitheque-wasm/types"
	"github.com/tbellembois/gochimitheque-wasm/widgets"
	"github.com/tbellembois/gochimitheque-wasm/widgets/themes"
	"github.com/tbellembois/gochimitheque/models"
)

// TODO: move me
func Select2SymbolTemplateResults(this js.Value, args []js.Value) interface{} {

	data := args[0]
	symbol := Symbol{Symbol: &models.Symbol{}}.FromJsJSONValue(data).(Symbol)

	if symbol.Symbol == nil {
		return jsutils.CreateJsHTMLElementFromString(widgets.NewDiv(widgets.DivAttributes{}).OuterHTML())
	}

	if symbol.Symbol == nil {
		return jsutils.CreateJsHTMLElementFromString(widgets.NewDiv(widgets.DivAttributes{}).OuterHTML())
	}

	image := widgets.NewImg(widgets.ImgAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
		},
		Height: "30",
		Width:  "30",
		Src:    fmt.Sprintf("%sstatic/img/%s.svg", ApplicationProxyPath, symbol.SymbolLabel),
		Alt:    symbol.SymbolLabel,
		Title:  symbol.SymbolLabel,
	})
	spanLabel := widgets.NewSpan(widgets.SpanAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
		},
		Text: symbol.SymbolLabel,
	})
	d := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
		},
	})

	d.AppendChild(image)
	d.AppendChild(spanLabel)

	return jsutils.CreateJsHTMLElementFromString(d.OuterHTML())

}

// TODO: move me
func Select2StoreLocationTemplateResults(this js.Value, args []js.Value) interface{} {

	var (
		iconCanStore dom.Node
	)

	data := args[0]

	storelocation := StoreLocation{StoreLocation: &models.StoreLocation{}}.FromJsJSONValue(data).(StoreLocation)

	if storelocation.StoreLocation == nil {
		return jsutils.CreateJsHTMLElementFromString(widgets.NewDiv(widgets.DivAttributes{}).OuterHTML())
	}

	iconColor := widgets.NewIcon(widgets.IconAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Attributes: map[string]string{
				"style": fmt.Sprintf("color: %s", *storelocation.StoreLocationColor),
			},
		},
		Icon: themes.NewMdiIcon(themes.MDI_COLOR, themes.MDI_24PX),
	})

	spanLabel := widgets.NewSpan(widgets.SpanAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
		},
		Text: storelocation.StoreLocationFullPath,
	})

	if storelocation.StoreLocationCanStore {
		iconCanStore = widgets.NewIcon(widgets.IconAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
				Attributes: map[string]string{
					"style": "float: right",
				},
			},
			Icon: themes.NewMdiIcon(themes.MDI_CHECK, themes.MDI_24PX),
		})
	} else {
		iconCanStore = widgets.NewIcon(widgets.IconAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
				Attributes: map[string]string{
					"style": "float: right",
				},
			},
			Icon: themes.NewMdiIcon(themes.MDI_NO_CHECK, themes.MDI_24PX),
		})
	}

	d := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
		},
	})
	d.AppendChild(iconColor)
	d.AppendChild(spanLabel)
	d.AppendChild(iconCanStore)

	return jsutils.CreateJsHTMLElementFromString(d.OuterHTML())

}
