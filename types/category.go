package types

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/select2"
)

type Categories struct {
	Rows  []*Category `json:"rows"`
	Total int         `json:"total"`
}

type Category struct {
	C             int            `json:"c"` // not stored in db but db:"c" set for sqlx
	CategoryID    sql.NullInt64  `json:"category_id"`
	CategoryLabel sql.NullString `json:"category_label"`
}

func (elems Categories) IsExactMatch() bool {

	for _, elem := range elems.Rows {
		if elem.C == 1 {
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

func (Categories) FromJsJSONValue(jsvalue js.Value) select2.Select2ResultAble {

	var (
		categories Categories
		err        error
	)

	jsCategoriesString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsCategoriesString), &categories); err != nil {
		fmt.Println(err)
	}

	return categories

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

func (c Categories) GetRows() []select2.Select2ItemAble {

	var select2ItemAble []select2.Select2ItemAble = make([]select2.Select2ItemAble, len(c.Rows))

	for i, row := range c.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (c Categories) GetTotal() int {

	return c.Total

}

func (c Category) GetSelect2Id() int {

	return int(c.CategoryID.Int64)

}

func (c Category) GetSelect2Text() string {

	return c.CategoryLabel.String

}
