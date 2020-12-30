package widgets

func FilterItem(title, value string) string {

	var (
		colTitle *div
		colValue *div
	)

	row := NewDiv(DivAttributes{
		BaseAttributes: BaseAttributes{
			Visible: true,
			Classes: []string{"row"},
		},
	})

	colTitle = NewDiv(DivAttributes{
		BaseAttributes: BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-2", "iconlabel"},
		},
	})
	colTitle.AppendChild(NewSpan(SpanAttributes{
		BaseAttributes: BaseAttributes{
			Visible: true,
			Classes: []string{"text-end"},
		},
		Text: title,
	}))

	var classes []string
	if title == "" {
		classes = []string{"col-sm-2", "badge", "badge-dark", "mt-sm-2"}
	} else {
		classes = []string{"col-sm-10"}
	}
	colValue = NewDiv(DivAttributes{
		BaseAttributes: BaseAttributes{
			Visible: true,
			Classes: classes,
		},
	})
	colValue.AppendChild(NewSpan(SpanAttributes{
		BaseAttributes: BaseAttributes{
			Visible: true,
			Classes: []string{"text-start"},
		},
		Text: value,
	}))

	if title != "" {
		row.AppendChild(colTitle)
	}
	row.AppendChild(colValue)

	return row.OuterHTML()

}
