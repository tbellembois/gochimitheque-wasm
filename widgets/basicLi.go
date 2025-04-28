//go:build go1.24 && js && wasm

package widgets

import (
	"honnef.co/go/js/dom/v2"
)

type li struct {
	Widget
}

type LiAttributes struct {
	BaseAttributes
	Text string
}

func NewLi(args LiAttributes) *li {

	o := &li{}
	o.HTMLElement = dom.GetWindow().Document().CreateElement("li").(*dom.HTMLLIElement)

	o.SetInnerHTML(args.Text)

	o.SetBaseAttributes(args.BaseAttributes)

	return o

}
