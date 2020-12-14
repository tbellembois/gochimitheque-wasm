package utils

import (
	"fmt"
	"syscall/js"

	. "github.com/tbellembois/gochimitheque-wasm/globals"
	. "github.com/tbellembois/gochimitheque-wasm/types"
	"github.com/tbellembois/gochimitheque-wasm/widgets"
	"github.com/tbellembois/gochimitheque-wasm/widgets/themes"
	"honnef.co/go/js/dom/v2"
)

func RedirectTo(href string) {

	js.Global().Get("window").Get("location").Set("href", href)

}

func CreateJsHTMLElementFromString(s string) js.Value {

	t := js.Global().Get("document").Call("createElement", "template")
	t.Set("innerHTML", s)
	return t.Get("content").Get("firstChild")

}

func LoadContent(viewName string, url string, callback func(args ...interface{}), args ...interface{}) {

	if viewName != "" {
		CurrentView = viewName
	}

	done := func(data js.Value) {
		Jq("div#content").SetHtml(data.String())
		if callback != nil {
			callback(args...)
		}
	}
	fail := func(data js.Value) {
		// TODO: improve this
		RedirectTo(fmt.Sprintf("%sdelete-token", ApplicationProxyPath))
	}

	Ajax{
		URL:    url,
		Method: "get",
		Done:   done,
		Fail:   fail,
	}.Send()

}

// TODO: merge with LoadContent
func LoadMenu(viewName string, url string, callback func(args ...interface{}), args ...interface{}) {

	done := func(data js.Value) {
		Jq("div#menu").SetHtml(data.String())
		callback(args...)
	}

	Ajax{
		URL:    url,
		Method: "get",
		Done:   done,
	}.Send()

}

// TODO: merge with LoadContent
func LoadSearch(viewName string, url string, callback func(args ...interface{}), args ...interface{}) {

	done := func(data js.Value) {
		Jq("div#searchbar").SetHtml(data.String())
		callback(args...)
	}

	Ajax{
		URL:    url,
		Method: "get",
		Done:   done,
	}.Send()

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
	icon := widgets.NewIcon(widgets.IconAttributes{Icon: themes.NewMdiIcon(themes.MDI_TOW_TRUCK, themes.MDI_48PX)})
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

	Jq("#list-collapse").Show()
	Jq("#edit-collapse").Hide()

	return nil

}

func DumpJsObject(object js.Value) {

	fmt.Println(js.Global().Get("JSON").Call("stringify", object).String())

}
