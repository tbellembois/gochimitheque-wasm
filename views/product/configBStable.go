package product

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
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
	"github.com/tbellembois/gochimitheque-wasm/views/storage"
	"github.com/tbellembois/gochimitheque-wasm/widgets"
	"github.com/tbellembois/gochimitheque-wasm/widgets/themes"
	"github.com/tbellembois/gochimitheque/models"
	"honnef.co/go/js/dom/v2"
)

func OperateEventsBookmark(this js.Value, args []js.Value) interface{} {

	row := args[2]
	globals.CurrentProduct = Product{Product: &models.Product{}}.ProductFromJsJSONValue(row)

	url := fmt.Sprintf("%sbookmarks/%d", ApplicationProxyPath, globals.CurrentProduct.ProductID)
	method := "put"

	done := func(data js.Value) {

		var (
			product Product
			err     error
		)

		if err = json.Unmarshal([]byte(data.String()), &product); err != nil {
			fmt.Println(err)
		}

		bstable.NewBootstraptable(jquery.Jq("#Product_table"), nil).Refresh(nil)

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

}

func OperateEventsStore(this js.Value, args []js.Value) interface{} {

	row := args[2]

	globals.CurrentProduct = Product{Product: &models.Product{}}.ProductFromJsJSONValue(row)

	BSTableQueryFilter.Lock()
	BSTableQueryFilter.QueryFilter.Product = strconv.Itoa(globals.CurrentProduct.ProductID)
	BSTableQueryFilter.QueryFilter.ProductFilterLabel = fmt.Sprintf("%s %s", globals.CurrentProduct.Name.NameLabel, globals.CurrentProduct.ProductSpecificity.String)
	BSTableQueryFilter.Unlock()

	href := fmt.Sprintf("%svc/storages", ApplicationProxyPath)
	jsutils.LoadContent("div#content", "storage", href, storage.Storage_createCallback, globals.CurrentProduct)

	return nil

}

func OperateEventsStorages(this js.Value, args []js.Value) interface{} {

	storageCallbackWrapper := func(args ...interface{}) {
		storage.Storage_listCallback(js.Null(), nil)
	}

	row := args[2]
	globals.CurrentProduct = Product{Product: &models.Product{}}.ProductFromJsJSONValue(row)

	BSTableQueryFilter.Lock()
	BSTableQueryFilter.QueryFilter.Product = strconv.Itoa(globals.CurrentProduct.ProductID)
	BSTableQueryFilter.QueryFilter.ProductFilterLabel = fmt.Sprintf("%s %s", globals.CurrentProduct.Name.NameLabel, globals.CurrentProduct.ProductSpecificity.String)

	href := fmt.Sprintf("%sv/storages", ApplicationProxyPath)
	jsutils.LoadContent("div#content", "storage", href, storageCallbackWrapper)

	return nil

}

func OperateEventsOStorages(this js.Value, args []js.Value) interface{} {

	row := args[2]
	globals.CurrentProduct = Product{Product: &models.Product{}}.ProductFromJsJSONValue(row)

	url := fmt.Sprintf("%sstorages/others?product=%d", ApplicationProxyPath, globals.CurrentProduct.ProductID)
	method := "get"

	done := func(data js.Value) {

		var (
			entities Entities
			err      error
		)

		if err = json.Unmarshal([]byte(data.String()), &entities); err != nil {
			jsutils.DisplayGenericErrorMessage()
			fmt.Println(err)
		}

		for _, entity := range entities.Rows {

			divEntity := widgets.NewDiv(widgets.DivAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"row"},
				},
			})

			managers := strings.Split(entity.EntityDescription, ",")

			divEntityName := widgets.NewDiv(widgets.DivAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"col-sm-auto"},
				},
			})
			spanEntityName := widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"iconlabel"},
				},
				Text: entity.EntityName,
			})
			divEntityName.AppendChild(spanEntityName)

			divEntity.AppendChild(divEntityName)

			for _, m := range managers {
				divEntityDescription := widgets.NewDiv(widgets.DivAttributes{
					BaseAttributes: widgets.BaseAttributes{
						Visible: true,
						Classes: []string{"col-sm-12"},
					},
				})
				spanEntityDescription := widgets.NewSpan(widgets.SpanAttributes{
					BaseAttributes: widgets.BaseAttributes{
						Visible: true,
						Classes: []string{"blockquote-footer"},
					},
					Text: m,
				})
				divEntityDescription.AppendChild(spanEntityDescription)

				divEntity.AppendChild(divEntityDescription)

			}

			jquery.Jq(fmt.Sprintf("#ostorages-collapse-%d", globals.CurrentProduct.ProductID)).SetHtml(divEntity.OuterHTML())

		}

		jquery.Jq(fmt.Sprintf("#ostorages-collapse-%d", globals.CurrentProduct.ProductID)).Show()

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

}

func OperateEventsEdit(this js.Value, args []js.Value) interface{} {

	row := args[2]
	globals.CurrentProduct = Product{Product: &models.Product{}}.ProductFromJsJSONValue(row)

	BSTableQueryFilter.Lock()
	BSTableQueryFilter.QueryFilter.Product = strconv.Itoa(globals.CurrentProduct.ProductID)
	BSTableQueryFilter.QueryFilter.ProductFilterLabel = fmt.Sprintf("%s %s", globals.CurrentProduct.Name.NameLabel, globals.CurrentProduct.ProductSpecificity.String)
	BSTableQueryFilter.Unlock()

	href := fmt.Sprintf("%svc/products", ApplicationProxyPath)
	jsutils.LoadContent("div#content", "storage", href, Product_createCallback, globals.CurrentProduct)

	return nil

}

func OperateEventsTotalStock(this js.Value, args []js.Value) interface{} {

	row := args[2]
	product := Product{Product: &models.Product{}}.ProductFromJsJSONValue(row)

	jquery.Jq(fmt.Sprintf("#totalstock-collapse-%d", product.ProductID)).Append(widgets.NewSpan(widgets.SpanAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"mdi", "mdi-loading", "mdi-spin", "mdi-36px"},
		},
	}).OuterHTML())

	jquery.Jq(fmt.Sprintf("#totalstock-collapse-%d", product.ProductID)).Show()

	url := fmt.Sprintf("%sentities/stocks/%d", ApplicationProxyPath, product.ProductID)
	method := "get"

	done := func(data js.Value) {

		var (
			storelocations []models.StoreLocation
			err            error
		)

		if err = json.Unmarshal([]byte(data.String()), &storelocations); err != nil {
			fmt.Println(err)
		}

		jquery.Jq(fmt.Sprintf("#totalstock-collapse-%d", product.ProductID)).SetHtml("")

		for _, storelocation := range storelocations {
			jsutils.ShowStockRecursive(&storelocation, 0, fmt.Sprintf("#totalstock-collapse-%d", product.ProductID))
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
}

func OperateEventsDelete(this js.Value, args []js.Value) interface{} {

	row := args[2]
	globals.CurrentProduct = Product{Product: &models.Product{}}.ProductFromJsJSONValue(row)

	confirm := widgets.NewLink(
		widgets.LinkAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Id:      fmt.Sprintf("delete%d", globals.CurrentProduct.ProductID),
				Classes: []string{"text-primary", "iconlabel"},
				Visible: true,
			},
			Href: "#",
			Label: widgets.NewSpan(
				widgets.SpanAttributes{
					BaseAttributes: widgets.BaseAttributes{
						Classes: []string{"mdi", themes.MDI_CONFIRM.ToString()},
						Visible: true,
					},
					Text: locales.Translate("confirm", HTTPHeaderAcceptLanguage),
				},
			),
		},
	).OuterHTML()

	jquery.Jq(fmt.Sprintf("div#confirm%d", globals.CurrentProduct.ProductID)).SetHtml(confirm)

	jquery.Jq(fmt.Sprintf("a#delete%d", globals.CurrentProduct.ProductID)).On("click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		url := fmt.Sprintf("%sproducts/%d", ApplicationProxyPath, globals.CurrentProduct.ProductID)
		method := "delete"

		done := func(data js.Value) {

			jsutils.DisplaySuccessMessage(locales.Translate("product_deleted_message", HTTPHeaderAcceptLanguage))
			bstable.NewBootstraptable(jquery.Jq("#Product_table"), nil).Refresh(nil)

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

	return nil

}

// Get remote bootstrap table data - params object defined here:
// https://bootstrap-table.com/docs/api/table-options/#queryparams
func GetTableData(this js.Value, args []js.Value) interface{} {

	row := args[0]
	params := bstable.QueryParamsFromJsJSONValue(row)

	go func() {

		u := url.URL{Path: ApplicationProxyPath + "products"}
		u.RawQuery = params.Data.ToRawQuery()

		ajax := ajax.Ajax{
			URL:    u.String(),
			Method: "get",
			Done: func(data js.Value) {

				var (
					products Products
					err      error
				)
				if err = json.Unmarshal([]byte(data.String()), &products); err != nil {
					fmt.Println(err)
				}

				if products.GetExportFn() != "" {
					jsutils.DisplaySuccessMessage(locales.Translate("export_done", HTTPHeaderAcceptLanguage))

					var icon widgets.Widget
					icon.HTMLElement = widgets.NewIcon(widgets.IconAttributes{
						BaseAttributes: widgets.BaseAttributes{
							Visible: true,
						},
						Text: locales.Translate("download_export", HTTPHeaderAcceptLanguage),
						Icon: themes.NewMdiIcon(themes.MDI_DOWNLOAD, themes.MDI_24PX),
					})

					downloadLink := widgets.NewLink(widgets.LinkAttributes{
						BaseAttributes: widgets.BaseAttributes{
							Visible: true,
							Classes: []string{"iconlabel"},
						},
						Onclick: "$('#export').collapse('hide')",
						Title:   locales.Translate("download_export", HTTPHeaderAcceptLanguage),
						Href:    fmt.Sprintf("%sdownload/%s", ApplicationProxyPath, products.GetExportFn()),
						Label:   icon,
					})

					jquery.Jq("#export-body").SetHtml(downloadLink.OuterHTML())
					jquery.Jq("#export").Show()
					jquery.Jq("button#export").SetProp("disabled", false)

				} else if products.GetTotal() != 0 {

					row.Call("success", js.ValueOf(js.Global().Get("JSON").Call("parse", data)))

				} else {

					// TODO: improve this
					jquery.Jq("span.loading-wrap").SetHtml(locales.Translate("no_result", globals.HTTPHeaderAcceptLanguage))

				}

			},
			Fail: func(jqXHR js.Value) {

				jsutils.DisplayGenericErrorMessage()

			},
		}

		ajax.Send()

	}()

	return nil

}

// TODO: factorise me with storage
func DataQueryParams(this js.Value, args []js.Value) interface{} {

	params := args[0]

	queryFilter := ajax.QueryFilterFromJsJSONValue(params)

	// Product_SaveCallback product id.
	queryFilter.Product = BSTableQueryFilter.Product
	queryFilter.ProductFilterLabel = BSTableQueryFilter.ProductFilterLabel
	queryFilter.ProductBookmark = BSTableQueryFilter.ProductBookmark
	queryFilter.Export = BSTableQueryFilter.Export
	BSTableQueryFilter.Export = false
	BSTableQueryFilter.Unlock()

	select2SProducerRef := select2.NewSelect2(jquery.Jq("select#s_producerref"), nil)
	if select2SProducerRef.Select2IsInitialized() {
		i := select2SProducerRef.Select2Data()
		if len(i) > 0 {
			queryFilter.ProducerRef = i[0].Id
			queryFilter.ProducerRefFilterLabel = i[0].Text
		}
	}

	select2SStoreLocation := select2.NewSelect2(jquery.Jq("select#s_storelocation"), nil)
	if select2SStoreLocation.Select2IsInitialized() {
		i := select2SStoreLocation.Select2Data()
		if len(i) > 0 {
			queryFilter.StoreLocation = i[0].Id
			queryFilter.StoreLocationFilterLabel = i[0].Text
		}
	}

	select2SName := select2.NewSelect2(jquery.Jq("select#s_name"), nil)
	if select2SName.Select2IsInitialized() {
		i := select2SName.Select2Data()
		if len(i) > 0 {
			queryFilter.Name = i[0].Id
			queryFilter.NameFilterLabel = i[0].Text
		}
	}

	select2SCasNumber := select2.NewSelect2(jquery.Jq("select#s_casnumber"), nil)
	if select2SCasNumber.Select2IsInitialized() {
		i := select2SCasNumber.Select2Data()
		if len(i) > 0 {
			queryFilter.CasNumber = i[0].Id
			queryFilter.CasNumberFilterLabel = i[0].Text
		}
	}

	select2SEmpiricalFormula := select2.NewSelect2(jquery.Jq("select#s_empiricalformula"), nil)
	if select2SEmpiricalFormula.Select2IsInitialized() {
		i := select2SEmpiricalFormula.Select2Data()
		if len(i) > 0 {
			queryFilter.EmpiricalFormula = i[0].Id
			queryFilter.EmpiricalFormulaFilterLabel = i[0].Text
		}
	}

	select2SCategory := select2.NewSelect2(jquery.Jq("select#s_category"), nil)
	if select2SCategory.Select2IsInitialized() {
		i := select2SCategory.Select2Data()
		if len(i) > 0 {
			queryFilter.Category = i[0].Id
			queryFilter.CategoryFilterLabel = i[0].Text
		}
	}

	select2SSignalWord := select2.NewSelect2(jquery.Jq("select#s_signalword"), nil)
	if select2SSignalWord.Select2IsInitialized() {
		i := select2SSignalWord.Select2Data()
		if len(i) > 0 {
			queryFilter.SignalWord = i[0].Id
			queryFilter.SignalWordFilterLabel = i[0].Text
		}
	}

	select2STags := select2.NewSelect2(jquery.Jq("select#s_tags"), nil)
	if select2STags.Select2IsInitialized() {
		i := select2STags.Select2Data()
		if len(i) > 0 {
			for _, tag := range i {
				queryFilter.Tags = append(queryFilter.Tags, tag.Id)
				queryFilter.TagsFilterLabel += fmt.Sprintf(" %s", tag.Text)
			}
		}
	}

	select2SHS := select2.NewSelect2(jquery.Jq("select#s_hazardstatements"), nil)
	if select2SHS.Select2IsInitialized() {
		i := select2SHS.Select2Data()
		if len(i) > 0 {
			for _, hs := range i {
				queryFilter.HazardStatements = append(queryFilter.HazardStatements, hs.Id)
				queryFilter.HazardStatementsFilterLabel += fmt.Sprintf(" %s", hs.Text)
			}
		}
	}

	select2SPS := select2.NewSelect2(jquery.Jq("select#s_precautionarystatements"), nil)
	if select2SPS.Select2IsInitialized() {
		i := select2SPS.Select2Data()
		if len(i) > 0 {
			for _, ps := range i {
				queryFilter.PrecautionaryStatements = append(queryFilter.PrecautionaryStatements, ps.Id)
				queryFilter.PrecautionaryStatementsFilterLabel += fmt.Sprintf(" %s", ps.Text)
			}
		}
	}

	select2SSymbols := select2.NewSelect2(jquery.Jq("select#s_symbols"), nil)
	if select2SSymbols.Select2IsInitialized() {
		i := select2SSymbols.Select2Data()
		if len(i) > 0 {
			for _, s := range i {
				queryFilter.Symbols = append(queryFilter.Symbols, s.Id)
				queryFilter.SignalWordFilterLabel += fmt.Sprintf(" %s", s.Text)
			}
		}
	}

	if jquery.Jq("#s_storage_batchnumber").GetVal().Truthy() {
		queryFilter.StorageBatchNumber = jquery.Jq("#s_storage_batchnumber").GetVal().String()
		queryFilter.StorageBatchNumberFilterLabel = jquery.Jq("#s_storage_batchnumber").GetVal().String()
	}
	if jquery.Jq("#s_storage_barecode").GetVal().Truthy() {
		queryFilter.StorageBarecode = jquery.Jq("#s_storage_barecode").GetVal().String()
		queryFilter.StorageBarecodeFilterLabel = jquery.Jq("#s_storage_barecode").GetVal().String()
	}
	if jquery.Jq("#s_custom_name_part_of").GetVal().Truthy() {
		queryFilter.CustomNamePartOf = jquery.Jq("#s_custom_name_part_of").GetVal().String()
		queryFilter.CustomNamePartOfFilterLabel = jquery.Jq("#s_custom_name_part_of").GetVal().String()
	}
	if jquery.Jq("#s_casnumber_cmr:checked").Object.Length() > 0 {
		queryFilter.CasNumberCMR = true
		queryFilter.CasNumberCMRFilterLabel = locales.Translate("s_casnumber_cmr", globals.HTTPHeaderAcceptLanguage)
	}
	if jquery.Jq("#s_borrowing:checked").Object.Length() > 0 {
		queryFilter.Borrowing = true
		queryFilter.BorrowingFilterLabel = locales.Translate("s_borrowing", globals.HTTPHeaderAcceptLanguage)
	}
	if jquery.Jq("#s_storage_to_destroy:checked").Object.Length() > 0 {
		queryFilter.StorageToDestroy = true
		queryFilter.StorageToDestroyFilterLabel = locales.Translate("s_storage_to_destroy", globals.HTTPHeaderAcceptLanguage)
	}

	if jquery.Jq("input#searchshowbio:checked").Object.Length() > 0 {
		queryFilter.ShowBio = true
	} else {
		queryFilter.ShowBio = false
	}
	if jquery.Jq("input#searchshowchem:checked").Object.Length() > 0 {
		queryFilter.ShowChem = true
	} else {
		queryFilter.ShowChem = false
	}
	if jquery.Jq("input#searchshowconsu:checked").Object.Length() > 0 {
		queryFilter.ShowConsu = true
	} else {
		queryFilter.ShowConsu = false
	}

	jsutils.DisplayFilter(queryFilter)

	return queryFilter.ToJsValue()

}

func DetailFormatter(this js.Value, args []js.Value) interface{} {

	row := args[1]
	globals.CurrentProduct = Product{Product: &models.Product{}}.ProductFromJsJSONValue(row)

	var synonyms strings.Builder
	for _, synonym := range globals.CurrentProduct.Synonyms {
		synonyms.WriteString(synonym.NameLabel)
		synonyms.WriteString(";")
	}

	//
	// 2dImage, synonyms, ID and person.
	//
	// 2dimage.
	col2dimage := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-auto m-1"},
		},
	})
	if globals.CurrentProduct.ProductTwoDFormula.Valid {
		col2dimage.AppendChild(widgets.NewImg(widgets.ImgAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Width: "200",
			Src:   globals.CurrentProduct.ProductTwoDFormula.String,
		}))
	}
	// Synonym.
	colSynonym := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-12 m-1"},
		},
	})
	colSynonym.AppendChild(
		widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
				Classes: []string{"iconlabel", "mr-sm-2"},
			},
			Text: locales.Translate("synonym_label_title", HTTPHeaderAcceptLanguage),
		}))
	colSynonym.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
		},
		Text: synonyms.String(),
	}))
	// ID.
	colID := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-auto m-1"},
		},
	})
	colID.AppendChild(
		widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
				Classes: []string{"iconlabel", "mr-sm-2"},
			},
			Text: "#",
		}))
	colID.AppendChild(
		widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Text: strconv.Itoa(globals.CurrentProduct.ProductID),
		}))
	// Person.
	colPerson := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-auto m-1"},
		},
	})
	colPerson.AppendChild(
		widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
				Classes: []string{"iconlabel"},
			},
			Text: globals.CurrentProduct.Person.PersonEmail,
		}))

	//
	// Category and tags.
	//
	// Category.
	colCategory := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-auto m-1"},
		},
	})
	if globals.CurrentProduct.Category.CategoryID.Valid {
		colCategory.AppendChild(
			widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"iconlabel", "mr-sm-2"},
				},
				Text: locales.Translate("category_label_title", HTTPHeaderAcceptLanguage),
			}))
		colCategory.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Text: globals.CurrentProduct.Category.CategoryLabel.String,
		}))
	}
	// Tags.
	colTags := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-auto m-1"},
		},
	})
	for _, tag := range globals.CurrentProduct.Tags {
		colTags.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
				Classes: []string{"mr-sm-2", "badge", "badge-secondary"},
			},
			Text: tag.TagLabel,
		}))
	}

	//
	// Suppliers and producer.
	//
	// Producer.
	colProducer := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-auto m-1"},
		},
	})
	if globals.CurrentProduct.ProducerRef.ProducerRefID.Valid {
		colProducer.AppendChild(
			widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"iconlabel", "mr-sm-2"},
				},
				Text: locales.Translate("producer_label_title", HTTPHeaderAcceptLanguage),
			}))
		colProducer.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Text: fmt.Sprintf("%s: %s", globals.CurrentProduct.Producer.ProducerLabel.String, globals.CurrentProduct.ProducerRef.ProducerRefLabel.String),
		}))
	}
	// Suppliers.
	colSuppliers := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-auto m-1"},
		},
	})
	if len(globals.CurrentProduct.SupplierRefs) > 0 {
		colSuppliers.AppendChild(
			widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"iconlabel", "mr-sm-4"},
				},
				Text: locales.Translate("supplier_label_title", HTTPHeaderAcceptLanguage),
			}))
		ul := widgets.NewUl(widgets.UlAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
		})
		for _, s := range globals.CurrentProduct.SupplierRefs {
			ul.AppendChild(widgets.NewLi(widgets.LiAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
				},
				Text: fmt.Sprintf("%s: %s", s.Supplier.SupplierLabel.String, s.SupplierRefLabel),
			}))
		}
		colSuppliers.AppendChild(ul)
	}

	//
	// Number per carton and per bag.
	//
	// Carton.
	colNumberPerCarton := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-auto m-1"},
		},
	})
	if globals.CurrentProduct.ProductNumberPerCarton.Valid && globals.CurrentProduct.ProductNumberPerCarton.Int64 > 0 {
		colNumberPerCarton.AppendChild(
			widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"iconlabel", "mr-sm-2"},
				},
				Text: locales.Translate("product_number_per_carton_title", HTTPHeaderAcceptLanguage),
			}))
		colNumberPerCarton.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Text: fmt.Sprintf("%d", globals.CurrentProduct.ProductNumberPerCarton.Int64),
		}))
	}
	// Bag.
	colNumberPerBag := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-auto m-1"},
		},
	})
	if globals.CurrentProduct.ProductNumberPerBag.Valid {
		colNumberPerBag.AppendChild(
			widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"iconlabel", "mr-sm-2"},
				},
				Text: locales.Translate("product_number_per_bag_title", HTTPHeaderAcceptLanguage),
			}))
		colNumberPerBag.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Text: fmt.Sprintf("%d", globals.CurrentProduct.ProductNumberPerBag.Int64),
		}))
	}

	//
	// Producer sheet.
	//
	// Producer sheet.
	colProducerSheet := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-auto m-1"},
		},
	})
	if globals.CurrentProduct.ProductSheet.Valid {
		colProducerSheet.AppendChild(
			widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"iconlabel", "mr-sm-2"},
				},
				Text: locales.Translate("product_sheet_title", HTTPHeaderAcceptLanguage),
			}))

		var icon widgets.Widget
		icon.HTMLElement = widgets.NewIcon(widgets.IconAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Icon: themes.NewMdiIcon(themes.MDI_LINK, themes.MDI_24PX),
		})

		colProducerSheet.AppendChild(widgets.NewLink(widgets.LinkAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Target: "_blank",
			Href:   globals.CurrentProduct.ProductSheet.String,
			Label:  icon,
		}))
	}

	//
	// Cas, Ce and MSDS.
	//
	// Cas.
	colCas := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-auto m-1"},
		},
	})
	if globals.CurrentProduct.CasNumber.CasNumberID.Valid {
		colCas.AppendChild(
			widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"iconlabel", "mr-sm-2"},
				},
				Text: locales.Translate("casnumber_label_title", HTTPHeaderAcceptLanguage),
			}))
		colCas.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Text: globals.CurrentProduct.CasNumber.CasNumberLabel.String,
		}))
	}
	if globals.CurrentProduct.CasNumber.CasNumberCMR.Valid {
		colCas.AppendChild(
			widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"iconlabel", "mr-sm-2"},
				},
				Text: locales.Translate("casnumber_cmr_title", HTTPHeaderAcceptLanguage),
			}))
		colCas.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Text: globals.CurrentProduct.CasNumber.CasNumberCMR.String,
		}))
	}
	for _, hs := range globals.CurrentProduct.HazardStatements {
		if hs.HazardStatementCMR.Valid {
			colCas.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
				},
				Text: hs.HazardStatementCMR.String,
			}))
		}
	}
	// Ce.
	colCe := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-auto m-1"},
		},
	})
	if globals.CurrentProduct.CeNumber.CeNumberID.Valid {
		colCe.AppendChild(
			widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"iconlabel", "mr-sm-2"},
				},
				Text: locales.Translate("cenumber_label_title", HTTPHeaderAcceptLanguage),
			}))
		colCe.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Text: globals.CurrentProduct.CeNumber.CeNumberLabel.String,
		}))
	}
	// MSDS.
	colMsds := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-auto m-1"},
		},
	})
	if globals.CurrentProduct.ProductMSDS.Valid {
		colMsds.AppendChild(
			widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"iconlabel", "mr-sm-2"},
				},
				Text: locales.Translate("product_msds_title", HTTPHeaderAcceptLanguage),
			}))

		var icon widgets.Widget
		icon.HTMLElement = widgets.NewIcon(widgets.IconAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Icon: themes.NewMdiIcon(themes.MDI_LINK, themes.MDI_16PX),
		})

		colMsds.AppendChild(widgets.NewLink(widgets.LinkAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Target: "_blank",
			Href:   globals.CurrentProduct.ProductMSDS.String,
			Label:  icon,
		}))
	}

	//
	// Recommended storage temperature
	//
	// Storage temperature.
	colStorageTemperature := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-auto m-1"},
		},
	})
	if globals.CurrentProduct.ProductTemperature.Valid {
		colStorageTemperature.AppendChild(
			widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"iconlabel", "mr-sm-2"},
				},
				Text: locales.Translate("product_temperature_title", HTTPHeaderAcceptLanguage),
			}))
		colStorageTemperature.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Text: fmt.Sprintf("%d%s", globals.CurrentProduct.ProductTemperature.Int64, globals.CurrentProduct.UnitTemperature.UnitLabel.String),
		}))
	}

	//
	// Empirical, linear and 3D formulas.
	//
	// Empirical formula.
	colEmpiricalFormula := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-auto m-1"},
		},
	})
	if globals.CurrentProduct.EmpiricalFormula.EmpiricalFormulaID.Valid {
		colEmpiricalFormula.AppendChild(
			widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"iconlabel", "mr-sm-2"},
				},
				Text: locales.Translate("empiricalformula_label_title", HTTPHeaderAcceptLanguage),
			}))
		colEmpiricalFormula.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Text: globals.CurrentProduct.EmpiricalFormula.EmpiricalFormulaLabel.String,
		}))
	}
	// Linear formula.
	colLinearFormula := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-auto m-1"},
		},
	})
	if globals.CurrentProduct.LinearFormula.LinearFormulaID.Valid {
		colLinearFormula.AppendChild(
			widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"iconlabel", "mr-sm-2"},
				},
				Text: locales.Translate("linearformula_label_title", HTTPHeaderAcceptLanguage),
			}))
		colLinearFormula.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Text: globals.CurrentProduct.LinearFormula.LinearFormulaLabel.String,
		}))
	}
	// 3D formula.
	colTreedFormula := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-auto m-1"},
		},
	})
	if globals.CurrentProduct.ProductThreeDFormula.Valid && globals.CurrentProduct.ProductThreeDFormula.String != "" {
		colTreedFormula.AppendChild(
			widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"iconlabel", "mr-sm-2"},
				},
				Text: locales.Translate("product_threedformula_title", HTTPHeaderAcceptLanguage),
			}))

		var icon widgets.Widget
		icon.HTMLElement = widgets.NewIcon(widgets.IconAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Icon: themes.NewMdiIcon(themes.MDI_LINK, themes.MDI_16PX),
		})

		colTreedFormula.AppendChild(widgets.NewLink(widgets.LinkAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Target: "_blank",
			Href:   globals.CurrentProduct.ProductThreeDFormula.String,
			Label:  icon,
		}))
	}

	//
	// Symbols, signal word and physical state.
	//
	// Symbols.
	colSymbols := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-12 m-1"},
		},
	})
	for _, symbol := range globals.CurrentProduct.Symbols {
		colSymbols.AppendChild(widgets.NewImg(widgets.ImgAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Alt:   symbol.SymbolLabel,
			Title: symbol.SymbolLabel,
			Src:   fmt.Sprintf("data:%s", symbol.SymbolImage),
		}))
	}
	// Signal word.
	colSignalWord := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-auto m-1"},
		},
	})
	if globals.CurrentProduct.SignalWord.SignalWordID.Valid {
		colSignalWord.AppendChild(
			widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"iconlabel", "mr-sm-4"},
				},
				Text: locales.Translate("signalword_label_title", HTTPHeaderAcceptLanguage),
			}))
		colSignalWord.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Text: globals.CurrentProduct.SignalWord.SignalWordLabel.String,
		}))
	}
	// Physical state.
	colPhysicalState := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-auto m-1"},
		},
	})
	if globals.CurrentProduct.PhysicalState.PhysicalStateID.Valid {
		colPhysicalState.AppendChild(
			widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"iconlabel", "mr-sm-4"},
				},
				Text: locales.Translate("physicalstate_label_title", HTTPHeaderAcceptLanguage),
			}))
		colPhysicalState.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Text: globals.CurrentProduct.PhysicalState.PhysicalStateLabel.String,
		}))
	}

	//
	// Hazard statements, precautionary statements, classes of compounds
	//
	// Hazard statements.
	colHs := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-12 m-1"},
		},
	})
	if len(globals.CurrentProduct.HazardStatements) > 0 {
		colHs.AppendChild(
			widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"iconlabel", "mr-sm-4"},
				},
				Text: locales.Translate("hazardstatement_label_title", HTTPHeaderAcceptLanguage),
			}))

		ul := widgets.NewUl(widgets.UlAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
		})
		for _, s := range globals.CurrentProduct.HazardStatements {
			ul.AppendChild(widgets.NewLi(widgets.LiAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
				},
				Text: fmt.Sprintf("%s: %s", s.HazardStatementReference, s.HazardStatementLabel),
			}))
		}
		colHs.AppendChild(ul)
	}
	// Precautionary statements.
	colPs := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-12 m-1"},
		},
	})
	if len(globals.CurrentProduct.PrecautionaryStatements) > 0 {
		colPs.AppendChild(
			widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"iconlabel", "mr-sm-4"},
				},
				Text: locales.Translate("precautionarystatement_label_title", HTTPHeaderAcceptLanguage),
			}))

		ul := widgets.NewUl(widgets.UlAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
		})
		for _, s := range globals.CurrentProduct.PrecautionaryStatements {
			ul.AppendChild(widgets.NewLi(widgets.LiAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
				},
				Text: fmt.Sprintf("%s: %s", s.PrecautionaryStatementReference, s.PrecautionaryStatementLabel),
			}))
		}
		colPs.AppendChild(ul)
	}
	// Classes of compounds.
	colCoc := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-auto m-1"},
		},
	})
	if len(globals.CurrentProduct.ClassOfCompound) > 0 {
		colCoc.AppendChild(
			widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"iconlabel", "mr-sm-4"},
				},
				Text: locales.Translate("classofcompound_label_title", HTTPHeaderAcceptLanguage),
			}))
		ul := widgets.NewUl(widgets.UlAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
		})
		for _, s := range globals.CurrentProduct.ClassOfCompound {
			ul.AppendChild(widgets.NewLi(widgets.LiAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
				},
				Text: s.ClassOfCompoundLabel,
			}))
		}
		colCoc.AppendChild(ul)
	}

	//
	// Disposal comment and remark.
	//
	// Disposal comment.
	colDisposalComment := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-auto m-1"},
		},
	})
	if globals.CurrentProduct.ProductDisposalComment.Valid && globals.CurrentProduct.ProductDisposalComment.String != "" {
		colDisposalComment.AppendChild(
			widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"iconlabel", "mr-sm-2"},
				},
				Text: locales.Translate("product_disposalcomment_title", HTTPHeaderAcceptLanguage),
			}))
		colDisposalComment.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Text: globals.CurrentProduct.ProductDisposalComment.String,
		}))
	}
	// Remark.
	colRemark := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-auto m-1"},
		},
	})
	if globals.CurrentProduct.ProductRemark.Valid && globals.CurrentProduct.ProductRemark.String != "" {
		colRemark.AppendChild(
			widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"iconlabel", "mr-sm-2"},
				},
				Text: locales.Translate("product_remark_title", HTTPHeaderAcceptLanguage),
			}))
		colRemark.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Text: globals.CurrentProduct.ProductRemark.String,
		}))
	}

	//
	// Radioactive, restricted.
	//
	// Radioactive.
	colRadioactive := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-auto m-1", "iconlabel"},
		},
	})
	if globals.CurrentProduct.ProductRadioactive.Valid && globals.CurrentProduct.ProductRadioactive.Bool {
		colRadioactive.AppendChild(widgets.NewIcon(widgets.IconAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Text:  locales.Translate("product_radioactive_title", HTTPHeaderAcceptLanguage),
			Title: locales.Translate("product_radioactive_title", HTTPHeaderAcceptLanguage),
			Icon:  themes.NewMdiIcon(themes.MDI_RADIOACTIVE, themes.MDI_24PX),
		}))
	}
	// Restricted.
	colRestricted := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-auto m-1", "iconlabel"},
		},
	})
	if globals.CurrentProduct.ProductRestricted.Valid && globals.CurrentProduct.ProductRestricted.Bool {
		colRestricted.AppendChild(widgets.NewIcon(widgets.IconAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Text:  locales.Translate("product_restricted_title", HTTPHeaderAcceptLanguage),
			Title: locales.Translate("product_restricted_title", HTTPHeaderAcceptLanguage),
			Icon:  themes.NewMdiIcon(themes.MDI_RESTRICTED, themes.MDI_24PX),
		}))
	}

	detailCard := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"card"},
		}})
	detailCardBody := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"card-body"},
		}})
	detailCard.AppendChild(detailCardBody)
	detailCardRow := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"row"},
		}})
	detailCardBody.AppendChild(detailCardRow)

	// colPin := widgets.NewDiv(widgets.DivAttributes{
	// 	BaseAttributes: widgets.BaseAttributes{
	// 		Visible: true,
	// 		Classes: []string{"offset-sm-5 col-sm-7 mb-2"},
	// 	},
	// })
	// colPin.AppendChild(widgets.NewSpan(
	// 	widgets.SpanAttributes{
	// 		BaseAttributes: widgets.BaseAttributes{
	// 			Visible: true,
	// 			Classes: []string{"mdi mdi-clipboard-list-outline mdi-36px"},
	// 		},
	// 	},
	// ))

	// detailCardRow.AppendChild(colPin)

	detailCardRow.AppendChild(colID)
	if globals.CurrentProduct.Category.CategoryID.Valid {
		detailCardRow.AppendChild(colCategory)
	}
	if len(globals.CurrentProduct.Tags) > 0 {
		detailCardRow.AppendChild(colTags)
	}
	if globals.CurrentProduct.ProductTwoDFormula.Valid && globals.CurrentProduct.ProductTwoDFormula.String != "" {
		detailCardRow.AppendChild(col2dimage)
	}
	if len(globals.CurrentProduct.Synonyms) > 0 {
		detailCardRow.AppendChild(colSynonym)
	}
	if globals.CurrentProduct.EmpiricalFormula.EmpiricalFormulaID.Valid {
		detailCardRow.AppendChild(colEmpiricalFormula)
	}
	if globals.CurrentProduct.LinearFormula.LinearFormulaID.Valid {
		detailCardRow.AppendChild(colLinearFormula)
	}
	if globals.CurrentProduct.ProductThreeDFormula.Valid && globals.CurrentProduct.ProductThreeDFormula.String != "" {
		detailCardRow.AppendChild(colTreedFormula)
	}
	if globals.CurrentProduct.CasNumber.CasNumberID.Valid {
		detailCardRow.AppendChild(colCas)
	}
	if globals.CurrentProduct.CeNumber.CeNumberID.Valid {
		detailCardRow.AppendChild(colCe)
	}
	if globals.CurrentProduct.ProductMSDS.Valid && globals.CurrentProduct.ProductMSDS.String != "" {
		detailCardRow.AppendChild(colMsds)
	}
	if globals.CurrentProduct.Producer.ProducerID.Valid {
		detailCardRow.AppendChild(colProducer)
	}
	if len(globals.CurrentProduct.SupplierRefs) > 0 {
		detailCardRow.AppendChild(colSuppliers)
	}
	if globals.CurrentProduct.ProductNumberPerCarton.Valid {
		detailCardRow.AppendChild(colNumberPerCarton)
	}
	if globals.CurrentProduct.ProductNumberPerBag.Valid {
		detailCardRow.AppendChild(colNumberPerBag)
	}
	if len(globals.CurrentProduct.Symbols) > 0 {
		detailCardRow.AppendChild(colSymbols)
	}
	if globals.CurrentProduct.SignalWord.SignalWordID.Valid {
		detailCardRow.AppendChild(colSignalWord)
	}
	if len(globals.CurrentProduct.HazardStatements) > 0 {
		detailCardRow.AppendChild(colHs)
	}
	if len(globals.CurrentProduct.PrecautionaryStatements) > 0 {
		detailCardRow.AppendChild(colPs)
	}
	if len(globals.CurrentProduct.ClassOfCompound) > 0 {
		detailCardRow.AppendChild(colCoc)
	}
	if globals.CurrentProduct.PhysicalState.PhysicalStateID.Valid {
		detailCardRow.AppendChild(colPhysicalState)
	}
	if globals.CurrentProduct.ProductSheet.Valid && globals.CurrentProduct.ProductSheet.String != "" {
		detailCardRow.AppendChild(colProducerSheet)
	}
	if globals.CurrentProduct.ProductTemperature.Valid {
		detailCardRow.AppendChild(colStorageTemperature)
	}
	if globals.CurrentProduct.ProductDisposalComment.Valid && globals.CurrentProduct.ProductDisposalComment.String != "" {
		detailCardRow.AppendChild(colDisposalComment)
	}
	if globals.CurrentProduct.ProductRemark.Valid && globals.CurrentProduct.ProductRemark.String != "" {
		detailCardRow.AppendChild(colRemark)
	}
	if globals.CurrentProduct.ProductRadioactive.Valid && globals.CurrentProduct.ProductRadioactive.Bool {
		detailCardRow.AppendChild(colRadioactive)
	}
	if globals.CurrentProduct.ProductRestricted.Valid && globals.CurrentProduct.ProductRestricted.Bool {
		detailCardRow.AppendChild(colRestricted)
	}
	detailCardRow.AppendChild(colPerson)

	return detailCard.OuterHTML()

}

func NameFormatter(this js.Value, args []js.Value) interface{} {

	row := args[1]
	globals.CurrentProduct = Product{Product: &models.Product{}}.ProductFromJsJSONValue(row)

	result := fmt.Sprintf("%s <i>%s</i>", globals.CurrentProduct.Name.NameLabel, globals.CurrentProduct.ProductSpecificity.String)

	for _, syn := range globals.CurrentProduct.Synonyms {
		if strings.HasPrefix(syn.NameLabel, "|") && strings.HasSuffix(syn.NameLabel, "|") {
			result += fmt.Sprintf(" %s", syn.NameLabel)
		}
	}

	if globals.CurrentProduct.ProductSL.Valid && globals.CurrentProduct.ProductSL.String != "" {
		result += fmt.Sprintf("<div><span class='text-white badge bg-secondary'>%s</span></div>", globals.CurrentProduct.ProductSL.String)
	}

	return result

}

func EmpiricalformulaFormatter(this js.Value, args []js.Value) interface{} {

	row := args[1]
	globals.CurrentProduct = Product{Product: &models.Product{}}.ProductFromJsJSONValue(row)

	if globals.CurrentProduct.EmpiricalFormulaID.Valid {
		return globals.CurrentProduct.EmpiricalFormulaLabel.String
	} else {
		return ""
	}

}

func TwodformulaFormatter(this js.Value, args []js.Value) interface{} {

	row := args[1]
	globals.CurrentProduct = Product{Product: &models.Product{}}.ProductFromJsJSONValue(row)

	var (
		twodformula string
	)

	if globals.CurrentProduct.ProductTwoDFormula.Valid {

		twodformula = widgets.NewImg(widgets.ImgAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Height:          "70",
			Src:             globals.CurrentProduct.ProductTwoDFormula.String,
			BackgroundColor: "white",
		}).OuterHTML()

	}

	return twodformula

}

func CasnumberFormatter(this js.Value, args []js.Value) interface{} {

	row := args[1]
	globals.CurrentProduct = Product{Product: &models.Product{}}.ProductFromJsJSONValue(row)

	var (
		spanCasNumber, spanCasCMR, spanHSCMR, imgSGH02, iconRestricted string
	)

	if globals.CurrentProduct.CasNumberID.Valid {
		spanCasNumber = widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Text: globals.CurrentProduct.CasNumberLabel.String,
		}).OuterHTML()
	}

	if globals.CurrentProduct.CasNumberCMR.Valid {
		spanCasCMR = widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
				Classes: []string{"text-danger", "font-italic"},
			},
			Text: globals.CurrentProduct.CasNumber.CasNumberCMR.String,
		}).OuterHTML()
	}

	for _, hs := range globals.CurrentProduct.HazardStatements {
		if hs.HazardStatementCMR.Valid {
			spanHSCMR = widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"text-danger", "font-italic"},
				},
				Text: hs.HazardStatementCMR.String,
			}).OuterHTML()
		}
	}

	for _, symbol := range globals.CurrentProduct.Symbols {
		if symbol.SymbolLabel == "SGH02" {
			imgSGH02 = widgets.NewImg(widgets.ImgAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
				},
				Src:   "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACYAAAAmCAYAAACoPemuAAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAAN1wAADdcBQiibeAAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAAIvSURBVFiFzdgxbI5BGMDx36uNJsJApFtFIh2QEIkBtYiYJFKsrDaJqYOofhKRMFhsFhNNMBgkTAaD0KUiBomN1EpYBHWGnvi8vu9r7/2eti65vHfv3T3P/57nuXvfuyqlJCRVVQuk1AqRl1LqP9NKpJxbETKjocLgYi0VaLl49wXBxUIFwsVDBcGFQWEAg1FwYZbCGMajLBfmPkzgUZRbw2IKFzGPrRFw/bpvD/bn8jUkXM719f3A9eu+k3iXA/92Bnub2yYx1NgDfbrvXIYZx8dcThjBExxvPOmGltqLIzmuEt63QSVczc+z/2whSw2ThpbajS+4UgOq59O4gYFSuGaByWb8zKvwN8RXXKiBPc7PLaWx3ARqY37O1CBe5/cvO1huVy+ZnfSX7y9MYxRTNeX32lZj+/sXWNfVnV3g1tT/aJeQ5vAGp3L9eXbjTFv7NzzM9VncSSnNF2lp4MqjNYvcxwEcy+0HcQg32/q8Kndl+YrcgM9Z4YdsrZ21PtvxHT9yv1vNgr8cbiIrnMUmbKu177PwVZjLgKPNt4sCOKzF0ww32aF9CA+yxSZKoTqDlVnucI6lMxhpg76OuxhrKr8oIENyXx/xxQKTE/hUkIdLJ1tlRd3TwtF/KtcuSalVVdUwdvQe+Fd6ljhfl9NzRKT5I8cvq/B+xi3vzFfk+FaqbEUPvEtVuipXBIspX9VLlW4Q/8U1VGe4EKgYsED3tefBgt271y7dUlV/ygHpF8bRglXiwx7BAAAAAElFTkSuQmCC",
				Alt:   locales.Translate("product_flammable", HTTPHeaderAcceptLanguage),
				Title: locales.Translate("product_flammable", HTTPHeaderAcceptLanguage),
			}).OuterHTML()
			break
		}
	}

	if globals.CurrentProduct.ProductRestricted.Valid && globals.CurrentProduct.ProductRestricted.Bool {
		iconRestricted = widgets.NewIcon(widgets.IconAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Title: locales.Translate("restricted", HTTPHeaderAcceptLanguage),
			Icon:  themes.NewMdiIcon(themes.MDI_RESTRICTED, ""),
		}).OuterHTML()
	}

	return spanCasNumber + spanCasCMR + spanHSCMR + imgSGH02 + iconRestricted

}

func Product_productSpecificityFormatter(this js.Value, args []js.Value) interface{} {

	row := args[1]
	globals.CurrentProduct = Product{Product: &models.Product{}}.ProductFromJsJSONValue(row)

	if globals.CurrentProduct.ProductSpecificity.Valid {
		return globals.CurrentProduct.ProductSpecificity.String
	} else {
		return ""
	}

}

func Product_productSlFormatter(this js.Value, args []js.Value) interface{} {

	row := args[1]
	globals.CurrentProduct = Product{Product: &models.Product{}}.ProductFromJsJSONValue(row)

	if globals.CurrentProduct.ProductSL.Valid {
		return globals.CurrentProduct.ProductSL.String
	} else {
		return ""
	}

}

func OperateFormatter(this js.Value, args []js.Value) interface{} {

	var (
		buttonStorages   string
		buttonOStorages  string
		buttonTotalStock string
		iconBookmark     themes.IconFace
		textBookmark     string
	)

	row := args[1]
	globals.CurrentProduct = Product{Product: &models.Product{}}.ProductFromJsJSONValue(row)

	if globals.CurrentProduct.Bookmark.BookmarkID.Valid {
		iconBookmark = themes.MDI_BOOKMARK
		textBookmark = locales.Translate("unbookmark", HTTPHeaderAcceptLanguage)
	} else {
		iconBookmark = themes.MDI_NO_BOOKMARK
		textBookmark = locales.Translate("bookmark", HTTPHeaderAcceptLanguage)
	}

	if globals.CurrentProduct.ProductSC != 0 || globals.CurrentProduct.ProductASC != 0 {

		buttonStorages = widgets.NewLink(
			widgets.LinkAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Id:      "storages" + strconv.Itoa(globals.CurrentProduct.ProductID),
					Classes: []string{"storages", "dropdown-item", "text-primary", "iconlabel"},
					Visible: false,
				},
				Href: "#",
				Label: widgets.NewSpan(
					widgets.SpanAttributes{
						BaseAttributes: widgets.BaseAttributes{
							Classes: []string{"mdi", themes.MDI_STORAGE.ToString()},
							Visible: true,
						},
						Text: fmt.Sprintf("%s %d  (%d)", locales.Translate("storages", HTTPHeaderAcceptLanguage), globals.CurrentProduct.ProductSC, globals.CurrentProduct.ProductASC),
					},
				),
			},
		).OuterHTML()

	}

	if globals.CurrentProduct.ProductTSC != 0 {

		buttonOStorages = widgets.NewLink(
			widgets.LinkAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Id:      "ostorages" + strconv.Itoa(globals.CurrentProduct.ProductID),
					Classes: []string{"ostorages", "dropdown-item", "text-primary", "iconlabel"},
					Visible: false,
				},
				Href: "#",
				Label: widgets.NewSpan(
					widgets.SpanAttributes{
						BaseAttributes: widgets.BaseAttributes{
							Classes: []string{"mdi", themes.MDI_OSTORAGE.ToString()},
							Visible: true,
						},
						Text: locales.Translate("ostorages", HTTPHeaderAcceptLanguage),
					},
				),
			},
		).OuterHTML()

	}

	buttonStore := widgets.NewLink(
		widgets.LinkAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Id:      "store" + strconv.Itoa(globals.CurrentProduct.ProductID),
				Classes: []string{"store", "dropdown-item", "text-primary", "iconlabel"},
				Visible: false,
			},
			Href: "#",
			Label: widgets.NewSpan(
				widgets.SpanAttributes{
					BaseAttributes: widgets.BaseAttributes{
						Classes: []string{"mdi", themes.MDI_STORE.ToString()},
						Visible: true,
					},
					Text: locales.Translate("store", HTTPHeaderAcceptLanguage),
				},
			),
		},
	).OuterHTML()

	buttonEdit := widgets.NewLink(
		widgets.LinkAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Id:      "edit" + strconv.Itoa(globals.CurrentProduct.ProductID),
				Classes: []string{"productedit", "dropdown-item", "text-primary", "iconlabel"},
				Visible: false,
			},
			Href: "#",
			Label: widgets.NewSpan(
				widgets.SpanAttributes{
					BaseAttributes: widgets.BaseAttributes{
						Classes: []string{"mdi", themes.MDI_EDIT.ToString()},
						Visible: true,
					},
					Text: locales.Translate("edit", HTTPHeaderAcceptLanguage),
				},
			),
		},
	).OuterHTML()

	buttonDelete := widgets.NewLink(
		widgets.LinkAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Id:      "delete" + strconv.Itoa(globals.CurrentProduct.ProductID),
				Classes: []string{"productdelete", "dropdown-item", "text-primary", "iconlabel"},
				Visible: false,
			},
			Href: "#",
			Label: widgets.NewSpan(
				widgets.SpanAttributes{
					BaseAttributes: widgets.BaseAttributes{
						Classes: []string{"mdi", themes.MDI_DELETE.ToString()},
						Visible: true,
					},
					Text: locales.Translate("delete", HTTPHeaderAcceptLanguage),
				},
			),
		},
	).OuterHTML()

	buttonBookmark := widgets.NewLink(
		widgets.LinkAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Id:         "bookmark" + strconv.Itoa(globals.CurrentProduct.ProductID),
				Classes:    []string{"bookmark", "dropdown-item", "text-primary", "iconlabel"},
				Visible:    false,
				Attributes: map[string]string{"pid": strconv.Itoa(globals.CurrentProduct.ProductID)},
			},
			Href: "#",
			Label: widgets.NewSpan(
				widgets.SpanAttributes{
					BaseAttributes: widgets.BaseAttributes{
						Classes: []string{"mdi", iconBookmark.ToString()},
						Visible: true,
					},
					Text: textBookmark,
				},
			),
		},
	).OuterHTML()

	if globals.CurrentProduct.ProductSC != 0 {

		buttonTotalStock = widgets.NewLink(
			widgets.LinkAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Id:      "totalstock" + strconv.Itoa(globals.CurrentProduct.ProductID),
					Classes: []string{"totalstock", "dropdown-item", "text-primary", "iconlabel"},
					Visible: false,
				},
				Href: "#",
				Label: widgets.NewSpan(
					widgets.SpanAttributes{
						BaseAttributes: widgets.BaseAttributes{
							Classes: []string{"mdi", themes.MDI_TOTALSTOCK.ToString()},
							Visible: true,
						},
						Text: locales.Translate("totalstock_text", HTTPHeaderAcceptLanguage),
					},
				),
			},
		).OuterHTML()

	}

	ostoragesDiv := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Id:      "ostorages-collapse-" + strconv.Itoa(globals.CurrentProduct.ProductID),
			Classes: []string{"collapse", "float-left"},
		},
	}).OuterHTML()

	totalstockDiv := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Id:      "totalstock-collapse-" + strconv.Itoa(globals.CurrentProduct.ProductID),
			Classes: []string{"collapse", "float-left"},
		},
	}).OuterHTML()

	finalDiv := `
<div class="dropdown">
  <button class="btn btn-secondary dropdown-toggle" type="button" id="productActions` + strconv.Itoa(globals.CurrentProduct.ProductID) + `" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
    <span class="mdi mdi-menu">&nbsp;</span>
  </button>
  <div class="dropdown-menu" aria-labelledby="productActions` + strconv.Itoa(globals.CurrentProduct.ProductID) + `">
  ` + buttonStorages + buttonOStorages + buttonStore + buttonEdit + buttonDelete + buttonBookmark + buttonTotalStock +
		`
  </div>
</div>
<div id="confirm` + strconv.Itoa(globals.CurrentProduct.ProductID) + `">
</div>
`
	return finalDiv + ostoragesDiv + totalstockDiv

}

func AddProducer(this js.Value, args []js.Value) interface{} {

	var (
		producer  Producer
		dataBytes []byte
		err       error
	)

	producerLabel := jquery.Jq("input#addproducer").GetVal().String()
	producerLabel = strings.Trim(producerLabel, " ")

	if producerLabel == "" {
		return nil
	}

	producer = Producer{
		Producer: &models.Producer{
			ProducerLabel: sql.NullString{
				String: producerLabel,
				Valid:  true,
			},
		},
	}
	if dataBytes, err = json.Marshal(producer); err != nil {
		fmt.Println(err)
		return nil
	}

	ajax.Ajax{
		URL:    ApplicationProxyPath + "products/producers",
		Method: "post",
		Data:   dataBytes,
		Done: func(data js.Value) {

			jsutils.DisplaySuccessMessage(locales.Translate("producer_added", HTTPHeaderAcceptLanguage))
			jquery.Jq("input#addproducer").SetVal("")

		},
		Fail: func(jqXHR js.Value) {

			jsutils.DisplayGenericErrorMessage()

		},
	}.Send()

	return nil

}

func AddSupplier(this js.Value, args []js.Value) interface{} {

	var (
		supplier  Supplier
		dataBytes []byte
		err       error
	)

	supplierLabel := jquery.Jq("input#addsupplier").GetVal().String()
	supplierLabel = strings.Trim(supplierLabel, " ")

	if supplierLabel == "" {
		return nil
	}

	supplier = Supplier{
		Supplier: &models.Supplier{
			SupplierLabel: sql.NullString{
				String: supplierLabel,
				Valid:  true,
			},
		},
	}
	if dataBytes, err = json.Marshal(supplier); err != nil {
		fmt.Println(err)
		return nil
	}

	ajax.Ajax{
		URL:    ApplicationProxyPath + "products/suppliers",
		Method: "post",
		Data:   dataBytes,
		Done: func(data js.Value) {

			jsutils.DisplaySuccessMessage(locales.Translate("supplier_added", HTTPHeaderAcceptLanguage))
			jquery.Jq("input#addsupplier").SetVal("")

		},
		Fail: func(jqXHR js.Value) {

			jsutils.DisplayGenericErrorMessage()

		},
	}.Send()

	return nil

}

func Consufy() {

	jquery.Jq(".chem").Not(".consu").Hide()
	jquery.Jq(".bio").Not(".consu").Hide()
	jquery.Jq(".consu").Show()

	validate.NewValidate(jquery.Jq("select#producerref"), nil).ValidateAddRequired()
	validate.NewValidate(jquery.Jq("input#product_batchnumber"), nil).ValidateRemoveRequired()

	jquery.Jq("span#producerref.badge").Show()
	jquery.Jq("span#product_batchnumber.badge").Hide()

	validate.NewValidate(jquery.Jq("select#empiricalformula"), nil).ValidateRemoveRequired()
	validate.NewValidate(jquery.Jq("select#casnumber"), nil).ValidateRemoveRequired()

	jquery.Jq("span#empiricalformula.badge").Hide()
	jquery.Jq("span#casnumber.badge").Hide()

	jquery.Jq("input#showconsu").SetProp("checked", "checked")

}

func Chemify() {

	jquery.Jq(".bio").Not(".chem").Hide()
	jquery.Jq(".consu").Not(".chem").Hide()
	jquery.Jq(".chem").Show()

	validate.NewValidate(jquery.Jq("select#producerref"), nil).ValidateRemoveRequired()
	validate.NewValidate(jquery.Jq("input#product_batchnumber"), nil).ValidateRemoveRequired()

	jquery.Jq("span#producerref.badge").Hide()
	jquery.Jq("span#product_batchnumber.badge").Hide()

	validate.NewValidate(jquery.Jq("select#empiricalformula"), nil).ValidateAddRequired()
	validate.NewValidate(jquery.Jq("select#casnumber"), nil).ValidateAddRequired()

	jquery.Jq("span#empiricalformula.badge").Show()
	jquery.Jq("span#casnumber.badge").Show()

	jquery.Jq("input#showchem").SetProp("checked", "checked")

}

func Biofy() {

	jquery.Jq(".chem").Not(".bio").Hide()
	jquery.Jq(".consu").Not(".bio").Hide()
	jquery.Jq(".bio").Show()

	validate.NewValidate(jquery.Jq("select#producerref"), nil).ValidateAddRequired()
	validate.NewValidate(jquery.Jq("input#product_batchnumber"), nil).ValidateAddRequired()

	jquery.Jq("span#producerref.badge").Show()
	jquery.Jq("span#product_batchnumber.badge").Show()

	validate.NewValidate(jquery.Jq("select#empiricalformula"), nil).ValidateRemoveRequired()
	validate.NewValidate(jquery.Jq("select#casnumber"), nil).ValidateRemoveRequired()

	jquery.Jq("span#empiricalformula.badge").Hide()
	jquery.Jq("span#casnumber.badge").Hide()

	jquery.Jq("input#showbio").SetProp("checked", "checked")

}

func ShowIfAuthorizedActionButtons(this js.Value, args []js.Value) interface{} {

	jquery.Jq(".bookmark").FadeIn()
	jsutils.HasPermission("storages", "", "post", func() {
		jquery.Jq(".store").FadeIn()
	}, func() {
	})
	jsutils.HasPermission("storages", "-2", "get", func() {
		jquery.Jq(".storages").FadeIn()
		jquery.Jq(".ostorages").FadeIn()
		jquery.Jq(".totalstock").FadeIn()

		jquery.Jq("#switchview").SetVisible()
	}, func() {
	})
	jsutils.HasPermission("products", "-2", "put", func() {
		jquery.Jq(".productedit").FadeIn()
	}, func() {
	})

	// Iterating other the button with the class "storelocation"
	// (we could choose "members" or "delete")
	// to retrieve once the product id.
	buttons := dom.GetWindow().Document().GetElementsByTagName("a")
	for _, button := range buttons {
		if button.Class().Contains("bookmark") {
			productId := button.GetAttribute("pid")

			jsutils.HasPermission("products", productId, "delete", func() {
				jquery.Jq("#delete" + productId).FadeIn()
			}, func() {
			})
		}
	}

	return nil

}
