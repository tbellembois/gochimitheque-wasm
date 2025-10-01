//go:build go1.24 && js && wasm

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

var (
	supplierrefToSupplier map[string]int64 // supplierref label -> supplier id
)

func init() {
	supplierrefToSupplier = make(map[string]int64)
}

func LinearToEmpirical(this js.Value, args []js.Value) interface{} {

	select2LinearFormula := select2.NewSelect2(jquery.Jq("select#linear_formula"), nil).Select2Data()

	if len(select2LinearFormula) == 0 {
		return ""
	}

	ajaxData := struct {
		EmpiricalFormula string `json:"empirical_formula"`
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
		URL:    fmt.Sprintf("%sformat/empiricalformula/", ApplicationProxyPath),
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

	validate.NewValidate(jquery.Jq("select#empirical_formula"), nil).ValidateRemoveRequired()
	jquery.Jq("span#empirical_formula.badge").Hide()

	return nil

}

func NoCas(this js.Value, args []js.Value) interface{} {

	validate.NewValidate(jquery.Jq("select#cas_number"), nil).ValidateRemoveRequired()
	jquery.Jq("span#cas_number.badge").Hide()

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
	rsymbols := regexp.MustCompile("(GHS[0-9]{2})")

	shs := rhs.FindAllStringSubmatch(magic, -1)
	sps := rps.FindAllStringSubmatch(magic, -1)
	ssymbols := rsymbols.FindAllStringSubmatch(magic, -1)

	var (
		processedH       map[string]string
		processedP       map[string]string
		processedSymbols map[string]string
		ok               bool
	)

	processedH = make(map[string]string)

	select2HS := select2.NewSelect2(jquery.Jq("select#hazard_statements"), nil)
	select2HS.Select2Clear()
	for _, h := range shs {

		if _, ok = processedH[h[1]]; !ok {
			processedH[h[1]] = ""

			for _, hs := range globals.DBHazardStatements {
				if hs.HazardStatementReference == h[1] {
					select2HS.Select2AppendOption(
						widgets.NewOption(widgets.OptionAttributes{
							Text:            h[0],
							Value:           strconv.Itoa(int(*hs.HazardStatementID)),
							DefaultSelected: true,
							Selected:        true,
						}).HTMLElement.OuterHTML())
					break
				}
			}

		}
	}

	processedP = make(map[string]string)

	select2PS := select2.NewSelect2(jquery.Jq("select#precautionary_statements"), nil)
	select2PS.Select2Clear()
	for _, p := range sps {

		if _, ok = processedP[p[1]]; !ok {
			processedP[p[1]] = ""

			for _, hs := range globals.DBPrecautionaryStatements {
				if hs.PrecautionaryStatementReference == p[1] {
					select2PS.Select2AppendOption(
						widgets.NewOption(widgets.OptionAttributes{
							Text:            p[0],
							Value:           strconv.Itoa(int(*hs.PrecautionaryStatementID)),
							DefaultSelected: true,
							Selected:        true,
						}).HTMLElement.OuterHTML())
					break
				}
			}

		}
	}

	processedSymbols = make(map[string]string)

	select2Symbols := select2.NewSelect2(jquery.Jq("select#symbols"), nil)
	select2Symbols.Select2Clear()
	for _, p := range ssymbols {

		if _, ok = processedSymbols[p[1]]; !ok {
			processedSymbols[p[1]] = ""

			for _, symbol := range globals.DBSymbols {
				if symbol.SymbolLabel == p[1] {
					select2Symbols.Select2AppendOption(
						widgets.NewOption(widgets.OptionAttributes{
							Text:            p[0],
							Value:           strconv.Itoa(int(*symbol.SymbolID)),
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
			"producer_ref": {
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
			"empirical_formula": {
				Required: js.FuncOf(func(this js.Value, args []js.Value) interface{} { return true }),
				Remote: validate.ValidateRemote{
					URL:        "",
					Type:       "post",
					BeforeSend: js.FuncOf(ValidateProductEmpiricalFormulaBeforeSend),
					Data: map[string]interface{}{
						"empirical_formula": js.FuncOf(ValidateProductEmpiricalFormulaData),
					},
				},
			},
			"cas_number": {
				Required: js.FuncOf(func(this js.Value, args []js.Value) interface{} { return true }),
				Remote: validate.ValidateRemote{
					URL:        "",
					Type:       "post",
					BeforeSend: js.FuncOf(ValidateProductCasNumberBeforeSend),
					Data: map[string]interface{}{
						"cas_number":          js.FuncOf(ValidateProductCasNumberData1),
						"product_specificity": js.FuncOf(ValidateProductCasNumberData2),
					},
				},
			},
			"ce_number": {
				Remote: validate.ValidateRemote{
					URL:        "",
					Type:       "post",
					BeforeSend: js.FuncOf(ValidateProductCeNumberBeforeSend),
					Data: map[string]interface{}{
						"ce_number": js.FuncOf(ValidateProductCeNumberData),
					},
				},
			},
		},
		Messages: map[string]validate.ValidateMessage{
			"name": {
				Required: locales.Translate("required_input", HTTPHeaderAcceptLanguage),
			},
			"empirical_formula": {
				Required: locales.Translate("required_input", HTTPHeaderAcceptLanguage),
			},
			"cas_number": {
				Required: locales.Translate("required_input", HTTPHeaderAcceptLanguage),
			},
			"producer_ref": {
				Required: locales.Translate("required_input", HTTPHeaderAcceptLanguage),
			},
		},
	}).Validate()

	// select2
	select2.NewSelect2(jquery.Jq("select#producer_ref"), &select2.Select2Config{
		Placeholder:       locales.Translate("product_producer_ref_placeholder", HTTPHeaderAcceptLanguage),
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
	jquery.Jq("select#producer_ref").On("change", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return validate.NewValidate(jquery.Jq(this), nil).Valid()
	}))
	jquery.Jq("select#producer_ref").On("select2:select", js.FuncOf(func(this js.Value, args []js.Value) interface{} {

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
				Text:            *producerref.Producer.ProducerLabel,
				Value:           strconv.Itoa(int(*producerref.Producer.ProducerID)),
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

	select2.NewSelect2(jquery.Jq("select#supplier_refs"), &select2.Select2Config{
		Placeholder:       locales.Translate("product_supplier_ref_placeholder", HTTPHeaderAcceptLanguage),
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

	select2.NewSelect2(jquery.Jq("select#unit_molecularweight"), &select2.Select2Config{
		Placeholder:    locales.Translate("product_unit_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(Unit{})),
		AllowClear:     true,
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "storages/units",
			DataType:       "json",
			Data:           js.FuncOf(Select2UnitMolecularWeightAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Units{})),
		},
	}).Select2ify()
	jquery.Jq("select#unit_molecularweight").On("change", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return validate.NewValidate(jquery.Jq(this), nil).Valid()
	}))

	select2.NewSelect2(jquery.Jq("select#cas_number"), &select2.Select2Config{
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
	jquery.Jq("select#cas_number").On("change", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return validate.NewValidate(jquery.Jq(this), nil).Valid()
	}))

	select2.NewSelect2(jquery.Jq("select#ce_number"), &select2.Select2Config{
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
	jquery.Jq("select#ce_number").On("change", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return validate.NewValidate(jquery.Jq(this), nil).Valid()
	}))

	select2.NewSelect2(jquery.Jq("select#physical_state"), &select2.Select2Config{
		Placeholder:    locales.Translate("product_physical_state_placeholder", HTTPHeaderAcceptLanguage),
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

	select2.NewSelect2(jquery.Jq("select#signal_word"), &select2.Select2Config{
		Placeholder:    locales.Translate("product_signal_word_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(SignalWord{})),
		AllowClear:     true,
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "products/signalwords/",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(SignalWords{})),
		},
	}).Select2ify()

	select2.NewSelect2(jquery.Jq("select#class_of_compound"), &select2.Select2Config{
		Placeholder:    locales.Translate("product_class_of_compound_placeholder", HTTPHeaderAcceptLanguage),
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

	select2.NewSelect2(jquery.Jq("select#empirical_formula"), &select2.Select2Config{
		Placeholder:    locales.Translate("product_empirical_formula_placeholder", HTTPHeaderAcceptLanguage),
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
	jquery.Jq("select#empirical_formula").On("change", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return validate.NewValidate(jquery.Jq(this), nil).Valid()
	}))

	select2.NewSelect2(jquery.Jq("select#linear_formula"), &select2.Select2Config{
		Placeholder:    locales.Translate("product_linear_formula_placeholder", HTTPHeaderAcceptLanguage),
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

	select2.NewSelect2(jquery.Jq("select#hazard_statements"), &select2.Select2Config{
		Placeholder:    locales.Translate("product_hazard_statements_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2HazardStatementTemplateResults),
		AllowClear:     true,
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "products/hazardstatements/",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(HazardStatements{})),
		},
	}).Select2ify()

	select2.NewSelect2(jquery.Jq("select#precautionary_statements"), &select2.Select2Config{
		Placeholder:    locales.Translate("product_precautionary_statements_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2PrecautionaryStatementTemplateResults),
		AllowClear:     true,
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "products/precautionarystatements/",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(PrecautionaryStatements{})),
		},
	}).Select2ify()

	jquery.Jq("#product_twod_formula").On("change", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
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

	//product_common()

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

	// Works only with no select2.
	jquery.Jq("input[name='searchpubchem']").On("keyup", js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		event := args[0]

		if !event.Get("keyCode").IsUndefined() && event.Get("keyCode").Int() == 13 {

			event.Call("preventDefault")
			PubchemSearch(js.Null(), nil)

		}

		return nil

	}))

}

func displaySection(section Section) {

	if section.Information != nil {
		for _, information := range *section.Information {

			jquery.Jq("#pubchemcompoundcontent").Append(
				widgets.NewSpan(widgets.SpanAttributes{
					BaseAttributes: widgets.BaseAttributes{
						Visible: true,
						Classes: []string{"iconlabel"},
					},
					Text: information.Name,
				}).OuterHTML())

			if information.Value.StringWithMarkup != nil {
				for _, value := range *information.Value.StringWithMarkup {
					jquery.Jq("#pubchemcompoundcontent").Append(
						widgets.NewSpan(widgets.SpanAttributes{
							BaseAttributes: widgets.BaseAttributes{
								Visible: true,
							},
							Text: value.String,
						}).OuterHTML())
				}
			}

			jquery.Jq("#pubchemcompoundcontent").Append(widgets.NewBr(widgets.BrAttributes{
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

func PubchemUpdateProduct(this js.Value, args []js.Value) interface{} {

	base64JsonPubchemProduct := args[0].String()
	product_id := args[1].String()

	var (
		jsonPubchemProduct []byte
		err                error
	)
	if jsonPubchemProduct, err = base64.StdEncoding.DecodeString(base64JsonPubchemProduct); err != nil {
		jsutils.DisplayGenericErrorMessage()
	}

	pubChemProduct := struct {
		Name string
	}{}

	err = json.Unmarshal(jsonPubchemProduct, &pubChemProduct)
	if err != nil {
		js.Global().Get("console").Call("log", fmt.Sprintf("%#v", string(jsonPubchemProduct)))

	}
	globals.CurrentProduct.Name = &models.Name{NameLabel: pubChemProduct.Name}

	ajax.Ajax{
		URL:    ApplicationProxyPath + "products/pubchemcreateproduct/" + product_id,
		Method: "post",
		Data:   jsonPubchemProduct,
		Done: func(data js.Value) {
			var (
				err        error
				product_id int64
			)

			if err = json.Unmarshal([]byte(data.String()), &product_id); err != nil {
				fmt.Println(err)
				return
			}

			globals.CurrentProduct.ProductID = &product_id

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

func PubchemCreateProduct(this js.Value, args []js.Value) interface{} {

	base64JsonPubchemProduct := args[0].String()

	var (
		jsonPubchemProduct []byte
		err                error
	)
	if jsonPubchemProduct, err = base64.StdEncoding.DecodeString(base64JsonPubchemProduct); err != nil {
		jsutils.DisplayGenericErrorMessage()
	}

	pubChemProduct := struct {
		Name string
	}{}

	err = json.Unmarshal(jsonPubchemProduct, &pubChemProduct)
	if err != nil {
		js.Global().Get("console").Call("log", fmt.Sprintf("%#v", string(jsonPubchemProduct)))

	}
	globals.CurrentProduct.Name = &models.Name{NameLabel: pubChemProduct.Name}

	ajax.Ajax{
		URL:    ApplicationProxyPath + "products/pubchemcreateproduct",
		Method: "post",
		Data:   jsonPubchemProduct,
		Done: func(data js.Value) {
			var (
				err        error
				product_id int64
			)

			if err = json.Unmarshal([]byte(data.String()), &product_id); err != nil {
				fmt.Println(err)
				return
			}

			globals.CurrentProduct.ProductID = &product_id

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

	jquery.Jq("#pubchemcompoundcontent").Empty()
	jquery.Jq("#pubchemcompoundcontent").Append(`
	<div class="spinner-border" role="status"><span class="sr-only">Loading...</span></div>
	`)

	ajax.Ajax{
		URL:    ApplicationProxyPath + "products/pubchemgetproductbyname/" + name,
		Method: "get",
		Done: func(data js.Value) {

			if err = json.Unmarshal([]byte(data.String()), &pubchemProduct); err != nil {
				fmt.Println(err)
			}

			if pubchemProduct.Cas != nil {

				ajax.Ajax{
					URL:    ApplicationProxyPath + "products?cas_number_string=" + *pubchemProduct.Cas,
					Method: "get",
					Done: func(data js.Value) {
						var (
							products Products
							err      error
						)
						if err = json.Unmarshal([]byte(data.String()), &products); err != nil {
							fmt.Println(err)
						}

						if products.GetTotal() != 0 {
							jquery.Jq("#pubchemcasexist").Empty()
							jquery.Jq("#pubchemcasexist").Append(`<div class="alert alert-danger" role="alert">` + locales.Translate("cas_number_validate_cas", HTTPHeaderAcceptLanguage) + `</div>`)
						}
					},
				}.Send()

			}
			// fmt.Println(data.String())

			base64JsonPubchem := base64.StdEncoding.EncodeToString([]byte(data.String()))

			jquery.Jq("#pubchemcompoundactions").Empty()
			jquery.Jq("#pubchemcompoundcontent").Empty()

			jsutils.HasPermission("products", "-2", "put", func() {

				// import button.
				jquery.Jq("#pubchemcompoundactions").Append(`<div id="import" class="row pb-3"></div>`)
				jquery.Jq("#pubchemcompoundactions #import").Append(`<button type="buton" class="btn btn-primary" href="#" onclick="Product_pubchemCreateProduct('` + base64JsonPubchem + `')">` + locales.Translate("import", HTTPHeaderAcceptLanguage) + `</button>`)

				if jquery.Jq("input[name='selected_product_id']").GetVal().String() != "" {
					// replace button.
					jquery.Jq("#pubchemcompoundactions").Append(`<div id="replace" class="row pb-3"></div>`)
					jquery.Jq("#pubchemcompoundactions #replace").Append(`<button type="buton" class="btn btn-primary" href="#" onclick="Product_pubchemUpdateProduct('` + base64JsonPubchem + `', '` + jquery.Jq("input[name='selected_product_id']").GetVal().String() + `')">` + locales.Translate("replace", HTTPHeaderAcceptLanguage) + " " + jquery.Jq("input[name='selected_product_name']").GetVal().String() + `</button>`)
				}

			}, func() {
			})

			// 2dpicture.
			jquery.Jq("#pubchemcompoundcontent").Append(`<div id="2dimage" class="row pb-3"></div>`)
			jquery.Jq("#pubchemcompoundcontent #2dimage").Append(`<img alt="2dpng" style="border: 1px solid grey;" title="2dpng" src="` + fmt.Sprintf("data:image/png;base64,%s", *pubchemProduct.Twodpicture) + `"></img>`)

			// Name.
			if pubchemProduct.Name != nil {
				jquery.Jq("#pubchemcompoundcontent").Append(`<div id="name" class="row pb-3"></div>`)
				jquery.Jq("#pubchemcompoundcontent #name").Append(`<div class="iconlabel col-sm-auto">name</div>`)
				jquery.Jq("#pubchemcompoundcontent #name").Append(`<div class="col-sm-auto">` + *pubchemProduct.Name + `</div>`)
			}

			// Molecular formula
			if pubchemProduct.MolecularFormula != nil {
				jquery.Jq("#pubchemcompoundcontent").Append(`<div id="molecular_formula" class="row pb-3"></div>`)
				jquery.Jq("#pubchemcompoundcontent #molecular_formula").Append(`<div class="iconlabel col-sm-auto">molecular_formula</div>`)
				jquery.Jq("#pubchemcompoundcontent #molecular_formula").Append(`<div class="col-sm-auto">` + *pubchemProduct.MolecularFormula + `</div>`)
			}

			// CAS
			if pubchemProduct.Cas != nil {
				jquery.Jq("#pubchemcompoundcontent").Append(`<div id="cas" class="row pb-3"></div>`)
				jquery.Jq("#pubchemcompoundcontent #cas").Append(`<div class="iconlabel col-sm-auto">cas</div>`)
				jquery.Jq("#pubchemcompoundcontent #cas").Append(`<div class="col-sm-auto">` + *pubchemProduct.Cas + `</div>`)
			}

			// EC
			if pubchemProduct.Ec != nil {
				jquery.Jq("#pubchemcompoundcontent").Append(`<div id="ec" class="row pb-3"></div>`)
				jquery.Jq("#pubchemcompoundcontent #ec").Append(`<div class="iconlabel col-sm-auto">ec</div>`)
				jquery.Jq("#pubchemcompoundcontent #ec").Append(`<div class="col-sm-auto">` + *pubchemProduct.Ec + `</div>`)
			}

			// Molecular weight
			if pubchemProduct.MolecularWeight != nil {

				mw := *pubchemProduct.MolecularWeight
				if pubchemProduct.MolecularWeightUnit != nil {
					mw += " " + *pubchemProduct.MolecularWeightUnit
				}

				jquery.Jq("#pubchemcompoundcontent").Append(`<div id="molecular_weight" class="row pb-3"></div>`)
				jquery.Jq("#pubchemcompoundcontent #molecular_weight").Append(`<div class="iconlabel col-sm-auto">molecular_weight</div>`)
				jquery.Jq("#pubchemcompoundcontent #molecular_weight").Append(`<div class="col-sm-auto">` + mw + `</div>`)
			}

			// Inchi
			if pubchemProduct.Inchi != nil {
				jquery.Jq("#pubchemcompoundcontent").Append(`<div id="inchi" class="row pb-3"></div>`)
				jquery.Jq("#pubchemcompoundcontent #inchi").Append(`<div class="iconlabel col-sm-auto">inchi</div>`)
				jquery.Jq("#pubchemcompoundcontent #inchi").Append(`<div class="col-sm-auto">` + *pubchemProduct.Inchi + `</div>`)
			}

			// Inchi key
			if pubchemProduct.InchiKey != nil {
				jquery.Jq("#pubchemcompoundcontent").Append(`<div id="inchi_key" class="row pb-3"></div>`)
				jquery.Jq("#pubchemcompoundcontent #inchi_key").Append(`<div class="iconlabel col-sm-auto">inchi_key</div>`)
				jquery.Jq("#pubchemcompoundcontent #inchi_key").Append(`<div class="col-sm-auto">` + *pubchemProduct.InchiKey + `</div>`)
			}

			// Canonical SMILES
			if pubchemProduct.CanonicalSmiles != nil {
				jquery.Jq("#pubchemcompoundcontent").Append(`<div id="canonical_smiles" class="row pb-3"></div>`)
				jquery.Jq("#pubchemcompoundcontent #canonical_smiles").Append(`<div class="iconlabel col-sm-auto">canonical_smiles</div>`)
				jquery.Jq("#pubchemcompoundcontent #canonical_smiles").Append(`<div class="col-sm-auto">` + *pubchemProduct.CanonicalSmiles + `</div>`)
			}

			// Symbols.
			if pubchemProduct.Symbols != nil {
				jquery.Jq("#pubchemcompoundcontent").Append(`<div id="symbols" class="row pb-3"></div>`)
				jquery.Jq("#pubchemcompoundcontent #symbols").Append(`<div class="iconlabel col-sm-auto">symbols</div>`)
				jquery.Jq("#pubchemcompoundcontent #symbols").Append(`<div id="symbols_content" class="col-sm-auto"></div>`)

				for _, sym := range *pubchemProduct.Symbols {
					jquery.Jq("#pubchemcompoundcontent #symbols #symbols_content").Append(`<img width="30" height="30" src='` + fmt.Sprintf("%sstatic/img/%s.svg", ApplicationProxyPath, sym) + `' alt="symbol">`)
				}
			}

			// Signal.
			if pubchemProduct.Signal != nil {
				jquery.Jq("#pubchemcompoundcontent").Append(`<div id="signal" class="row pb-3"></div>`)
				jquery.Jq("#pubchemcompoundcontent #signal").Append(`<div class="iconlabel col-sm-auto">signal</div>`)
				jquery.Jq("#pubchemcompoundcontent #signal").Append(`<div id="signal_content" class="col-sm-auto"></div>`)

				for _, syn := range *pubchemProduct.Signal {
					jquery.Jq("#pubchemcompoundcontent #signal #signal_content").Append(`<li>` + syn + `</li>`)
				}
			}

			// HS.
			if pubchemProduct.Hs != nil {
				jquery.Jq("#pubchemcompoundcontent").Append(`<div id="hazard_statement" class="row pb-3"></div>`)
				jquery.Jq("#pubchemcompoundcontent #hazard_statement").Append(`<div class="iconlabel col-sm-auto">hazard statements</div>`)
				jquery.Jq("#pubchemcompoundcontent #hazard_statement").Append(`<div id="hazard_statement_content" class="col-sm-auto"></div>`)

				for _, hs := range *pubchemProduct.Hs {
					jquery.Jq("#pubchemcompoundcontent #hazard_statement #hazard_statement_content").Append(`<span class="badge badge-secondary mr-1">` + hs + `</span>`)
				}
			}

			// PS.
			if pubchemProduct.Ps != nil {
				jquery.Jq("#pubchemcompoundcontent").Append(`<div id="precautionary_statement" class="row pb-3"></div>`)
				jquery.Jq("#pubchemcompoundcontent #precautionary_statement").Append(`<div class="iconlabel col-sm-auto">precautionary statements</div>`)
				jquery.Jq("#pubchemcompoundcontent #precautionary_statement").Append(`<div id="precautionary_statement_content" class="col-sm-auto"></div>`)

				for _, ps := range *pubchemProduct.Ps {
					jquery.Jq("#pubchemcompoundcontent #precautionary_statement #precautionary_statement_content").Append(`<span class="badge badge-secondary mr-1">` + ps + `</span>`)
				}
			}

			// Synonyms.
			// if pubchemProduct.Synonyms != nil {
			// 	jquery.Jq("#pubchemcompoundcontent").Append(`<div id="synonyms" class="row pb-3"></div>`)
			// 	jquery.Jq("#pubchemcompoundcontent #synonyms").Append(`<div class="iconlabel col-sm-auto">synonyms</div>`)
			// 	jquery.Jq("#pubchemcompoundcontent #synonyms").Append(`<div id="synonyms_content" class="col-sm-auto"></div>`)

			// 	for _, syn := range *pubchemProduct.Synonyms {
			// 		jquery.Jq("#pubchemcompoundcontent #synonyms #synonyms_content").Append(`<li>` + syn) + `</li>`)
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

			jquery.Jq("#pubchemcompoundcontent").Empty()

			for _, pccompound := range compounds.PCCompounds {
				jquery.Jq("#pubchemcompoundcontent").Append(
					widgets.NewImg(widgets.ImgAttributes{
						BaseAttributes: widgets.BaseAttributes{
							Visible: true,
						},
						Src:   fmt.Sprintf("%s", compounds.Base64Png),
						Alt:   "2dpng",
						Title: "2dpng",
					}).OuterHTML())
				jquery.Jq("#pubchemcompoundcontent").Append(widgets.NewBr(widgets.BrAttributes{
					BaseAttributes: widgets.BaseAttributes{
						Visible: true,
					},
				}).OuterHTML())

				jquery.Jq("#pubchemcompoundcontent").Append(
					widgets.NewSpan(widgets.SpanAttributes{
						BaseAttributes: widgets.BaseAttributes{
							Visible: true,
							Classes: []string{"iconlabel"},
						},
						Text: "cid ",
					}).OuterHTML())
				jquery.Jq("#pubchemcompoundcontent").Append(
					widgets.NewSpan(widgets.SpanAttributes{
						BaseAttributes: widgets.BaseAttributes{
							Visible: true,
						},
						Text: fmt.Sprint(pccompound.ID.ID.CID),
					}).OuterHTML())
				jquery.Jq("#pubchemcompoundcontent").Append(widgets.NewBr(widgets.BrAttributes{
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

					jquery.Jq("#pubchemcompoundcontent").Append(
						widgets.NewSpan(widgets.SpanAttributes{
							BaseAttributes: widgets.BaseAttributes{
								Visible: true,
								Classes: []string{"iconlabel"},
							},
							Text: prop.URN.Name + " " + prop.URN.Label + " ",
						}).OuterHTML())
					jquery.Jq("#pubchemcompoundcontent").Append(
						widgets.NewSpan(widgets.SpanAttributes{
							BaseAttributes: widgets.BaseAttributes{
								Visible: true,
							},
							Text: propval,
						}).OuterHTML())
					jquery.Jq("#pubchemcompoundcontent").Append(widgets.NewBr(widgets.BrAttributes{
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
			jquery.Jq("#pubchemcompoundcontent").Empty()

			for _, compound := range autocomplete.DictionaryTerms.Compound {
				jquery.Jq("#pubchemcompoundcontent").Append(`<div class="row"><div class="col-sm-auto mx-auto"><a href="#" onclick="Product_pubchemGetProductByName('` + compound + `')">` + compound + `</a></div></div>`)
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
	BSTableQueryFilter.QueryFilter.Id = int(args[0].(int64))

	if CurrentProduct.ProductSpecificity != nil {
		BSTableQueryFilter.QueryFilter.ProductFilterLabel = fmt.Sprintf("#%d %s %s", *CurrentProduct.ProductID, CurrentProduct.Name.NameLabel, *CurrentProduct.ProductSpecificity)
	} else {
		BSTableQueryFilter.QueryFilter.ProductFilterLabel = fmt.Sprintf("#%d %s", *CurrentProduct.ProductID, CurrentProduct.Name.NameLabel)
	}

	bstable.NewBootstraptable(jquery.Jq("#Product_table"), nil).Refresh(nil)

	// jquery.Jq("#Product_table").Bootstraptable(nil).Refresh(&BootstraptableRefreshQuery{
	// 	Query: QueryFilter{
	// 		Product: strconv.Itoa(args[0].(int)),
	// 	},
	// })
	jquery.Jq("#Product_table").On("load-success.bs.table", js.FuncOf(ShowIfAuthorizedActionButtons))

	//product_common()

	jquery.Jq("#searchbar").Show()
	jquery.Jq("#actions").Show()

}
