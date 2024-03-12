package types

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/barweiss/go-tuple"
	"github.com/tbellembois/gochimitheque-wasm/models"
	"github.com/tbellembois/gochimitheque-wasm/select2"
)

type Select2Tags struct {
	Rows  []*Tag `json:"rows"`
	Total int    `json:"total"`
}

type Tag struct {
	*models.Tag
}

func (elems Select2Tags) GetRowConcreteTypeName() string {

	return "Tag"

}

func (elems Select2Tags) IsExactMatch() bool {

	for _, elem := range elems.Rows {
		if elem.MatchExactSearch {
			return true
		}
	}

	return false

}

func (t *Tag) MarshalJSON() ([]byte, error) {
	type Copy Tag
	return json.Marshal(&struct {
		Id   int    `json:"id"`
		Text string `json:"text"`

		*Copy
	}{
		Id:   t.GetSelect2Id(),
		Text: t.GetSelect2Text(),
		Copy: (*Copy)(t),
	})
}

func (Select2Tags) FromJsJSONValue(jsvalue js.Value) select2.Select2ResultAble {
	var (
		tagsAjaxResponse tuple.T2[[]struct {
			MatchExactSearch bool   `json:"match_exact_search"`
			TagID            int64  `json:"tag_id"`
			TagLabel         string `json:"tag_label"`
		}, int]
		select2Tags Select2Tags
		err         error
	)

	jsTagsString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsTagsString), &tagsAjaxResponse); err != nil {
		fmt.Println("(Select2Tags) FromJsJSONValue:" + err.Error())
	}

	for _, tag := range tagsAjaxResponse.V1 {
		select2Tags.Rows = append(select2Tags.Rows, &Tag{
			&models.Tag{
				MatchExactSearch: tag.MatchExactSearch,
				TagID:            int(tag.TagID),
				TagLabel:         tag.TagLabel,
			},
		})
	}

	select2Tags.Total = tagsAjaxResponse.V2

	return select2Tags
}

func (t Tag) FromJsJSONValue(jsvalue js.Value) select2.Select2ItemAble {

	var (
		tag Tag
		err error
	)

	jsTagString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsTagString), &tag); err != nil {
		fmt.Println(err)
	}

	return tag

}

func (t Select2Tags) GetRows() []select2.Select2ItemAble {

	var select2ItemAble []select2.Select2ItemAble = make([]select2.Select2ItemAble, len(t.Rows))

	for i, row := range t.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (t Select2Tags) GetTotal() int {

	return t.Total

}

func (t Tag) GetSelect2Id() int {

	return t.TagID

}

func (t Tag) GetSelect2Text() string {

	if t.Tag != nil {
		return t.TagLabel
	}

	return ""

}
