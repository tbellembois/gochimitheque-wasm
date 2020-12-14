package widgets

import (
	"honnef.co/go/js/dom/v2"
)

type Link struct {
	Widget
}

type LinkAttributes struct {
	BaseAttributes
	Title   string
	Href    string
	Onclick string
	Label   dom.HTMLElement
}

func NewLink(args LinkAttributes) *Link {

	l := &Link{}
	l.HTMLElement = dom.GetWindow().Document().CreateElement("a").(*dom.HTMLAnchorElement)

	l.SetAttribute("href", args.Href)
	l.SetAttribute("onclick", args.Onclick)
	l.SetTitle(args.Title)
	l.SetInnerHTML(args.Label.OuterHTML())

	l.SetBaseAttributes(args.BaseAttributes)

	return l

}
