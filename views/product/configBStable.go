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
	"honnef.co/go/js/dom/v2"
)

func OperateEventsBookmark(this js.Value, args []js.Value) interface{} {

	row := args[2]
	globals.CurrentProduct = Product{}.ProductFromJsJSONValue(row)

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

	globals.CurrentProduct = Product{}.ProductFromJsJSONValue(row)

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
	globals.CurrentProduct = Product{}.ProductFromJsJSONValue(row)

	BSTableQueryFilter.Lock()
	BSTableQueryFilter.QueryFilter.Product = strconv.Itoa(globals.CurrentProduct.ProductID)
	BSTableQueryFilter.QueryFilter.ProductFilterLabel = fmt.Sprintf("%s %s", globals.CurrentProduct.Name.NameLabel, globals.CurrentProduct.ProductSpecificity.String)

	href := fmt.Sprintf("%sv/storages", ApplicationProxyPath)
	jsutils.LoadContent("div#content", "storage", href, storageCallbackWrapper)

	return nil

}

func OperateEventsOStorages(this js.Value, args []js.Value) interface{} {

	row := args[2]
	globals.CurrentProduct = Product{}.ProductFromJsJSONValue(row)

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

			spanEntityName := widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"iconlabel"},
				},
				Text: entity.EntityName,
			})
			spanEntityDescription := widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"blockquote-footer"},
				},
				Text: entity.EntityDescription,
			})
			divEntity := widgets.NewDiv(widgets.DivAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
				},
			})
			divEntity.AppendChild(spanEntityName)
			divEntity.AppendChild(spanEntityDescription)

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
	globals.CurrentProduct = Product{}.ProductFromJsJSONValue(row)

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
	product := Product{}.ProductFromJsJSONValue(row)

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
			storelocations []StoreLocation
			err            error
		)

		if err = json.Unmarshal([]byte(data.String()), &storelocations); err != nil {
			fmt.Println(err)
		}

		jquery.Jq(fmt.Sprintf("#totalstock-collapse-%d", product.ProductID)).SetHtml("")

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
		// 				"onclick": fmt.Sprintf("$('#totalstock-collapse-%d').html('')", product.ProductID),
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

		// jquery.Jq(fmt.Sprintf("#totalstock-collapse-%d", product.ProductID)).Append(rowButtonClose.OuterHTML())

		for _, storelocation := range storelocations {
			showStockRecursive(&storelocation, 0, fmt.Sprintf("#totalstock-collapse-%d", product.ProductID))
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
	globals.CurrentProduct = Product{}.ProductFromJsJSONValue(row)

	jquery.Jq(fmt.Sprintf("button#delete%d", globals.CurrentProduct.ProductID)).On("click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {

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

	buttonTitle := widgets.NewIcon(widgets.IconAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"iconlabel"},
		},
		Icon: themes.NewMdiIcon(themes.MDI_CONFIRM, ""),
		Text: locales.Translate("confirm", HTTPHeaderAcceptLanguage),
	})
	jquery.Jq(fmt.Sprintf("button#delete%d", globals.CurrentProduct.ProductID)).SetHtml("")
	jquery.Jq(fmt.Sprintf("button#delete%d", globals.CurrentProduct.ProductID)).Append(buttonTitle.OuterHTML())

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
	globals.CurrentProduct = Product{}.ProductFromJsJSONValue(row)

	var synonyms strings.Builder
	for _, synonym := range globals.CurrentProduct.Synonyms {
		synonyms.WriteString(synonym.NameLabel)
		synonyms.WriteString("<br/>")
	}

	// JSMol div.
	widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Id:      fmt.Sprintf("jsmol%d", globals.CurrentProduct.ProductID),
		},
	})

	//
	// 2dImage, synonyms, ID and person.
	//
	rowSynonymAndID := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"row", "mt-sm-3"},
		},
	})
	// 2dimage.
	col2dimage := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-4"},
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
			Classes: []string{"col-sm-4"},
		},
	})
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
			Classes: []string{"col-sm-2"},
		},
	})
	colID.AppendChild(
		widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
				Classes: []string{"iconlabel", "mr-sm-2"},
			},
			Text: "id",
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
			Classes: []string{"col-sm-2"},
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
	rowSynonymAndID.AppendChild(col2dimage)
	rowSynonymAndID.AppendChild(colSynonym)
	rowSynonymAndID.AppendChild(colID)
	rowSynonymAndID.AppendChild(colPerson)

	//
	// Category and tags.
	//
	rowCategoryAndTags := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"row", "mt-sm-3"},
		},
	})
	// Category.
	colCategory := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-10"},
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
			Classes: []string{"col-sm-2"},
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

	rowCategoryAndTags.AppendChild(colCategory)
	rowCategoryAndTags.AppendChild(colTags)

	//
	// Suppliers and producer.
	//
	rowSupplierAndProducer := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"row", "mt-sm-3"},
		},
	})
	// Producer.
	colProducer := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-6"},
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
			Classes: []string{"col-sm-6"},
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

	rowSupplierAndProducer.AppendChild(colProducer)
	rowSupplierAndProducer.AppendChild(colSuppliers)

	//
	// Number per carton and per bag.
	//
	rowNumberPerCartonAndBag := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"row", "mt-sm-3"},
		},
	})
	// Carton.
	colNumberPerCarton := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-6"},
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
			Classes: []string{"col-sm-6"},
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

	rowNumberPerCartonAndBag.AppendChild(colNumberPerCarton)
	rowNumberPerCartonAndBag.AppendChild(colNumberPerBag)

	//
	// Producer sheet.
	//
	rowProducerSheet := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"row", "mt-sm-3"},
		},
	})
	// Producer sheet.
	colProducerSheet := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-12"},
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

	rowProducerSheet.AppendChild(colProducerSheet)

	//
	// Cas, Ce and MSDS.
	//
	rowCasCeMsds := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"row", "mt-sm-3"},
		},
	})
	// Cas.
	colCas := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-4"},
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
			Classes: []string{"col-sm-4"},
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
			Classes: []string{"col-sm-4"},
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
			Icon: themes.NewMdiIcon(themes.MDI_LINK, themes.MDI_24PX),
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

	rowCasCeMsds.AppendChild(colCas)
	rowCasCeMsds.AppendChild(colCe)
	rowCasCeMsds.AppendChild(colMsds)

	//
	// Recommended storage temperature
	//
	rowStorageTemperature := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"row", "mt-sm-3"},
		},
	})
	// Storage temperature.
	colStorageTemperature := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-12"},
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

	rowStorageTemperature.AppendChild(colStorageTemperature)

	//
	// Empirical, linear and 3D formulas.
	//
	rowFormulas := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"row", "mt-sm-3"},
		},
	})
	// Empirical formula.
	colEmpiricalFormula := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-4"},
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
			Classes: []string{"col-sm-4"},
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
			Classes: []string{"col-sm-4"},
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
			Icon: themes.NewMdiIcon(themes.MDI_LINK, themes.MDI_24PX),
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

	rowFormulas.AppendChild(colEmpiricalFormula)
	rowFormulas.AppendChild(colLinearFormula)
	rowFormulas.AppendChild(colTreedFormula)

	//
	// Symbols, signal word and physical state.
	//
	rowSymbolsSignalWordPhysicalState := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"row", "mt-sm-3"},
		},
	})
	// Symbols.
	colSymbols := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-4"},
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
			Classes: []string{"col-sm-4"},
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
			Classes: []string{"col-sm-4"},
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

	rowSymbolsSignalWordPhysicalState.AppendChild(colSymbols)
	rowSymbolsSignalWordPhysicalState.AppendChild(colSignalWord)
	rowSymbolsSignalWordPhysicalState.AppendChild(colPhysicalState)

	//
	// Hazard statements, precautionary statements, classes of compounds
	//
	rowHsPsCoc := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"row", "mt-sm-3"},
		},
	})
	// Hazard statements.
	colHs := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-4"},
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
			Classes: []string{"col-sm-4"},
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
			Classes: []string{"col-sm-4"},
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

	rowHsPsCoc.AppendChild(colHs)
	rowHsPsCoc.AppendChild(colPs)
	rowHsPsCoc.AppendChild(colCoc)

	//
	// Disposal comment and remark.
	//
	rowDisposalCommentRemark := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"row", "mt-sm-3"},
		},
	})
	// Disposal comment.
	colDisposalComment := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-6"},
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
			Classes: []string{"col-sm-6"},
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

	rowDisposalCommentRemark.AppendChild(colDisposalComment)
	rowDisposalCommentRemark.AppendChild(colRemark)

	//
	// Radioactive, restricted.
	//
	rowRadioactiveRestricted := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"row", "mt-sm-3"},
		},
	})
	// Radioactive.
	colRadioactive := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-1", "iconlabel"},
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
			Classes: []string{"col-sm-1", "iconlabel"},
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

	rowRadioactiveRestricted.AppendChild(colRadioactive)
	rowRadioactiveRestricted.AppendChild(colRestricted)

	return rowSynonymAndID.OuterHTML() +
		rowCategoryAndTags.OuterHTML() +
		rowSupplierAndProducer.OuterHTML() +
		rowProducerSheet.OuterHTML() +
		rowNumberPerCartonAndBag.OuterHTML() +
		rowStorageTemperature.OuterHTML() +
		rowCasCeMsds.OuterHTML() +
		rowFormulas.OuterHTML() +
		rowSymbolsSignalWordPhysicalState.OuterHTML() +
		rowHsPsCoc.OuterHTML() +
		rowDisposalCommentRemark.OuterHTML() +
		rowRadioactiveRestricted.OuterHTML()

}

func NameFormatter(this js.Value, args []js.Value) interface{} {

	row := args[1]
	globals.CurrentProduct = Product{}.ProductFromJsJSONValue(row)

	result := fmt.Sprintf("%s <i>%s</i>", globals.CurrentProduct.Name.NameLabel, globals.CurrentProduct.ProductSpecificity.String)
	if globals.CurrentProduct.ProductSL.Valid && globals.CurrentProduct.ProductSL.String != "" {
		result += fmt.Sprintf("<div><span class='text-white badge bg-secondary'>%s</span></div>", globals.CurrentProduct.ProductSL.String)
	}

	return result

}

func EmpiricalformulaFormatter(this js.Value, args []js.Value) interface{} {

	row := args[1]
	globals.CurrentProduct = Product{}.ProductFromJsJSONValue(row)

	if globals.CurrentProduct.EmpiricalFormulaID.Valid {
		return globals.CurrentProduct.EmpiricalFormulaLabel.String
	} else {
		return ""
	}

}

func CasnumberFormatter(this js.Value, args []js.Value) interface{} {

	row := args[1]
	globals.CurrentProduct = Product{}.ProductFromJsJSONValue(row)

	if globals.CurrentProduct.CasNumberID.Valid {
		return globals.CurrentProduct.CasNumberLabel.String
	} else {
		return ""
	}

}

func Product_productSpecificityFormatter(this js.Value, args []js.Value) interface{} {

	row := args[1]
	globals.CurrentProduct = Product{}.ProductFromJsJSONValue(row)

	if globals.CurrentProduct.ProductSpecificity.Valid {
		return globals.CurrentProduct.ProductSpecificity.String
	} else {
		return ""
	}

}

func Product_productSlFormatter(this js.Value, args []js.Value) interface{} {

	row := args[1]
	globals.CurrentProduct = Product{}.ProductFromJsJSONValue(row)

	if globals.CurrentProduct.ProductSL.Valid {
		return globals.CurrentProduct.ProductSL.String
	} else {
		return ""
	}

}

func OperateFormatter(this js.Value, args []js.Value) interface{} {

	var (
		imgSGH02         string
		spanCasCMR       string
		spanHSCMR        string
		iconRestricted   string
		buttonStorages   string
		buttonOStorages  string
		buttonTotalStock string
		iconBookmark     themes.IconFace
		textBookmark     string
	)

	row := args[1]
	globals.CurrentProduct = Product{}.ProductFromJsJSONValue(row)

	if globals.CurrentProduct.Bookmark.BookmarkID.Valid {
		iconBookmark = themes.MDI_BOOKMARK
		textBookmark = locales.Translate("unbookmark", HTTPHeaderAcceptLanguage)
	} else {
		iconBookmark = themes.MDI_NO_BOOKMARK
		textBookmark = locales.Translate("bookmark", HTTPHeaderAcceptLanguage)
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

	if globals.CurrentProduct.ProductSC != 0 || globals.CurrentProduct.ProductASC != 0 {
		buttonStorages = widgets.NewBSButtonWithIcon(
			widgets.ButtonAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Id:      "storages" + strconv.Itoa(globals.CurrentProduct.ProductID),
					Classes: []string{"storages"},
					Visible: false,
				},
				Title: locales.Translate("storages", HTTPHeaderAcceptLanguage),
			},
			widgets.IconAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"iconlabel"},
				},
				Text: fmt.Sprintf("%s %d  (%d)", locales.Translate("storages", HTTPHeaderAcceptLanguage), globals.CurrentProduct.ProductSC, globals.CurrentProduct.ProductASC),
				Icon: themes.NewMdiIcon(themes.MDI_STORAGE, ""),
			},
			[]themes.BSClass{themes.BS_BTN, themes.BS_BNT_LINK},
		).OuterHTML()
	}

	if globals.CurrentProduct.ProductTSC != 0 {
		buttonOStorages = widgets.NewBSButtonWithIcon(
			widgets.ButtonAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Id:      "ostorages" + strconv.Itoa(globals.CurrentProduct.ProductID),
					Classes: []string{"ostorages"},
					Visible: false,
				},
				Title: locales.Translate("ostorages", HTTPHeaderAcceptLanguage),
			},
			widgets.IconAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"iconlabel"},
				},
				Text: locales.Translate("ostorages", HTTPHeaderAcceptLanguage),
				Icon: themes.NewMdiIcon(themes.MDI_OSTORAGE, ""),
			},
			[]themes.BSClass{themes.BS_BTN, themes.BS_BNT_LINK},
		).OuterHTML()
	}

	buttonStore := widgets.NewBSButtonWithIcon(
		widgets.ButtonAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Id:      "store" + strconv.Itoa(globals.CurrentProduct.ProductID),
				Classes: []string{"store"},
				Visible: false,
			},
			Title: locales.Translate("store", HTTPHeaderAcceptLanguage),
		},
		widgets.IconAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
				Classes: []string{"iconlabel"},
			},
			Text: locales.Translate("store", HTTPHeaderAcceptLanguage),
			Icon: themes.NewMdiIcon(themes.MDI_STORE, ""),
		},
		[]themes.BSClass{themes.BS_BTN, themes.BS_BNT_LINK},
	).OuterHTML()

	buttonEdit := widgets.NewBSButtonWithIcon(
		widgets.ButtonAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Id:      "edit" + strconv.Itoa(globals.CurrentProduct.ProductID),
				Classes: []string{"productedit"},
				Visible: false,
			},
			Title: locales.Translate("edit", HTTPHeaderAcceptLanguage),
		},
		widgets.IconAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
				Classes: []string{"iconlabel"},
			},
			Text: locales.Translate("edit", HTTPHeaderAcceptLanguage),
			Icon: themes.NewMdiIcon(themes.MDI_EDIT, ""),
		},
		[]themes.BSClass{themes.BS_BTN, themes.BS_BNT_LINK},
	).OuterHTML()

	buttonDelete := widgets.NewBSButtonWithIcon(
		widgets.ButtonAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Id:      "delete" + strconv.Itoa(globals.CurrentProduct.ProductID),
				Classes: []string{"productdelete"},
				Visible: false,
			},
			Title: locales.Translate("delete", HTTPHeaderAcceptLanguage),
		},
		widgets.IconAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
				Classes: []string{"iconlabel"},
			},
			Text: locales.Translate("delete", HTTPHeaderAcceptLanguage),
			Icon: themes.NewMdiIcon(themes.MDI_DELETE, ""),
		},
		[]themes.BSClass{themes.BS_BTN, themes.BS_BNT_LINK},
	).OuterHTML()

	buttonBookmark := widgets.NewBSButtonWithIcon(
		widgets.ButtonAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Id:         "bookmark" + strconv.Itoa(globals.CurrentProduct.ProductID),
				Classes:    []string{"bookmark"},
				Visible:    false,
				Attributes: map[string]string{"pid": strconv.Itoa(globals.CurrentProduct.ProductID)},
			},
			Title: textBookmark,
		},
		widgets.IconAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
				Classes: []string{"iconlabel"},
			},
			Text: textBookmark,
			Icon: themes.NewMdiIcon(iconBookmark, ""),
		},
		[]themes.BSClass{themes.BS_BTN, themes.BS_BNT_LINK},
	).OuterHTML()

	if globals.CurrentProduct.ProductSC != 0 {
		buttonTotalStock = widgets.NewBSButtonWithIcon(
			widgets.ButtonAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Id:         "totalstock" + strconv.Itoa(globals.CurrentProduct.ProductID),
					Classes:    []string{"totalstock"},
					Visible:    false,
					Attributes: map[string]string{"pid": strconv.Itoa(globals.CurrentProduct.ProductID)},
				},
				Title: locales.Translate("totalstock_text", HTTPHeaderAcceptLanguage),
			},
			widgets.IconAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"iconlabel"},
				},
				Text: locales.Translate("totalstock_text", HTTPHeaderAcceptLanguage),
				Icon: themes.NewMdiIcon(themes.MDI_TOTALSTOCK, ""),
			},
			[]themes.BSClass{themes.BS_BTN, themes.BS_BNT_LINK},
		).OuterHTML()
	}

	ostoragesDiv := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Id:      "ostorages-collapse-" + strconv.Itoa(globals.CurrentProduct.ProductID),
			Classes: []string{"collapse"},
		},
	}).OuterHTML()

	totalstockDiv := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Id:      "totalstock-collapse-" + strconv.Itoa(globals.CurrentProduct.ProductID),
			Classes: []string{"collapse", "p-sm-3", "float-right"},
		},
	}).OuterHTML()

	if globals.CurrentProduct.ProductRestricted.Valid && globals.CurrentProduct.ProductRestricted.Bool {
		iconRestricted = widgets.NewIcon(widgets.IconAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Title: locales.Translate("restricted", HTTPHeaderAcceptLanguage),
			Icon:  themes.NewMdiIcon(themes.MDI_RESTRICTED, ""),
		}).OuterHTML()
	}

	return buttonStorages + buttonOStorages + buttonStore + buttonEdit + buttonDelete + buttonBookmark + buttonTotalStock + ostoragesDiv + totalstockDiv + iconRestricted + spanCasCMR + spanHSCMR + imgSGH02

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
		ProducerLabel: sql.NullString{
			String: producerLabel,
			Valid:  true,
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
		SupplierLabel: sql.NullString{
			String: supplierLabel,
			Valid:  true,
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
	buttons := dom.GetWindow().Document().GetElementsByTagName("button")
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
