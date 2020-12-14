package widgets

import "honnef.co/go/js/dom/v2"

type img struct {
	Widget
}

type ImgAttributes struct {
	BaseAttributes
	Src    string
	Alt    string
	Title  string
	Height string
	Width  string
}

func NewImg(args ImgAttributes) *img {

	i := &img{}
	i.HTMLElement = dom.GetWindow().Document().CreateElement("img").(*dom.HTMLImageElement)

	i.SetTitle(args.Title)
	i.SetAttribute("alt", args.Alt)
	i.SetAttribute("src", args.Src)

	if args.Height != "" {
		i.SetAttribute("height", args.Height)
	}
	if args.Width != "" {
		i.SetAttribute("width", args.Width)
	}

	i.SetBaseAttributes(args.BaseAttributes)

	return i

}
