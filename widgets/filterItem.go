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
			Classes: []string{"col-sm-auto", "iconlabel"},
		},
	})
	colTitle.AppendChild(NewSpan(SpanAttributes{
		BaseAttributes: BaseAttributes{
			Visible: true,
		},
		Text: title,
	}))

	classes := []string{"col-sm-auto"}
	if title == "" {
		classes = append(classes, "badge", "badge-secondary")
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
		},
		Text: value,
	}))

	if title != "" {
		row.AppendChild(colTitle)
	}
	row.AppendChild(colValue)

	return row.OuterHTML()

}
