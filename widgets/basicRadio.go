//go:build go1.24 && js && wasm

package widgets

import (
	"honnef.co/go/js/dom/v2"
)

type radio struct {
	Widget
}

type RadioAttributes struct {
	BaseAttributes
	Checked bool
}

func NewRadio(args RadioAttributes) *radio {

	r := &radio{}

	htmlElement := dom.GetWindow().Document().CreateElement("input").(*dom.HTMLInputElement)
	htmlElement.SetAttribute("type", "radio")
	htmlElement.SetChecked(args.Checked)

	r.HTMLElement = htmlElement

	r.SetBaseAttributes(args.BaseAttributes)

	return r

}
