//go:build go1.24 && js && wasm

package widgets

import (
	"honnef.co/go/js/dom/v2"
)

type ul struct {
	Widget
}

type UlAttributes struct {
	BaseAttributes
}

func NewUl(args UlAttributes) *ul {

	o := &ul{}
	o.HTMLElement = dom.GetWindow().Document().CreateElement("ul").(*dom.HTMLUListElement)

	o.SetBaseAttributes(args.BaseAttributes)

	return o

}
