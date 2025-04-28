//go:build go1.24 && js && wasm

package widgets

import (
	"honnef.co/go/js/dom/v2"
)

// InlineRadio return a radio inline block
func InlineRadio(inputattr map[string]string) *dom.HTMLDivElement {

	Doc := dom.GetWindow().Document()

	// main div
	maindiv := Doc.CreateElement("div").(*dom.HTMLDivElement)
	maindiv.SetClass("form-check form-check-inline")
	// input
	input := Doc.CreateElement("input").(*dom.HTMLInputElement)
	input.SetAttribute("type", "radio")
	input.SetClass("form-check-input")
	// ...setting up additional attributes
	for a := range inputattr {
		input.SetAttribute(a, inputattr[a])
	}
	// label
	label := Doc.CreateElement("label").(*dom.HTMLLabelElement)
	label.SetClass("form-check-label")
	label.SetAttribute("for", inputattr["id"])
	label.SetInnerHTML(inputattr["label"])

	// building the final result
	maindiv.AppendChild(input)
	maindiv.AppendChild(label)
	return maindiv
}
