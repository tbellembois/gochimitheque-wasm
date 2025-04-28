//go:build go1.24 && js && wasm

package widgets

import (
	"github.com/tbellembois/gochimitheque-wasm/widgets/themes"

	"honnef.co/go/js/dom/v2"
)

type icon struct {
	Widget
}

type IconAttributes struct {
	BaseAttributes
	Text  string
	Title string
	Icon  themes.MDIcon
}

func NewIcon(args IconAttributes) *icon {

	i := &icon{}
	i.HTMLElement = dom.GetWindow().Document().CreateElement("span").(*dom.HTMLSpanElement)

	i.SetInnerHTML(args.Text)
	i.SetTitle(args.Title)
	// Appending mateial design icon to classes.
	args.BaseAttributes.Classes = append(args.BaseAttributes.Classes, args.Icon.ToString())
	args.BaseAttributes.Visible = true

	i.SetBaseAttributes(args.BaseAttributes)

	return i

}
