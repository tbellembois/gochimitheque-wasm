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

type Select2SignalWords struct {
	Rows  []*SignalWord `json:"rows"`
	Total int           `json:"total"`
}

type SignalWord struct {
	*models.SignalWord
}

func (elems Select2SignalWords) GetRowConcreteTypeName() string {

	return "SignalWord"

}

func (elems Select2SignalWords) IsExactMatch() bool {

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

func (Select2SignalWords) FromJsJSONValue(jsvalue js.Value) select2.Select2ResultAble {
	var (
		signalwordsAjaxResponse tuple.T2[[]struct {
			// MatchExactSearch bool   `json:"match_exact_search"`
			SignalWordID    int64  `json:"signalword_id"`
			SignalWordLabel string `json:"signalword_label"`
		}, int]
		select2SignalWords Select2SignalWords
		err                error
	)

	jsSignalWordsString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsSignalWordsString), &signalwordsAjaxResponse); err != nil {
		fmt.Println("(Select2SignalWords) FromJsJSONValue:" + err.Error())
	}

	for _, signalword := range signalwordsAjaxResponse.V1 {
		select2SignalWords.Rows = append(select2SignalWords.Rows, &SignalWord{
			&models.SignalWord{
				// MatchExactSearch: signalword.MatchExactSearch,
				SignalWordID:    sql.NullInt64{Int64: signalword.SignalWordID, Valid: true},
				SignalWordLabel: sql.NullString{String: signalword.SignalWordLabel, Valid: true},
			},
		})
	}

	select2SignalWords.Total = signalwordsAjaxResponse.V2

	return select2SignalWords
}

func (s Select2SignalWords) GetRows() []select2.Select2ItemAble {

	var select2ItemAble []select2.Select2ItemAble = make([]select2.Select2ItemAble, len(s.Rows))

	for i, row := range s.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (s Select2SignalWords) GetTotal() int {

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

	return int(s.SignalWordID.Int64)

}

func (s SignalWord) GetSelect2Text() string {

	if s.SignalWord != nil {
		return s.SignalWordLabel.String
	}

	return ""

}
