package product

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/ajax"
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
	"github.com/tbellembois/gochimitheque/models"
)

var supplierrefToSupplier map[string]int64 // supplierref label -> supplier id

func init() {
	supplierrefToSupplier = make(map[string]int64)
}

func LinearToEmpirical(this js.Value, args []js.Value) interface{} {
	select2LinearFormula := select2.NewSelect2(jquery.Jq("select#linearformula"), nil).Select2Data()

	if len(select2LinearFormula) == 0 {
		return ""
	}

	ajaxData := struct {
		EmpiricalFormula string `json:"empiricalformula"`
	}{
		EmpiricalFormula: select2LinearFormula[0].Text,
	}

	var (
		ajaxDataJson []byte
		err          error
	)
	if ajaxDataJson, err = json.Marshal(ajaxData); err != nil {
		fmt.Println(err)
		return ""
	}

	ajax.Ajax{
		URL:    fmt.Sprintf("%sformat/product/empiricalformula/", ApplicationProxyPath),
		Method: "post",
		Data:   ajaxDataJson,
		Done: func(data js.Value) {
			jquery.Jq("#convertedEmpiricalFormula").SetHtml(data)
		},
		Fail: func(jqXHR js.Value) {
			jsutils.DisplayErrorMessage(locales.Translate("empirical_formula_convert_failed", HTTPHeaderAcceptLanguage))
		},
	}.Send()

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
		default:
			Consufy()
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
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Select2ProducerRefs{})),
		},
	}).Select2ify()
	jquery.Jq("select#producerref").On("change", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return validate.NewValidate(jquery.Jq(this), nil).Valid()
	}))
	jquery.Jq("select#producerref").On("select2:select", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		select2ProducerrefSelected := args[0].Get("params").Get("data")
		producerrefSelected := ProducerRef{ProducerRef: &models.ProducerRef{}}.FromJsJSONValue(select2ProducerrefSelected)

		producerref := producerrefSelected.(ProducerRef)

		// If we create a new producerref
		if (producerref == ProducerRef{} || producerref.Producer == nil) {
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
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Select2Producers{})),
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
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Select2SupplierRefs{})),
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
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Select2Suppliers{})),
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
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Select2Tags{})),
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
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Select2Categories{})),
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
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Select2Units{})),
		},
	}).Select2ify()
	jquery.Jq("select#unit_temperature").On("change", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return validate.NewValidate(jquery.Jq(this), nil).Valid()
	}))

	select2.NewSelect2(jquery.Jq("select#unit_molecularweight"), &select2.Select2Config{
		Placeholder:    locales.Translate("product_unit_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(Unit{})),
		AllowClear:     true,
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "storages/units",
			DataType:       "json",
			Data:           js.FuncOf(Select2UnitTemperatureAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Select2Units{})),
		},
	}).Select2ify()
	jquery.Jq("select#unit_molecularweight").On("change", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
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
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Select2CasNumbers{})),
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
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Select2CeNumbers{})),
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
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Select2PhysicalStates{})),
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
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Select2SignalWords{})),
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
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Select2ClassesOfCompound{})),
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
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Select2Names{})),
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
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Select2EmpiricalFormulas{})),
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
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Select2LinearFormulas{})),
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
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Select2Names{})),
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
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Select2Symbols{})),
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
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Select2HazardStatements{})),
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
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Select2PrecautionaryStatements{})),
		},
	}).Select2ify()

	jquery.Jq("#product_twodformula").On("change", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		js.Global().Call("load2dimage")
		return nil
	}))
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
	// product_common()

	bstable.NewBootstraptable(jquery.Jq("#Product_table"), &bstable.BootstraptableParams{Ajax: "Product_getTableData"})
	jquery.Jq("#Product_table").On("load-success.bs.table", js.FuncOf(ShowIfAuthorizedActionButtons))

	jquery.Jq("#searchbar").Show()
	jquery.Jq("#actions").Show()
	jquery.Jq("#s_storage_archive_button").SetInvisible()
	jquery.Jq("#s_storage_stock_button").SetInvisible()
	jquery.Jq("#stock").SetHtml("")

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

func Product_pubchemCallback(args ...interface{}) {
	jquery.Jq("#searchbar").Hide()
	jquery.Jq("#actions").Hide()

	jquery.Jq("#search input").On("keyup", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
		if !event.Get("which").IsUndefined() && event.Get("which").Int() == 13 {

			event.Call("preventDefault")
		}

		return nil
	}))
}

func displaySection(section Section) {
	if section.Information != nil {
		for _, information := range *section.Information {

			jquery.Jq("#pubchemcompound").Append(
				widgets.NewSpan(widgets.SpanAttributes{
					BaseAttributes: widgets.BaseAttributes{
						Visible: true,
						Classes: []string{"iconlabel"},
					},
					Text: information.Name,
				}).OuterHTML())

			if information.Value.StringWithMarkup != nil {
				for _, value := range *information.Value.StringWithMarkup {
					jquery.Jq("#pubchemcompound").Append(
						widgets.NewSpan(widgets.SpanAttributes{
							BaseAttributes: widgets.BaseAttributes{
								Visible: true,
							},
							Text: value.String,
						}).OuterHTML())
				}
			}

			jquery.Jq("#pubchemcompound").Append(widgets.NewBr(widgets.BrAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
				},
			}).OuterHTML())

		}
	}

	if section.Section != nil {
		for _, sectionChild := range *section.Section {
			displaySection(sectionChild)
		}
	}
}

func PubchemCreateProduct(this js.Value, args []js.Value) interface{} {

	base64JsonPubchemProduct := args[0].String()

	var (
		jsonPubchemProduct []byte
		err                error
	)
	if jsonPubchemProduct, err = base64.StdEncoding.DecodeString(base64JsonPubchemProduct); err != nil {
		jsutils.DisplayGenericErrorMessage()
	}

	ajax.Ajax{
		URL:    ApplicationProxyPath + "products/pubchemcreateproduct",
		Method: "post",
		Data:   jsonPubchemProduct,
		Done: func(data js.Value) {
			var (
				err        error
				product_id int
			)

			if err = json.Unmarshal([]byte(data.String()), &product_id); err != nil {
				fmt.Println(err)
				return
			}

			href := fmt.Sprintf("%sv/products", ApplicationProxyPath)
			jsutils.ClearSearch(js.Null(), nil)
			jsutils.LoadContent("div#content", "product", href, Product_SaveCallback, product_id)
		},
		Fail: func(jqXHR js.Value) {

			jsutils.DisplayGenericErrorMessage()

		},
	}.Send()

	return nil
}

func PubchemGetProductByName(this js.Value, args []js.Value) interface{} {
	var (
		pubchemProduct PubchemProduct
		err            error
	)

	name := args[0].String()

	if name == "" {
		return nil
	}

	jquery.Jq("#pubchemcompound").Empty()
	jquery.Jq("#pubchemcompound").Append(`
	<div class="spinner-border" role="status"><span class="sr-only">Loading...</span></div>
	`)

	ajax.Ajax{
		URL:    ApplicationProxyPath + "products/pubchemgetproductbyname/" + name,
		Method: "get",
		Done: func(data js.Value) {
			if err = json.Unmarshal([]byte(data.String()), &pubchemProduct); err != nil {
				fmt.Println(err)
			}

			fmt.Println(data.String())

			base64JsonPubchem := base64.StdEncoding.EncodeToString([]byte(data.String()))

			jquery.Jq("#pubchemcompound").Empty()

			// import button.
			jquery.Jq("#pubchemcompound").Append(`<div id="import" class="row pb-3"></div>`)
			jquery.Jq("#pubchemcompound #import").Append(`<a href="#" onclick="Product_pubchemCreateProduct('` + base64JsonPubchem + `')">` + locales.Translate("import", HTTPHeaderAcceptLanguage) + `</a>`)

			// 2dpicture.
			jquery.Jq("#pubchemcompound").Append(`<div id="2dimage" class="row pb-3"></div>`)
			jquery.Jq("#pubchemcompound #2dimage").Append(`<img alt="2dpng" style="border: 1px solid grey;" title="2dpng" src="` + fmt.Sprintf("data:image/png;base64,%s", *pubchemProduct.Twodpicture) + `"></img>`)

			// Name.
			if pubchemProduct.Name != nil {
				jquery.Jq("#pubchemcompound").Append(`<div id="name" class="row pb-3"></div>`)
				jquery.Jq("#pubchemcompound #name").Append(`<div class="iconlabel col-sm-auto">name</div>`)
				jquery.Jq("#pubchemcompound #name").Append(`<div class="col-sm-auto">` + *pubchemProduct.Name + `</div>`)
			}

			// Molecular formula
			if pubchemProduct.MolecularFormula != nil {
				jquery.Jq("#pubchemcompound").Append(`<div id="molecular_formula" class="row pb-3"></div>`)
				jquery.Jq("#pubchemcompound #molecular_formula").Append(`<div class="iconlabel col-sm-auto">molecular_formula</div>`)
				jquery.Jq("#pubchemcompound #molecular_formula").Append(`<div class="col-sm-auto">` + *pubchemProduct.MolecularFormula + `</div>`)
			}

			// CAS
			if pubchemProduct.Cas != nil {
				jquery.Jq("#pubchemcompound").Append(`<div id="cas" class="row pb-3"></div>`)
				jquery.Jq("#pubchemcompound #cas").Append(`<div class="iconlabel col-sm-auto">cas</div>`)
				jquery.Jq("#pubchemcompound #cas").Append(`<div class="col-sm-auto">` + *pubchemProduct.Cas + `</div>`)
			}

			// EC
			if pubchemProduct.Ec != nil {
				jquery.Jq("#pubchemcompound").Append(`<div id="ec" class="row pb-3"></div>`)
				jquery.Jq("#pubchemcompound #ec").Append(`<div class="iconlabel col-sm-auto">ec</div>`)
				jquery.Jq("#pubchemcompound #ec").Append(`<div class="col-sm-auto">` + *pubchemProduct.Ec + `</div>`)
			}

			// Molecular weight
			if pubchemProduct.MolecularWeight != nil {

				mw := *pubchemProduct.MolecularWeight
				if pubchemProduct.MolecularWeightUnit != nil {
					mw += " " + *pubchemProduct.MolecularWeightUnit
				}

				jquery.Jq("#pubchemcompound").Append(`<div id="molecular_weight" class="row pb-3"></div>`)
				jquery.Jq("#pubchemcompound #molecular_weight").Append(`<div class="iconlabel col-sm-auto">molecular_weight</div>`)
				jquery.Jq("#pubchemcompound #molecular_weight").Append(`<div class="col-sm-auto">` + mw + `</div>`)
			}

			// Inchi
			if pubchemProduct.Inchi != nil {
				jquery.Jq("#pubchemcompound").Append(`<div id="inchi" class="row pb-3"></div>`)
				jquery.Jq("#pubchemcompound #inchi").Append(`<div class="iconlabel col-sm-auto">inchi</div>`)
				jquery.Jq("#pubchemcompound #inchi").Append(`<div class="col-sm-auto">` + *pubchemProduct.Inchi + `</div>`)
			}

			// Inchi key
			if pubchemProduct.InchiKey != nil {
				jquery.Jq("#pubchemcompound").Append(`<div id="inchi_key" class="row pb-3"></div>`)
				jquery.Jq("#pubchemcompound #inchi_key").Append(`<div class="iconlabel col-sm-auto">inchi_key</div>`)
				jquery.Jq("#pubchemcompound #inchi_key").Append(`<div class="col-sm-auto">` + *pubchemProduct.InchiKey + `</div>`)
			}

			// Canonical SMILES
			if pubchemProduct.CanonicalSmiles != nil {
				jquery.Jq("#pubchemcompound").Append(`<div id="canonical_smiles" class="row pb-3"></div>`)
				jquery.Jq("#pubchemcompound #canonical_smiles").Append(`<div class="iconlabel col-sm-auto">canonical_smiles</div>`)
				jquery.Jq("#pubchemcompound #canonical_smiles").Append(`<div class="col-sm-auto">` + *pubchemProduct.CanonicalSmiles + `</div>`)
			}

			// Symbols.
			if pubchemProduct.Symbols != nil {
				jquery.Jq("#pubchemcompound").Append(`<div id="symbols" class="row pb-3"></div>`)
				jquery.Jq("#pubchemcompound #symbols").Append(`<div class="iconlabel col-sm-auto">symbols</div>`)
				jquery.Jq("#pubchemcompound #symbols").Append(`<div id="symbols_content" class="col-sm-auto"></div>`)

				for _, sym := range *pubchemProduct.Symbols {
					jquery.Jq("#pubchemcompound #symbols #symbols_content").Append(`<img src='data:` + globals.SymbolImages[sym] + `' alt="symbol">`)
				}
			}

			// Signal.
			if pubchemProduct.Signal != nil {
				jquery.Jq("#pubchemcompound").Append(`<div id="signal" class="row pb-3"></div>`)
				jquery.Jq("#pubchemcompound #signal").Append(`<div class="iconlabel col-sm-auto">signal</div>`)
				jquery.Jq("#pubchemcompound #signal").Append(`<div id="signal_content" class="col-sm-auto"></div>`)

				for _, syn := range *pubchemProduct.Signal {
					jquery.Jq("#pubchemcompound #signal #signal_content").Append(`<li>` + syn + `</li>`)
				}
			}

			// HS.
			if pubchemProduct.Hs != nil {
				jquery.Jq("#pubchemcompound").Append(`<div id="hazard_statement" class="row pb-3"></div>`)
				jquery.Jq("#pubchemcompound #hazard_statement").Append(`<div class="iconlabel col-sm-auto">hazard statements</div>`)
				jquery.Jq("#pubchemcompound #hazard_statement").Append(`<div id="hazard_statement_content" class="col-sm-auto"></div>`)

				for _, hs := range *pubchemProduct.Hs {
					jquery.Jq("#pubchemcompound #hazard_statement #hazard_statement_content").Append(`<span class="badge badge-secondary mr-1">` + hs + `</span>`)
				}
			}

			// PS.
			if pubchemProduct.Ps != nil {
				jquery.Jq("#pubchemcompound").Append(`<div id="precautionary_statement" class="row pb-3"></div>`)
				jquery.Jq("#pubchemcompound #precautionary_statement").Append(`<div class="iconlabel col-sm-auto">precautionary statements</div>`)
				jquery.Jq("#pubchemcompound #precautionary_statement").Append(`<div id="precautionary_statement_content" class="col-sm-auto"></div>`)

				for _, ps := range *pubchemProduct.Ps {
					jquery.Jq("#pubchemcompound #precautionary_statement #precautionary_statement_content").Append(`<span class="badge badge-secondary mr-1">` + ps + `</span>`)
				}
			}

			// Synonyms.
			// if pubchemProduct.Synonyms != nil {
			// 	jquery.Jq("#pubchemcompound").Append(`<div id="synonyms" class="row pb-3"></div>`)
			// 	jquery.Jq("#pubchemcompound #synonyms").Append(`<div class="iconlabel col-sm-auto">synonyms</div>`)
			// 	jquery.Jq("#pubchemcompound #synonyms").Append(`<div id="synonyms_content" class="col-sm-auto"></div>`)

			// 	for _, syn := range *pubchemProduct.Synonyms {
			// 		jquery.Jq("#pubchemcompound #synonyms #synonyms_content").Append(`<li>` + syn) + `</li>`)
			// 	}
			// }
		},
		Fail: func(jqXHR js.Value) {
			jsutils.DisplayGenericErrorMessage()
		},
	}.Send()
	return nil
}

func PubchemGetCompoundByName(this js.Value, args []js.Value) interface{} {
	var (
		compounds Compounds
		err       error
	)

	name := args[0].String()

	if name == "" {
		return nil
	}

	ajax.Ajax{
		URL:    ApplicationProxyPath + "products/pubchemgetcompoundbyname/" + name,
		Method: "get",
		Done: func(data js.Value) {
			if err = json.Unmarshal([]byte(data.String()), &compounds); err != nil {
				fmt.Println(err)
			}

			jquery.Jq("#pubchemcompound").Empty()

			for _, pccompound := range compounds.PCCompounds {
				jquery.Jq("#pubchemcompound").Append(
					widgets.NewImg(widgets.ImgAttributes{
						BaseAttributes: widgets.BaseAttributes{
							Visible: true,
						},
						Src:   fmt.Sprintf("data:image/png;base64,%s", compounds.Base64Png),
						Alt:   "2dpng",
						Title: "2dpng",
					}).OuterHTML())
				jquery.Jq("#pubchemcompound").Append(widgets.NewBr(widgets.BrAttributes{
					BaseAttributes: widgets.BaseAttributes{
						Visible: true,
					},
				}).OuterHTML())

				jquery.Jq("#pubchemcompound").Append(
					widgets.NewSpan(widgets.SpanAttributes{
						BaseAttributes: widgets.BaseAttributes{
							Visible: true,
							Classes: []string{"iconlabel"},
						},
						Text: "cid ",
					}).OuterHTML())
				jquery.Jq("#pubchemcompound").Append(
					widgets.NewSpan(widgets.SpanAttributes{
						BaseAttributes: widgets.BaseAttributes{
							Visible: true,
						},
						Text: fmt.Sprint(pccompound.ID.ID.CID),
					}).OuterHTML())
				jquery.Jq("#pubchemcompound").Append(widgets.NewBr(widgets.BrAttributes{
					BaseAttributes: widgets.BaseAttributes{
						Visible: true,
					},
				}).OuterHTML())

				for _, section := range compounds.Record.Record.Section {
					displaySection(section)
				}

				for _, prop := range pccompound.Props {

					propval := ""
					if prop.Value.Ival != nil {
						propval = fmt.Sprint(*prop.Value.Ival)
					} else if prop.Value.Fval != nil {
						propval = strconv.FormatFloat(*prop.Value.Fval, 'f', 64, 64)
					} else if prop.Value.Sval != nil {
						propval = *prop.Value.Sval
					} else {
						propval = *prop.Value.Binary
					}

					jquery.Jq("#pubchemcompound").Append(
						widgets.NewSpan(widgets.SpanAttributes{
							BaseAttributes: widgets.BaseAttributes{
								Visible: true,
								Classes: []string{"iconlabel"},
							},
							Text: prop.URN.Name + " " + prop.URN.Label + " ",
						}).OuterHTML())
					jquery.Jq("#pubchemcompound").Append(
						widgets.NewSpan(widgets.SpanAttributes{
							BaseAttributes: widgets.BaseAttributes{
								Visible: true,
							},
							Text: propval,
						}).OuterHTML())
					jquery.Jq("#pubchemcompound").Append(widgets.NewBr(widgets.BrAttributes{
						BaseAttributes: widgets.BaseAttributes{
							Visible: true,
						},
					}).OuterHTML())
				}
			}
		},
		Fail: func(jqXHR js.Value) {
			jsutils.DisplayGenericErrorMessage()
		},
	}.Send()
	return nil
}

func PubchemSearch(this js.Value, args []js.Value) interface{} {
	var (
		autocomplete Autocomplete
		err          error
	)

	search := jquery.Jq("input#searchpubchem").GetVal().String()
	search = strings.Trim(search, " ")

	if search == "" {
		return nil
	}

	jquery.Jq("#pubchemsearchresult").Append(`
	<div class="spinner-border" role="status"><span class="sr-only">Loading...</span></div>
	`)

	ajax.Ajax{
		URL:    ApplicationProxyPath + "products/pubchemautocomplete/" + search,
		Method: "get",
		Done: func(data js.Value) {
			if err = json.Unmarshal([]byte(data.String()), &autocomplete); err != nil {
				fmt.Println(err)
			}

			jquery.Jq("#pubchemsearchresult").Empty()
			jquery.Jq("#pubchemcompound").Empty()

			for _, compound := range autocomplete.DictionaryTerms.Compound {
				jquery.Jq("#pubchemcompound").Append(`<div class="row"><div class="col-sm-auto mx-auto"><a href="#" onclick="Product_pubchemGetProductByName('` + compound + `')">` + compound + `</a></div></div>`)
			}
		},
		Fail: func(jqXHR js.Value) {
			jsutils.DisplayGenericErrorMessage()
		},
	}.Send()
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

		jquery.Jq("#searchbar").Hide()
		jquery.Jq("#actions").Hide()
	}

	// Chemical product by default on creation.
	if jquery.Jq("input#product_id").GetVal().String() == "" {
		Chemify()
	} else {
		jquery.Jq("input#showchem").SetProp("disabled", "disabled")
		jquery.Jq("input#showbio").SetProp("disabled", "disabled")
	}

	jquery.Jq("#searchbar").Hide()
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

	// product_common()

	jquery.Jq("#searchbar").Show()
	jquery.Jq("#actions").Show()
}
