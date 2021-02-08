package storage

import (
	"fmt"
	"reflect"
	"strconv"
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
			"storelocation": {
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
			"storelocation": {
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

	select2.NewSelect2(jquery.Jq("select#storelocation"), &select2.Select2Config{
		Placeholder:    locales.Translate("storage_storelocation_placeholder", HTTPHeaderAcceptLanguage),
		TemplateResult: js.FuncOf(Select2StoreLocationTemplateResults),
		AllowClear:     true,
		Ajax: select2.Select2Ajax{
			URL:            ApplicationProxyPath + "storelocations",
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

		jquery.Jq("input#storage_nbitem").SetProp("disabled", "disabled")
		jquery.Jq("input#storage_identicalbarecode").SetProp("disabled", "disabled")

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

	jquery.Jq("#search").Hide()
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

	jquery.Jq("#search").Show()
	jquery.Jq("#actions").Show()
	jquery.Jq("#s_storage_archive_button").SetVisible()

	if BSTableQueryFilter.QueryFilter.Product != "" {
		jquery.Jq("#s_storage_stock_button").SetVisible()
	}

	changeSwitchButtonToProduct()

	return nil

}

func Storage_SaveCallback(args ...interface{}) {

	BSTableQueryFilter.Lock()
	BSTableQueryFilter.QueryFilter.Storage = strconv.Itoa(args[0].(int))
	BSTableQueryFilter.QueryFilter.StorageFilterLabel = fmt.Sprintf("#%d", CurrentStorage.StorageID.Int64)
	BSTableQueryFilter.QueryFilter.ProductFilterLabel = fmt.Sprintf("%s %s", globals.CurrentProduct.Name.NameLabel, globals.CurrentProduct.ProductSpecificity.String)
	bstable.NewBootstraptable(jquery.Jq("#Storage_table"), nil).Refresh(nil)

	// jquery.bstable.NewBootstraptable(jquery.Jq("#Storage_table"),nil).Refresh(&BootstraptableRefreshQuery{
	// 	Query: QueryFilter{
	// 		Storage: strconv.Itoa(args[0].(int)),
	// 	},
	// })
	jquery.Jq("#Storage_table").On("load-success.bs.table", js.FuncOf(ShowIfAuthorizedActionButtons))

	//storage_common()

	jquery.Jq("#search").Show()
	jquery.Jq("#actions").Show()

	changeSwitchButtonToProduct()

}
