package types

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"syscall/js"
)

type Producers struct {
	Rows  []*Producer `json:"rows"`
	Total int         `json:"total"`
}

type Producer struct {
	C             int            `json:"c"` // not stored in db but db:"c" set for sqlx
	ProducerID    sql.NullInt64  `json:"producer_id"`
	ProducerLabel sql.NullString `json:"producer_label"`
}

func (elems Producers) IsExactMatch() bool {

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

func (Producers) FromJsJSONValue(jsvalue js.Value) Select2ResultAble {

	var (
		producers Producers
		err       error
	)

	jsProducersString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsProducersString), &producers); err != nil {
		fmt.Println(err)
	}

	return producers

}

func (r Producers) GetRows() []Select2ItemAble {

	var select2ItemAble []Select2ItemAble = make([]Select2ItemAble, len(r.Rows))

	for i, row := range r.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (r Producers) GetTotal() int {

	return r.Total

}

func (r Producer) FromJsJSONValue(jsvalue js.Value) Select2ItemAble {

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

	return r.ProducerLabel.String

}
