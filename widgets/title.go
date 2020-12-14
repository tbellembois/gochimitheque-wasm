package widgets

import (
	"syscall/js"

	"honnef.co/go/js/dom/v2"
)

func TitleWrapper(this js.Value, args []js.Value) interface{} {

	t := Title(args[0].String(), args[1].String())
	return js.ValueOf(t)

}

// Title returns a page title with a customized icon
func Title(msgText string, msgType string) *dom.HTMLDivElement {

	Doc := dom.GetWindow().Document()

	t := Doc.CreateElement("div").(*dom.HTMLDivElement)
	s := Doc.CreateElement("span").(*dom.HTMLSpanElement)
	icon := Doc.CreateElement("span").(*dom.HTMLSpanElement)
	t.Class().SetString("mt-md-3 mb-md-3 row text-right")
	s.Class().SetString("col-sm-11 align-bottom")
	s.SetTextContent(msgText)

	switch msgType {
	case "history":
		icon.Class().SetString("mdi mdi-24px mdi-alarm")
	case "bookmark":
		icon.Class().SetString("mdi mdi-24px mdi-bookmark")
	case "entity":
		icon.Class().SetString("mdi mdi-24px mdi-store")
	case "storelocation":
		icon.Class().SetString("mdi mdi-24px mdi-Docker")
	case "product":
		icon.Class().SetString("mdi mdi-24px mdi-tag")
	case "storage":
		icon.Class().SetString("mdi mdi-24px mdi-cube-unfolded")
	default:
		icon.Class().SetString("mdi mdi-24px mdi-menu-right-outline")
	}

	t.AppendChild(s)
	t.AppendChild(icon)

	return t
}
