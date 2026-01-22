//go:build go1.24 && js && wasm

package product

import (
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

	url := fmt.Sprintf("%sbookmarks/%d", BackProxyPath, *globals.CurrentProduct.ProductID)
	method := "get"

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
	BSTableQueryFilter.QueryFilter.Product = strconv.Itoa(int(*globals.CurrentProduct.ProductID))
	if globals.CurrentProduct.ProductSpecificity != nil {
		BSTableQueryFilter.QueryFilter.ProductFilterLabel = fmt.Sprintf("%s %s", globals.CurrentProduct.Name.NameLabel, *globals.CurrentProduct.ProductSpecificity)
	} else {
		BSTableQueryFilter.QueryFilter.ProductFilterLabel = globals.CurrentProduct.Name.NameLabel
	}
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
	BSTableQueryFilter.QueryFilter.Product = strconv.Itoa(int(*globals.CurrentProduct.ProductID))
	if globals.CurrentProduct.ProductSpecificity != nil {
		BSTableQueryFilter.QueryFilter.ProductFilterLabel = fmt.Sprintf("%s %s", globals.CurrentProduct.Name.NameLabel, *globals.CurrentProduct.ProductSpecificity)
	} else {
		BSTableQueryFilter.QueryFilter.ProductFilterLabel = globals.CurrentProduct.Name.NameLabel
	}
	href := fmt.Sprintf("%sv/storages", ApplicationProxyPath)
	jsutils.LoadContent("div#content", "storage", href, storageCallbackWrapper)

	return nil

}

func OperateEventsOStorages(this js.Value, args []js.Value) interface{} {

	row := args[2]
	globals.CurrentProduct = Product{Product: &models.Product{}}.ProductFromJsJSONValue(row)

	url := fmt.Sprintf("%sstorages/others?product=%d", BackProxyPath, globals.CurrentProduct.ProductID)
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
	BSTableQueryFilter.QueryFilter.Product = strconv.Itoa(int(*globals.CurrentProduct.ProductID))
	if globals.CurrentProduct.ProductSpecificity != nil {
		BSTableQueryFilter.QueryFilter.ProductFilterLabel = fmt.Sprintf("%s %s", globals.CurrentProduct.Name.NameLabel, *globals.CurrentProduct.ProductSpecificity)
	} else {
		BSTableQueryFilter.QueryFilter.ProductFilterLabel = globals.CurrentProduct.Name.NameLabel
	}
	BSTableQueryFilter.Unlock()

	href := fmt.Sprintf("%svc/products", ApplicationProxyPath)
	jsutils.LoadContent("div#content", "storage", href, Product_createCallback, globals.CurrentProduct)

	return nil

}

func OperateEventsTotalStock(this js.Value, args []js.Value) interface{} {

	row := args[2]
	product := Product{Product: &models.Product{}}.ProductFromJsJSONValue(row)

	// jquery.Jq(fmt.Sprintf("#totalstock-collapse-%d", product.ProductID)).Append(widgets.NewSpan(widgets.SpanAttributes{
	// 	BaseAttributes: widgets.BaseAttributes{
	// 		Visible: true,
	// 		Classes: []string{"mdi", "mdi-loading", "mdi-spin", "mdi-36px"},
	// 	},
	// }).OuterHTML())

	jquery.Jq(fmt.Sprintf("#totalstock-collapse-%d", *product.ProductID)).Show()

	url := fmt.Sprintf("%sstocks/%d", BackProxyPath, *product.ProductID)
	method := "get"

	done := func(data js.Value) {

		var (
			stocks []models.Stock
			err    error
		)

		if err = json.Unmarshal([]byte(data.String()), &stocks); err != nil {
			jsutils.DisplayGenericErrorMessage()
			fmt.Println(err)
		}

		jquery.Jq("div#stock").SetHtml("")

		for _, stock := range stocks {
			jquery.Jq("div#stock").Append(fmt.Sprintf("<div class='col-sm-auto'><span>%s</span></div>", stock.StoreLocation.StoreLocationFullPath))

			unit := ""
			if stock.Unit != nil {
				unit = *stock.Unit.UnitLabel
			}
			jquery.Jq("div#stock").Append(fmt.Sprintf("<div class='col-sm-auto'><span>%.2f %s</span></div>", stock.Quantity, unit))
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

func OperateEventsSelect(this js.Value, args []js.Value) interface{} {

	row := args[2]
	globals.CurrentProduct = Product{Product: &models.Product{}}.ProductFromJsJSONValue(row)

	jquery.Jq("input[name=selected_product_id]").SetVal(*globals.CurrentProduct.ProductID)
	jquery.Jq("input[name=selected_product_name]").SetVal(globals.CurrentProduct.Name.NameLabel)

	jsutils.DisplaySuccessMessage(locales.Translate("selected", HTTPHeaderAcceptLanguage) + ": #" + strconv.Itoa(int(*globals.CurrentProduct.ProductID)) + " " + globals.CurrentProduct.Name.NameLabel)

	return nil
}

func OperateEventsDelete(this js.Value, args []js.Value) interface{} {

	row := args[2]
	globals.CurrentProduct = Product{Product: &models.Product{}}.ProductFromJsJSONValue(row)

	confirm := widgets.NewLink(
		widgets.LinkAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Id:      fmt.Sprintf("delete%d", *globals.CurrentProduct.ProductID),
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

	jquery.Jq(fmt.Sprintf("div#confirm%d", *globals.CurrentProduct.ProductID)).SetHtml(confirm)

	jquery.Jq(fmt.Sprintf("a#delete%d", *globals.CurrentProduct.ProductID)).On("click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		url := fmt.Sprintf("%sproducts/%d", BackProxyPath, *globals.CurrentProduct.ProductID)
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

	if params.Data.Sort == "name.name_label" {
		params.Data.Sort = "name"
	}
	if params.Data.Sort == "empirical_formula.empirical_formula_label" {
		params.Data.Sort = "empirical_formula"
	}
	if params.Data.Sort == "cas_number.cas_number_label" {
		params.Data.Sort = "cas_number"
	}
	// jsutils.ConsoleLog(fmt.Sprintf("%#v", params))

	go func() {

		var u *url.URL

		if params.Data.Export {
			u, _ = url.Parse(BackProxyPath + "products/export")

			params.Data.Offset = 0
			params.Data.Limit = 999999999999999999
		} else {
			u, _ = url.Parse(BackProxyPath + "products_old")

		}
		u.RawQuery = params.Data.ToRawQuery()

		// if params.Data.Export {
		// 	jsutils.RedirectTo(u.String())
		// 	return
		// }

		ajax := ajax.Ajax{
			URL:    u.String(),
			Method: "get",
			Done: func(data js.Value) {

				// jsutils.ConsoleLog(fmt.Sprintf("%#v", data.String()))

				if params.Data.Export {

					text := data.String()

					// Create Blob from text
					array := js.Global().Get("Array").New()
					array.Call("push", text)

					blob := js.Global().Get("Blob").New(
						array,
						map[string]any{
							"type": "text/plain;charset=utf-8",
						},
					)

					// Create object URL
					url := js.Global().Get("URL").Call("createObjectURL", blob)

					// Create <a> element
					document := js.Global().Get("document")
					a := document.Call("createElement", "a")
					a.Set("href", url)
					a.Set("download", "chimitheque_export.csv")

					document.Get("body").Call("appendChild", a)
					a.Call("click")
					document.Get("body").Call("removeChild", a)

					// Cleanup
					js.Global().Get("URL").Call("revokeObjectURL", url)

					return

					// jsutils.DisplaySuccessMessage(locales.Translate("export_done", HTTPHeaderAcceptLanguage))

					// var icon widgets.Widget
					// icon.HTMLElement = widgets.NewIcon(widgets.IconAttributes{
					// 	BaseAttributes: widgets.BaseAttributes{
					// 		Visible: true,
					// 	},
					// 	Text: locales.Translate("download_export", HTTPHeaderAcceptLanguage),
					// 	Icon: themes.NewMdiIcon(themes.MDI_DOWNLOAD, themes.MDI_24PX),
					// })

					// downloadLink := widgets.NewLink(widgets.LinkAttributes{
					// 	BaseAttributes: widgets.BaseAttributes{
					// 		Visible: true,
					// 		Classes: []string{"iconlabel"},
					// 	},
					// 	Onclick: "$('#export').collapse('hide')",
					// 	Title:   locales.Translate("download_export", HTTPHeaderAcceptLanguage),
					// 	Href:    fmt.Sprintf("%sdownload/%s", ApplicationProxyPath, products.GetExportFn()),
					// 	Label:   icon,
					// })

					// jquery.Jq("#export-body").SetHtml(downloadLink.OuterHTML())
					// jquery.Jq("#export").Show()
					// jquery.Jq("button#export").SetProp("disabled", false)

				} else {
					var (
						products Products
						err      error
					)
					if err = json.Unmarshal([]byte(data.String()), &products); err != nil {
						fmt.Println(err)
					}

					if products.GetTotal() != 0 {

						row.Call("success", js.ValueOf(js.Global().Get("JSON").Call("parse", data)))

					} else {

						// TODO: improve this
						jquery.Jq("span.loading-wrap").SetHtml(locales.Translate("no_result", globals.HTTPHeaderAcceptLanguage))

					}

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
	queryFilter.Id = BSTableQueryFilter.Id
	queryFilter.ProductFilterLabel = BSTableQueryFilter.ProductFilterLabel
	queryFilter.ProductBookmark = BSTableQueryFilter.ProductBookmark
	queryFilter.Export = BSTableQueryFilter.Export
	BSTableQueryFilter.Export = false
	BSTableQueryFilter.Unlock()

	select2SProducerRef := select2.NewSelect2(jquery.Jq("select#s_producer_ref"), nil)
	if select2SProducerRef.Select2IsInitialized() {
		i := select2SProducerRef.Select2Data()
		if len(i) > 0 {
			queryFilter.ProducerRef = i[0].Id
			queryFilter.ProducerRefFilterLabel = i[0].Text
		}
	}

	select2SStoreLocation := select2.NewSelect2(jquery.Jq("select#s_store_location"), nil)
	if select2SStoreLocation.Select2IsInitialized() {
		i := select2SStoreLocation.Select2Data()
		if len(i) > 0 {
			queryFilter.StoreLocation = i[0].Id
			queryFilter.StoreLocationFilterLabel = i[0].Text
		}
	}

	select2SEntity := select2.NewSelect2(jquery.Jq("select#s_entity"), nil)
	if select2SEntity.Select2IsInitialized() {
		i := select2SEntity.Select2Data()
		if len(i) > 0 {
			queryFilter.Entity = i[0].Id
			queryFilter.EntityFilterLabel = i[0].Text
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

	select2SCasNumber := select2.NewSelect2(jquery.Jq("select#s_cas_number"), nil)
	if select2SCasNumber.Select2IsInitialized() {
		i := select2SCasNumber.Select2Data()
		if len(i) > 0 {
			queryFilter.CasNumber = i[0].Id
			queryFilter.CasNumberFilterLabel = i[0].Text
		}
	}

	select2SEmpiricalFormula := select2.NewSelect2(jquery.Jq("select#s_empirical_formula"), nil)
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

	select2SSignalWord := select2.NewSelect2(jquery.Jq("select#s_signal_word"), nil)
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

	select2SHS := select2.NewSelect2(jquery.Jq("select#s_hazard_statements"), nil)
	if select2SHS.Select2IsInitialized() {
		i := select2SHS.Select2Data()
		if len(i) > 0 {
			for _, hs := range i {
				queryFilter.HazardStatements = append(queryFilter.HazardStatements, hs.Id)
				queryFilter.HazardStatementsFilterLabel += fmt.Sprintf(" %s", hs.Text)
			}
		}
	}

	select2SPS := select2.NewSelect2(jquery.Jq("select#s_precautionary_statements"), nil)
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

	if jquery.Jq("#s_storage_batch_number").GetVal().Truthy() {
		queryFilter.StorageBatchNumber = jquery.Jq("#s_storage_batch_number").GetVal().String()
		queryFilter.StorageBatchNumberFilterLabel = jquery.Jq("#s_storage_batch_number").GetVal().String()
	}
	if jquery.Jq("#s_storage").GetVal().Truthy() {
		queryFilter.Storage = jquery.Jq("#s_storage").GetVal().String()
		queryFilter.StorageFilterLabel = jquery.Jq("#s_storage").GetVal().String()
	}
	if jquery.Jq("#s_storage_barecode").GetVal().Truthy() {
		queryFilter.StorageBarecode = jquery.Jq("#s_storage_barecode").GetVal().String()
		queryFilter.StorageBarecodeFilterLabel = jquery.Jq("#s_storage_barecode").GetVal().String()
	}
	if jquery.Jq("#s_custom_name_part_of").GetVal().Truthy() {
		queryFilter.CustomNamePartOf = jquery.Jq("#s_custom_name_part_of").GetVal().String()
		queryFilter.CustomNamePartOfFilterLabel = jquery.Jq("#s_custom_name_part_of").GetVal().String()
	}
	if jquery.Jq("#s_cas_number_cmr:checked").Object.Length() > 0 {
		queryFilter.IsCMR = true
		queryFilter.CasNumberCMRFilterLabel = locales.Translate("s_cas_number_cmr", globals.HTTPHeaderAcceptLanguage)
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

	// jsutils.ConsoleLog(fmt.Sprintf("%#v", globals.CurrentProduct.Product))

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
	if globals.CurrentProduct.ProductTwoDFormula != nil {
		col2dimage.AppendChild(widgets.NewImg(widgets.ImgAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Width: "200",
			Src:   *globals.CurrentProduct.ProductTwoDFormula,
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
				Classes: []string{"iconlabel"},
			},
			Text: "#",
		}))
	colID.AppendChild(
		widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Text: strconv.Itoa(int(*globals.CurrentProduct.ProductID)) + globals.CurrentProduct.ProductType,
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
	if globals.CurrentProduct.Category != nil && globals.CurrentProduct.Category.CategoryID != nil {
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
			Text: *globals.CurrentProduct.Category.CategoryLabel,
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
	if globals.CurrentProduct.ProducerRef != nil && globals.CurrentProduct.ProducerRef.ProducerRefID != nil {
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
			Text: fmt.Sprintf("%s: %s", *globals.CurrentProduct.Producer.ProducerLabel, *globals.CurrentProduct.ProducerRef.ProducerRefLabel),
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
				Text: fmt.Sprintf("%s: %s", *s.Supplier.SupplierLabel, s.SupplierRefLabel),
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
	if globals.CurrentProduct.ProductNumberPerCarton != nil && *globals.CurrentProduct.ProductNumberPerCarton > 0 {
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
			Text: fmt.Sprintf("%d", *globals.CurrentProduct.ProductNumberPerCarton),
		}))
	}
	// Bag.
	colNumberPerBag := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-auto m-1"},
		},
	})
	if globals.CurrentProduct.ProductNumberPerBag != nil && *globals.CurrentProduct.ProductNumberPerBag > 0 {
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
			Text: fmt.Sprintf("%d", *globals.CurrentProduct.ProductNumberPerBag),
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
	if globals.CurrentProduct.ProductSheet != nil {
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
			Href:   *globals.CurrentProduct.ProductSheet,
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
	// if globals.CurrentProduct.CasNumber.CasNumberID!= nil {
	if globals.CurrentProduct.CasNumber != nil && globals.CurrentProduct.CasNumber.CasNumberID != nil {
		colCas.AppendChild(
			widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"iconlabel", "mr-sm-2"},
				},
				Text: locales.Translate("cas_number_label_title", HTTPHeaderAcceptLanguage),
			}))
		colCas.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Text: *globals.CurrentProduct.CasNumber.CasNumberLabel,
		}))
	}
	if globals.CurrentProduct.CasNumber != nil && globals.CurrentProduct.CasNumber.CasNumberCMR != nil {
		colCas.AppendChild(
			widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"iconlabel", "mr-sm-2"},
				},
				Text: locales.Translate("cas_number_cmr_title", HTTPHeaderAcceptLanguage),
			}))
		colCas.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Text: *globals.CurrentProduct.CasNumber.CasNumberCMR,
		}))
	}
	for _, hs := range globals.CurrentProduct.HazardStatements {
		if hs.HazardStatementCMR != nil {
			colCas.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
				},
				Text: *hs.HazardStatementCMR,
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
	if globals.CurrentProduct.CeNumber != nil && globals.CurrentProduct.CeNumber.CeNumberID != nil {
		colCe.AppendChild(
			widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"iconlabel", "mr-sm-2"},
				},
				Text: locales.Translate("ce_number_label_title", HTTPHeaderAcceptLanguage),
			}))
		colCe.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Text: *globals.CurrentProduct.CeNumber.CeNumberLabel,
		}))
	}
	// MSDS.
	colMsds := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-auto m-1"},
		},
	})
	if globals.CurrentProduct.ProductMSDS != nil {
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
			Href:   *globals.CurrentProduct.ProductMSDS,
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
	if globals.CurrentProduct.ProductTemperature != nil {
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
			Text: fmt.Sprintf("%f%s", *globals.CurrentProduct.ProductTemperature, *globals.CurrentProduct.UnitTemperature.UnitLabel),
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
	if globals.CurrentProduct.EmpiricalFormula != nil && globals.CurrentProduct.EmpiricalFormula.EmpiricalFormulaID != nil {
		colEmpiricalFormula.AppendChild(
			widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"iconlabel", "mr-sm-2"},
				},
				Text: locales.Translate("empirical_formula_label_title", HTTPHeaderAcceptLanguage),
			}))
		colEmpiricalFormula.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Text: *globals.CurrentProduct.EmpiricalFormula.EmpiricalFormulaLabel,
		}))
	}
	// Linear formula.
	colLinearFormula := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-auto m-1"},
		},
	})
	if globals.CurrentProduct.LinearFormula != nil && globals.CurrentProduct.LinearFormula.LinearFormulaID != nil {
		colLinearFormula.AppendChild(
			widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"iconlabel", "mr-sm-2"},
				},
				Text: locales.Translate("linear_formula_label_title", HTTPHeaderAcceptLanguage),
			}))
		colLinearFormula.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Text: *globals.CurrentProduct.LinearFormula.LinearFormulaLabel,
		}))
	}
	// 3D formula.
	colTreedFormula := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-auto m-1"},
		},
	})
	if globals.CurrentProduct.ProductThreeDFormula != nil && *globals.CurrentProduct.ProductThreeDFormula != "" {
		colTreedFormula.AppendChild(
			widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"iconlabel", "mr-sm-2"},
				},
				Text: locales.Translate("product_threed_formula_title", HTTPHeaderAcceptLanguage),
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
			Href:   *globals.CurrentProduct.ProductThreeDFormula,
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
			Height: "30",
			Width:  "30",
			Alt:    symbol.SymbolLabel,
			Title:  symbol.SymbolLabel,
			Src:    fmt.Sprintf("%sstatic/img/%s.svg", ApplicationProxyPath, symbol.SymbolLabel),
		}))
	}
	// Signal word.
	colSignalWord := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-auto m-1"},
		},
	})
	if globals.CurrentProduct.SignalWord != nil && globals.CurrentProduct.SignalWord.SignalWordID != nil {
		colSignalWord.AppendChild(
			widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"iconlabel", "mr-sm-4"},
				},
				Text: locales.Translate("signal_word_label_title", HTTPHeaderAcceptLanguage),
			}))
		colSignalWord.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Text: *globals.CurrentProduct.SignalWord.SignalWordLabel,
		}))
	}
	// Physical state.
	colPhysicalState := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-auto m-1"},
		},
	})
	if globals.CurrentProduct.PhysicalState != nil && globals.CurrentProduct.PhysicalState.PhysicalStateID != nil {
		colPhysicalState.AppendChild(
			widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"iconlabel", "mr-sm-4"},
				},
				Text: locales.Translate("physical_state_label_title", HTTPHeaderAcceptLanguage),
			}))
		colPhysicalState.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Text: *globals.CurrentProduct.PhysicalState.PhysicalStateLabel,
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
				Text: locales.Translate("hazard_statement_label_title", HTTPHeaderAcceptLanguage),
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
				Text: locales.Translate("precautionary_statement_label_title", HTTPHeaderAcceptLanguage),
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
				Text: locales.Translate("class_of_compound_label_title", HTTPHeaderAcceptLanguage),
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
	if globals.CurrentProduct.ProductDisposalComment != nil && *globals.CurrentProduct.ProductDisposalComment != "" {
		colDisposalComment.AppendChild(
			widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"iconlabel", "mr-sm-2"},
				},
				Text: locales.Translate("product_disposal_comment_title", HTTPHeaderAcceptLanguage),
			}))
		colDisposalComment.AppendChild(widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Text: *globals.CurrentProduct.ProductDisposalComment,
		}))
	}
	// Remark.
	colRemark := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-auto m-1"},
		},
	})
	if globals.CurrentProduct.ProductRemark != nil && *globals.CurrentProduct.ProductRemark != "" {
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
			Text: *globals.CurrentProduct.ProductRemark,
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
	if globals.CurrentProduct.ProductRadioactive {
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
	if globals.CurrentProduct.ProductRestricted {
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
	if globals.CurrentProduct.Category != nil && globals.CurrentProduct.Category.CategoryID != nil {
		detailCardRow.AppendChild(colCategory)
	}
	if len(globals.CurrentProduct.Tags) > 0 {
		detailCardRow.AppendChild(colTags)
	}
	if globals.CurrentProduct.ProductTwoDFormula != nil && *globals.CurrentProduct.ProductTwoDFormula != "" {
		detailCardRow.AppendChild(col2dimage)
	}
	if len(globals.CurrentProduct.Synonyms) > 0 {
		detailCardRow.AppendChild(colSynonym)
	}
	if globals.CurrentProduct.EmpiricalFormula != nil && globals.CurrentProduct.EmpiricalFormula.EmpiricalFormulaID != nil {
		detailCardRow.AppendChild(colEmpiricalFormula)
	}
	if globals.CurrentProduct.LinearFormula != nil && globals.CurrentProduct.LinearFormula.LinearFormulaID != nil {
		detailCardRow.AppendChild(colLinearFormula)
	}
	if globals.CurrentProduct.ProductThreeDFormula != nil && *globals.CurrentProduct.ProductThreeDFormula != "" {
		detailCardRow.AppendChild(colTreedFormula)
	}
	if globals.CurrentProduct.ProductInchi != nil {
		div := dom.GetWindow().Document().CreateElement("div")
		inchi_div := div.(*dom.HTMLDivElement)
		inchi_div.SetClass("col-sm-12 m-1")

		span_title := dom.GetWindow().Document().CreateElement("span")
		inchi_title_html := span_title.(*dom.HTMLSpanElement)
		inchi_title_html.SetInnerHTML(locales.Translate("product_inchi_title", HTTPHeaderAcceptLanguage))
		inchi_title_html.SetClass("iconlabel mr-sm-2")

		span := dom.GetWindow().Document().CreateElement("span")
		inchi_html := span.(*dom.HTMLSpanElement)
		inchi_html.SetInnerHTML(*globals.CurrentProduct.ProductInchi)

		inchi_div.AppendChild(inchi_title_html)
		inchi_div.AppendChild(inchi_html)
		detailCardRow.AppendChild(inchi_div)
	}
	if globals.CurrentProduct.ProductInchikey != nil {
		div := dom.GetWindow().Document().CreateElement("div")
		inchikey_div := div.(*dom.HTMLDivElement)
		inchikey_div.SetClass("col-sm-12 m-1")

		span_title := dom.GetWindow().Document().CreateElement("span")
		inchikey_title_html := span_title.(*dom.HTMLSpanElement)
		inchikey_title_html.SetInnerHTML(locales.Translate("product_inchi_key_title", HTTPHeaderAcceptLanguage))
		inchikey_title_html.SetClass("iconlabel mr-sm-2")

		span := dom.GetWindow().Document().CreateElement("span")
		inchikey_html := span.(*dom.HTMLSpanElement)
		inchikey_html.SetInnerHTML(*globals.CurrentProduct.ProductInchikey)

		inchikey_div.AppendChild(inchikey_title_html)
		inchikey_div.AppendChild(inchikey_html)
		detailCardRow.AppendChild(inchikey_div)
	}
	if globals.CurrentProduct.ProductCanonicalSmiles != nil {
		div := dom.GetWindow().Document().CreateElement("div")
		smiles_div := div.(*dom.HTMLDivElement)
		smiles_div.SetClass("col-sm-12 m-1")

		span_title := dom.GetWindow().Document().CreateElement("span")
		smiles_title_html := span_title.(*dom.HTMLSpanElement)
		smiles_title_html.SetInnerHTML(locales.Translate("product_smiles_title", HTTPHeaderAcceptLanguage))
		smiles_title_html.SetClass("iconlabel mr-sm-2")

		span := dom.GetWindow().Document().CreateElement("span")
		smiles_html := span.(*dom.HTMLSpanElement)
		smiles_html.SetInnerHTML(*globals.CurrentProduct.ProductCanonicalSmiles)

		smiles_div.AppendChild(smiles_title_html)
		smiles_div.AppendChild(smiles_html)
		detailCardRow.AppendChild(smiles_div)
	}
	if globals.CurrentProduct.ProductMolecularWeight != nil {
		div := dom.GetWindow().Document().CreateElement("div")
		molecularweight_div := div.(*dom.HTMLDivElement)
		molecularweight_div.SetClass("col-sm-12 m-1")

		span_title := dom.GetWindow().Document().CreateElement("span")
		molecularweight_title_html := span_title.(*dom.HTMLSpanElement)
		molecularweight_title_html.SetInnerHTML(locales.Translate("product_molecular_weight_title", HTTPHeaderAcceptLanguage))
		molecularweight_title_html.SetClass("iconlabel mr-sm-2")

		span := dom.GetWindow().Document().CreateElement("span")
		molecularweight_html := span.(*dom.HTMLSpanElement)

		if globals.CurrentProduct.UnitMolecularWeight != nil {
			molecularweight_html.SetInnerHTML(fmt.Sprintf("%f %s", *globals.CurrentProduct.ProductMolecularWeight, *globals.CurrentProduct.UnitMolecularWeight.UnitLabel))
		} else {
			molecularweight_html.SetInnerHTML(fmt.Sprintf("%f", *globals.CurrentProduct.ProductMolecularWeight))

		}
		molecularweight_div.AppendChild(molecularweight_title_html)
		molecularweight_div.AppendChild(molecularweight_html)
		detailCardRow.AppendChild(molecularweight_div)
	}

	// if globals.CurrentProduct.CasNumber.CasNumberID!= nil {
	if globals.CurrentProduct.CasNumber != nil && globals.CurrentProduct.CasNumber.CasNumberID != nil {
		detailCardRow.AppendChild(colCas)
	}
	if globals.CurrentProduct.CeNumber != nil && globals.CurrentProduct.CeNumber.CeNumberID != nil {
		detailCardRow.AppendChild(colCe)
	}
	if globals.CurrentProduct.ProductMSDS != nil && *globals.CurrentProduct.ProductMSDS != "" {
		detailCardRow.AppendChild(colMsds)
	}
	if globals.CurrentProduct.ProducerRef != nil && globals.CurrentProduct.Producer.ProducerID != nil {
		detailCardRow.AppendChild(colProducer)
	}
	if len(globals.CurrentProduct.SupplierRefs) > 0 {
		detailCardRow.AppendChild(colSuppliers)
	}
	if globals.CurrentProduct.ProductNumberPerCarton != nil && *globals.CurrentProduct.ProductNumberPerCarton > 0 {
		detailCardRow.AppendChild(colNumberPerCarton)
	}
	if globals.CurrentProduct.ProductNumberPerBag != nil && *globals.CurrentProduct.ProductNumberPerBag > 0 {
		detailCardRow.AppendChild(colNumberPerBag)
	}
	if len(globals.CurrentProduct.Symbols) > 0 {
		detailCardRow.AppendChild(colSymbols)
	}
	if globals.CurrentProduct.SignalWord != nil && globals.CurrentProduct.SignalWord.SignalWordID != nil {
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
	if globals.CurrentProduct.PhysicalState != nil && globals.CurrentProduct.PhysicalState.PhysicalStateID != nil {
		detailCardRow.AppendChild(colPhysicalState)
	}
	if globals.CurrentProduct.ProductSheet != nil && *globals.CurrentProduct.ProductSheet != "" {
		detailCardRow.AppendChild(colProducerSheet)
	}
	if globals.CurrentProduct.ProductTemperature != nil {
		detailCardRow.AppendChild(colStorageTemperature)
	}
	if globals.CurrentProduct.ProductDisposalComment != nil && *globals.CurrentProduct.ProductDisposalComment != "" {
		detailCardRow.AppendChild(colDisposalComment)
	}
	if globals.CurrentProduct.ProductRemark != nil && *globals.CurrentProduct.ProductRemark != "" {
		detailCardRow.AppendChild(colRemark)
	}
	if globals.CurrentProduct.ProductRadioactive {
		detailCardRow.AppendChild(colRadioactive)
	}
	if globals.CurrentProduct.ProductRestricted {
		detailCardRow.AppendChild(colRestricted)
	}
	detailCardRow.AppendChild(colPerson)

	if globals.CurrentProduct.ProductAvailability != nil {
		div := dom.GetWindow().Document().CreateElement("div")
		availability_div := div.(*dom.HTMLDivElement)
		availability_div.SetClass("col-sm-12 m-1")

		span_title := dom.GetWindow().Document().CreateElement("span")
		availability_title_html := span_title.(*dom.HTMLSpanElement)
		availability_title_html.SetInnerHTML(locales.Translate("product_availability_title", HTTPHeaderAcceptLanguage))
		availability_title_html.SetClass("iconlabel mr-sm-2")

		span := dom.GetWindow().Document().CreateElement("span")
		availability_html := span.(*dom.HTMLSpanElement)

		var span_content string
		for _, entity := range *globals.CurrentProduct.ProductAvailability {
			span_content = span_content + fmt.Sprintf(" %s", entity.EntityName)
			span_content = span_content + " ("
			for _, manager := range *entity.Managers {
				span_content = span_content + fmt.Sprintf("%s ", manager.PersonEmail)
			}
			span_content = span_content + ")"
		}
		availability_html.SetInnerHTML(span_content)

		availability_div.AppendChild(availability_title_html)
		availability_div.AppendChild(availability_html)
		detailCardRow.AppendChild(availability_div)
	}

	return detailCard.OuterHTML()

}

func NameFormatter(this js.Value, args []js.Value) interface{} {

	row := args[1]
	globals.CurrentProduct = Product{Product: &models.Product{}}.ProductFromJsJSONValue(row)

	specificity := ""
	if globals.CurrentProduct.ProductSpecificity != nil {
		specificity = *globals.CurrentProduct.ProductSpecificity
	}
	result := fmt.Sprintf("%s <i>%s</i>", globals.CurrentProduct.Name.NameLabel, specificity)

	for _, syn := range globals.CurrentProduct.Synonyms {
		if strings.HasPrefix(syn.NameLabel, "|") && strings.HasSuffix(syn.NameLabel, "|") {
			result += fmt.Sprintf(" %s", syn.NameLabel)
		}
	}

	result += "<div>"

	if globals.CurrentProduct.ProductSL != nil && *globals.CurrentProduct.ProductSL != "" {
		result += fmt.Sprintf("<span class='text-white badge bg-secondary'>%s</span>", *globals.CurrentProduct.ProductSL)
	}

	if globals.CurrentProduct.ProductHasBookmark {
		result += fmt.Sprintf("<span class='mdi mdi-bookmark mdi-24px iconlabel'>%s</span>", "")
	}

	result += "</div>"

	return result

}

func EmpiricalformulaFormatter(this js.Value, args []js.Value) interface{} {

	row := args[1]
	globals.CurrentProduct = Product{Product: &models.Product{}}.ProductFromJsJSONValue(row)

	if globals.CurrentProduct.EmpiricalFormula != nil {
		return *globals.CurrentProduct.EmpiricalFormulaLabel
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

	if globals.CurrentProduct.ProductTwoDFormula != nil {

		twodformula = widgets.NewImg(widgets.ImgAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Height:          "70",
			Src:             *globals.CurrentProduct.ProductTwoDFormula,
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

	if globals.CurrentProduct.CasNumber != nil {
		spanCasNumber = widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Text: *globals.CurrentProduct.CasNumberLabel,
		}).OuterHTML()
	}

	if globals.CurrentProduct.CasNumber != nil && globals.CurrentProduct.CasNumberCMR != nil {
		spanCasCMR = widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
				Classes: []string{"text-danger", "font-italic"},
			},
			Text: *globals.CurrentProduct.CasNumber.CasNumberCMR,
		}).OuterHTML()
	}

	for _, hs := range globals.CurrentProduct.HazardStatements {
		if hs.HazardStatementCMR != nil {
			spanHSCMR = widgets.NewSpan(widgets.SpanAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
					Classes: []string{"text-danger", "font-italic"},
				},
				Text: *hs.HazardStatementCMR,
			}).OuterHTML()
		}
	}

	for _, symbol := range globals.CurrentProduct.Symbols {
		if symbol.SymbolLabel == "GHS02" {
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

	if globals.CurrentProduct.ProductRestricted {
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

	if globals.CurrentProduct.ProductSpecificity != nil {
		return *globals.CurrentProduct.ProductSpecificity
	} else {
		return ""
	}

}

func Product_productSlFormatter(this js.Value, args []js.Value) interface{} {

	row := args[1]
	globals.CurrentProduct = Product{Product: &models.Product{}}.ProductFromJsJSONValue(row)

	if globals.CurrentProduct.ProductSL != nil {
		return *globals.CurrentProduct.ProductSL
	} else {
		return ""
	}

}

func OperateFormatter(this js.Value, args []js.Value) interface{} {

	var (
		buttonStorages   string
		buttonTotalStock string
		iconBookmark     themes.IconFace
		textBookmark     string
	)

	row := args[1]
	globals.CurrentProduct = Product{Product: &models.Product{}}.ProductFromJsJSONValue(row)

	if globals.CurrentProduct.ProductHasBookmark {
		iconBookmark = themes.MDI_BOOKMARK
		textBookmark = locales.Translate("unbookmark", HTTPHeaderAcceptLanguage)
	} else {
		iconBookmark = themes.MDI_NO_BOOKMARK
		textBookmark = locales.Translate("bookmark", HTTPHeaderAcceptLanguage)
	}

	if globals.CurrentProduct.ProductSC != nil || globals.CurrentProduct.ProductASC != nil {

		var _product_sc = 0
		var _product_asc = 0

		if globals.CurrentProduct.ProductSC != nil {
			_product_sc = *globals.CurrentProduct.ProductSC
		}
		if globals.CurrentProduct.ProductASC != nil {
			_product_asc = *globals.CurrentProduct.ProductASC
		}

		buttonStorages = widgets.NewLink(
			widgets.LinkAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Id:      "storages" + strconv.Itoa(int(*globals.CurrentProduct.ProductID)),
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
						Text: fmt.Sprintf("%s %d  (%d)", locales.Translate("storages", HTTPHeaderAcceptLanguage), _product_sc, _product_asc),
					},
				),
			},
		).OuterHTML()

	}

	buttonStore := widgets.NewLink(
		widgets.LinkAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Id:      "store" + strconv.Itoa(int(*globals.CurrentProduct.ProductID)),
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
				Id:      "edit" + strconv.Itoa(int(*globals.CurrentProduct.ProductID)),
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
				Id:      "delete" + strconv.Itoa(int(*globals.CurrentProduct.ProductID)),
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
				Id:         "bookmark" + strconv.Itoa(int(*globals.CurrentProduct.ProductID)),
				Classes:    []string{"bookmark", "dropdown-item", "text-primary", "iconlabel"},
				Visible:    false,
				Attributes: map[string]string{"pid": strconv.Itoa(int(*globals.CurrentProduct.ProductID))},
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

	buttonSelect := widgets.NewLink(
		widgets.LinkAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Id:         "select" + strconv.Itoa(int(*globals.CurrentProduct.ProductID)),
				Classes:    []string{"productselect", "dropdown-item", "text-primary", "iconlabel"},
				Visible:    false,
				Attributes: map[string]string{"pid": strconv.Itoa(int(*globals.CurrentProduct.ProductID))},
			},
			Href: "#",
			Label: widgets.NewSpan(
				widgets.SpanAttributes{
					BaseAttributes: widgets.BaseAttributes{
						Classes: []string{"mdi", themes.MDI_PUBCHEM.ToString()},
						Visible: true,
					},
					Text: locales.Translate("select", HTTPHeaderAcceptLanguage),
				},
			),
		},
	).OuterHTML()

	if globals.CurrentProduct.ProductSC != nil {

		buttonTotalStock = widgets.NewLink(
			widgets.LinkAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Id:      "totalstock" + strconv.Itoa(int(*globals.CurrentProduct.ProductID)),
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

	totalstockDiv := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Id:      "totalstock-collapse-" + strconv.Itoa(int(*globals.CurrentProduct.ProductID)),
			Classes: []string{"collapse", "float-left"},
		},
	}).OuterHTML()

	finalDiv := `
<div class="dropdown">
  <button class="btn btn-secondary dropdown-toggle" type="button" id="productActions` + strconv.Itoa(int(*globals.CurrentProduct.ProductID)) + `" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
    <span class="mdi mdi-menu">&nbsp;</span>
  </button>
  <div class="dropdown-menu" aria-labelledby="productActions` + strconv.Itoa(int(*globals.CurrentProduct.ProductID)) + `">
  ` + buttonStorages + buttonStore + buttonEdit + buttonDelete + buttonBookmark + buttonTotalStock + buttonSelect +
		`
  </div>
</div>
<div id="confirm` + strconv.Itoa(int(*globals.CurrentProduct.ProductID)) + `">
</div>
`
	return finalDiv + totalstockDiv

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
			ProducerLabel: &producerLabel,
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
			SupplierLabel: &supplierLabel,
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

	validate.NewValidate(jquery.Jq("select#producer_ref"), nil).ValidateAddRequired()
	validate.NewValidate(jquery.Jq("input#storage_batch_number"), nil).ValidateRemoveRequired()

	jquery.Jq("span#producerref.badge").Show()
	jquery.Jq("span#storage_batch_number.badge").Hide()

	validate.NewValidate(jquery.Jq("select#empirical_formula"), nil).ValidateRemoveRequired()
	validate.NewValidate(jquery.Jq("select#cas_number"), nil).ValidateRemoveRequired()

	jquery.Jq("span#empirical_formula.badge").Hide()
	jquery.Jq("span#cas_number.badge").Hide()

	jquery.Jq("input#showconsu").SetProp("checked", "checked")

}

func Chemify() {

	jquery.Jq(".bio").Not(".chem").Hide()
	jquery.Jq(".consu").Not(".chem").Hide()
	jquery.Jq(".chem").Show()

	validate.NewValidate(jquery.Jq("select#producer_ref"), nil).ValidateRemoveRequired()
	validate.NewValidate(jquery.Jq("input#storage_batch_number"), nil).ValidateRemoveRequired()

	jquery.Jq("span#producerref.badge").Hide()
	jquery.Jq("span#storage_batch_number.badge").Hide()

	validate.NewValidate(jquery.Jq("select#empirical_formula"), nil).ValidateAddRequired()
	validate.NewValidate(jquery.Jq("select#cas_number"), nil).ValidateAddRequired()

	jquery.Jq("span#empirical_formula.badge").Show()
	jquery.Jq("span#cas_number.badge").Show()

	jquery.Jq("input#showchem").SetProp("checked", "checked")

}

func Biofy() {

	jquery.Jq(".chem").Not(".bio").Hide()
	jquery.Jq(".consu").Not(".bio").Hide()
	jquery.Jq(".bio").Show()

	validate.NewValidate(jquery.Jq("select#producer_ref"), nil).ValidateAddRequired()
	validate.NewValidate(jquery.Jq("select#producer"), nil).ValidateAddRequired()
	validate.NewValidate(jquery.Jq("input#storage_batch_number"), nil).ValidateAddRequired()

	jquery.Jq("span#producerref.badge").Show()
	jquery.Jq("span#storage_batch_number.badge").Show()

	validate.NewValidate(jquery.Jq("select#empirical_formula"), nil).ValidateRemoveRequired()
	validate.NewValidate(jquery.Jq("select#cas_number"), nil).ValidateRemoveRequired()

	jquery.Jq("span#empirical_formula.badge").Hide()
	jquery.Jq("span#cas_number.badge").Hide()

	jquery.Jq("input#showbio").SetProp("checked", "checked")

}

func ShowIfAuthorizedActionButtons(this js.Value, args []js.Value) interface{} {

	jquery.Jq(".bookmark").FadeIn()
	jsutils.HasPermission("storages", "", "post", func() {
		jquery.Jq(".store").FadeIn()
	}, func() {
	})
	jsutils.HasPermission("storages", "", "get", func() {
		jquery.Jq(".storages").FadeIn()
		jquery.Jq(".ostorages").FadeIn()
		jquery.Jq(".totalstock").FadeIn()

		jquery.Jq("#switchview").SetVisible()
	}, func() {
	})
	jsutils.HasPermission("products", "", "post", func() {
		jquery.Jq(".productselect").FadeIn()
		jquery.Jq(".productedit").FadeIn()
	}, func() {
	})

	// Iterating other the button with the class "store_location"
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
