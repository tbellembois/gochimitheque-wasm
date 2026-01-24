//go:build go1.24 && js && wasm

package search

import (
	"encoding/json"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/bstable"
	. "github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/jquery"
	"github.com/tbellembois/gochimitheque-wasm/jsutils"
	"github.com/tbellembois/gochimitheque-wasm/locales"
	"github.com/tbellembois/gochimitheque-wasm/select2"
	. "github.com/tbellembois/gochimitheque-wasm/types"
	"github.com/tbellembois/gochimitheque-wasm/views/storage"
	"github.com/tbellembois/gochimitheque-wasm/widgets"
	"github.com/tbellembois/gochimitheque-wasm/widgets/themes"
)

func Search_listCallback(args ...interface{}) {

	window := js.Global()
	var keycloak js.Value
	keycloak = window.Get("keycloak")
	token := keycloak.Get("token").String()
	marshalToken, _ := json.Marshal(map[string]string{"Authorization": "Bearer " + token})

	select2.NewSelect2(jquery.Jq("select#s_tags"), &select2.Select2Config{
		Placeholder:    locales.Translate("s_tags", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(Tag{})),
		AllowClear:     true,
		Ajax: select2.Select2Ajax{
			URL:            BackProxyPath + "products/tags_old",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Tags{})),
			Headers:        js.Global().Get("JSON").Call("parse", string(marshalToken)),
		},
	}).Select2ify()

	select2.NewSelect2(jquery.Jq("select#s_category"), &select2.Select2Config{
		Placeholder:    locales.Translate("s_category", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(Category{})),
		AllowClear:     true,
		Ajax: select2.Select2Ajax{
			URL:            BackProxyPath + "products/categories_old",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Categories{})),
			Headers:        js.Global().Get("JSON").Call("parse", string(marshalToken)),
		},
	}).Select2ify()

	select2.NewSelect2(jquery.Jq("select#s_entity"), &select2.Select2Config{
		Placeholder:    locales.Translate("s_entity", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(Entity{})),
		AllowClear:     true,
		Ajax: select2.Select2Ajax{
			URL:            BackProxyPath + "entities_old",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Entities{})),
			Headers:        js.Global().Get("JSON").Call("parse", string(marshalToken)),
		},
	}).Select2ify()

	select2.NewSelect2(jquery.Jq("select#s_store_location"), &select2.Select2Config{
		Placeholder:    locales.Translate("s_store_location", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2StoreLocationTemplateResults),
		AllowClear:     true,
		Ajax: select2.Select2Ajax{
			URL:            BackProxyPath + "store_locations_old",
			DataType:       "json",
			Data:           js.FuncOf(Select2StoreLocationAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(StoreLocations{})),
			Headers:        js.Global().Get("JSON").Call("parse", string(marshalToken)),
		},
	}).Select2ify()

	select2.NewSelect2(jquery.Jq("select#s_cas_number"), &select2.Select2Config{
		Placeholder:    locales.Translate("s_cas_number", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(CasNumber{})),
		AllowClear:     true,
		Ajax: select2.Select2Ajax{
			URL:            BackProxyPath + "products/casnumbers_old",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(CasNumbers{})),
			Headers:        js.Global().Get("JSON").Call("parse", string(marshalToken)),
		},
	}).Select2ify()

	select2.NewSelect2(jquery.Jq("select#s_name"), &select2.Select2Config{
		Placeholder:    locales.Translate("s_name", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(Name{})),
		AllowClear:     true,
		Ajax: select2.Select2Ajax{
			URL:            BackProxyPath + "products/names_old",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Names{})),
			Headers:        js.Global().Get("JSON").Call("parse", string(marshalToken)),
		},
	}).Select2ify()

	select2.NewSelect2(jquery.Jq("select#s_empirical_formula"), &select2.Select2Config{
		Placeholder:    locales.Translate("s_empirical_formula", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(EmpiricalFormula{})),
		AllowClear:     true,
		Ajax: select2.Select2Ajax{
			URL:            BackProxyPath + "products/empiricalformulas_old",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(EmpiricalFormulas{})),
			Headers:        js.Global().Get("JSON").Call("parse", string(marshalToken)),
		},
	}).Select2ify()

	select2.NewSelect2(jquery.Jq("select#s_producer_ref"), &select2.Select2Config{
		Placeholder:    locales.Translate("s_producer_ref", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(ProducerRef{})),
		AllowClear:     true,
		Ajax: select2.Select2Ajax{
			URL:            BackProxyPath + "products/producerrefs_old",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(ProducerRefs{})),
			Headers:        js.Global().Get("JSON").Call("parse", string(marshalToken)),
		},
	}).Select2ify()

	select2.NewSelect2(jquery.Jq("select#s_signal_word"), &select2.Select2Config{
		Placeholder:    locales.Translate("s_signal_word", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(SignalWord{})),
		AllowClear:     true,
		Ajax: select2.Select2Ajax{
			URL:            BackProxyPath + "products/signalwords_old",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(SignalWords{})),
			Headers:        js.Global().Get("JSON").Call("parse", string(marshalToken)),
		},
	}).Select2ify()

	select2.NewSelect2(jquery.Jq("select#s_symbols"), &select2.Select2Config{
		Placeholder:    locales.Translate("s_symbols", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2SymbolTemplateResults),
		AllowClear:     true,
		Ajax: select2.Select2Ajax{
			URL:            BackProxyPath + "products/symbols_old",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Symbols{})),
			Headers:        js.Global().Get("JSON").Call("parse", string(marshalToken)),
		},
	}).Select2ify()

	select2.NewSelect2(jquery.Jq("select#s_hazard_statements"), &select2.Select2Config{
		Placeholder:    locales.Translate("s_hazard_statements", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(HazardStatement{})),
		AllowClear:     true,
		Ajax: select2.Select2Ajax{
			URL:            BackProxyPath + "products/hazardstatements_old",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			Headers:        js.Global().Get("JSON").Call("parse", string(marshalToken)),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(HazardStatements{})),
		},
	}).Select2ify()

	select2.NewSelect2(jquery.Jq("select#s_precautionary_statements"), &select2.Select2Config{
		Placeholder:    locales.Translate("s_precautionary_statements", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(PrecautionaryStatement{})),
		AllowClear:     true,
		Ajax: select2.Select2Ajax{
			URL:            BackProxyPath + "products/precautionarystatements_old",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			Headers:        js.Global().Get("JSON").Call("parse", string(marshalToken)),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(PrecautionaryStatements{})),
		},
	}).Select2ify()

	// Works only with no select2.
	jquery.Jq("#search input").On("keyup", js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		event := args[0]
		if !event.Get("keyCode").IsUndefined() && event.Get("keyCode").Int() == 13 {

			event.Call("preventDefault")
			jsutils.Search(js.Null(), nil)

		}

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
			BSTableQueryFilter.QueryFilter.StorageArchive = new(bool)
			*BSTableQueryFilter.QueryFilter.StorageArchive = true
			btnIcon = themes.MDI_HIDE_DELETED
			btnLabel = locales.Translate("hidedeleted_text", HTTPHeaderAcceptLanguage)
		} else {
			BSTableQueryFilter.QueryFilter.StorageArchive = new(bool)
			*BSTableQueryFilter.QueryFilter.StorageArchive = false
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
