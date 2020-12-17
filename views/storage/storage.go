package storage

import (
	"fmt"
	"reflect"
	"strconv"
	"syscall/js"

	. "github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/locales"
	. "github.com/tbellembois/gochimitheque-wasm/types"
	"github.com/tbellembois/gochimitheque-wasm/widgets"
	"github.com/tbellembois/gochimitheque-wasm/widgets/themes"
)

func storage_common() {

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
	// create form
	//
	// validate
	Jq("#product").Validate(ValidateConfig{
		ErrorClass: "alert alert-danger",
		Rules: map[string]ValidateRule{
			"storelocation": {
				Required: js.FuncOf(func(this js.Value, args []js.Value) interface{} { return true }),
			},
			"unit_concentration": {
				Required: js.FuncOf(func(this js.Value, args []js.Value) interface{} {
					return Jq("#storage_concentration").GetVal().Truthy()
				}),
			},
		},
		Messages: map[string]ValidateMessage{
			"storelocation": {
				Required: locales.Translate("required_input", HTTPHeaderAcceptLanguage),
			},
		},
	})
	Jq("#borrowing").Validate(ValidateConfig{
		ErrorClass: "alert alert-danger",
		Rules: map[string]ValidateRule{
			"borrower": {
				Required: js.FuncOf(func(this js.Value, args []js.Value) interface{} { return true }),
			},
		},
		Messages: map[string]ValidateMessage{
			"borrower": {
				Required: locales.Translate("required_input", HTTPHeaderAcceptLanguage),
			},
		},
	})

	// select2
	Jq("select#unit_concentration").Select2(Select2Config{
		Placeholder:    locales.Translate("product_unit_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2GenericTemplateResults(Unit{})),
		AllowClear:     true,
		Ajax: Select2Ajax{
			URL:            ApplicationProxyPath + "storages/units",
			DataType:       "json",
			Data:           js.FuncOf(Select2UnitConcentrationAjaxData),
			ProcessResults: js.FuncOf(Select2GenericAjaxProcessResults(Units{})),
		},
	})

	Jq("select#storelocation").Select2(Select2Config{
		Placeholder:    locales.Translate("storage_storelocation_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2StoreLocationTemplateResults),
		AllowClear:     true,
		Ajax: Select2Ajax{
			URL:            ApplicationProxyPath + "storelocations",
			DataType:       "json",
			Data:           js.FuncOf(Select2StoreLocationAjaxData),
			ProcessResults: js.FuncOf(Select2GenericAjaxProcessResults(StoreLocations{})),
		},
	})

	Jq("select#unit_quantity").Select2(Select2Config{
		Placeholder:    locales.Translate("product_unit_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2GenericTemplateResults(Unit{})),
		AllowClear:     true,
		Ajax: Select2Ajax{
			URL:            ApplicationProxyPath + "storages/units",
			DataType:       "json",
			Data:           js.FuncOf(Select2UnitQuantityAjaxData),
			ProcessResults: js.FuncOf(Select2GenericAjaxProcessResults(Units{})),
		},
	})

	Jq("select#supplier").Select2(Select2Config{
		Placeholder:    locales.Translate("storage_supplier_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2GenericTemplateResults(Supplier{})),
		AllowClear:     true,
		Tags:           true,
		CreateTag:      js.FuncOf(Select2GenericCreateTag(Supplier{})),
		Ajax: Select2Ajax{
			URL:            ApplicationProxyPath + "products/suppliers/",
			DataType:       "json",
			Data:           js.FuncOf(Select2GenericAjaxData),
			ProcessResults: js.FuncOf(Select2GenericAjaxProcessResults(Suppliers{})),
		},
	})
	Jq("select#borrower").Select2(Select2Config{
		Placeholder:    locales.Translate("storage_borrower_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2GenericTemplateResults(Person{})),
		AllowClear:     true,
		Ajax: Select2Ajax{
			URL:            ApplicationProxyPath + "people",
			DataType:       "json",
			Data:           js.FuncOf(Select2GenericAjaxData),
			ProcessResults: js.FuncOf(Select2GenericAjaxProcessResults(People{})),
		},
	})

}

func Storage_createCallback(args ...interface{}) {

	var (
		productId   int
		productName string
	)

	storage_common()

	switch reflect.TypeOf(args[0]) {
	case reflect.TypeOf(Product{}):

		product := args[0].(Product)
		productId = product.ProductID
		productName = fmt.Sprintf("%s %s", product.Name.NameLabel, product.ProductSpecificity.String)

	case reflect.TypeOf(Storage{}):

		storage := args[0].(Storage)
		productId = storage.Product.ProductID
		productName = fmt.Sprintf("%s %s", storage.Product.Name.NameLabel, storage.Product.ProductSpecificity.String)

		FillInStorageForm(storage, "storage")

	}

	title := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"alert", "alert-light"},
			Attributes: map[string]string{
				"role": "alert",
			},
		},
	})
	title.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
		},
		Text: productName,
	}))

	Jq("input#product_id").SetVal(productId)
	Jq("#filter-product").SetHtml(title.OuterHTML())

	Jq("#search").Hide()
	Jq("#actions").Hide()

}

func Storage_listCallback(this js.Value, args []js.Value) interface{} {

	storage_common()

	Jq("#Storage_table").Bootstraptable(&BootstraptableParams{Ajax: "Storage_getTableData"})
	Jq("#Storage_table").On("load-success.bs.table", js.FuncOf(ShowIfAuthorizedActionButtons))

	Jq("#search").Show()
	Jq("#actions").Show()
	Jq("#s_storage_archive_button").Show()

	btnLabel := locales.Translate("switchproductview_text", HTTPHeaderAcceptLanguage)
	buttonTitle := widgets.NewIcon(widgets.IconAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"iconlabel"},
		},
		Icon: themes.NewMdiIcon(themes.MDI_PRODUCT, ""),
		Text: btnLabel,
	})
	Jq("#switchview").SetProp("title", btnLabel)
	Jq("#switchview").SetHtml("")
	Jq("#switchview").Append(buttonTitle.OuterHTML())

	return nil

}

func Storage_SaveCallback(args ...interface{}) {

	BSTableQueryFilter.Lock()
	BSTableQueryFilter.QueryFilter.Storage = strconv.Itoa(args[0].(int))
	BSTableQueryFilter.QueryFilter.StorageFilterLabel = fmt.Sprintf("id: %d", args[0].(int))
	Jq("#Storage_table").Bootstraptable(nil).Refresh(nil)

	// Jq("#Storage_table").Bootstraptable(nil).Refresh(&BootstraptableRefreshQuery{
	// 	Query: QueryFilter{
	// 		Storage: strconv.Itoa(args[0].(int)),
	// 	},
	// })
	Jq("#Storage_table").On("load-success.bs.table", js.FuncOf(ShowIfAuthorizedActionButtons))

	storage_common()

}
