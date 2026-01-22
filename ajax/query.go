//go:build go1.24 && js && wasm

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
	ShowBio                            bool   `json:"show_bio,omitempty"`
	ShowChem                           bool   `json:"show_chem,omitempty"`
	ShowConsu                          bool   `json:"show_consu,omitempty"`
	Product                            string `json:"product,omitempty"`
	ProductFilterLabel                 string
	ProductBookmark                    bool   `json:"bookmark,omitempty"`
	ProducerRef                        string `json:"producer_ref,omitempty"`
	ProducerRefFilterLabel             string
	Storage                            string `json:"storage,omitempty"`
	StorageFilterLabel                 string
	Storages                           []int `json:"ids,omitempty"`
	StoragesFilterLabel                string
	UnitType                           string `json:"unit_type,omitempty"`
	Supplier                           string `json:"supplier,omitempty"`
	Producer                           string `json:"producer,omitempty"`
	StoreLocation                      string `json:"store_location,omitempty"`
	StoreLocationCanStore              bool   `json:"store_location_can_store,omitempty"`
	StoreLocationFilterLabel           string
	Entity                             string `json:"entity,omitempty"`
	EntityFilterLabel                  string
	Name                               string `json:"name,omitempty"`
	NameFilterLabel                    string
	CasNumber                          string `json:"cas_number,omitempty"`
	CasNumberFilterLabel               string
	EmpiricalFormula                   string `json:"empirical_formula,omitempty"`
	EmpiricalFormulaFilterLabel        string
	SignalWord                         string `json:"signal_word,omitempty"`
	SignalWordFilterLabel              string
	HazardStatements                   []string `json:"hazard_statements,omitempty"`
	HazardStatementsFilterLabel        string
	PrecautionaryStatements            []string `json:"precautionary_statements,omitempty"`
	PrecautionaryStatementsFilterLabel string
	Symbols                            []string `json:"symbols,omitempty"`
	SymbolsFilterLabel                 string
	Tags                               []string `json:"tags,omitempty"`
	TagsFilterLabel                    string
	Category                           string `json:"category,omitempty"`
	CategoryFilterLabel                string
	StorageBarecode                    string `json:"storage_barecode,omitempty"`
	StorageBarecodeFilterLabel         string
	CustomNamePartOf                   string `json:"custom_name_part_of,omitempty"`
	CustomNamePartOfFilterLabel        string
	IsCMR                              bool `json:"is_cmr,omitempty"`
	CasNumberCMRFilterLabel            string
	Borrowing                          bool `json:"borrowing,omitempty"`
	BorrowingFilterLabel               string
	StorageToDestroy                   bool `json:"storage_to_destroy,omitempty"`
	StorageToDestroyFilterLabel        string
	StorageArchive                     *bool  `json:"storage_archive,omitempty"`
	StorageHistory                     bool   `json:"storage_history,omitempty"`
	StorageBatchNumber                 string `json:"storage_batch_number,omitempty"`
	StorageBatchNumberFilterLabel      string
	Export                             bool `json:"export,omitempty"`

	Search string `json:"search,omitempty"`
	Page   int    `json:"page,omitempty"`
	Sort   string `json:"sort,omitempty"`
	Order  string `json:"order,omitempty"`
	Offset int    `json:"offset,omitempty"`
	Limit  int    `json:"limit,omitempty"`
	Id     int    `json:"id,omitempty"`
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

	// js.Global().Get("console").Call("log", fmt.Sprintf("%#v", q))

	values := url.Values{}
	if q.Sort != "" {
		values.Set("order_by", q.Sort)
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
	if q.Id != 0 {
		values.Set("id", strconv.Itoa(q.Id))
	}

	if q.Product != "" {
		values.Set("product", q.Product)
	}
	if q.Storage != "" {
		values.Set("storage", q.Storage)
	}
	if len(q.Storages) > 0 {
		for _, storage := range q.Storages {
			values.Add("ids", strconv.Itoa(storage))
		}
	}
	if q.Supplier != "" {
		values.Set("supplier", q.Supplier)
	}
	if q.Producer != "" {
		values.Set("producer", q.Producer)
	}
	if q.StoreLocation != "" {
		values.Set("store_location", q.StoreLocation)
	}
	if q.ProducerRef != "" {
		values.Set("producer_ref", q.ProducerRef)
	}
	if q.Entity != "" {
		values.Set("entity", q.Entity)
	}
	if q.Name != "" {
		values.Set("name", q.Name)
	}
	if q.CasNumber != "" {
		values.Set("cas_number", q.CasNumber)
	}
	if q.EmpiricalFormula != "" {
		values.Set("empirical_formula", q.EmpiricalFormula)
	}
	if q.SignalWord != "" {
		values.Set("signal_word", q.SignalWord)
	}
	if q.Category != "" {
		values.Set("category", q.Category)
	}
	if len(q.HazardStatements) > 0 {
		hs_ids := ""
		for _, hs := range q.HazardStatements {
			hs_ids = hs + ","
		}
		values.Set("hazard_statements", hs_ids)
	}
	if len(q.PrecautionaryStatements) > 0 {
		ps_ids := ""
		for _, ps := range q.PrecautionaryStatements {
			ps_ids = ps + ","
		}
		values.Set("precautionary_statements", ps_ids)

	}
	if len(q.Symbols) > 0 {
		symbol_ids := ""
		for _, s := range q.Symbols {
			symbol_ids = s + ","
		}
		values.Set("symbols", symbol_ids)
	}
	if len(q.Tags) > 0 {
		tag_ids := ""
		for _, tag := range q.Tags {
			tag_ids = tag + ","
		}
		values.Set("tags", tag_ids)
	}

	if q.UnitType != "" {
		values.Set("unit_type", q.UnitType)
	}
	if q.StorageBatchNumber != "" {
		values.Set("storage_batch_number", q.StorageBatchNumber)
	}
	if q.StorageBarecode != "" {
		values.Set("storage_barecode", q.StorageBarecode)
	}
	if q.CustomNamePartOf != "" {
		values.Set("custom_name_part_of", q.CustomNamePartOf)
	}
	if q.IsCMR {
		values.Set("is_cmr", strconv.FormatBool(true))
	}
	if q.Borrowing {
		values.Set("borrowing", strconv.FormatBool(true))
	}
	if q.StorageToDestroy {
		values.Set("storage_to_destroy", strconv.FormatBool(true))
	}
	if q.StorageArchive != nil && *q.StorageArchive {
		values.Set("storage_archive", strconv.FormatBool(true))
	}
	if q.StorageArchive != nil && !*q.StorageArchive {
		values.Set("storage_archive", strconv.FormatBool(false))
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

	if q.ShowBio {
		values.Set("show_bio", strconv.FormatBool(true))
	}
	if q.ShowChem {
		values.Set("show_chem", strconv.FormatBool(true))
	}
	if q.ShowConsu {
		values.Set("show_consu", strconv.FormatBool(true))
	}

	return values.Encode()

}
