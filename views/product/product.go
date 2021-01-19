package product

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"syscall/js"

	. "github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/locales"
	"github.com/tbellembois/gochimitheque-wasm/types"
	. "github.com/tbellembois/gochimitheque-wasm/types"
	"github.com/tbellembois/gochimitheque-wasm/utils"
	"github.com/tbellembois/gochimitheque-wasm/views/search"
	"github.com/tbellembois/gochimitheque-wasm/widgets"
	"github.com/tbellembois/gochimitheque-wasm/widgets/themes"
)

var (
	supplierrefToSupplier map[string]int64 // supplierref label -> supplier id
	err                   error
)

func init() {
	supplierrefToSupplier = make(map[string]int64)
}

func LinearToEmpirical(this js.Value, args []js.Value) interface{} {

	select2LinearFormula := Jq("select#linearformula").Select2Data()

	if len(select2LinearFormula) == 0 {
		return ""
	}

	Jq("#convertedEmpiricalFormula").Append(utils.LinearToEmpiricalFormula(select2LinearFormula[0].Text))

	return nil

}

func NoEmpiricalFormula(this js.Value, args []js.Value) interface{} {

	Jq("select#empiricalformula").ValidateRemoveRequired()
	Jq("span#empiricalformula.badge").Hide()

	return nil

}

func NoCas(this js.Value, args []js.Value) interface{} {

	Jq("select#casnumber").ValidateRemoveRequired()
	Jq("span#casnumber.badge").Hide()

	return nil

}

func HowToMagicalSelector(this js.Value, args []js.Value) interface{} {

	Win.Call("open", fmt.Sprintf("%simg/magicalselector.webm", ApplicationProxyPath), "_blank")

	return nil

}

func Magic(this js.Value, args []js.Value) interface{} {

	magic := Jq("textarea#magical").GetVal().String()

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
	Jq("select#hazardstatements").Select2Clear()
	for _, h := range shs {

		if _, ok = processedH[h[1]]; !ok {
			processedH[h[1]] = ""

			for _, hs := range types.DBHazardStatements {
				if hs.HazardStatementReference == h[1] {
					Jq("select#hazardstatements").Select2AppendOption(
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
	Jq("select#precautionarystatements").Select2Clear()
	for _, p := range sps {

		if _, ok = processedP[p[1]]; !ok {
			processedP[p[1]] = ""

			for _, hs := range types.DBPrecautionaryStatements {
				if hs.PrecautionaryStatementReference == p[1] {
					Jq("select#precautionarystatements").Select2AppendOption(
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

// TODO: factorise with storage
func product_common() {

	//
	// search form
	//
	Jq("select#s_storelocation").Select2(Select2Config{
		Placeholder:    locales.Translate("s_storelocation", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2StoreLocationTemplateResults),
		AllowClear:     true,
		Ajax: Select2Ajax{
			URL:            ApplicationProxyPath + "storelocations",
			DataType:       "json",
			Data:           js.FuncOf(Select2GenericAjaxData),
			ProcessResults: js.FuncOf(Select2GenericAjaxProcessResults(StoreLocations{})),
		},
	})

	Jq("select#s_casnumber").Select2(Select2Config{
		Placeholder:    locales.Translate("s_casnumber", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2GenericTemplateResults(CasNumber{})),
		AllowClear:     true,
		Ajax: Select2Ajax{
			URL:            ApplicationProxyPath + "products/casnumbers/",
			DataType:       "json",
			Data:           js.FuncOf(Select2GenericAjaxData),
			ProcessResults: js.FuncOf(Select2GenericAjaxProcessResults(CasNumbers{})),
		},
	})

	Jq("select#s_name").Select2(Select2Config{
		Placeholder:    locales.Translate("s_name", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2GenericTemplateResults(Name{})),
		AllowClear:     true,
		Ajax: Select2Ajax{
			URL:            ApplicationProxyPath + "products/names/",
			DataType:       "json",
			Data:           js.FuncOf(Select2GenericAjaxData),
			ProcessResults: js.FuncOf(Select2GenericAjaxProcessResults(Names{})),
		},
	})

	Jq("select#s_empiricalformula").Select2(Select2Config{
		Placeholder:    locales.Translate("s_empiricalformula", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2GenericTemplateResults(EmpiricalFormula{})),
		AllowClear:     true,
		Ajax: Select2Ajax{
			URL:            ApplicationProxyPath + "products/empiricalformulas/",
			DataType:       "json",
			Data:           js.FuncOf(Select2GenericAjaxData),
			ProcessResults: js.FuncOf(Select2GenericAjaxProcessResults(EmpiricalFormulas{})),
		},
	})

	Jq("select#s_signalword").Select2(Select2Config{
		Placeholder:    locales.Translate("s_signalword", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2GenericTemplateResults(SignalWord{})),
		AllowClear:     true,
		Ajax: Select2Ajax{
			URL:            ApplicationProxyPath + "products/signalwords/",
			DataType:       "json",
			Data:           js.FuncOf(Select2GenericAjaxData),
			ProcessResults: js.FuncOf(Select2GenericAjaxProcessResults(SignalWords{})),
		},
	})

	Jq("select#s_symbols").Select2(Select2Config{
		Placeholder:    locales.Translate("s_symbols", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2SymbolTemplateResults),
		AllowClear:     true,
		Ajax: Select2Ajax{
			URL:            ApplicationProxyPath + "products/symbols/",
			DataType:       "json",
			Data:           js.FuncOf(Select2GenericAjaxData),
			ProcessResults: js.FuncOf(Select2GenericAjaxProcessResults(Symbols{})),
		},
	})

	Jq("select#s_hazardstatements").Select2(Select2Config{
		Placeholder:    locales.Translate("s_hazardstatements", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2GenericTemplateResults(HazardStatement{})),
		AllowClear:     true,
		Ajax: Select2Ajax{
			URL:            ApplicationProxyPath + "products/hazardstatements/",
			DataType:       "json",
			Data:           js.FuncOf(Select2GenericAjaxData),
			ProcessResults: js.FuncOf(Select2GenericAjaxProcessResults(HazardStatements{})),
		},
	})

	Jq("select#s_precautionarystatements").Select2(Select2Config{
		Placeholder:    locales.Translate("s_precautionarystatements", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2GenericTemplateResults(PrecautionaryStatement{})),
		AllowClear:     true,
		Ajax: Select2Ajax{
			URL:            ApplicationProxyPath + "products/precautionarystatements/",
			DataType:       "json",
			Data:           js.FuncOf(Select2GenericAjaxData),
			ProcessResults: js.FuncOf(Select2GenericAjaxProcessResults(PrecautionaryStatements{})),
		},
	})

	//
	// Type chooser.
	//
	Jq("input[type=radio][name=typechooser]").On("change", js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		switch Jq("input[name=typechooser]:checked").GetVal().String() {
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
	Jq("#product").Validate(ValidateConfig{
		ErrorClass: "alert alert-danger",
		Rules: map[string]ValidateRule{
			"name": {
				Required: js.FuncOf(func(this js.Value, args []js.Value) interface{} { return true }),
			},
			"producerref": {
				Required: js.FuncOf(func(this js.Value, args []js.Value) interface{} { return true }),
			},
			"unit_temperature": {
				Required: js.FuncOf(func(this js.Value, args []js.Value) interface{} {
					return Jq("#product_temperature").GetVal().Truthy()
				}),
			},
			"empiricalformula": {
				Required: js.FuncOf(func(this js.Value, args []js.Value) interface{} { return true }),
				Remote: ValidateRemote{
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
				Remote: ValidateRemote{
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
				Remote: ValidateRemote{
					URL:        "",
					Type:       "post",
					BeforeSend: js.FuncOf(ValidateProductCeNumberBeforeSend),
					Data: map[string]interface{}{
						"cenumber": js.FuncOf(ValidateProductCeNumberData),
					},
				},
			},
		},
		Messages: map[string]ValidateMessage{
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
	})

	// select2
	Jq("select#producerref").Select2(Select2Config{
		Placeholder:       locales.Translate("product_producerref_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult:    js.FuncOf(Select2GenericTemplateResults(ProducerRef{})),
		TemplateSelection: js.FuncOf(Select2ProducerRefTemplateSelection),
		Tags:              true,
		AllowClear:        true,
		CreateTag:         js.FuncOf(Select2ProducerRefCreateTag),
		Ajax: Select2Ajax{
			URL:            ApplicationProxyPath + "products/producerrefs/",
			DataType:       "json",
			Data:           js.FuncOf(Select2ProducerRefAjaxData),
			ProcessResults: js.FuncOf(Select2GenericAjaxProcessResults(ProducerRefs{})),
		},
	})
	Jq("select#producerref").On("change", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return Jq(this).Valid()
	}))
	Jq("select#producerref").On("select2:select", js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		select2ProducerrefSelected := args[0].Get("params").Get("data")
		producerrefSelected := ProducerRef{}.FromJsJSONValue(select2ProducerrefSelected)

		producerref := producerrefSelected.(ProducerRef)

		// If we create a new producerref
		if producerref.Producer == nil {
			return nil
		}

		Jq("select#producer").Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            producerref.Producer.ProducerLabel.String,
				Value:           strconv.Itoa(int(producerref.Producer.ProducerID.Int64)),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())

		return nil

	}))
	Jq("select#producer").Select2(Select2Config{
		Placeholder:    locales.Translate("product_producer_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2GenericTemplateResults(Producer{})),
		AllowClear:     true,
		Ajax: Select2Ajax{
			URL:            ApplicationProxyPath + "products/producers/",
			DataType:       "json",
			Data:           js.FuncOf(Select2GenericAjaxData),
			ProcessResults: js.FuncOf(Select2GenericAjaxProcessResults(Producers{})),
		},
	})

	Jq("select#supplierrefs").Select2(Select2Config{
		Placeholder:       locales.Translate("product_supplierref_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult:    js.FuncOf(Select2GenericTemplateResults(SupplierRef{})),
		TemplateSelection: js.FuncOf(Select2SupplierRefTemplateSelection),
		Tags:              true,
		AllowClear:        true,
		CreateTag:         js.FuncOf(Select2SupplierRefCreateTag),
		Ajax: Select2Ajax{
			URL:            ApplicationProxyPath + "products/supplierrefs/",
			DataType:       "json",
			Data:           js.FuncOf(Select2SupplierRefAjaxData),
			ProcessResults: js.FuncOf(Select2GenericAjaxProcessResults(SupplierRefs{})),
		},
	})
	Jq("select#supplier").Select2(Select2Config{
		Placeholder:    locales.Translate("product_supplier_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2GenericTemplateResults(Supplier{})),
		AllowClear:     true,
		Ajax: Select2Ajax{
			URL:            ApplicationProxyPath + "products/suppliers/",
			DataType:       "json",
			Data:           js.FuncOf(Select2GenericAjaxData),
			ProcessResults: js.FuncOf(Select2GenericAjaxProcessResults(Suppliers{})),
		},
	})

	Jq("select#tags").Select2(Select2Config{
		Placeholder:    locales.Translate("product_tag_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2GenericTemplateResults(Tag{})),
		CreateTag:      js.FuncOf(Select2GenericCreateTag(Tag{})),
		AllowClear:     true,
		Tags:           true,
		Ajax: Select2Ajax{
			URL:            ApplicationProxyPath + "products/tags/",
			DataType:       "json",
			Data:           js.FuncOf(Select2GenericAjaxData),
			ProcessResults: js.FuncOf(Select2GenericAjaxProcessResults(Tags{})),
		},
	})

	Jq("select#category").Select2(Select2Config{
		Placeholder:    locales.Translate("product_category_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2GenericTemplateResults(Category{})),
		AllowClear:     true,
		Tags:           true,
		CreateTag:      js.FuncOf(Select2GenericCreateTag(Category{})),
		Ajax: Select2Ajax{
			URL:            ApplicationProxyPath + "products/categories/",
			DataType:       "json",
			Data:           js.FuncOf(Select2GenericAjaxData),
			ProcessResults: js.FuncOf(Select2GenericAjaxProcessResults(Categories{})),
		},
	})

	Jq("select#unit_temperature").Select2(Select2Config{
		Placeholder:    locales.Translate("product_unit_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2GenericTemplateResults(Unit{})),
		AllowClear:     true,
		Ajax: Select2Ajax{
			URL:            ApplicationProxyPath + "storages/units",
			DataType:       "json",
			Data:           js.FuncOf(Select2UnitTemperatureAjaxData),
			ProcessResults: js.FuncOf(Select2GenericAjaxProcessResults(Units{})),
		},
	})
	Jq("select#unit_temperature").On("change", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return Jq(this).Valid()
	}))

	Jq("select#casnumber").Select2(Select2Config{
		Placeholder:    locales.Translate("product_cas_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2GenericTemplateResults(CasNumber{})),
		AllowClear:     true,
		Tags:           true,
		CreateTag:      js.FuncOf(Select2GenericCreateTag(CasNumber{})),
		Ajax: Select2Ajax{
			URL:            ApplicationProxyPath + "products/casnumbers/",
			DataType:       "json",
			Data:           js.FuncOf(Select2GenericAjaxData),
			ProcessResults: js.FuncOf(Select2GenericAjaxProcessResults(CasNumbers{})),
		},
	})
	Jq("select#casnumber").On("change", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return Jq(this).Valid()
	}))

	Jq("select#cenumber").Select2(Select2Config{
		Placeholder:    locales.Translate("product_ce_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2GenericTemplateResults(CeNumber{})),
		AllowClear:     true,
		Tags:           true,
		CreateTag:      js.FuncOf(Select2GenericCreateTag(CeNumber{})),
		Ajax: Select2Ajax{
			URL:            ApplicationProxyPath + "products/cenumbers/",
			DataType:       "json",
			Data:           js.FuncOf(Select2GenericAjaxData),
			ProcessResults: js.FuncOf(Select2GenericAjaxProcessResults(CeNumbers{})),
		},
	})
	Jq("select#cenumber").On("change", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return Jq(this).Valid()
	}))

	Jq("select#physicalstate").Select2(Select2Config{
		Placeholder:    locales.Translate("product_physicalstate_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2GenericTemplateResults(PhysicalState{})),
		AllowClear:     true,
		Tags:           true,
		CreateTag:      js.FuncOf(Select2GenericCreateTag(PhysicalState{})),
		Ajax: Select2Ajax{
			URL:            ApplicationProxyPath + "products/physicalstates/",
			DataType:       "json",
			Data:           js.FuncOf(Select2GenericAjaxData),
			ProcessResults: js.FuncOf(Select2GenericAjaxProcessResults(PhysicalStates{})),
		},
	})

	Jq("select#signalword").Select2(Select2Config{
		Placeholder:    locales.Translate("product_signalword_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2GenericTemplateResults(SignalWord{})),
		AllowClear:     true,
		Ajax: Select2Ajax{
			URL:            ApplicationProxyPath + "products/signalwords/",
			DataType:       "json",
			Data:           js.FuncOf(Select2GenericAjaxData),
			ProcessResults: js.FuncOf(Select2GenericAjaxProcessResults(SignalWords{})),
		},
	})

	Jq("select#classofcompound").Select2(Select2Config{
		Placeholder:    locales.Translate("product_classofcompound_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2GenericTemplateResults(ClassOfCompound{})),
		AllowClear:     true,
		Tags:           true,
		CreateTag:      js.FuncOf(Select2GenericCreateTag(ClassOfCompound{})),
		Ajax: Select2Ajax{
			URL:            ApplicationProxyPath + "products/classofcompounds/",
			DataType:       "json",
			Data:           js.FuncOf(Select2GenericAjaxData),
			ProcessResults: js.FuncOf(Select2GenericAjaxProcessResults(ClassesOfCompound{})),
		},
	})

	Jq("select#name").Select2(Select2Config{
		Placeholder:    locales.Translate("product_name_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2GenericTemplateResults(Name{})),
		AllowClear:     true,
		Tags:           true,
		CreateTag:      js.FuncOf(Select2GenericCreateTag(Name{})),
		Ajax: Select2Ajax{
			URL:            ApplicationProxyPath + "products/names/",
			DataType:       "json",
			Data:           js.FuncOf(Select2GenericAjaxData),
			ProcessResults: js.FuncOf(Select2GenericAjaxProcessResults(Names{})),
		},
	})
	Jq("select#name").On("change", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return Jq(this).Valid()
	}))

	Jq("select#empiricalformula").Select2(Select2Config{
		Placeholder:    locales.Translate("product_empiricalformula_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2GenericTemplateResults(EmpiricalFormula{})),
		AllowClear:     true,
		Tags:           true,
		CreateTag:      js.FuncOf(Select2GenericCreateTag(EmpiricalFormula{})),
		Ajax: Select2Ajax{
			URL:            ApplicationProxyPath + "products/empiricalformulas/",
			DataType:       "json",
			Data:           js.FuncOf(Select2GenericAjaxData),
			ProcessResults: js.FuncOf(Select2GenericAjaxProcessResults(EmpiricalFormulas{})),
		},
	})
	Jq("select#empiricalformula").On("change", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return Jq(this).Valid()
	}))

	Jq("select#linearformula").Select2(Select2Config{
		Placeholder:    locales.Translate("product_linearformula_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2GenericTemplateResults(LinearFormula{})),
		AllowClear:     true,
		Tags:           true,
		CreateTag:      js.FuncOf(Select2GenericCreateTag(LinearFormula{})),
		Ajax: Select2Ajax{
			URL:            ApplicationProxyPath + "products/linearformulas/",
			DataType:       "json",
			Data:           js.FuncOf(Select2GenericAjaxData),
			ProcessResults: js.FuncOf(Select2GenericAjaxProcessResults(LinearFormulas{})),
		},
	})

	Jq("select#synonyms").Select2(Select2Config{
		Placeholder:    locales.Translate("product_synonyms_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2GenericTemplateResults(Name{})),
		AllowClear:     true,
		Tags:           true,
		CreateTag:      js.FuncOf(Select2GenericCreateTag(Name{})),
		Ajax: Select2Ajax{
			URL:            ApplicationProxyPath + "products/synonyms/",
			DataType:       "json",
			Data:           js.FuncOf(Select2GenericAjaxData),
			ProcessResults: js.FuncOf(Select2GenericAjaxProcessResults(Names{})),
		},
	})

	Jq("select#symbols").Select2(Select2Config{
		Placeholder:    locales.Translate("product_symbols_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2SymbolTemplateResults),
		AllowClear:     true,
		Ajax: Select2Ajax{
			URL:            ApplicationProxyPath + "products/symbols/",
			DataType:       "json",
			Data:           js.FuncOf(Select2GenericAjaxData),
			ProcessResults: js.FuncOf(Select2GenericAjaxProcessResults(Symbols{})),
		},
	})

	Jq("select#hazardstatements").Select2(Select2Config{
		Placeholder:    locales.Translate("product_hazardstatements_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2HazardStatementTemplateResults),
		AllowClear:     true,
		Ajax: Select2Ajax{
			URL:            ApplicationProxyPath + "products/hazardstatements/",
			DataType:       "json",
			Data:           js.FuncOf(Select2GenericAjaxData),
			ProcessResults: js.FuncOf(Select2GenericAjaxProcessResults(HazardStatements{})),
		},
	})

	Jq("select#precautionarystatements").Select2(Select2Config{
		Placeholder:    locales.Translate("product_precautionarystatements_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2PrecautionaryStatementTemplateResults),
		AllowClear:     true,
		Ajax: Select2Ajax{
			URL:            ApplicationProxyPath + "products/precautionarystatements/",
			DataType:       "json",
			Data:           js.FuncOf(Select2GenericAjaxData),
			ProcessResults: js.FuncOf(Select2GenericAjaxProcessResults(PrecautionaryStatements{})),
		},
	})

	Jq("#product_twodformula").On("change", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		js.Global().Call("load2dimage")
		return nil
	}))

}

func ShowIfAuthorizedMenuItems(args ...interface{}) {

	utils.HasPermission("products", "-2", "get", func() {
		Jq("#menu_scan_qrcode").FadeIn()
		Jq("#menu_list_products").FadeIn()
		Jq("#menu_list_bookmarks").FadeIn()
	}, func() {
	})

	utils.HasPermission("products", "", "post", func() {
		Jq("#menu_create_product").FadeIn()
	}, func() {
	})

	utils.HasPermission("entities", "-2", "get", func() {
		Jq("#menu_entities").FadeIn()
	}, func() {
	})

	utils.HasPermission("entities", "", "post", func() {
		Jq("#menu_create_entity").FadeIn()
	}, func() {
	})

	utils.HasPermission("entities", "-2", "put", func() {
		Jq("#menu_update_welcomeannounce").FadeIn()
	}, func() {
	})

	utils.HasPermission("storages", "-2", "get", func() {
		Jq("#menu_storelocations").FadeIn()
	}, func() {
	})

	utils.HasPermission("storelocations", "", "post", func() {
		Jq("#menu_create_storelocation").FadeIn()
	}, func() {
	})

	utils.HasPermission("people", "-2", "get", func() {
		Jq("#menu_people").FadeIn()
	}, func() {
	})

	utils.HasPermission("people", "", "post", func() {
		Jq("#menu_create_person").FadeIn()
	}, func() {
	})

}

func LoadMenu() {

	href := fmt.Sprintf("%smenu", ApplicationProxyPath)
	utils.LoadMenu("menu", href, ShowIfAuthorizedMenuItems, nil)

	MenuLoaded = true

}

func LoadUser() {

	// Can not read Email and ID from container as those values
	// are not set before login.
	cookie := js.Global().Get("document").Get("cookie").String()
	regex := regexp.MustCompile(`(?P<token>token=\S*)\s{0,1}(?P<email>email=\S*)\s{0,1}(?P<id>id=\S*)\s{0,1}`)
	match := regex.FindStringSubmatch(cookie)

	if len(match) > 0 {

		result := make(map[string]string)
		for i, name := range regex.SubexpNames() {
			if i != 0 && name != "" {
				result[name] = match[i]
			}
		}
		ConnectedUserEmail = strings.TrimRight(result["email"], ";")[6:]
		if ConnectedUserID, err = strconv.Atoi(strings.TrimRight(result["id"], ";")[3:]); err != nil {
			panic(err)
		}

	}

	Jq("#logged").SetHtml(ConnectedUserEmail)

}

func LoadSearch() {

	bindSearchButtons := func(args ...interface{}) {

		// Works only with no select2.
		Jq("input").On("keyup", js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			event := args[0]
			if event.Get("which").Int() == 13 {

				event.Call("preventDefault")
				search.Search(js.Null(), nil)

			}

			return nil

		}))

		// Show/Hide archives.
		Jq("#s_storage_archive_button").On("click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			var (
				btnIcon  themes.IconFace
				btnLabel string
			)

			BSTableQueryFilter.Lock()

			if Jq("#s_storage_archive_button > span").HasClass(themes.MDI_SHOW_DELETED.ToString()) {
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

			Jq("#s_storage_archive_button").SetProp("title", btnLabel)
			Jq("#s_storage_archive_button").SetHtml("")
			Jq("#s_storage_archive_button").Append(buttonTitle.OuterHTML())

			Jq("#Storage_table").Bootstraptable(nil).Refresh(nil)
			Jq("#Storage_table").On("load-success.bs.table", js.FuncOf(ShowIfAuthorizedActionButtons))

			return nil

		}))

	}

	href := fmt.Sprintf("%ssearch", ApplicationProxyPath)
	utils.LoadSearch("search", href, bindSearchButtons, nil)

	SearchLoaded = true

}

func Product_listBookmarkCallback(this js.Value, args []js.Value) interface{} {

	BSTableQueryFilter.Clean()
	BSTableQueryFilter.ProductBookmark = true

	productCallbackWrapper := func(args ...interface{}) {
		Product_listCallback(js.Null(), nil)
	}

	utils.LoadContent("product", fmt.Sprintf("%sv/products", ApplicationProxyPath), productCallbackWrapper, nil)

	return nil

}

func Product_listCallback(this js.Value, args []js.Value) interface{} {

	product_common()

	Jq("#Product_table").Bootstraptable(&BootstraptableParams{Ajax: "Product_getTableData"})
	Jq("#Product_table").On("load-success.bs.table", js.FuncOf(ShowIfAuthorizedActionButtons))

	// TODO: move this
	if !MenuLoaded {
		LoadMenu()
	}
	if !SearchLoaded {
		LoadSearch()
	}
	if !UserLoaded {
		LoadUser()
	}

	Jq("#search").Show()
	Jq("#actions").Show()
	Jq("#s_storage_archive_button").Hide()

	btnLabel := locales.Translate("switchstorageview_text", HTTPHeaderAcceptLanguage)
	buttonTitle := widgets.NewIcon(widgets.IconAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"iconlabel"},
		},
		Icon: themes.NewMdiIcon(themes.MDI_STORAGE, ""),
		Text: btnLabel,
	})
	Jq("#switchview").SetProp("title", btnLabel)
	Jq("#switchview").SetHtml("")
	Jq("#switchview").Append(buttonTitle.OuterHTML())

	return nil

}

func Product_createCallback(this js.Value, args []js.Value) interface{} {

	product_common()

	// Chemical product by default on creation.
	if Jq("input#product_id").Object.Length() == 0 {
		Chemify()
	}

	Jq("#search").Hide()
	Jq("#actions").Hide()

	return nil

}

func Product_SaveCallback(args ...interface{}) {

	BSTableQueryFilter.Lock()
	BSTableQueryFilter.QueryFilter.Product = strconv.Itoa(args[0].(int))
	BSTableQueryFilter.QueryFilter.ProductFilterLabel = fmt.Sprintf("%s %s", CurrentProduct.Name.NameLabel, CurrentProduct.ProductSpecificity.String)
	Jq("#Product_table").Bootstraptable(nil).Refresh(nil)

	// Jq("#Product_table").Bootstraptable(nil).Refresh(&BootstraptableRefreshQuery{
	// 	Query: QueryFilter{
	// 		Product: strconv.Itoa(args[0].(int)),
	// 	},
	// })
	Jq("#Product_table").On("load-success.bs.table", js.FuncOf(ShowIfAuthorizedActionButtons))

	product_common()

	Jq("#search").Show()
	Jq("#actions").Show()

}
