package types

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/select2"
	"github.com/tbellembois/gochimitheque/models"
)

type Producers struct {
	Rows  []*Producer `json:"rows"`
	Total int         `json:"total"`
}

type Producer struct {
	*models.Producer
}

func (elems Producers) GetRowConcreteTypeName() string {

	return "Producer"

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

func (Producers) FromJsJSONValue(jsvalue js.Value) select2.Select2ResultAble {

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

func (r Producers) GetRows() []select2.Select2ItemAble {

	var select2ItemAble []select2.Select2ItemAble = make([]select2.Select2ItemAble, len(r.Rows))

	for i, row := range r.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (r Producers) GetTotal() int {

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
