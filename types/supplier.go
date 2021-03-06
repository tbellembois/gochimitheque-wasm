package types

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"syscall/js"
)

type Suppliers struct {
	Rows  []*Supplier `json:"rows"`
	Total int         `json:"total"`
}

type Supplier struct {
	C             int            `json:"c"` // not stored in db but db:"c" set for sqlx
	SupplierID    sql.NullInt64  `json:"supplier_id"`
	SupplierLabel sql.NullString `json:"supplier_label"`
}

func (elems Suppliers) IsExactMatch() bool {

	return false

}

func (s *Supplier) MarshalJSON() ([]byte, error) {

	type Copy Supplier
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

func (Suppliers) FromJsJSONValue(jsvalue js.Value) Select2ResultAble {

	var (
		suppliers Suppliers
		err       error
	)

	jsSuppliersString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsSuppliersString), &suppliers); err != nil {
		fmt.Println(err)
	}

	return suppliers

}

func (s Suppliers) GetRows() []Select2ItemAble {

	var select2ItemAble []Select2ItemAble = make([]Select2ItemAble, len(s.Rows))

	for i, row := range s.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (r Suppliers) GetTotal() int {

	return r.Total

}

func (s Supplier) FromJsJSONValue(jsvalue js.Value) Select2ItemAble {

	var (
		supplier Supplier
		err      error
	)

	jsSupplierString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsSupplierString), &supplier); err != nil {
		fmt.Println(err)
	}

	return supplier

}

func (s Supplier) GetSelect2Id() int {

	return int(s.SupplierID.Int64)

}

func (s Supplier) GetSelect2Text() string {

	return s.SupplierLabel.String

}
