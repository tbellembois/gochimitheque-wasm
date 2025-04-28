//go:build go1.24 && js && wasm

package widgets

import (
	"honnef.co/go/js/dom/v2"
)

type button struct {
	Widget
}

type ButtonAttributes struct {
	BaseAttributes
	Label string
	Title string
}

func NewButton(args ButtonAttributes) *button {

	b := &button{}
	b.HTMLElement = dom.GetWindow().Document().CreateElement("button").(*dom.HTMLButtonElement)

	b.SetTitle(args.Title)
	if args.Label != "" {
		b.SetAttribute("label", args.Label)
	}

	b.SetBaseAttributes(args.BaseAttributes)

	return b

}
