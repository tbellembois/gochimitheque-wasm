package search

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/ajax"
	"github.com/tbellembois/gochimitheque-wasm/bstable"
	"github.com/tbellembois/gochimitheque-wasm/globals"
	. "github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/jquery"
	"github.com/tbellembois/gochimitheque-wasm/jsutils"
	"github.com/tbellembois/gochimitheque-wasm/locales"
	"github.com/tbellembois/gochimitheque-wasm/models"
	"github.com/tbellembois/gochimitheque-wasm/select2"
	. "github.com/tbellembois/gochimitheque-wasm/types"
	"github.com/tbellembois/gochimitheque-wasm/views/storage"
	"github.com/tbellembois/gochimitheque-wasm/widgets"
	"github.com/tbellembois/gochimitheque-wasm/widgets/themes"
)

func Search_listCallback(args ...interface{}) {

	select2.NewSelect2(jquery.Jq("select#s_tags"), &select2.Select2Config{
		Placeholder:    locales.Translate("s_tags", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(Tag{})),
		AllowClear:     true,
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "products/tags/",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Tags{})),
		},
	}).Select2ify()

	select2.NewSelect2(jquery.Jq("select#s_category"), &select2.Select2Config{
		Placeholder:    locales.Translate("s_category", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(Category{})),
		AllowClear:     true,
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "products/categories/",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Categories{})),
		},
	}).Select2ify()

	select2.NewSelect2(jquery.Jq("select#s_entity"), &select2.Select2Config{
		Placeholder:    locales.Translate("s_entity", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(Entity{})),
		AllowClear:     true,
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "entities",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Entities{})),
		},
	}).Select2ify()

	select2.NewSelect2(jquery.Jq("select#s_storelocation"), &select2.Select2Config{
		Placeholder:    locales.Translate("s_storelocation", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2StoreLocationTemplateResults),
		AllowClear:     true,
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "storelocations",
			DataType:       "json",
			Data:           js.FuncOf(Select2StoreLocationAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(StoreLocations{})),
		},
	}).Select2ify()

	select2.NewSelect2(jquery.Jq("select#s_casnumber"), &select2.Select2Config{
		Placeholder:    locales.Translate("s_casnumber", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(CasNumber{})),
		AllowClear:     true,
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "products/casnumbers/",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(CasNumbers{})),
		},
	}).Select2ify()

	select2.NewSelect2(jquery.Jq("select#s_name"), &select2.Select2Config{
		Placeholder:    locales.Translate("s_name", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(Name{})),
		AllowClear:     true,
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "products/names/",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Names{})),
		},
	}).Select2ify()

	select2.NewSelect2(jquery.Jq("select#s_empiricalformula"), &select2.Select2Config{
		Placeholder:    locales.Translate("s_empiricalformula", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(EmpiricalFormula{})),
		AllowClear:     true,
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "products/empiricalformulas/",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(EmpiricalFormulas{})),
		},
	}).Select2ify()

	select2.NewSelect2(jquery.Jq("select#s_producerref"), &select2.Select2Config{
		Placeholder:    locales.Translate("s_producerref", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(ProducerRef{})),
		AllowClear:     true,
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "products/producerrefs/",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Select2ProducerRefs{})),
		},
	}).Select2ify()

	select2.NewSelect2(jquery.Jq("select#s_signalword"), &select2.Select2Config{
		Placeholder:    locales.Translate("s_signalword", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(SignalWord{})),
		AllowClear:     true,
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "products/signalwords/",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(SignalWords{})),
		},
	}).Select2ify()

	select2.NewSelect2(jquery.Jq("select#s_symbols"), &select2.Select2Config{
		Placeholder:    locales.Translate("s_symbols", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2SymbolTemplateResults),
		AllowClear:     true,
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "products/symbols/",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Symbols{})),
		},
	}).Select2ify()

	select2.NewSelect2(jquery.Jq("select#s_hazardstatements"), &select2.Select2Config{
		Placeholder:    locales.Translate("s_hazardstatements", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(HazardStatement{})),
		AllowClear:     true,
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "products/hazardstatements/",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(HazardStatements{})),
		},
	}).Select2ify()

	select2.NewSelect2(jquery.Jq("select#s_precautionarystatements"), &select2.Select2Config{
		Placeholder:    locales.Translate("s_precautionarystatements", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(PrecautionaryStatement{})),
		AllowClear:     true,
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "products/precautionarystatements/",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(PrecautionaryStatements{})),
		},
	}).Select2ify()

	// Works only with no select2.
	jquery.Jq("#search input").On("keyup", js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		event := args[0]
		if !event.Get("which").IsUndefined() && event.Get("which").Int() == 13 {

			event.Call("preventDefault")
			jsutils.Search(js.Null(), nil)

		}

		return nil

	}))

	// Stock.
	jquery.Jq("#s_storage_stock_button").On("click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		jquery.Jq("#stock").Append(widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
				Classes: []string{"mdi", "mdi-loading", "mdi-spin", "mdi-36px"},
			},
		}).OuterHTML())

		url := fmt.Sprintf("%sentities/stocks/%d", ApplicationProxyPath, globals.CurrentProduct.ProductID)
		method := "get"

		done := func(data js.Value) {

			var (
				storelocations []models.StoreLocation
				err            error
			)

			if err = json.Unmarshal([]byte(data.String()), &storelocations); err != nil {
				fmt.Println(err)
			}

			jquery.Jq("#stock").SetHtml("")

			// rowButtonClose := widgets.NewDiv(widgets.DivAttributes{
			// 	BaseAttributes: widgets.BaseAttributes{
			// 		Visible: true,
			// 		Classes: []string{"row"},
			// 	},
			// })
			// buttonClose := widgets.NewBSButtonWithIcon(
			// 	widgets.ButtonAttributes{
			// 		BaseAttributes: widgets.BaseAttributes{
			// 			Visible: true,
			// 			Attributes: map[string]string{
			// 				"onclick": "$('#stock').html('')",
			// 			},
			// 		},
			// 		Title: locales.Translate("close", HTTPHeaderAcceptLanguage),
			// 	},
			// 	widgets.IconAttributes{
			// 		BaseAttributes: widgets.BaseAttributes{
			// 			Visible: true,
			// 			Classes: []string{"iconlabel"},
			// 		},
			// 		Text: locales.Translate("close", HTTPHeaderAcceptLanguage),
			// 		Icon: themes.NewMdiIcon(themes.MDI_CLOSE, ""),
			// 	},
			// 	[]themes.BSClass{themes.BS_BTN, themes.BS_BNT_LINK},
			// )
			// rowButtonClose.AppendChild(buttonClose)

			// rowProduct := widgets.NewDiv(widgets.DivAttributes{
			// 	BaseAttributes: widgets.BaseAttributes{
			// 		Visible: true,
			// 		Classes: []string{"row", "iconlabel"},
			// 	},
			// })
			// rowProduct.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
			// 	BaseAttributes: widgets.BaseAttributes{
			// 		Visible: true,
			// 		Classes: []string{"col", "iconlabel"},
			// 	},
			// 	Text: CurrentProduct.Name.NameLabel,
			// }))

			//jquery.Jq("#stock").Append(rowButtonClose.OuterHTML())
			// jquery.Jq("#stock").Append(rowProduct.OuterHTML())

			for _, storelocation := range storelocations {
				jsutils.ShowStockRecursive(&storelocation, 0, "#stock")
			}

		}
		fail := func(data js.Value) {

			jsutils.DisplayGenericErrorMessage()

		}

		ajax.Ajax{
			Method: method,
			URL:    url,
			Done:   done,
			Fail:   fail,
		}.Send()

		return nil

	}))

	// Show/Hide archives.
	jquery.Jq("#s_storage_archive_button").On("click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		var (
			btnIcon  themes.IconFace
			btnLabel string
		)

		BSTableQueryFilter.Lock()

		if jquery.Jq("#s_storage_archive_button > span").HasClass(themes.MDI_SHOW_DELETED.ToString()) {
			BSTableQueryFilter.QueryFilter.StorageArchive = true
			btnIcon = themes.MDI_HIDE_DELETED
			btnLabel = locales.Translate("hidedeleted_text", HTTPHeaderAcceptLanguage)
		} else {
			BSTableQueryFilter.QueryFilter.StorageArchive = false
			btnIcon = themes.MDI_SHOW_DELETED
			btnLabel = locales.Translate("showdeleted_text", HTTPHeaderAcceptLanguage)
		}

		buttonTitle := widgets.NewIcon(widgets.IconAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
				Classes: []string{"iconlabel"},
			},
			Icon: themes.NewMdiIcon(btnIcon, ""),
			Text: btnLabel,
		})

		jquery.Jq("#s_storage_archive_button").SetProp("title", btnLabel)
		jquery.Jq("#s_storage_archive_button").SetHtml("")
		jquery.Jq("#s_storage_archive_button").Append(buttonTitle.OuterHTML())

		bstable.NewBootstraptable(jquery.Jq("#Storage_table"), nil).Refresh(nil)
		jquery.Jq("#Storage_table").On("load-success.bs.table", js.FuncOf(storage.ShowIfAuthorizedActionButtons))

		return nil

	}))

}
