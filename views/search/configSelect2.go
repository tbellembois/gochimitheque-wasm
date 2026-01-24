//go:build go1.24 && js && wasm

package search

import (
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/ajax"
	"github.com/tbellembois/gochimitheque-wasm/jquery"
	"github.com/tbellembois/gochimitheque-wasm/jsutils"
	"github.com/tbellembois/gochimitheque-wasm/select2"
	. "github.com/tbellembois/gochimitheque-wasm/types"
	"github.com/tbellembois/gochimitheque-wasm/widgets"
	"github.com/tbellembois/gochimitheque/models"
)

func Select2StoreLocationAjaxData(this js.Value, args []js.Value) interface{} {

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

	select2SEntity := select2.NewSelect2(jquery.Jq("select#s_entity"), nil)
	if select2SEntity.Select2IsInitialized() {
		i := select2SEntity.Select2Data()
		if len(i) > 0 {
			return ajax.QueryFilter{
				Entity: i[0].Id,
				Search: search,
				Offset: offset,
				Page:   page,
				Limit:  limit,
			}.ToJsValue()
		}
	}

	return ajax.QueryFilter{
		Search: search,
		Offset: offset,
		Page:   page,
		Limit:  limit,
	}.ToJsValue()

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
