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
	ShowBio                            bool   `json:"showbio,omitempty"`
	ShowChem                           bool   `json:"showchem,omitempty"`
	ShowConsu                          bool   `json:"showconsu,omitempty"`
	Product                            string `json:"product,omitempty"`
	ProductFilterLabel                 string
	ProductBookmark                    bool   `json:"bookmark,omitempty"`
	ProducerRef                        string `json:"producerref,omitempty"`
	ProducerRefFilterLabel             string
	Storage                            string `json:"storage,omitempty"`
	StorageFilterLabel                 string
	Storages                           []int `json:"ids,omitempty"`
	StoragesFilterLabel                string
	UnitType                           string `json:"unit_type,omitempty"`
	Supplier                           string `json:"supplier,omitempty"`
	Producer                           string `json:"producer,omitempty"`
	StoreLocation                      string `json:"storelocation,omitempty"`
	StoreLocationCanStore              bool   `json:"storelocation_canstore,omitempty"`
	StoreLocationFilterLabel           string
	Entity                             string `json:"entity,omitempty"`
	EntityFilterLabel                  string
	Name                               string `json:"name,omitempty"`
	NameFilterLabel                    string
	CasNumber                          string `json:"casnumber,omitempty"`
	CasNumberFilterLabel               string
	EmpiricalFormula                   string `json:"empiricalformula,omitempty"`
	EmpiricalFormulaFilterLabel        string
	SignalWord                         string `json:"signalword,omitempty"`
	SignalWordFilterLabel              string
	HazardStatements                   []string `json:"hazardstatements,omitempty"`
	HazardStatementsFilterLabel        string
	PrecautionaryStatements            []string `json:"precautionarystatements,omitempty"`
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
	CasNumberCMR                       bool `json:"casnumber_cmr,omitempty"`
	CasNumberCMRFilterLabel            string
	Borrowing                          bool `json:"borrowing,omitempty"`
	BorrowingFilterLabel               string
	StorageToDestroy                   bool `json:"storage_to_destroy,omitempty"`
	StorageToDestroyFilterLabel        string
	StorageArchive                     bool   `json:"storage_archive,omitempty"`
	StorageHistory                     bool   `json:"storage_history,omitempty"`
	StorageBatchNumber                 string `json:"storage_batchnumber,omitempty"`
	StorageBatchNumberFilterLabel      string
	Export                             bool `json:"export,omitempty"`

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
		values.Set("storelocation", q.StoreLocation)
	}
	if q.ProducerRef != "" {
		values.Set("producerref", q.ProducerRef)
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
	if q.Category != "" {
		values.Set("category", q.Category)
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
	if len(q.Tags) > 0 {
		for _, tag := range q.Tags {
			values.Set("tags[]", tag)
		}
	}

	if q.UnitType != "" {
		values.Set("unit_type", q.UnitType)
	}
	if q.StorageBatchNumber != "" {
		values.Set("storage_batchnumber", q.StorageBatchNumber)
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

	if q.ShowBio {
		values.Set("showbio", strconv.FormatBool(true))
	}
	if q.ShowChem {
		values.Set("showchem", strconv.FormatBool(true))
	}
	if q.ShowConsu {
		values.Set("showconsu", strconv.FormatBool(true))
	}

	return values.Encode()

}
