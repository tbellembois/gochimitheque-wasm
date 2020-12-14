package widgets

import (
	"honnef.co/go/js/dom/v2"
)

type Option struct {
	Widget
}

type OptionAttributes struct {
	BaseAttributes
	Text            string
	Value           string
	DefaultSelected bool
	Selected        bool
}

func NewOption(args OptionAttributes) *Option {

	o := &Option{}

	htmlElement := dom.GetWindow().Document().CreateElement("option").(*dom.HTMLOptionElement)
	htmlElement.SetTextContent(args.Text)
	htmlElement.SetAttribute("value", args.Value)
	htmlElement.SetSelected(args.Selected)
	htmlElement.SetDefaultSelected(args.DefaultSelected)

	o.HTMLElement = htmlElement

	o.SetBaseAttributes(args.BaseAttributes)

	return o

}
