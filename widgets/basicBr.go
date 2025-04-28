//go:build go1.24 && js && wasm

package widgets

import (
	"honnef.co/go/js/dom/v2"
)

type br struct {
	Widget
}

type BrAttributes struct {
	BaseAttributes
}

func NewBr(args BrAttributes) *br {

	b := &br{}
	b.HTMLElement = dom.GetWindow().Document().CreateElement("br").(*dom.HTMLBRElement)

	b.SetBaseAttributes(args.BaseAttributes)

	return b

}
