package product

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"syscall/js"

	. "github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/locales"
	. "github.com/tbellembois/gochimitheque-wasm/types"
	"github.com/tbellembois/gochimitheque-wasm/utils"
	"github.com/tbellembois/gochimitheque-wasm/views/storage"
	"github.com/tbellembois/gochimitheque-wasm/widgets"
	"github.com/tbellembois/gochimitheque-wasm/widgets/themes"
	"honnef.co/go/js/dom/v2"
)

func OperateEventsBookmark(this js.Value, args []js.Value) interface{} {

	row := args[2]
	product := Product{}.ProductFromJsJSONValue(row)

	url := fmt.Sprintf("%sbookmarks/%d", ApplicationProxyPath, product.ProductID)
	method := "put"

	done := func(data js.Value) {

		var (
			product Product
			err     error
		)

		if err = json.Unmarshal([]byte(data.String()), &product); err != nil {
			fmt.Println(err)
		}

		if product.Bookmark != nil {
			Jq(fmt.Sprintf("button#bookmark%d > span", product.ProductID)).RemoveClass(themes.MDI_NO_BOOKMARK.ToString())
			Jq(fmt.Sprintf("button#bookmark%d > span", product.ProductID)).AddClass(themes.MDI_BOOKMARK.ToString())
		} else {
			Jq(fmt.Sprintf("button#bookmark%d > span", product.ProductID)).AddClass(themes.MDI_NO_BOOKMARK.ToString())
			Jq(fmt.Sprintf("button#bookmark%d > span", product.ProductID)).RemoveClass(themes.MDI_BOOKMARK.ToString())
		}

	}
	fail := func(data js.Value) {

		utils.DisplayGenericErrorMessage()

	}

	Ajax{
		Method: method,
		URL:    url,
		Done:   done,
		Fail:   fail,
	}.Send()

	return nil

}

func OperateEventsStore(this js.Value, args []js.Value) interface{} {

	row := args[2]
	product := Product{}.FromJsJSONValue(row)

	href := fmt.Sprintf("%svc/storages", ApplicationProxyPath)
	utils.LoadContent("storage", href, storage.Storage_createCallback, product)

	return nil

}

func OperateEventsStorages(this js.Value, args []js.Value) interface{} {

	// TODO: stock

	storageCallbackWrapper := func(args ...interface{}) {
		storage.Storage_listCallback(js.Null(), nil)
	}

	row := args[2]
	product := Product{}.ProductFromJsJSONValue(row)

	BSTableQueryFilter.Lock()
	BSTableQueryFilter.QueryFilter.Product = strconv.Itoa(product.ProductID)
	BSTableQueryFilter.QueryFilter.ProductFilterLabel = fmt.Sprintf("%s %s", product.Name.NameLabel, product.ProductSpecificity.String)

	href := fmt.Sprintf("%sv/storages", ApplicationProxyPath)
	utils.LoadContent("storage", href, storageCallbackWrapper, product)

	return nil

}

func OperateEventsOStorages(this js.Value, args []js.Value) interface{} {

	row := args[2]
	product := Product{}.ProductFromJsJSONValue(row)

	url := fmt.Sprintf("%sstorages/others?product=%d", ApplicationProxyPath, product.ProductID)
	method := "get"

	done := func(data js.Value) {

		var (
			entities Entities
			err      error
		)

		if err = json.Unmarshal([]byte(data.String()), &entities); err != nil {
			utils.DisplayGenericErrorMessage()
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

			Jq(fmt.Sprintf("#ostorages-collapse-%d", product.ProductID)).Append(divEntity.OuterHTML())

		}

		Jq(fmt.Sprintf("#ostorages-collapse-%d", product.ProductID)).Show()

	}
	fail := func(data js.Value) {

		utils.DisplayGenericErrorMessage()

	}

	Ajax{
		Method: method,
		URL:    url,
		Done:   done,
		Fail:   fail,
	}.Send()

	return nil

}

func OperateEventsEdit(this js.Value, args []js.Value) interface{} {

	row := args[2]
	index := args[3].Int()
	product := Product{}.ProductFromJsJSONValue(row)

	url := fmt.Sprintf("%sproducts/%d", ApplicationProxyPath, product.ProductID)
	method := "get"

	done := func(data js.Value) {

		var (
			product Product
			err     error
		)

		if err = json.Unmarshal([]byte(data.String()), &product); err != nil {
			fmt.Println(err)
		}

		FillInProductForm(product, "edit-collapse")

		Jq("input#index").SetVal(index)

		Jq("#edit-collapse").Show()
		Jq("#search").Hide()
		Jq("#actions").Hide()

	}
	fail := func(data js.Value) {

		utils.DisplayGenericErrorMessage()

	}

	Ajax{
		Method: method,
		URL:    url,
		Done:   done,
		Fail:   fail,
	}.Send()

	return nil

}

func OperateEventsDelete(this js.Value, args []js.Value) interface{} {

	row := args[2]
	product := Product{}.ProductFromJsJSONValue(row)

	url := fmt.Sprintf("%sproducts/%d", ApplicationProxyPath, product.ProductID)
	method := "delete"

	done := func(data js.Value) {

		utils.DisplaySuccessMessage(locales.Translate("product_deleted_message", HTTPHeaderAcceptLanguage))
		Jq("#Product_table").Bootstraptable(nil).Refresh(nil)

	}
	fail := func(data js.Value) {

		utils.DisplayGenericErrorMessage()

	}

	Ajax{
		Method: method,
		URL:    url,
		Done:   done,
		Fail:   fail,
	}.Send()

	return nil

}

// Get remote bootstrap table data - params object defined here:
// https://bootstrap-table.com/docs/api/table-options/#queryparams
func GetTableData(this js.Value, args []js.Value) interface{} {

	row := args[0]
	params := QueryParamsFromJsJSONValue(row)

	go func() {

		u := url.URL{Path: ApplicationProxyPath + "products"}
		u.RawQuery = params.Data.ToRawQuery()

		ajax := Ajax{
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
					utils.DisplaySuccessMessage(locales.Translate("export_done", HTTPHeaderAcceptLanguage))

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
						Onclick: "$('#export').hide()",
						Title:   locales.Translate("download_export", HTTPHeaderAcceptLanguage),
						Href:    fmt.Sprintf("%sdownload/%s", ApplicationProxyPath, products.GetExportFn()),
						Label:   icon,
					})

					Jq("#export-body").SetHtml("")
					Jq("#export-body").Append(downloadLink.OuterHTML())
					Jq("#export").Show()

				} else if products.GetTotal() != 0 {
					row.Call("success", js.ValueOf(js.Global().Get("JSON").Call("parse", data)))
				} else {
					// FIXME: Does not work
					//Jq("#Product_table").Bootstraptable(nil).RemoveAll()
					utils.DisplayErrorMessage(locales.Translate("no_result", HTTPHeaderAcceptLanguage))
				}

			},
			Fail: func(jqXHR js.Value) {

				utils.DisplayGenericErrorMessage()

			},
		}

		ajax.Send()

	}()

	return nil

}

// TODO: factorise me with storage
func DataQueryParams(this js.Value, args []js.Value) interface{} {

	params := args[0]

	queryFilter := QueryFilterFromJsJSONValue(params)

	// Product_SaveCallback product id.
	queryFilter.Product = BSTableQueryFilter.Product
	queryFilter.ProductFilterLabel = BSTableQueryFilter.ProductFilterLabel
	queryFilter.ProductBookmark = BSTableQueryFilter.ProductBookmark
	queryFilter.Export = BSTableQueryFilter.Export
	BSTableQueryFilter.Export = false
	BSTableQueryFilter.Unlock()

	if Jq("select#s_storelocation").Select2IsInitialized() {
		i := Jq("select#s_storelocation").Select2Data()
		if len(i) > 0 {
			queryFilter.StoreLocation = i[0].Id
			queryFilter.StoreLocationFilterLabel = i[0].Text
		}
	}
	if Jq("select#s_name").Select2IsInitialized() {
		i := Jq("select#s_name").Select2Data()
		if len(i) > 0 {
			queryFilter.Name = i[0].Id
			queryFilter.NameFilterLabel = i[0].Text
		}
	}
	if Jq("select#s_casnumber").Select2IsInitialized() {
		i := Jq("select#s_casnumber").Select2Data()
		if len(i) > 0 {
			queryFilter.CasNumber = i[0].Id
			queryFilter.CasNumberFilterLabel = i[0].Text
		}
	}
	if Jq("select#s_empiricalformula").Select2IsInitialized() {
		i := Jq("select#s_empiricalformula").Select2Data()
		if len(i) > 0 {
			queryFilter.EmpiricalFormula = i[0].Id
			queryFilter.EmpiricalFormulaFilterLabel = i[0].Text
		}
	}
	if Jq("select#s_signalword").Select2IsInitialized() {
		i := Jq("select#s_signalword").Select2Data()
		if len(i) > 0 {
			queryFilter.SignalWord = i[0].Id
		}
	}
	if Jq("select#s_hazardstatements").Select2IsInitialized() {
		i := Jq("select#s_hazardstatements").Select2Data()
		if len(i) > 0 {
			for _, hs := range i {
				queryFilter.HazardStatements = append(queryFilter.HazardStatements, hs.Id)
			}
		}
	}
	if Jq("select#s_precautionarystatements").Select2IsInitialized() {
		i := Jq("select#s_precautionarystatements").Select2Data()
		if len(i) > 0 {
			for _, ps := range i {
				queryFilter.PrecautionaryStatements = append(queryFilter.PrecautionaryStatements, ps.Id)
			}
		}
	}
	if Jq("select#s_symbols").Select2IsInitialized() {
		i := Jq("select#s_symbols").Select2Data()
		if len(i) > 0 {
			for _, s := range i {
				queryFilter.Symbols = append(queryFilter.Symbols, s.Id)
			}
		}
	}

	if Jq("#s_storage_barecode").GetVal().Truthy() {
		queryFilter.StorageBarecode = Jq("#s_storage_barecode").GetVal().String()
	}
	if Jq("#s_custom_name_part_of").GetVal().Truthy() {
		queryFilter.CustomNamePartOf = Jq("#s_custom_name_part_of").GetVal().String()
	}
	if Jq("#s_casnumber_cmr:checked").Object.Length() > 0 {
		queryFilter.CasNumberCMR = true
	}
	if Jq("#s_borrowing:checked").Object.Length() > 0 {
		queryFilter.Borrowing = true
	}
	if Jq("#s_storage_to_destroy:checked").Object.Length() > 0 {
		queryFilter.StorageToDestroy = true
	}

	queryFilter.DisplayFilter()

	return queryFilter.ToJsValue()

}

func DetailFormatter(this js.Value, args []js.Value) interface{} {

	row := args[1]
	product := Product{}.ProductFromJsJSONValue(row)

	var synonyms strings.Builder
	for _, synonym := range product.Synonyms {
		synonyms.WriteString(synonym.NameLabel)
		synonyms.WriteString("<br/>")
	}

	// JSMol div.
	widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Id:      fmt.Sprintf("jsmol%d", product.ProductID),
		},
	})

	//
	// Synonyms, ID and person.
	//
	rowSynonymAndID := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"row", "mt-sm-3"},
		},
	})
	// Synonym.
	colSynonym := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-8"},
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
			Text: strconv.Itoa(product.ProductID),
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
			Text: product.Person.PersonEmail,
		}))

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
	if product.Category.CategoryID.Valid {
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
			Text: product.Category.CategoryLabel.String,
		}))
	}
	// Tags.
	colTags := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-2"},
		},
	})
	for _, tag := range product.Tags {
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
	// Producer
	colProducer := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-6"},
		},
	})
	if product.ProducerRef.ProducerRefID.Valid {
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
			Text: fmt.Sprintf("%s: %s", product.Producer.ProducerLabel.String, product.ProducerRef.ProducerRefLabel.String),
		}))
	}
	// Suppliers
	colSuppliers := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-6"},
		},
	})
	if len(product.SupplierRefs) > 0 {
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
		for _, s := range product.SupplierRefs {
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
	if product.CasNumber.CasNumberID.Valid {
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
			Text: product.CasNumber.CasNumberLabel.String,
		}))
	}
	if product.CasNumber.CasNumberCMR.Valid {
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
			Text: product.CasNumber.CasNumberCMR.String,
		}))
	}
	for _, hs := range product.HazardStatements {
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
	if product.CeNumber.CeNumberID.Valid {
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
			Text: product.CeNumber.CeNumberLabel.String,
		}))
	}
	// MSDS.
	colMsds := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-4"},
		},
	})
	if product.ProductMSDS.Valid {
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
			Href:  product.ProductMSDS.String,
			Label: icon,
		}))
	}

	rowCasCeMsds.AppendChild(colCas)
	rowCasCeMsds.AppendChild(colCe)
	rowCasCeMsds.AppendChild(colMsds)

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
	if product.EmpiricalFormula.EmpiricalFormulaID.Valid {
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
			Text: product.EmpiricalFormula.EmpiricalFormulaLabel.String,
		}))
	}
	// Linear formula.
	colLinearFormula := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-4"},
		},
	})
	if product.LinearFormula.LinearFormulaID.Valid {
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
			Text: product.LinearFormula.LinearFormulaLabel.String,
		}))
	}
	// 3D formula.
	colTreedFormula := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-4"},
		},
	})
	if product.ProductThreeDFormula.Valid && product.ProductThreeDFormula.String != "" {
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
			Href:  product.ProductThreeDFormula.String,
			Label: icon,
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
	for _, symbol := range product.Symbols {
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
	if product.SignalWord.SignalWordID.Valid {
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
			Text: product.SignalWord.SignalWordLabel.String,
		}))
	}
	// Physical state.
	colPhysicalState := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-4"},
		},
	})
	if product.PhysicalState.PhysicalStateID.Valid {
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
			Text: product.PhysicalState.PhysicalStateLabel.String,
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
	if len(product.HazardStatements) > 0 {
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
		for _, s := range product.HazardStatements {
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
	if len(product.PrecautionaryStatements) > 0 {
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
		for _, s := range product.PrecautionaryStatements {
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
	if len(product.ClassOfCompound) > 0 {
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
		for _, s := range product.ClassOfCompound {
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
	if product.ProductDisposalComment.Valid && product.ProductDisposalComment.String != "" {
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
			Text: product.ProductDisposalComment.String,
		}))
	}
	// Remark.
	colRemark := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"col-sm-6"},
		},
	})
	if product.ProductRemark.Valid && product.ProductRemark.String != "" {
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
			Text: product.ProductRemark.String,
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
	if product.ProductRadioactive.Valid && product.ProductRadioactive.Bool {
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
	if product.ProductRestricted.Valid && product.ProductRestricted.Bool {
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
		rowCasCeMsds.OuterHTML() +
		rowFormulas.OuterHTML() +
		rowSymbolsSignalWordPhysicalState.OuterHTML() +
		rowHsPsCoc.OuterHTML() +
		rowDisposalCommentRemark.OuterHTML() +
		rowRadioactiveRestricted.OuterHTML()

}

func EmpiricalformulaFormatter(this js.Value, args []js.Value) interface{} {

	row := args[1]
	product := Product{}.ProductFromJsJSONValue(row)

	if product.EmpiricalFormulaID.Valid {
		return product.EmpiricalFormulaLabel.String
	} else {
		return ""
	}

}

func CasnumberFormatter(this js.Value, args []js.Value) interface{} {

	row := args[1]
	product := Product{}.ProductFromJsJSONValue(row)

	if product.CasNumberID.Valid {
		return product.CasNumberLabel.String
	} else {
		return ""
	}

}

func Product_productSpecificityFormatter(this js.Value, args []js.Value) interface{} {

	row := args[1]
	product := Product{}.ProductFromJsJSONValue(row)

	if product.ProductSpecificity.Valid {
		return product.ProductSpecificity.String
	} else {
		return ""
	}

}

func Product_productSlFormatter(this js.Value, args []js.Value) interface{} {

	row := args[1]
	product := Product{}.ProductFromJsJSONValue(row)

	if product.ProductSL.Valid {
		return product.ProductSL.String
	} else {
		return ""
	}

}

func OperateFormatter(this js.Value, args []js.Value) interface{} {

	var (
		imgSGH02        string
		spanCasCMR      string
		spanHSCMR       string
		iconRestricted  string
		buttonStorages  string
		buttonOStorages string
		iconBookmark    themes.IconFace
	)

	row := args[1]
	product := Product{}.ProductFromJsJSONValue(row)

	if product.Bookmark.BookmarkID.Valid {
		iconBookmark = themes.MDI_BOOKMARK
	} else {
		iconBookmark = themes.MDI_NO_BOOKMARK
	}

	for _, symbol := range product.Symbols {
		if symbol.SymbolLabel == "SGH02" {
			imgSGH02 = widgets.NewImg(widgets.ImgAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Visible: true,
				},
				Src:   "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACYAAAAmCAYAAACoPemuAAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAAN1wAADdcBQiibeAAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAAIvSURBVFiFzdgxbI5BGMDx36uNJsJApFtFIh2QEIkBtYiYJFKsrDaJqYOofhKRMFhsFhNNMBgkTAaD0KUiBomN1EpYBHWGnvi8vu9r7/2eti65vHfv3T3P/57nuXvfuyqlJCRVVQuk1AqRl1LqP9NKpJxbETKjocLgYi0VaLl49wXBxUIFwsVDBcGFQWEAg1FwYZbCGMajLBfmPkzgUZRbw2IKFzGPrRFw/bpvD/bn8jUkXM719f3A9eu+k3iXA/92Bnub2yYx1NgDfbrvXIYZx8dcThjBExxvPOmGltqLIzmuEt63QSVczc+z/2whSw2ThpbajS+4UgOq59O4gYFSuGaByWb8zKvwN8RXXKiBPc7PLaWx3ARqY37O1CBe5/cvO1huVy+ZnfSX7y9MYxRTNeX32lZj+/sXWNfVnV3g1tT/aJeQ5vAGp3L9eXbjTFv7NzzM9VncSSnNF2lp4MqjNYvcxwEcy+0HcQg32/q8Kndl+YrcgM9Z4YdsrZ21PtvxHT9yv1vNgr8cbiIrnMUmbKu177PwVZjLgKPNt4sCOKzF0ww32aF9CA+yxSZKoTqDlVnucI6lMxhpg76OuxhrKr8oIENyXx/xxQKTE/hUkIdLJ1tlRd3TwtF/KtcuSalVVdUwdvQe+Fd6ljhfl9NzRKT5I8cvq/B+xi3vzFfk+FaqbEUPvEtVuipXBIspX9VLlW4Q/8U1VGe4EKgYsED3tefBgt271y7dUlV/ygHpF8bRglXiwx7BAAAAAElFTkSuQmCC",
				Alt:   locales.Translate("product_flammable", HTTPHeaderAcceptLanguage),
				Title: locales.Translate("product_flammable", HTTPHeaderAcceptLanguage),
			}).OuterHTML()
		}
		break
	}

	if product.CasNumberCMR.Valid {
		spanCasCMR = widgets.NewSpan(widgets.SpanAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
				Classes: []string{"text-danger", "font-italic"},
			},
			Text: product.CasNumber.CasNumberCMR.String,
		}).OuterHTML()
	}

	for _, hs := range product.HazardStatements {
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

	if product.ProductSC != 0 || product.ProductASC != 0 {
		buttonStorages = widgets.NewBSButtonWithIcon(
			widgets.ButtonAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Id:      "storages" + strconv.Itoa(product.ProductID),
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
				Text: fmt.Sprintf("%s %d  (%d)", locales.Translate("storages", HTTPHeaderAcceptLanguage), product.ProductSC, product.ProductASC),
				Icon: themes.NewMdiIcon(themes.MDI_STORAGE, ""),
			},
			[]themes.BSClass{themes.BS_BTN, themes.BS_BNT_LINK},
		).OuterHTML()
	}

	if product.ProductTSC != 0 {
		buttonOStorages = widgets.NewBSButtonWithIcon(
			widgets.ButtonAttributes{
				BaseAttributes: widgets.BaseAttributes{
					Id:      "ostorages" + strconv.Itoa(product.ProductID),
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
				Id:      "store" + strconv.Itoa(product.ProductID),
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
				Id:      "edit" + strconv.Itoa(product.ProductID),
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
				Id:      "delete" + strconv.Itoa(product.ProductID),
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
				Id:         "bookmark" + strconv.Itoa(product.ProductID),
				Classes:    []string{"bookmark"},
				Visible:    false,
				Attributes: map[string]string{"pid": strconv.Itoa(product.ProductID)},
			},
			Title: locales.Translate("bookmark", HTTPHeaderAcceptLanguage),
		},
		widgets.IconAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
				Classes: []string{"iconlabel"},
			},
			Text: locales.Translate("bookmark", HTTPHeaderAcceptLanguage),
			Icon: themes.NewMdiIcon(iconBookmark, ""),
		},
		[]themes.BSClass{themes.BS_BTN, themes.BS_BNT_LINK},
	).OuterHTML()

	collapseDiv := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Id:      "ostorages-collapse-" + strconv.Itoa(product.ProductID),
			Classes: []string{"collapse"},
		},
	}).OuterHTML()

	if product.ProductRestricted.Valid && product.ProductRestricted.Bool {
		iconRestricted = widgets.NewIcon(widgets.IconAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
			},
			Title: locales.Translate("restricted", HTTPHeaderAcceptLanguage),
			Icon:  themes.NewMdiIcon(themes.MDI_RESTRICTED, ""),
		}).OuterHTML()
	}

	return spanCasCMR + spanHSCMR + imgSGH02 + buttonStorages + buttonOStorages + buttonStore + buttonEdit + buttonDelete + buttonBookmark + collapseDiv + iconRestricted

}

func AddProducer(this js.Value, args []js.Value) interface{} {

	var (
		producer  Producer
		dataBytes []byte
		err       error
	)

	producerLabel := Jq("input#addproducer").GetVal().String()
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

	Ajax{
		URL:    ApplicationProxyPath + "products/producers",
		Method: "post",
		Data:   dataBytes,
		Done: func(data js.Value) {

			utils.DisplaySuccessMessage(locales.Translate("producer_added", HTTPHeaderAcceptLanguage))
			Jq("input#addproducer").SetVal("")

		},
		Fail: func(jqXHR js.Value) {

			utils.DisplayGenericErrorMessage()

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

	supplierLabel := Jq("input#addsupplier").GetVal().String()
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

	Ajax{
		URL:    ApplicationProxyPath + "products/suppliers",
		Method: "post",
		Data:   dataBytes,
		Done: func(data js.Value) {

			utils.DisplaySuccessMessage(locales.Translate("supplier_added", HTTPHeaderAcceptLanguage))
			Jq("input#addsupplier").SetVal("")

		},
		Fail: func(jqXHR js.Value) {

			utils.DisplayGenericErrorMessage()

		},
	}.Send()

	return nil

}

func Chemify() {

	Jq(".chem").Show()
	Jq(".bio").Hide()

	Jq("select#producerref").ValidateRemoveRequired()
	Jq("input#product_batchnumber").ValidateRemoveRequired()
	Jq("span#producerref.badge").Hide()
	Jq("span#product_batchnumber.badge").Hide()

	Jq("select#empiricalformula").ValidateAddRequired()
	Jq("select#casnumber").ValidateAddRequired()
	Jq("span#empiricalformula.badge").Show()
	Jq("span#casnumber.badge").Show()

	Jq("input#showchem").SetProp("checked", "checked")

}

func Biofy() {

	Jq(".bio").Show()
	Jq(".chem").Hide()

	Jq("select#producerref").ValidateAddRequired()
	Jq("input#product_batchnumber").ValidateAddRequired()
	Jq("span#producerref.badge").Show()
	Jq("span#product_batchnumber.badge").Show()

	Jq("select#empiricalformula").ValidateRemoveRequired()
	Jq("select#casnumber").ValidateRemoveRequired()
	Jq("span#empiricalformula.badge").Hide()
	Jq("span#casnumber.badge").Hide()

	Jq("input#showbio").SetProp("checked", "checked")

}

func ShowIfAuthorizedActionButtons(this js.Value, args []js.Value) interface{} {

	Jq(".bookmark").FadeIn()
	utils.HasPermission("storages", "", "post", func() {
		Jq(".store").FadeIn()
	}, func() {
	})
	utils.HasPermission("storages", "-2", "get", func() {
		Jq(".storages").FadeIn()
		Jq(".ostorages").FadeIn()

		Jq("#switchview").RemoveClass("invisible")
	}, func() {
	})
	utils.HasPermission("products", "-2", "put", func() {
		Jq(".productedit").FadeIn()
	}, func() {
	})

	// Iterating other the button with the class "storelocation"
	// (we could choose "members" or "delete")
	// to retrieve once the product id.
	buttons := dom.GetWindow().Document().GetElementsByTagName("button")
	for _, button := range buttons {
		if button.Class().Contains("bookmark") {
			productId := button.GetAttribute("pid")

			utils.HasPermission("products", productId, "delete", func() {
				Jq("#delete" + productId).FadeIn()
			}, func() {
			})
		}
	}

	return nil

}
