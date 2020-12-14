package widgets

import (
	"honnef.co/go/js/dom/v2"
)

type span struct {
	Widget
}

type SpanAttributes struct {
	BaseAttributes
	Text string
}

func NewSpan(args SpanAttributes) *span {

	s := &span{}
	s.HTMLElement = dom.GetWindow().Document().CreateElement("span").(*dom.HTMLSpanElement)

	s.SetInnerHTML(args.Text)

	s.SetBaseAttributes(args.BaseAttributes)

	return s

}
