package types

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/models"
	"github.com/tbellembois/gochimitheque-wasm/select2"
)

type Products struct {
	Rows     []*Product `json:"rows"`
	Total    int        `json:"total"`
	ExportFn string     `json:"exportfn"`
}

func (elems Products) IsExactMatch() bool {

	return false

}

func (elems Products) GetRowConcreteTypeName() string {

	return "Product"

}

type Bookmark struct {
	BookmarkID sql.NullInt64 `json:"bookmark_id"`
	Person     `json:"person"`
	Product    `json:"product"`
}

type Product struct {
	*models.Product
}

// ProductFromJsJSONValue converts a JS JSON into a
// Go product.
func (p Product) FromJsJSONValue(jsvalue js.Value) select2.Select2ItemAble {

	var (
		product Product
		err     error
	)

	jsProductString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsProductString), &product); err != nil {
		fmt.Println(err)
	}

	return product

}

// ProductFromJsJSONValue converts a JS JSON into a
// Go product.
func (p Product) ProductFromJsJSONValue(jsvalue js.Value) Product {

	var (
		product Product
		err     error
	)

	jsProductString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsProductString), &product); err != nil {
		fmt.Println(err)
	}

	return product

}

func (Products) FromJsJSONValue(jsvalue js.Value) select2.Select2ResultAble {

	var (
		products Products
		err      error
	)

	jsProductsString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsProductsString), &products); err != nil {
		fmt.Println(err)
	}

	return products

}

func (p Products) GetRows() []select2.Select2ItemAble {

	var select2ItemAble []select2.Select2ItemAble = make([]select2.Select2ItemAble, len(p.Rows))

	for i, row := range p.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (p Products) GetTotal() int {

	return p.Total

}

func (p Products) GetExportFn() string {

	return p.ExportFn

}

func (p Product) GetSelect2Id() int {

	return p.ProductID

}

func (p Product) GetSelect2Text() string {

	if p.Product != nil {
		return p.Name.NameLabel
	}

	return ""

}
