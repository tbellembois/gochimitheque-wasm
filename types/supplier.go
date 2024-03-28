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

type Select2Suppliers struct {
	Rows  []*Supplier `json:"rows"`
	Total int         `json:"total"`
}

type Supplier struct {
	*models.Supplier
}

func (elems Select2Suppliers) GetRowConcreteTypeName() string {

	return "Supplier"

}

func (elems Select2Suppliers) IsExactMatch() bool {

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

func (Select2Suppliers) FromJsJSONValue(jsvalue js.Value) select2.Select2ResultAble {

	var (
		suppliersAjaxResponse tuple.T2[[]struct {
			MatchExactSearch bool   `json:"match_exact_search"`
			SupplierID       int64  `json:"supplier_id"`
			SupplierLabel    string `json:"supplier_label"`
		}, int]
		select2Suppliers Select2Suppliers
		err              error
	)

	jsSuppliersString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsSuppliersString), &suppliersAjaxResponse); err != nil {
		fmt.Println("(Select2Suppliers) FromJsJSONValue:" + err.Error())
	}

	for _, supplier := range suppliersAjaxResponse.V1 {
		select2Suppliers.Rows = append(select2Suppliers.Rows, &Supplier{
			&models.Supplier{
				MatchExactSearch: supplier.MatchExactSearch,
				SupplierID:       sql.NullInt64{Int64: supplier.SupplierID, Valid: true},
				SupplierLabel:    sql.NullString{String: supplier.SupplierLabel, Valid: true},
			},
		})
	}

	select2Suppliers.Total = suppliersAjaxResponse.V2

	return select2Suppliers

}

func (s Select2Suppliers) GetRows() []select2.Select2ItemAble {

	var select2ItemAble []select2.Select2ItemAble = make([]select2.Select2ItemAble, len(s.Rows))

	for i, row := range s.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (r Select2Suppliers) GetTotal() int {

	return r.Total

}

func (s Supplier) FromJsJSONValue(jsvalue js.Value) select2.Select2ItemAble {

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

	if s.Supplier != nil {
		return s.SupplierLabel.String
	}

	return ""

}
