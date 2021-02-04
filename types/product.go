package types

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"syscall/js"

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
	ProductID              int            `json:"product_id"`
	ProductSpecificity     sql.NullString `json:"product_specificity"`
	ProductMSDS            sql.NullString `json:"product_msds"`
	ProductRestricted      sql.NullBool   `json:"product_restricted"`
	ProductRadioactive     sql.NullBool   `json:"product_radioactive"`
	ProductThreeDFormula   sql.NullString `json:"product_threedformula"`
	ProductTwoDFormula     sql.NullString `json:"product_twodformula"`
	ProductMolFormula      sql.NullString `json:"product_molformula"`
	ProductDisposalComment sql.NullString `json:"product_disposalcomment"`
	ProductRemark          sql.NullString `json:"product_remark"`
	ProductTemperature     sql.NullInt64  `json:"product_temperature"`
	ProductSheet           sql.NullString `json:"product_sheet"`
	EmpiricalFormula       `json:"empiricalformula"`
	LinearFormula          `json:"linearformula"`
	PhysicalState          `json:"physicalstate"`
	SignalWord             `json:"signalword"`
	Person                 `json:"person"`
	CasNumber              `json:"casnumber"`
	CeNumber               `json:"cenumber"`
	Name                   `json:"name"`
	ProducerRef            `json:"producerref"`
	Category               `json:"category"`
	UnitTemperature        Unit `json:"unit_temperature"`

	ClassOfCompound         []ClassOfCompound        `json:"classofcompound"`
	Synonyms                []Name                   `json:"synonyms"`
	Symbols                 []Symbol                 `json:"symbols"`
	HazardStatements        []HazardStatement        `json:"hazardstatements"`
	PrecautionaryStatements []PrecautionaryStatement `json:"precautionarystatements"`
	SupplierRefs            []SupplierRef            `json:"supplierrefs"`
	Tags                    []Tag                    `json:"tags"`

	Bookmark *Bookmark `json:"bookmark"` // not in db but sqlx requires the "db" entry

	// archived storage count in the logged user entity(ies)
	ProductASC int `json:"product_asc"` // not in db but sqlx requires the "db" entry
	// total storage count
	ProductTSC int `json:"product_tsc"` // not in db but sqlx requires the "db" entry
	// storage count in the logged user entity(ies)
	ProductSC int `json:"product_sc"` // not in db but sqlx requires the "db" entry
	// storage barecode concatenation
	ProductSL sql.NullString `json:"product_sl"` // not in db but sqlx requires the "db" entry
	// hazard statement CMR concatenation
	HazardStatementCMR sql.NullString `json:"hazardstatement_cmr"` // not in db but sqlx requires the "db" entry

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

	return p.Name.NameLabel

}
