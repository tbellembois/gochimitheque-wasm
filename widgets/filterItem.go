package widgets

import (
	"fmt"

	"github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/locales"
	"github.com/tbellembois/gochimitheque-wasm/widgets/themes"
)

func FilterItem(title, value string) string {

	var (
		colValue *div
	)

	row := NewDiv(DivAttributes{
		BaseAttributes: BaseAttributes{
			Visible: true,
			Classes: []string{"row"},
		},
	})

	colRemoveFilter := NewDiv(DivAttributes{
		BaseAttributes: BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-auto"},
		},
	})

	colRemoveFilter.AppendChild(NewBSButtonWithIcon(
		ButtonAttributes{
			BaseAttributes: BaseAttributes{
				Visible: true,
				Id:      fmt.Sprintf("removefilter%s", title),
				Classes: []string{"btn-sm", "btn", "btn-outline-primary"},
			},
			Title: locales.Translate("remove_filter", globals.HTTPHeaderAcceptLanguage),
		},
		IconAttributes{
			BaseAttributes: BaseAttributes{
				Visible: true,
			},
			Icon:  themes.NewMdiIcon(themes.MDI_REMOVEFILTER, themes.MDI_16PX),
			Title: locales.Translate("remove_filter", globals.HTTPHeaderAcceptLanguage),
		},
		[]themes.BSClass{themes.BS_BTN, themes.BS_BNT_LINK},
	))

	colRemoveFilter.AppendChild(NewSpan(SpanAttributes{
		BaseAttributes: BaseAttributes{
			Visible: true,
			Classes: []string{"iconlabel", "text-primary"},
		},
		Text: locales.Translate(title, globals.HTTPHeaderAcceptLanguage),
	}))

	colValue = NewDiv(DivAttributes{
		BaseAttributes: BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-auto"},
		},
	})
	// switch title {
	// case "s_cas_number_cmr", "s_borrowing", "s_storage_to_destroy":
	// 	colValue.AppendChild(NewIcon(IconAttributes{
	// 		BaseAttributes: BaseAttributes{
	// 			Visible: true,
	// 			Classes: []string{"iconlabel"},
	// 		},
	// 		Icon: themes.NewMdiIcon(themes.MDI_OK, ""),
	// 	}))
	// default:
	// 	colValue.AppendChild(NewSpan(SpanAttributes{
	// 		BaseAttributes: BaseAttributes{
	// 			Visible: true,
	// 			Classes: []string{"text-start"},
	// 		},
	// 		Text: value,
	// 	}))
	// }

	colValue.AppendChild(NewSpan(SpanAttributes{
		BaseAttributes: BaseAttributes{
			Visible: true,
			Classes: []string{"text-start"},
		},
		Text: value,
	}))

	row.AppendChild(colRemoveFilter)
	row.AppendChild(colValue)

	return row.OuterHTML()

}
