package product

import (
	"fmt"
	"strconv"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/ajax"
	. "github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/jquery"
	"github.com/tbellembois/gochimitheque-wasm/jsutils"
	"github.com/tbellembois/gochimitheque-wasm/locales"
	"github.com/tbellembois/gochimitheque-wasm/select2"
	. "github.com/tbellembois/gochimitheque-wasm/types"
	"github.com/tbellembois/gochimitheque-wasm/widgets"
	"github.com/tbellembois/gochimitheque-wasm/widgets/themes"
	"github.com/tbellembois/gochimitheque/models"
	"honnef.co/go/js/dom/v2"
)

// TODO: factorise with storage
func Select2StoreLocationTemplateResults(this js.Value, args []js.Value) interface{} {

	var (
		iconCanStore dom.Node
	)

	data := args[0]

	storelocation := StoreLocation{StoreLocation: &models.StoreLocation{}}.FromJsJSONValue(data).(StoreLocation)

	if storelocation.StoreLocation == nil {
		return jsutils.CreateJsHTMLElementFromString(widgets.NewDiv(widgets.DivAttributes{}).OuterHTML())
	}

	iconColor := widgets.NewIcon(widgets.IconAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Attributes: map[string]string{
				"style": fmt.Sprintf("color: %s", storelocation.StoreLocationColor.String),
			},
		},
		Icon: themes.NewMdiIcon(themes.MDI_COLOR, themes.MDI_24PX),
	})

	spanLabel := widgets.NewSpan(widgets.SpanAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
		},
		Text: storelocation.StoreLocationFullPath,
	})

	if storelocation.StoreLocationCanStore.Valid && storelocation.StoreLocationCanStore.Bool {
		iconCanStore = widgets.NewIcon(widgets.IconAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
				Attributes: map[string]string{
					"style": "float: right",
				},
			},
			Icon: themes.NewMdiIcon(themes.MDI_CHECK, themes.MDI_24PX),
		})
	} else {
		iconCanStore = widgets.NewIcon(widgets.IconAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible: true,
				Attributes: map[string]string{
					"style": "float: right",
				},
			},
			Icon: themes.NewMdiIcon(themes.MDI_NO_CHECK, themes.MDI_24PX),
		})
	}

	d := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
		},
	})
	d.AppendChild(iconColor)
	d.AppendChild(spanLabel)
	d.AppendChild(iconCanStore)

	return jsutils.CreateJsHTMLElementFromString(d.OuterHTML())

}

// TODO: factorise with storage
func Select2SymbolTemplateResults(this js.Value, args []js.Value) interface{} {

	data := args[0]
	symbol := Symbol{Symbol: &models.Symbol{}}.FromJsJSONValue(data).(Symbol)

	if symbol.Symbol == nil {
		return jsutils.CreateJsHTMLElementFromString(widgets.NewDiv(widgets.DivAttributes{}).OuterHTML())
	}

	image := widgets.NewImg(widgets.ImgAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
		},
		Src:   fmt.Sprintf("data:%s", symbol.SymbolImage),
		Alt:   symbol.SymbolLabel,
		Title: symbol.SymbolLabel,
	})
	spanLabel := widgets.NewSpan(widgets.SpanAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
		},
		Text: symbol.SymbolLabel,
	})
	d := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
		},
	})

	d.AppendChild(image)
	d.AppendChild(spanLabel)

	return jsutils.CreateJsHTMLElementFromString(d.OuterHTML())

}

func Select2HazardStatementTemplateResults(this js.Value, args []js.Value) interface{} {

	data := args[0]
	hs := HazardStatement{HazardStatement: &models.HazardStatement{}}.FromJsJSONValue(data).(HazardStatement)

	if hs.HazardStatement == nil {
		return jsutils.CreateJsHTMLElementFromString(widgets.NewDiv(widgets.DivAttributes{}).OuterHTML())
	}

	spanLabel := widgets.NewSpan(widgets.SpanAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"iconlabel"},
		},
		Text: hs.HazardStatementLabel,
	})
	spanReference := widgets.NewSpan(widgets.SpanAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
		},
		Text: hs.HazardStatementReference,
	})
	d := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
		},
	})

	d.AppendChild(spanReference)
	d.AppendChild(spanLabel)

	return jsutils.CreateJsHTMLElementFromString(d.OuterHTML())

}

func Select2PrecautionaryStatementTemplateResults(this js.Value, args []js.Value) interface{} {

	data := args[0]
	ps := PrecautionaryStatement{PrecautionaryStatement: &models.PrecautionaryStatement{}}.FromJsJSONValue(data).(PrecautionaryStatement)

	if ps.PrecautionaryStatement == nil {
		return jsutils.CreateJsHTMLElementFromString(widgets.NewDiv(widgets.DivAttributes{}).OuterHTML())
	}

	spanLabel := widgets.NewSpan(widgets.SpanAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
			Classes: []string{"iconlabel"},
		},
		Text: ps.PrecautionaryStatementLabel,
	})
	spanReference := widgets.NewSpan(widgets.SpanAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
		},
		Text: ps.PrecautionaryStatementReference,
	})
	d := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
		},
	})

	d.AppendChild(spanReference)
	d.AppendChild(spanLabel)

	return jsutils.CreateJsHTMLElementFromString(d.OuterHTML())

}

func Select2ProducerRefCreateTag(this js.Value, args []js.Value) interface{} {

	// var (
	// 	producerId int
	// 	err        error
	// )

	params := args[0]

	if jquery.Jq("input#exactMatchProducerRefs").GetVal().String() == "true" {
		return nil
	}

	if len(select2.NewSelect2(jquery.Jq("select#producer"), nil).Select2Data()) == 0 {
		jsutils.DisplayErrorMessage(locales.Translate("producerref_create_needproducer", HTTPHeaderAcceptLanguage))
		return nil
	}

	// select2ProducerId := jquery.Jq("select#producer").Select2Data()[0].Id
	// select2ProducerText := jquery.Jq("select#producer").Select2Data()[0].Text

	// if producerId, err = strconv.Atoi(select2ProducerId); err != nil {
	// 	fmt.Println(err)
	// 	return nil
	// }

	return select2.Select2Item{
		Id:   params.Get("term").String(),
		Text: params.Get("term").String(),
	}.ToJsValue()

	// return ProducerRef{
	// 	ProducerRefID:    sql.NullInt64{Int64: int64(0), Valid: true},
	// 	ProducerRefLabel: sql.NullString{String: term, Valid: true},
	// 	Producer: &Producer{
	// 		ProducerID:    sql.NullInt64{Int64: int64(producerId), Valid: true},
	// 		ProducerLabel: sql.NullString{String: select2ProducerText, Valid: true},
	// 	},
	// }.ToJsValue()

}

func Select2ProducerRefTemplateSelection(this js.Value, args []js.Value) interface{} {

	var text string

	data := args[0]

	producerRef := ProducerRef{ProducerRef: &models.ProducerRef{}}.FromJsJSONValue(data).(ProducerRef)

	if producerRef.ProducerRef == nil {
		return data.Get("text")
		// return jsutils.CreateJsHTMLElementFromString(widgets.NewDiv(widgets.DivAttributes{}).OuterHTML())
	}

	if !producerRef.ProducerRefID.Valid {
		return data.Get("text")
	}

	if producerRef.Producer != nil {
		text = fmt.Sprintf("%s (%s)", producerRef.ProducerRefLabel.String, producerRef.Producer.ProducerLabel.String)
	} else {
		text = producerRef.ProducerRefLabel.String
	}

	spanLabel := widgets.NewSpan(widgets.SpanAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
		},
		Text: text,
	})
	d := widgets.NewDiv(widgets.DivAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
		},
	})

	d.AppendChild(spanLabel)

	return jsutils.CreateJsHTMLElementFromString(d.OuterHTML())

}

func Select2ProducerRefAjaxData(this js.Value, args []js.Value) interface{} {

	params := args[0]

	var producerId string

	select2Producer := select2.NewSelect2(jquery.Jq("select#producer"), nil)
	if len(select2Producer.Select2Data()) > 0 {
		select2ItemProducer := select2Producer.Select2Data()[0]
		if !select2ItemProducer.IsEmpty() {
			producerId = select2ItemProducer.Id
		}
	}

	search := ""
	if params.Get("term").Truthy() {
		search = params.Get("term").String()
	}
	page := 1
	if params.Get("page").Truthy() {
		page = params.Get("page").Int()
	}
	offset := (page - 1) * 10
	limit := 10

	if offset < 0 {
		offset = 0
	}

	return ajax.QueryFilter{
		Producer: producerId,
		Search:   search,
		Offset:   offset,
		Page:     page,
		Limit:    limit,
	}.ToJsValue()

}

func Select2SupplierRefTemplateSelection(this js.Value, args []js.Value) interface{} {

	var text string

	data := args[0]

	supplierRef := SupplierRef{SupplierRef: &models.SupplierRef{}}.FromJsJSONValue(data).(SupplierRef)

	if supplierRef.SupplierRef == nil {
		return data.Get("text")
		// return jsutils.CreateJsHTMLElementFromString(widgets.NewDiv(widgets.DivAttributes{}).OuterHTML())
	}

	if supplierRef.SupplierRefID == 0 {
		// Autofill.
		return data.Get("text")
	}

	if supplierRef.Supplier != nil {
		// Selection.
		text = fmt.Sprintf("%s@%s", supplierRef.SupplierRefLabel, supplierRef.Supplier.SupplierLabel.String)
		supplierrefToSupplier[supplierRef.SupplierRefLabel] = supplierRef.Supplier.SupplierID.Int64
	} else {
		text = supplierRef.SupplierRefLabel
	}

	spanLabel := widgets.NewSpan(widgets.SpanAttributes{
		BaseAttributes: widgets.BaseAttributes{
			Visible: true,
		},
		Text: text,
	})

	return jsutils.CreateJsHTMLElementFromString(spanLabel.OuterHTML())

}

func Select2SupplierRefCreateTag(this js.Value, args []js.Value) interface{} {

	var (
		err               error
		select2SupplierId int
	)

	params := args[0]

	if jquery.Jq("input#exactMatchSupplierRefs").GetVal().String() == "true" {
		return nil
	}

	select2Supplier := select2.NewSelect2(jquery.Jq("select#supplier"), nil)

	if len(select2Supplier.Select2Data()) == 0 {
		jsutils.DisplayErrorMessage(locales.Translate("supplierref_create_needsupplier", HTTPHeaderAcceptLanguage))
		return nil
	}

	if select2SupplierId, err = strconv.Atoi(select2Supplier.Select2Data()[0].Id); err != nil {
		fmt.Println(err)
		return nil
	}

	supplierrefToSupplier[params.Get("term").String()] = int64(select2SupplierId)

	return select2.Select2Item{
		Id:   params.Get("term").String(),
		Text: fmt.Sprintf("%s@%s", params.Get("term").String(), select2Supplier.Select2Data()[0].Text),
	}.ToJsValue()

}

func Select2SupplierRefAjaxData(this js.Value, args []js.Value) interface{} {

	params := args[0]

	var supplierId string

	select2Supplier := select2.NewSelect2(jquery.Jq("select#supplier"), nil)
	if len(select2Supplier.Select2Data()) > 0 {
		select2ItemSupplier := select2Supplier.Select2Data()[0]
		if !select2ItemSupplier.IsEmpty() {
			supplierId = select2ItemSupplier.Id
		}
	}

	search := ""
	if params.Get("term").Truthy() {
		search = params.Get("term").String()
	}
	page := 1
	if params.Get("page").Truthy() {
		page = params.Get("page").Int()
	}
	offset := (page - 1) * 10
	limit := 10

	if offset < 0 {
		offset = 0
	}

	return ajax.QueryFilter{
		Supplier: supplierId,
		Search:   search,
		Offset:   offset,
		Page:     page,
		Limit:    limit,
	}.ToJsValue()

}

// TODO; factorise with storage
func Select2UnitTemperatureAjaxData(this js.Value, args []js.Value) interface{} {

	params := args[0]

	search := ""
	if params.Get("term").Truthy() {
		search = params.Get("term").String()
	}
	page := 1
	if params.Get("page").Truthy() {
		page = params.Get("page").Int()
	}
	offset := (page - 1) * 10
	limit := 10

	if offset < 0 {
		offset = 0
	}

	return ajax.QueryFilter{
		UnitType: "temperature",
		Search:   search,
		Offset:   offset,
		Page:     page,
		Limit:    limit,
	}.ToJsValue()

}
