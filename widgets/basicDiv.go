//go:build go1.24 && js && wasm

package widgets

import (
	"honnef.co/go/js/dom/v2"
)

type div struct {
	Widget
}

type DivAttributes struct {
	BaseAttributes
}

func NewDiv(args DivAttributes) *div {

	d := &div{}
	d.HTMLElement = dom.GetWindow().Document().CreateElement("div").(*dom.HTMLDivElement)

	d.SetBaseAttributes(args.BaseAttributes)

	return d

}
