//go:build go1.24 && js && wasm

package storage

import (
	"fmt"
	"reflect"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/bstable"
	"github.com/tbellembois/gochimitheque-wasm/globals"
	. "github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/jquery"
	"github.com/tbellembois/gochimitheque-wasm/locales"
	"github.com/tbellembois/gochimitheque-wasm/select2"
	. "github.com/tbellembois/gochimitheque-wasm/types"
	"github.com/tbellembois/gochimitheque-wasm/validate"
	"github.com/tbellembois/gochimitheque-wasm/widgets"
	"github.com/tbellembois/gochimitheque-wasm/widgets/themes"
)

func storage_common() {

	//
	// create form
	//
	// validate
	validate.NewValidate(jquery.Jq("#storage"), &validate.ValidateConfig{
		ErrorClass: "alert alert-danger",
		Rules: map[string]validate.ValidateRule{
			// "storage_number_of_unit": {
			// 	Required: js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			// 		return jquery.Jq("#storage_number_of_bag").GetVal().String() == "" && jquery.Jq("#storage_number_of_carton").GetVal().String() == ""
			// 	}),
			// 	Remote: validate.ValidateRemote{
			// 		BeforeSend: js.FuncOf(func(this js.Value, args []js.Value) interface{} { return false }),
			// 	},
			// },
			// "storage_number_of_bag": {
			// 	Required: js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			// 		return jquery.Jq("#storage_number_of_unit").GetVal().String() == "" && jquery.Jq("#storage_number_of_carton").GetVal().String() == ""
			// 	}),
			// 	Remote: validate.ValidateRemote{
			// 		BeforeSend: js.FuncOf(func(this js.Value, args []js.Value) interface{} { return false }),
			// 	},
			// },
			// "storage_number_of_carton": {
			// 	Required: js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			// 		return jquery.Jq("#storage_number_of_bag").GetVal().String() == "" && jquery.Jq("#storage_number_of_unit").GetVal().String() == ""
			// 	}),
			// 	Remote: validate.ValidateRemote{
			// 		BeforeSend: js.FuncOf(func(this js.Value, args []js.Value) interface{} { return false }),
			// 	},
			// },
			"storage_batch_number": {
				Required: js.FuncOf(func(this js.Value, args []js.Value) interface{} {
					return true
				}),
				Remote: validate.ValidateRemote{
					BeforeSend: js.FuncOf(func(this js.Value, args []js.Value) interface{} { return false }),
				},
			},
			"store_location": {
				Required: js.FuncOf(func(this js.Value, args []js.Value) interface{} { return true }),
				Remote: validate.ValidateRemote{
					BeforeSend: js.FuncOf(func(this js.Value, args []js.Value) interface{} { return false }),
				},
			},
			"unit_concentration": {
				Required: js.FuncOf(func(this js.Value, args []js.Value) interface{} {
					return jquery.Jq("#storage_concentration").GetVal().Truthy()
				}),
				Remote: validate.ValidateRemote{
					BeforeSend: js.FuncOf(func(this js.Value, args []js.Value) interface{} { return false }),
				},
			},
		},
		Messages: map[string]validate.ValidateMessage{
			"store_location": {
				Required: locales.Translate("required_input", HTTPHeaderAcceptLanguage),
			},
			"unit_concentration": {
				Required: locales.Translate("required_input", HTTPHeaderAcceptLanguage),
			},
			"storage_number_of_unit": {
				Required: locales.Translate("storage_one_number_required", HTTPHeaderAcceptLanguage),
			},
			"storage_number_of_bag": {
				Required: locales.Translate("storage_one_number_required", HTTPHeaderAcceptLanguage),
			},
			"storage_number_of_carton": {
				Required: locales.Translate("storage_one_number_required", HTTPHeaderAcceptLanguage),
			},
			"storage_batch_number": {
				Required: locales.Translate("required_input", HTTPHeaderAcceptLanguage),
			},
		},
	}).Validate()

	// select2
	select2.NewSelect2(jquery.Jq("select#unit_concentration"), &select2.Select2Config{
		Placeholder:    locales.Translate("product_unit_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(Unit{})),
		AllowClear:     true,
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "storages/units",
			DataType:       "json",
			Data:           js.FuncOf(Select2UnitConcentrationAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Units{})),
		},
	}).Select2ify()

	select2.NewSelect2(jquery.Jq("select#store_location"), &select2.Select2Config{
		Placeholder:    locales.Translate("storage_store_location_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2StoreLocationTemplateResults),
		AllowClear:     true,
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "store_locations",
			DataType:       "json",
			Data:           js.FuncOf(Select2StoreLocationAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(StoreLocations{})),
		},
	}).Select2ify()

	select2.NewSelect2(jquery.Jq("select#unit_quantity"), &select2.Select2Config{
		Placeholder:    locales.Translate("product_unit_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(Unit{})),
		AllowClear:     true,
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "storages/units",
			DataType:       "json",
			Data:           js.FuncOf(Select2UnitQuantityAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Units{})),
		},
	}).Select2ify()

	select2.NewSelect2(jquery.Jq("select#supplier"), &select2.Select2Config{
		Placeholder:    locales.Translate("storage_supplier_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(Supplier{})),
		AllowClear:     true,
		Tags:           true,
		CreateTag:      js.FuncOf(select2.Select2GenericCreateTag(Supplier{})),
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "products/suppliers/",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(Suppliers{})),
		},
	}).Select2ify()

	storage_borrower()

}

func storage_borrower() {

	validate.NewValidate(jquery.Jq("#borrowing"), &validate.ValidateConfig{
		ErrorClass: "alert alert-danger",
		Rules: map[string]validate.ValidateRule{
			"borrower": {
				Required: js.FuncOf(func(this js.Value, args []js.Value) interface{} { return true }),
				Remote: validate.ValidateRemote{
					BeforeSend: js.FuncOf(func(this js.Value, args []js.Value) interface{} { return false }),
				},
			},
		},
		Messages: map[string]validate.ValidateMessage{
			"borrower": {
				Required: locales.Translate("required_input", HTTPHeaderAcceptLanguage),
			},
		},
	}).Validate()

	select2.NewSelect2(jquery.Jq("select#borrower"), &select2.Select2Config{
		Placeholder:    locales.Translate("storage_borrower_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(select2.Select2GenericTemplateResults(Person{})),
		AllowClear:     true,
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "people",
			DataType:       "json",
			Data:           js.FuncOf(select2.Select2GenericAjaxData),
			ProcessResults: js.FuncOf(select2.Select2GenericAjaxProcessResults(People{})),
		},
	}).Select2ify()

}

func Storage_createCallback(args ...interface{}) {

	var (
		productId   int64
		productName string
	)

	storage_common()

	switch reflect.TypeOf(args[0]) {
	case reflect.TypeOf(Product{}):

		product := args[0].(Product)
		productId = *product.ProductID

		if product.ProductSpecificity == nil {
			productName = product.Name.NameLabel
		} else {
			productName = fmt.Sprintf("%s %s", product.Name.NameLabel, *product.ProductSpecificity)
		}

		// Chem/Bio/Consu detection.
		switch product.ProductType {
		case "cons":
			Consufy()
		case "bio":
			Biofy()
		default:
			Chemify()
		}

		globals.CurrentProduct = product

	case reflect.TypeOf(Storage{}):

		storage := args[0].(Storage)
		productId = *storage.Product.ProductID
		if storage.Product.ProductSpecificity == nil {
			productName = storage.Product.Name.NameLabel
		} else {
			productName = fmt.Sprintf("%s %s", storage.Product.Name.NameLabel, *storage.Product.ProductSpecificity)
		}

		FillInStorageForm(storage, "storage")

		// Chem/Bio/Consu detection.
		switch storage.Product.ProductType {
		case "cons":
			Consufy()
		case "bio":
			Biofy()
		default:
			Chemify()
		}

		if !(len(args) > 1 && args[1] == "clone") {
			jquery.Jq("input#storage_nbitem").SetProp("disabled", "disabled")
			jquery.Jq("input#storage_identicalbarecode").SetProp("disabled", "disabled")
		}

		globals.CurrentProduct = Product{Product: &storage.Product}
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

	jquery.Jq("input#product_id").SetVal(productId)
	jquery.Jq("#filter-product").SetHtml(title.OuterHTML())

	if globals.CurrentProduct.ProductNumberPerBag == nil || !(*globals.CurrentProduct.ProductNumberPerBag > 0) {
		jquery.Jq("input#storage_number_of_bag").SetProp("disabled", true)
	}
	if globals.CurrentProduct.ProductNumberPerCarton == nil || !(*globals.CurrentProduct.ProductNumberPerCarton > 0) {
		jquery.Jq("input#storage_number_of_carton").SetProp("disabled", true)
	}
	jquery.Jq("#searchbar").Hide()
	jquery.Jq("#actions").Hide()

}

func changeSwitchButtonToProduct() {

	btnLabel := locales.Translate("switchproductview_text", HTTPHeaderAcceptLanguage)
	buttonTitle := widgets.NewIcon(widgets.IconAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"iconlabel"},
		},
		Icon: themes.NewMdiIcon(themes.MDI_PRODUCT, ""),
		Text: btnLabel,
	})
	jquery.Jq("#switchview").SetProp("title", btnLabel)
	jquery.Jq("#switchview").SetHtml("")
	jquery.Jq("#switchview").Append(buttonTitle.OuterHTML())

}

func Storage_listCallback(this js.Value, args []js.Value) interface{} {

	//storage_common()
	storage_borrower()

	bstable.NewBootstraptable(jquery.Jq("#Storage_table"), &bstable.BootstraptableParams{Ajax: "Storage_getTableData"})
	jquery.Jq("#Storage_table").On("load-success.bs.table", js.FuncOf(ShowIfAuthorizedActionButtons))

	jquery.Jq("#searchbar").Show()
	jquery.Jq("#actions").Show()
	jquery.Jq("#s_storage_archive_button").SetVisible()

	if BSTableQueryFilter.QueryFilter.Product != "" {
		jquery.Jq("#s_storage_stock_button").SetVisible()
	}

	changeSwitchButtonToProduct()

	return nil

}

func Storage_SaveCallback(args ...interface{}) {

	var (
		// currentStorageIds []int
		filterLabel string
	)
	// for _, storage := range CurrentStorages {
	// 	currentStorageIds = append(currentStorageIds, int(*storage.StorageID))
	// 	filterLabel += fmt.Sprintf("#%d ", storage.StorageID)
	// }

	BSTableQueryFilter.Lock()
	BSTableQueryFilter.QueryFilter.Storages = globals.CurrentStorageIds
	BSTableQueryFilter.QueryFilter.StoragesFilterLabel = filterLabel

	if globals.CurrentProduct.ProductSpecificity != nil {
		BSTableQueryFilter.QueryFilter.ProductFilterLabel = fmt.Sprintf("%s %s", globals.CurrentProduct.Name.NameLabel, *globals.CurrentProduct.ProductSpecificity)
	} else {
		BSTableQueryFilter.QueryFilter.ProductFilterLabel = globals.CurrentProduct.Name.NameLabel
	}

	bstable.NewBootstraptable(jquery.Jq("#Storage_table"), nil).Refresh(nil)

	// jquery.bstable.NewBootstraptable(jquery.Jq("#Storage_table"),nil).Refresh(&BootstraptableRefreshQuery{
	// 	Query: QueryFilter{
	// 		Storage: strconv.Itoa(args[0].(int)),
	// 	},
	// })
	jquery.Jq("#Storage_table").On("load-success.bs.table", js.FuncOf(ShowIfAuthorizedActionButtons))

	//storage_common()

	jquery.Jq("#searchbar").Show()
	jquery.Jq("#actions").Show()

	changeSwitchButtonToProduct()

}
