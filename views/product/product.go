package product

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-utils/convert"
	"github.com/tbellembois/gochimitheque-wasm/bstable"
	"github.com/tbellembois/gochimitheque-wasm/globals"
	. "github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/jquery"
	"github.com/tbellembois/gochimitheque-wasm/jsutils"
	"github.com/tbellembois/gochimitheque-wasm/locales"
	"github.com/tbellembois/gochimitheque-wasm/select2"
	. "github.com/tbellembois/gochimitheque-wasm/types"
	"github.com/tbellembois/gochimitheque-wasm/validate"
	"github.com/tbellembois/gochimitheque-wasm/widgets"
	"github.com/tbellembois/gochimitheque-wasm/widgets/themes"
)

var (
	supplierrefToSupplier map[string]int64 // supplierref label -> supplier id
)

func init() {
	supplierrefToSupplier = make(map[string]int64)
}

func LinearToEmpirical(this js.Value, args []js.Value) interface{} {

	select2LinearFormula := select2.NewSelect2(jquery.Jq("select#linearformula"), nil).Select2Data()

	if len(select2LinearFormula) == 0 {
		return ""
	}

	jquery.Jq("#convertedEmpiricalFormula").Append(convert.LinearToEmpiricalFormula(select2LinearFormula[0].Text))

	return nil

}

func NoEmpiricalFormula(this js.Value, args []js.Value) interface{} {

	validate.NewValidate(jquery.Jq("select#empiricalformula"), nil).ValidateRemoveRequired()
	jquery.Jq("span#empiricalformula.badge").Hide()

	return nil

}

func NoCas(this js.Value, args []js.Value) interface{} {

	validate.NewValidate(jquery.Jq("select#casnumber"), nil).ValidateRemoveRequired()
	jquery.Jq("span#casnumber.badge").Hide()

	return nil

}

func HowToMagicalSelector(this js.Value, args []js.Value) interface{} {

	js.Global().Get("window").Call("open", fmt.Sprintf("%simg/magicalselector.webm", ApplicationProxyPath), "_blank")

	return nil

}

func Magic(this js.Value, args []js.Value) interface{} {

	magic := jquery.Jq("textarea#magical").GetVal().String()

	rhs := regexp.MustCompile("((?:EU){0,1}H[0-9]{3}[FfDdAi]{0,2})")
	rps := regexp.MustCompile("(P[0-9]{3})")

	shs := rhs.FindAllStringSubmatch(magic, -1)
	sps := rps.FindAllStringSubmatch(magic, -1)

	var (
		processedH map[string]string
		processedP map[string]string
		ok         bool
	)

	processedH = make(map[string]string)

	select2HS := select2.NewSelect2(jquery.Jq("select#hazardstatements"), nil)
	select2HS.Select2Clear()
	for _, h := range shs {

		if _, ok = processedH[h[1]]; !ok {
			processedH[h[1]] = ""

			for _, hs := range globals.DBHazardStatements {
				if hs.HazardStatementReference == h[1] {
					select2HS.Select2AppendOption(
						widgets.NewOption(widgets.OptionAttributes{
							Text:            h[0],
							Value:           strconv.Itoa(hs.HazardStatementID),
							DefaultSelected: true,
							Selected:        true,
						}).HTMLElement.OuterHTML())
					break
				}
			}

		}
	}

	processedP = make(map[string]string)

	select2PS := select2.NewSelect2(jquery.Jq("select#precautionarystatements"), nil)
	select2PS.Select2Clear()
	for _, p := range sps {

		if _, ok = processedP[p[1]]; !ok {
			processedP[p[1]] = ""

			for _, hs := range globals.DBPrecautionaryStatements {
				if hs.PrecautionaryStatementReference == p[1] {
					select2PS.Select2AppendOption(
						widgets.NewOption(widgets.OptionAttributes{
							Text:            p[0],
							Value:           strconv.Itoa(hs.PrecautionaryStatementID),
							DefaultSelected: true,
							Selected:        true,
						}).HTMLElement.OuterHTML())
					break
				}
			}

		}
	}

	return nil

}

func product_common() {

	//
	// Type chooser.
	//
	jquery.Jq("input[type=radio][name=typechooser]").On("change", js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		switch jquery.Jq("input[name=typechooser]:checked").GetVal().String() {
		case "chem":
			Chemify()
		case "bio":
			Biofy()
		}

		return nil

	}))

	//
	// create form
	//
	// validate
	validate.NewValidate(jquery.Jq("#product"), &validate.ValidateConfig{
		ErrorClass: "alert alert-danger",
		Rules: map[string]validate.ValidateRule{
			"name": {
				Required: js.FuncOf(func(this js.Value, args []js.Value) interface{} { return true }),
				Remote: validate.ValidateRemote{
					BeforeSend: js.FuncOf(func(this js.Value, args []js.Value) interface{} { return false }),
				},
			},
			"producerref": {
				Required: js.FuncOf(func(this js.Value, args []js.Value) interface{} { return true }),
				Remote: validate.ValidateRemote{
					BeforeSend: js.FuncOf(func(this js.Value, args []js.Value) interface{} { return false }),
				},
			},
			"unit_temperature": {
				Required: js.FuncOf(func(this js.Value, args []js.Value) interface{} {
					return jquery.Jq("#product_temperature").GetVal().Truthy()
				}),
				Remote: validate.ValidateRemote{
					BeforeSend: js.FuncOf(func(this js.Value, args []js.Value) interface{} { return false }),
				},
			},
			"empiricalformula": {
				Required: js.FuncOf(func(this js.Value, args []js.Value) interface{} { return true }),
				Remote: validate.ValidateRemote{
					URL:        "",
					Type:       "post",
					BeforeSend: js.FuncOf(ValidateProductEmpiricalFormulaBeforeSend),
					Data: map[string]interface{}{
						"empiricalformula": js.FuncOf(ValidateProductEmpiricalFormulaData),
					},
				},
			},
			"casnumber": {
				Required: js.FuncOf(func(this js.Value, args []js.Value) interface{} { return true }),
				Remote: validate.ValidateRemote{
					URL:        "",
					Type:       "post",
					BeforeSend: js.FuncOf(ValidateProductCasNumberBeforeSend),
					Data: map[string]interface{}{
						"casnumber":           js.FuncOf(ValidateProductCasNumberData1),
						"product_specificity": js.FuncOf(ValidateProductCasNumberData2),
					},
				},
			},
			"cenumber": {
				Remote: validate.ValidateRemote{
					URL:        "",
					Type:       "post",
					BeforeSend: js.FuncOf(ValidateProductCeNumberBeforeSend),
					Data: map[string]interface{}{
						"cenumber": js.FuncOf(ValidateProductCeNumberData),
					},
				},
			},
		},
		Messages: map[string]validate.ValidateMessage{
			"name": {
				Required: locales.Translate("required_input", HTTPHeaderAcceptLanguage),
			},
			"empiricalformula": {
				Required: locales.Translate("required_input", HTTPHeaderAcceptLanguage),
			},
			"casnumber": {
				Required: locales.Translate("required_input", HTTPHeaderAcceptLanguage),
			},
			"producerref": {
				Required: locales.Translate("required_input", HTTPHeaderAcceptLanguage),
			},
		},
	}).Validate()

	// select2
	select2.NewSelect2(jquery.Jq("select#producerref"), &select2.Select2Config{
		Placeholder:       locales.Translate("product_producerref_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult:    js.FuncOf(select2.Select2GenericTemplateResults(ProducerRef{})),
		TemplateSelection: js.FuncOf(Select2ProducerRefTemplateSelection),
		Tags:              true,
		AllowClear:        true,
		CreateTag:         js.FuncOf(Select2ProducerRefCreateTag),
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "products/producerrefs/",
			DataType:       "json",
			Data:           js.FuncOf(Select2ProducerRefAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(ProducerRefs{})),
		},
	}).Select2ify()
	jquery.Jq("select#producerref").On("change", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return validate.NewValidate(jquery.Jq(this), nil).Valid()
	}))
	jquery.Jq("select#producerref").On("select2:select", js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		select2ProducerrefSelected := args[0].Get("params").Get("data")
		producerrefSelected := ProducerRef{}.FromJsJSONValue(select2ProducerrefSelected)

		producerref := producerrefSelected.(ProducerRef)

		// If we create a new producerref
		if producerref.Producer == nil {
			return nil
		}

		select2Producer := select2.NewSelect2(jquery.Jq("select#producer"), nil)
		select2Producer.Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            producerref.Producer.ProducerLabel.String,
				Value:           strconv.Itoa(int(producerref.Producer.ProducerID.Int64)),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())

		return nil

	}))
	select2.NewSelect2(jquery.Jq("select#producer"), &select2.Select2Config{
		Placeholder:    locales.Translate("product_producer_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(Producer{})),
		AllowClear:     true,
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "products/producers/",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Producers{})),
		},
	}).Select2ify()

	select2.NewSelect2(jquery.Jq("select#supplierrefs"), &select2.Select2Config{
		Placeholder:       locales.Translate("product_supplierref_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult:    js.FuncOf(select2.Select2GenericTemplateResults(SupplierRef{})),
		TemplateSelection: js.FuncOf(Select2SupplierRefTemplateSelection),
		Tags:              true,
		AllowClear:        true,
		CreateTag:         js.FuncOf(Select2SupplierRefCreateTag),
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "products/supplierrefs/",
			DataType:       "json",
			Data:           js.FuncOf(Select2SupplierRefAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(SupplierRefs{})),
		},
	}).Select2ify()

	select2.NewSelect2(jquery.Jq("select#supplier"), &select2.Select2Config{
		Placeholder:    locales.Translate("product_supplier_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(Supplier{})),
		AllowClear:     true,
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "products/suppliers/",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Suppliers{})),
		},
	}).Select2ify()

	select2.NewSelect2(jquery.Jq("select#tags"), &select2.Select2Config{
		Placeholder:    locales.Translate("product_tag_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(Tag{})),
		CreateTag:      js.FuncOf(select2.Select2GenericCreateTag(Tag{})),
		AllowClear:     true,
		Tags:           true,
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "products/tags/",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Tags{})),
		},
	}).Select2ify()

	select2.NewSelect2(jquery.Jq("select#category"), &select2.Select2Config{
		Placeholder:    locales.Translate("product_category_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(Category{})),
		AllowClear:     true,
		Tags:           true,
		CreateTag:      js.FuncOf(select2.Select2GenericCreateTag(Category{})),
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "products/categories/",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Categories{})),
		},
	}).Select2ify()

	select2.NewSelect2(jquery.Jq("select#unit_temperature"), &select2.Select2Config{
		Placeholder:    locales.Translate("product_unit_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(Unit{})),
		AllowClear:     true,
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "storages/units",
			DataType:       "json",
			Data:           js.FuncOf(Select2UnitTemperatureAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Units{})),
		},
	}).Select2ify()
	jquery.Jq("select#unit_temperature").On("change", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return validate.NewValidate(jquery.Jq(this), nil).Valid()
	}))

	select2.NewSelect2(jquery.Jq("select#casnumber"), &select2.Select2Config{
		Placeholder:    locales.Translate("product_cas_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(CasNumber{})),
		AllowClear:     true,
		Tags:           true,
		CreateTag:      js.FuncOf(select2.Select2GenericCreateTag(CasNumber{})),
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "products/casnumbers/",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(CasNumbers{})),
		},
	}).Select2ify()
	jquery.Jq("select#casnumber").On("change", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return validate.NewValidate(jquery.Jq(this), nil).Valid()
	}))

	select2.NewSelect2(jquery.Jq("select#cenumber"), &select2.Select2Config{
		Placeholder:    locales.Translate("product_ce_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(CeNumber{})),
		AllowClear:     true,
		Tags:           true,
		CreateTag:      js.FuncOf(select2.Select2GenericCreateTag(CeNumber{})),
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "products/cenumbers/",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(CeNumbers{})),
		},
	}).Select2ify()
	jquery.Jq("select#cenumber").On("change", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return validate.NewValidate(jquery.Jq(this), nil).Valid()
	}))

	select2.NewSelect2(jquery.Jq("select#physicalstate"), &select2.Select2Config{
		Placeholder:    locales.Translate("product_physicalstate_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(PhysicalState{})),
		AllowClear:     true,
		Tags:           true,
		CreateTag:      js.FuncOf(select2.Select2GenericCreateTag(PhysicalState{})),
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "products/physicalstates/",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(PhysicalStates{})),
		},
	}).Select2ify()

	select2.NewSelect2(jquery.Jq("select#signalword"), &select2.Select2Config{
		Placeholder:    locales.Translate("product_signalword_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(SignalWord{})),
		AllowClear:     true,
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "products/signalwords/",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(SignalWords{})),
		},
	}).Select2ify()

	select2.NewSelect2(jquery.Jq("select#classofcompound"), &select2.Select2Config{
		Placeholder:    locales.Translate("product_classofcompound_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(ClassOfCompound{})),
		AllowClear:     true,
		Tags:           true,
		CreateTag:      js.FuncOf(select2.Select2GenericCreateTag(ClassOfCompound{})),
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "products/classofcompounds/",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(ClassesOfCompound{})),
		},
	}).Select2ify()

	select2.NewSelect2(jquery.Jq("select#name"), &select2.Select2Config{
		Placeholder:    locales.Translate("product_name_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(Name{})),
		AllowClear:     true,
		Tags:           true,
		CreateTag:      js.FuncOf(select2.Select2GenericCreateTag(Name{})),
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "products/names/",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Names{})),
		},
	}).Select2ify()
	jquery.Jq("select#name").On("change", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return validate.NewValidate(jquery.Jq(this), nil).Valid()
	}))

	select2.NewSelect2(jquery.Jq("select#empiricalformula"), &select2.Select2Config{
		Placeholder:    locales.Translate("product_empiricalformula_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(EmpiricalFormula{})),
		AllowClear:     true,
		Tags:           true,
		CreateTag:      js.FuncOf(select2.Select2GenericCreateTag(EmpiricalFormula{})),
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "products/empiricalformulas/",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(EmpiricalFormulas{})),
		},
	}).Select2ify()
	jquery.Jq("select#empiricalformula").On("change", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return validate.NewValidate(jquery.Jq(this), nil).Valid()
	}))

	select2.NewSelect2(jquery.Jq("select#linearformula"), &select2.Select2Config{
		Placeholder:    locales.Translate("product_linearformula_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(LinearFormula{})),
		AllowClear:     true,
		Tags:           true,
		CreateTag:      js.FuncOf(select2.Select2GenericCreateTag(LinearFormula{})),
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "products/linearformulas/",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(LinearFormulas{})),
		},
	}).Select2ify()

	select2.NewSelect2(jquery.Jq("select#synonyms"), &select2.Select2Config{
		Placeholder:    locales.Translate("product_synonyms_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(Name{})),
		AllowClear:     true,
		Tags:           true,
		CreateTag:      js.FuncOf(select2.Select2GenericCreateTag(Name{})),
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "products/synonyms/",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Names{})),
		},
	}).Select2ify()

	select2.NewSelect2(jquery.Jq("select#symbols"), &select2.Select2Config{
		Placeholder:    locales.Translate("product_symbols_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2SymbolTemplateResults),
		AllowClear:     true,
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "products/symbols/",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Symbols{})),
		},
	}).Select2ify()

	select2.NewSelect2(jquery.Jq("select#hazardstatements"), &select2.Select2Config{
		Placeholder:    locales.Translate("product_hazardstatements_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2HazardStatementTemplateResults),
		AllowClear:     true,
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "products/hazardstatements/",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(HazardStatements{})),
		},
	}).Select2ify()

	select2.NewSelect2(jquery.Jq("select#precautionarystatements"), &select2.Select2Config{
		Placeholder:    locales.Translate("product_precautionarystatements_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2PrecautionaryStatementTemplateResults),
		AllowClear:     true,
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "products/precautionarystatements/",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(PrecautionaryStatements{})),
		},
	}).Select2ify()

	jquery.Jq("#product_twodformula").On("change", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		js.Global().Call("load2dimage")
		return nil
	}))

}

func ShowIfAuthorizedMenuItems(args ...interface{}) {

	jsutils.HasPermission("products", "-2", "get", func() {
		jquery.Jq("#menu_scan_qrcode").FadeIn()
		jquery.Jq("#menu_list_products").FadeIn()
		jquery.Jq("#menu_list_bookmarks").FadeIn()
	}, func() {
	})

	jsutils.HasPermission("products", "", "post", func() {
		jquery.Jq("#menu_create_product").FadeIn()
	}, func() {
	})

	jsutils.HasPermission("entities", "-2", "get", func() {
		jquery.Jq("#menu_entities").FadeIn()
	}, func() {
	})

	jsutils.HasPermission("entities", "", "post", func() {
		jquery.Jq("#menu_create_entity").FadeIn()
	}, func() {
	})

	jsutils.HasPermission("entities", "-2", "put", func() {
		jquery.Jq("#menu_update_welcomeannounce").FadeIn()
	}, func() {
	})

	jsutils.HasPermission("storages", "-2", "get", func() {
		jquery.Jq("#menu_storelocations").FadeIn()
	}, func() {
	})

	jsutils.HasPermission("storelocations", "", "post", func() {
		jquery.Jq("#menu_create_storelocation").FadeIn()
	}, func() {
	})

	jsutils.HasPermission("people", "-2", "get", func() {
		jquery.Jq("#menu_people").FadeIn()
	}, func() {
	})

	jsutils.HasPermission("people", "", "post", func() {
		jquery.Jq("#menu_create_person").FadeIn()
	}, func() {
	})

}

func showStockRecursive(storelocation *StoreLocation, depth int) {

	// Checking if there is a stock or not for the store location.
	hasStock := false
	for _, stock := range storelocation.Stocks {
		if stock.Total != 0 || stock.Current != 0 {
			hasStock = true
			break
		}
	}

	if hasStock {

		// Depth.
		depthSep := ""
		for i := 0; i <= depth; i++ {
			depthSep += "<span class='mdi mdi-microsoft'></span>"
		}

		rowStorelocation := widgets.NewDiv(widgets.DivAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
				Classes: []string{"row", "iconlabel"},
			},
		})
		rowStorelocation.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
				Classes: []string{"col", "iconlabel"},
			},
			Text: fmt.Sprintf("%s %s", depthSep, storelocation.StoreLocationName.String),
		}))

		jquery.Jq("#stock").Append(rowStorelocation.OuterHTML())

		for _, stock := range storelocation.Stocks {

			rowStocks := widgets.NewDiv(widgets.DivAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"row", "iconlabel"},
				},
			})

			if !(stock.Total == 0 && stock.Current == 0) {

				rowStock := widgets.NewDiv(widgets.DivAttributes{
					BaseAttributes: widgets.BaseAttributes{
						Visible: true,
						Classes: []string{"col-sm-12", "iconlabel"},
					},
				})

				rowStock.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
					BaseAttributes: widgets.BaseAttributes{
						Visible: true,
						Classes: []string{"iconlabel"},
					},
					Text: fmt.Sprintf("%s %s: %f %s %s: %f %s",
						depthSep,
						locales.Translate("stock_storelocation_title", HTTPHeaderAcceptLanguage),
						stock.Current,
						stock.Unit.UnitLabel.String,
						locales.Translate("stock_storelocation_sub_title", HTTPHeaderAcceptLanguage),
						stock.Total,
						stock.Unit.UnitLabel.String),
				}))

				rowStocks.AppendChild(rowStock)

				jquery.Jq("#stock").Append(rowStocks.OuterHTML())

			}

		}

	}

	if len(storelocation.Children) > 0 {
		depth++
		for _, child := range storelocation.Children {

			showStockRecursive(child, depth)

		}
	}

}

func Product_listBookmarkCallback(this js.Value, args []js.Value) interface{} {

	BSTableQueryFilter.Clean()
	BSTableQueryFilter.ProductBookmark = true

	productCallbackWrapper := func(args ...interface{}) {
		Product_listCallback(js.Null(), nil)
	}

	jsutils.LoadContent("div#content", "product", fmt.Sprintf("%sv/products", ApplicationProxyPath), productCallbackWrapper, nil)

	return nil

}

func Product_listCallback(this js.Value, args []js.Value) interface{} {

	//product_common()

	bstable.NewBootstraptable(jquery.Jq("#Product_table"), &bstable.BootstraptableParams{Ajax: "Product_getTableData"})
	jquery.Jq("#Product_table").On("load-success.bs.table", js.FuncOf(ShowIfAuthorizedActionButtons))

	jquery.Jq("#search").Show()
	jquery.Jq("#actions").Show()
	jquery.Jq("#s_storage_archive_button").Hide()
	jquery.Jq("#s_storage_stock_button").Hide()

	btnLabel := locales.Translate("switchstorageview_text", HTTPHeaderAcceptLanguage)
	buttonTitle := widgets.NewIcon(widgets.IconAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"iconlabel"},
		},
		Icon: themes.NewMdiIcon(themes.MDI_STORAGE, ""),
		Text: btnLabel,
	})
	jquery.Jq("#switchview").SetProp("title", btnLabel)
	jquery.Jq("#switchview").SetHtml("")
	jquery.Jq("#switchview").Append(buttonTitle.OuterHTML())

	return nil

}

var ProductCreateCallbackWrapper = func(this js.Value, args []js.Value) interface{} {
	Product_createCallback(nil)
	return nil
}

func Product_createCallback(args ...interface{}) {

	product_common()

	switch reflect.TypeOf(args[0]) {
	case reflect.TypeOf(Product{}):

		product := args[0].(Product)

		FillInProductForm(product, "product")

		jquery.Jq("#search").Hide()
		jquery.Jq("#actions").Hide()

	}

	// Chemical product by default on creation.
	if jquery.Jq("input#product_id").GetVal().String() == "" {
		Chemify()
	} else {
		jquery.Jq("input#showchem").SetProp("disabled", "disabled")
		jquery.Jq("input#showbio").SetProp("disabled", "disabled")
	}

	jquery.Jq("#search").Hide()
	jquery.Jq("#actions").Hide()

}

func Product_SaveCallback(args ...interface{}) {

	BSTableQueryFilter.Lock()
	BSTableQueryFilter.QueryFilter.Product = strconv.Itoa(args[0].(int))
	BSTableQueryFilter.QueryFilter.ProductFilterLabel = fmt.Sprintf("%s %s", CurrentProduct.Name.NameLabel, CurrentProduct.ProductSpecificity.String)
	bstable.NewBootstraptable(jquery.Jq("#Product_table"), nil).Refresh(nil)

	// jquery.Jq("#Product_table").Bootstraptable(nil).Refresh(&BootstraptableRefreshQuery{
	// 	Query: QueryFilter{
	// 		Product: strconv.Itoa(args[0].(int)),
	// 	},
	// })
	jquery.Jq("#Product_table").On("load-success.bs.table", js.FuncOf(ShowIfAuthorizedActionButtons))

	//product_common()

	jquery.Jq("#search").Show()
	jquery.Jq("#actions").Show()

}
