package types

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/barweiss/go-tuple"
	"github.com/tbellembois/gochimitheque-wasm/models"
	"github.com/tbellembois/gochimitheque-wasm/select2"
)

type Select2Producers struct {
	Rows  []*Producer `json:"rows"`
	Total int         `json:"total"`
}

type Producer struct {
	*models.Producer
}

func (elems Select2Producers) GetRowConcreteTypeName() string {

	return "Producer"

}

func (elems Select2Producers) IsExactMatch() bool {

	return false

}

func (r *Producer) MarshalJSON() ([]byte, error) {

	type Copy Producer
	return json.Marshal(&struct {
		Id   int    `json:"id"`
		Text string `json:"text"`

		*Copy
	}{
		Id:   r.GetSelect2Id(),
		Text: r.GetSelect2Text(),
		Copy: (*Copy)(r),
	})

}

func (Select2Producers) FromJsJSONValue(jsvalue js.Value) select2.Select2ResultAble {

	var (
		producersAjaxResponse tuple.T2[[]struct {
			MatchExactSearch bool   `json:"match_exact_search"`
			ProducerID       int64  `json:"producer_id"`
			ProducerLabel    string `json:"producer_label"`
		}, int]
		select2Producers Select2Producers
		err              error
	)

	jsProducersString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsProducersString), &producersAjaxResponse); err != nil {
		fmt.Println("(Select2Producers) FromJsJSONValue:" + err.Error())
	}

	for _, producer := range producersAjaxResponse.V1 {
		select2Producers.Rows = append(select2Producers.Rows, &Producer{
			&models.Producer{
				MatchExactSearch: producer.MatchExactSearch,
				ProducerID:       sql.NullInt64{Int64: producer.ProducerID, Valid: true},
				ProducerLabel:    sql.NullString{String: producer.ProducerLabel, Valid: true},
			},
		})
	}

	select2Producers.Total = producersAjaxResponse.V2

	return select2Producers

}

func (r Select2Producers) GetRows() []select2.Select2ItemAble {

	var select2ItemAble []select2.Select2ItemAble = make([]select2.Select2ItemAble, len(r.Rows))

	for i, row := range r.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (r Select2Producers) GetTotal() int {

	return r.Total

}

func (r Producer) FromJsJSONValue(jsvalue js.Value) select2.Select2ItemAble {

	var (
		Producer Producer
		err      error
	)

	jsProducerString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsProducerString), &Producer); err != nil {
		fmt.Println(err)
	}

	return Producer

}

func (r Producer) GetSelect2Id() int {

	return int(r.ProducerID.Int64)

}

func (r Producer) GetSelect2Text() string {

	if r.Producer != nil {
		return r.ProducerLabel.String
	}

	return ""

}
