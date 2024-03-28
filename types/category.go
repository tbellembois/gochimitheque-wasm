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

type Select2Categories struct {
	Rows  []*Category `json:"rows"`
	Total int         `json:"total"`
}

type Category struct {
	*models.Category
}

func (elems Select2Categories) GetRowConcreteTypeName() string {

	return "Category"

}

func (elems Select2Categories) IsExactMatch() bool {

	for _, elem := range elems.Rows {
		if elem.MatchExactSearch {
			return true
		}
	}

	return false

}

func (c *Category) MarshalJSON() ([]byte, error) {
	type Copy Category
	return json.Marshal(&struct {
		Id   int    `json:"id"`
		Text string `json:"text"`

		*Copy
	}{
		Id:   c.GetSelect2Id(),
		Text: c.GetSelect2Text(),

		Copy: (*Copy)(c),
	})
}

func (Select2Categories) FromJsJSONValue(jsvalue js.Value) select2.Select2ResultAble {
	var (
		categoriesAjaxResponse tuple.T2[[]struct {
			MatchExactSearch bool   `json:"match_exact_search"`
			CategoryID       int64  `json:"category_id"`
			CategoryLabel    string `json:"category_label"`
		}, int]
		select2Categories Select2Categories
		err               error
	)

	jsCategoriesString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsCategoriesString), &categoriesAjaxResponse); err != nil {
		fmt.Println("(Select2Categories) FromJsJSONValue:" + err.Error())
	}

	for _, category := range categoriesAjaxResponse.V1 {
		select2Categories.Rows = append(select2Categories.Rows, &Category{
			&models.Category{
				MatchExactSearch: category.MatchExactSearch,
				CategoryID:       sql.NullInt64{Int64: category.CategoryID, Valid: true},
				CategoryLabel:    sql.NullString{String: category.CategoryLabel, Valid: true},
			},
		})
	}

	select2Categories.Total = categoriesAjaxResponse.V2

	return select2Categories
}

func (c Category) FromJsJSONValue(jsvalue js.Value) select2.Select2ItemAble {

	var (
		category Category
		err      error
	)

	jsCategoryString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsCategoryString), &category); err != nil {
		fmt.Println(err)
	}

	return category

}

func (c Select2Categories) GetRows() []select2.Select2ItemAble {

	var select2ItemAble []select2.Select2ItemAble = make([]select2.Select2ItemAble, len(c.Rows))

	for i, row := range c.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (c Select2Categories) GetTotal() int {

	return c.Total

}

func (c Category) GetSelect2Id() int {

	return int(c.CategoryID.Int64)

}

func (c Category) GetSelect2Text() string {

	if c.Category != nil {
		return c.CategoryLabel.String
	}

	return ""

}
