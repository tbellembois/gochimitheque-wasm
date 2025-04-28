//go:build go1.24 && js && wasm

package widgets

import (
	"fmt"
	"syscall/js"

	. "github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/locales"
	"honnef.co/go/js/dom/v2"
)

func PermissionWrapper(this js.Value, args []js.Value) interface{} {

	p := Permission(args[0].Int(), args[1].String(), args[2].Bool())
	return js.ValueOf(p)

}

// Permission return a widget to setup people permissions
func Permission(entityID int, entityName string, ismanager bool) *dom.HTMLDivElement {

	Doc := dom.GetWindow().Document()

	// create main widget div
	widgetdiv := Doc.CreateElement("div").(*dom.HTMLDivElement)
	widgetdiv.SetID(fmt.Sprintf("perm%d", entityID))
	widgetdiv.SetClass("col-sm-12")
	title := Doc.CreateElement("div").(*dom.HTMLDivElement)
	title.SetClass("d-flex")
	title.SetInnerHTML("<span class='mdi mdi-home-group mdi-24px'/>" + entityName)

	widgetdiv.AppendChild(title)

	if ismanager {
		s := Doc.CreateElement("span").(*dom.HTMLSpanElement)
		s.SetClass("mdi mdi-36px mdi-account-star")
		s.SetAttribute("title", "manager")
		coldiv := Doc.CreateElement("div").(*dom.HTMLDivElement)
		coldiv.SetClass("col-sm-2")
		coldiv.AppendChild(s)

		widgetdiv.AppendChild(coldiv)
		return widgetdiv
	}

	for _, i := range PermItems {
		// products permissions widget is static
		if i != "products" && i != "rproducts" {
			//println(i)
			// building main row
			mainrowdiv := Doc.CreateElement("div").(*dom.HTMLDivElement)
			mainrowdiv.SetClass("form-group row d-flex")
			// building first col for table name
			label := Doc.CreateElement("div").(*dom.HTMLDivElement)
			label.SetClass("iconlabel text-right")
			label.SetInnerHTML(locales.Translate(fmt.Sprintf("permission_%s", i), HTTPHeaderAcceptLanguage))
			firstcoldiv := Doc.CreateElement("div").(*dom.HTMLDivElement)
			firstcoldiv.SetClass("col-sm-2")
			firstcoldiv.AppendChild(label)
			// building second col for radios
			noneradioattrs := map[string]string{
				"id":        fmt.Sprintf("permn%s%d", i, entityID),
				"name":      fmt.Sprintf("perm%s%d", i, entityID),
				"value":     "none",
				"checked":   "checked",
				"label":     fmt.Sprintf("<label class=\"form-check-label ml-sm-1 pr-sm-1 pl-sm-1 text-secondary border border-secondary rounded\" for=\"permn%s%d\" title=\"%s\"><span class=\"mdi mdi-close\"></span></label>", i, entityID, locales.Translate("permission_none", HTTPHeaderAcceptLanguage)),
				"perm_name": "n",
				"item_name": i,
				"entity_id": fmt.Sprintf("%d", entityID),
				"class":     fmt.Sprintf("perm permn permn%s", i)}
			readradioattrs := map[string]string{
				"id":        fmt.Sprintf("permr%s%d", i, entityID),
				"name":      fmt.Sprintf("perm%s%d", i, entityID),
				"value":     "r",
				"label":     fmt.Sprintf("<label class=\"form-check-label ml-sm-1 pr-sm-1 pl-sm-1 text-secondary border border-secondary rounded\" for=\"permn%s%d\" title=\"%s\"><span class=\"mdi mdi-eye\"></span></label>", i, entityID, locales.Translate("permission_read", HTTPHeaderAcceptLanguage)),
				"perm_name": "r",
				"item_name": i,
				"entity_id": fmt.Sprintf("%d", entityID),
				"class":     fmt.Sprintf("perm permr permr%s", i)}
			writeradioattrs := map[string]string{
				"id":        fmt.Sprintf("permw%s%d", i, entityID),
				"name":      fmt.Sprintf("perm%s%d", i, entityID),
				"value":     "w",
				"label":     fmt.Sprintf("<label class=\"form-check-label ml-sm-1 pr-sm-1 pl-sm-1 text-secondary border border-secondary rounded\" for=\"permn%s%d\" title=\"%s\"><span class=\"mdi mdi-eye\"></span><span class=\"mdi mdi-creation\"></span><span class=\"mdi mdi-pencil-outline\"></span><span class=\"mdi mdi-delete\"></span></label>", i, entityID, locales.Translate("permission_crud", HTTPHeaderAcceptLanguage)),
				"perm_name": "w",
				"item_name": i,
				"entity_id": fmt.Sprintf("%d", entityID),
				"class":     fmt.Sprintf("perm permw permw%s", i)}
			secondcoldiv := Doc.CreateElement("div").(*dom.HTMLDivElement)
			secondcoldiv.SetClass("col-sm-4")
			secondcoldiv.AppendChild(InlineRadio(noneradioattrs))
			secondcoldiv.AppendChild(InlineRadio(readradioattrs))
			secondcoldiv.AppendChild(InlineRadio(writeradioattrs))

			// appending to final div
			mainrowdiv.AppendChild(firstcoldiv)
			mainrowdiv.AppendChild(secondcoldiv)
			widgetdiv.AppendChild(mainrowdiv)
		}
	}
	return widgetdiv
}
