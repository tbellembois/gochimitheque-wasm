package ajax

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"syscall/js"
)

// QueryParams is the data sent while requesting
// remote data as defined
// https://api.jquery.com/jquery.ajax/
// and
// https://bootstrap-table.com/docs/api/table-options/#queryparams
type QueryParams struct {
	Data QueryFilter `js:"data"`
}

// QueryFilter contains the parameters sent
// while doing AJAX requests to retrieve multiple
// results ("/entities", "/people"...).
// It is especially used by select2 and bootstraptable.
type QueryFilter struct {
	Product                     string `json:"product,omitempty"`
	ProductFilterLabel          string
	ProductBookmark             bool   `json:"bookmark,omitempty"`
	Storage                     string `json:"storage,omitempty"`
	StorageFilterLabel          string
	UnitType                    string `json:"unit_type,omitempty"`
	Supplier                    string `json:"supplier,omitempty"`
	Producer                    string `json:"producer,omitempty"`
	StoreLocation               string `json:"storelocation,omitempty"`
	StoreLocationCanStore       bool   `json:"storelocation_canstore,omitempty"`
	StoreLocationFilterLabel    string
	Entity                      string `json:"entity,omitempty"`
	Name                        string `json:"name,omitempty"`
	NameFilterLabel             string
	CasNumber                   string `json:"casnumber,omitempty"`
	CasNumberFilterLabel        string
	EmpiricalFormula            string `json:"empiricalformula,omitempty"`
	EmpiricalFormulaFilterLabel string
	SignalWord                  string   `json:"signalword,omitempty"`
	HazardStatements            []string `json:"hazardstatements,omitempty"`
	PrecautionaryStatements     []string `json:"precautionarystatements,omitempty"`
	Symbols                     []string `json:"symbols,omitempty"`
	StorageBarecode             string   `json:"storage_barecode,omitempty"`
	CustomNamePartOf            string   `json:"custom_name_part_of,omitempty"`
	CasNumberCMR                bool     `json:"casnumber_cmr,omitempty"`
	Borrowing                   bool     `json:"borrowing,omitempty"`
	StorageToDestroy            bool     `json:"storage_to_destroy,omitempty"`
	StorageArchive              bool     `json:"storage_archive,omitempty"`
	StorageHistory              bool     `json:"storage_history,omitempty"`
	Export                      bool     `json:"export,omitempty"`

	Search string `json:"search,omitempty"`
	Page   int    `json:"page,omitempty"`
	Sort   string `json:"sort,omitempty"`
	Order  string `json:"order,omitempty"`
	Offset int    `json:"offset,omitempty"`
	Limit  int    `json:"limit,omitempty"`
}

// QueryFilterFromJsJSONValue converts a JS JSON into a
// Go queryFilter.
func QueryFilterFromJsJSONValue(jsvalue js.Value) QueryFilter {

	var (
		queryFilter QueryFilter
		err         error
	)

	jsQueryFilterString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsQueryFilterString), &queryFilter); err != nil {
		fmt.Println(err)
	}

	return queryFilter

}

// ToJsValue converts a Go QueryFilter
// into a JS JSON.
func (q QueryFilter) ToJsValue() js.Value {

	var (
		marshalJson []byte
		err         error
	)

	if marshalJson, err = json.Marshal(q); err != nil {
		fmt.Println(err)
		return js.Null()
	}

	return js.Global().Get("JSON").Call("parse", string(marshalJson))

}

// ToRawQuery converts a QueryFilter into
// an url.RawQuery
func (q QueryFilter) ToRawQuery() string {

	values := url.Values{}
	if q.Sort != "" {
		values.Set("sort", q.Sort)
	}
	if q.Order != "" {
		values.Set("order", q.Order)
	}
	if q.Search != "" {
		values.Set("search", q.Search)
	}
	values.Set("page", strconv.Itoa(q.Page))
	values.Set("offset", strconv.Itoa(q.Offset))
	values.Set("limit", strconv.Itoa(q.Limit))

	if q.Product != "" {
		values.Set("product", q.Product)
	}
	if q.Storage != "" {
		values.Set("storage", q.Storage)
	}
	if q.Supplier != "" {
		values.Set("supplier", q.Supplier)
	}
	if q.Producer != "" {
		values.Set("producer", q.Producer)
	}
	if q.StoreLocation != "" {
		values.Set("storelocation", q.StoreLocation)
	}
	if q.Entity != "" {
		values.Set("entity", q.Entity)
	}
	if q.Name != "" {
		values.Set("name", q.Name)
	}
	if q.CasNumber != "" {
		values.Set("casnumber", q.CasNumber)
	}
	if q.EmpiricalFormula != "" {
		values.Set("empiricalformula", q.EmpiricalFormula)
	}
	if q.SignalWord != "" {
		values.Set("signalword", q.SignalWord)
	}
	if len(q.HazardStatements) > 0 {
		for _, hs := range q.HazardStatements {
			values.Set("hazardstatements[]", hs)
		}
	}
	if len(q.PrecautionaryStatements) > 0 {
		for _, ps := range q.PrecautionaryStatements {
			values.Set("precautionarystatements[]", ps)
		}
	}
	if len(q.Symbols) > 0 {
		for _, s := range q.Symbols {
			values.Set("symbols[]", s)
		}
	}

	if q.UnitType != "" {
		values.Set("unit_type", q.UnitType)
	}
	if q.StorageBarecode != "" {
		values.Set("storage_barecode", q.StorageBarecode)
	}
	if q.CustomNamePartOf != "" {
		values.Set("custom_name_part_of", q.CustomNamePartOf)
	}
	if q.CasNumberCMR {
		values.Set("casnumber_cmr", strconv.FormatBool(true))
	}
	if q.Borrowing {
		values.Set("borrowing", strconv.FormatBool(true))
	}
	if q.StorageToDestroy {
		values.Set("storage_to_destroy", strconv.FormatBool(true))
	}
	if q.StorageArchive {
		values.Set("storage_archive", strconv.FormatBool(true))
	}
	if q.StorageHistory {
		values.Set("history", strconv.FormatBool(true))
	}
	if q.ProductBookmark {
		values.Set("bookmark", strconv.FormatBool(true))
	}
	if q.Export {
		values.Set("export", strconv.FormatBool(true))
	}

	return values.Encode()

}