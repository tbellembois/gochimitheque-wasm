package types

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/barweiss/go-tuple"
	"github.com/tbellembois/gochimitheque-wasm/select2"
	"github.com/tbellembois/gochimitheque/models"
)

type Select2ProducerRefs struct {
	Rows  []*ProducerRef `json:"rows"`
	Total int            `json:"total"`
}

type ProducerRef struct {
	*models.ProducerRef
}

func (elems Select2ProducerRefs) GetRowConcreteTypeName() string {

	return "ProducerRef"

}

func (p ProducerRef) ToJsValue() js.Value {

	var (
		marshalJson []byte
		err         error
	)

	if marshalJson, err = json.Marshal(p); err != nil {
		fmt.Println(err)
		return js.Null()
	}

	return js.Global().Get("JSON").Call("parse", string(marshalJson))

}

func (elems Select2ProducerRefs) IsExactMatch() bool {

	for _, elem := range elems.Rows {
		if elem.MatchExactSearch {
			return true
		}
	}

	return false

}

func (r ProducerRef) MarshalJSON() ([]byte, error) {

	type Copy ProducerRef
	return json.Marshal(&struct {
		Id   int    `json:"id"`
		Text string `json:"text"`
		Copy
	}{
		Id:   r.GetSelect2Id(),
		Text: r.GetSelect2Text(),
		Copy: (Copy)(r),
	})

}

func (Select2ProducerRefs) FromJsJSONValue(jsvalue js.Value) select2.Select2ResultAble {

	var (
		producersAjaxResponse tuple.T2[[]struct {
			MatchExactSearch bool   `json:"match_exact_search"`
			ProducerRefID    int64  `json:"producerref_id"`
			ProducerRefLabel string `json:"producerref_label"`
		}, int]
		select2ProducerRefs Select2ProducerRefs
		err                 error
	)

	jsProducersString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsProducersString), &producersAjaxResponse); err != nil {
		fmt.Println("(Select2ProducerRefs) FromJsJSONValue:" + err.Error())
	}

	for _, producerref := range producersAjaxResponse.V1 {
		select2ProducerRefs.Rows = append(select2ProducerRefs.Rows, &ProducerRef{
			&models.ProducerRef{
				MatchExactSearch: producerref.MatchExactSearch,
				ProducerRefID:    sql.NullInt64{Int64: producerref.ProducerRefID, Valid: true},
				ProducerRefLabel: sql.NullString{String: producerref.ProducerRefLabel, Valid: true},
			},
		})
	}

	select2ProducerRefs.Total = producersAjaxResponse.V2

	return select2ProducerRefs

}

func (r Select2ProducerRefs) GetRows() []select2.Select2ItemAble {

	var select2ItemAble []select2.Select2ItemAble = make([]select2.Select2ItemAble, len(r.Rows))

	for i, row := range r.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (r Select2ProducerRefs) GetTotal() int {

	return r.Total

}

func (r ProducerRef) FromJsJSONValue(jsvalue js.Value) select2.Select2ItemAble {

	var (
		producerRef ProducerRef
		err         error
	)

	jsProducerRefString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsProducerRefString), &producerRef); err != nil {
		fmt.Println(err)
	}

	return producerRef

}

func (r ProducerRef) GetSelect2Id() int {

	return int(r.ProducerRefID.Int64)

}

func (r ProducerRef) GetSelect2Text() string {

	if r.ProducerRef != nil {
		return r.ProducerRefLabel.String
	}

	return ""

}
