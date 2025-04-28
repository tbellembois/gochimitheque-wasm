//go:build go1.24 && js && wasm

package types

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/select2"
	"github.com/tbellembois/gochimitheque/models"
)

type SignalWords struct {
	Rows  []*SignalWord `json:"rows"`
	Total int           `json:"total"`
}

type SignalWord struct {
	*models.SignalWord
}

func (elems SignalWords) GetRowConcreteTypeName() string {

	return "SignalWord"

}

func (elems SignalWords) IsExactMatch() bool {

	return false

}

func (s *SignalWord) MarshalJSON() ([]byte, error) {
	type Copy SignalWord
	return json.Marshal(&struct {
		Id   int    `json:"id"`
		Text string `json:"text"`

		*Copy
	}{
		Id:   s.GetSelect2Id(),
		Text: s.GetSelect2Text(),
		Copy: (*Copy)(s),
	})
}

func (SignalWords) FromJsJSONValue(jsvalue js.Value) select2.Select2ResultAble {

	var (
		signalWords SignalWords
		err         error
	)

	jsSignalWordsString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsSignalWordsString), &signalWords); err != nil {
		fmt.Println(err)
	}

	return signalWords

}

func (s SignalWords) GetRows() []select2.Select2ItemAble {

	var select2ItemAble []select2.Select2ItemAble = make([]select2.Select2ItemAble, len(s.Rows))

	for i, row := range s.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (s SignalWords) GetTotal() int {

	return s.Total

}

func (s SignalWord) FromJsJSONValue(jsvalue js.Value) select2.Select2ItemAble {

	var (
		signalWord SignalWord
		err        error
	)

	jsSignalWordString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsSignalWordString), &signalWord); err != nil {
		fmt.Println(err)
	}

	return signalWord

}

func (s SignalWord) GetSelect2Id() int {

	return int(*s.SignalWordID)

}

func (s SignalWord) GetSelect2Text() string {

	if s.SignalWord != nil {
		return *s.SignalWordLabel
	}

	return ""

}
