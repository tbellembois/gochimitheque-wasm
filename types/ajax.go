package types

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/locales"
	"github.com/tbellembois/gochimitheque-wasm/widgets"
)

// TODO: move this
var (
	BSTableQueryFilter SafeQueryFilter
	// CurrentProduct is the current viewed or edited product, or the product of
	// the listed storages.
	CurrentProduct Product
	// CurrentStorage is the current viewed or edited storage.
	CurrentStorage            Storage
	DBPrecautionaryStatements []PrecautionaryStatement // for magic selector
	DBHazardStatements        []HazardStatement        // for magic selector
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

func (q QueryFilter) DisplayFilter() {

	var isFilter bool

	Jq("#filter-item").SetHtml("")

	if q.Product != "" {
		isFilter = true
		Jq("#filter-item").Append(widgets.FilterItem(locales.Translate("storage_product_table_header", globals.HTTPHeaderAcceptLanguage), q.ProductFilterLabel))
	}
	if q.Storage != "" {
		isFilter = true
		Jq("#filter-item").Append(widgets.FilterItem(locales.Translate("storage", globals.HTTPHeaderAcceptLanguage), q.StorageFilterLabel))
	}

	if q.CustomNamePartOf != "" {
		isFilter = true
		Jq("#filter-item").Append(widgets.FilterItem(locales.Translate("s_custom_name_part_of", globals.HTTPHeaderAcceptLanguage), q.CustomNamePartOf))
	}
	if q.CasNumber != "" {
		isFilter = true
		Jq("#filter-item").Append(widgets.FilterItem(locales.Translate("s_casnumber", globals.HTTPHeaderAcceptLanguage), q.CasNumberFilterLabel))
	}
	if q.EmpiricalFormula != "" {
		isFilter = true
		Jq("#filter-item").Append(widgets.FilterItem(locales.Translate("s_empiricalformula", globals.HTTPHeaderAcceptLanguage), q.EmpiricalFormulaFilterLabel))
	}
	if q.StorageBarecode != "" {
		isFilter = true
		Jq("#filter-item").Append(widgets.FilterItem(locales.Translate("s_storage_barecode", globals.HTTPHeaderAcceptLanguage), q.StorageBarecode))
	}
	if q.StoreLocation != "" {
		isFilter = true
		Jq("#filter-item").Append(widgets.FilterItem(locales.Translate("s_storelocation", globals.HTTPHeaderAcceptLanguage), q.StoreLocationFilterLabel))
	}
	if q.Name != "" {
		isFilter = true
		Jq("#filter-item").Append(widgets.FilterItem(locales.Translate("s_name", globals.HTTPHeaderAcceptLanguage), q.NameFilterLabel))
	}

	if q.ProductBookmark {
		isFilter = true
		Jq("#filter-item").Append(widgets.FilterItem("", locales.Translate("menu_bookmark", globals.HTTPHeaderAcceptLanguage)))
	}
	if q.StorageArchive {
		isFilter = true
		Jq("#filter-item").Append(widgets.FilterItem("", locales.Translate("archives", globals.HTTPHeaderAcceptLanguage)))
	}
	if q.StorageHistory {
		isFilter = true
		Jq("#filter-item").Append(widgets.FilterItem("", locales.Translate("storage_history", globals.HTTPHeaderAcceptLanguage)))
	}

	if !isFilter {
		Jq("#filter-item").Append(widgets.FilterItem("", locales.Translate("no_filter", globals.HTTPHeaderAcceptLanguage)))
	}

}

// SafeQueryFilter is a concurent
// query filter used in  Ajax requests.
type SafeQueryFilter struct {
	lock   sync.Mutex
	locked bool
	QueryFilter
}

func (s *SafeQueryFilter) Lock() {
	s.locked = true
	s.lock.Lock()
}

func (s *SafeQueryFilter) Unlock() {
	if s.locked {
		s.locked = false
		s.lock.Unlock()
	}
}

func (s *SafeQueryFilter) Clean() {
	s.QueryFilter = QueryFilter{}
}

func (s *SafeQueryFilter) CleanExceptProduct() {

	backupProduct := s.QueryFilter.Product
	backupProductFilterLabel := s.QueryFilter.ProductFilterLabel
	s.Clean()
	s.QueryFilter.Product = backupProduct
	s.QueryFilter.ProductFilterLabel = backupProductFilterLabel

}

// Response contains the data retrieved
// from the query above.
type Response struct {
	Rows  js.Value `json:"rows"`
	Total js.Value `json:"total"`
}

// Ajax represents a JQuery Ajax method.
type Ajax struct {
	URL    string
	Method string
	Data   []byte
	Done   AjaxDone
	Fail   AjaxFail
}

type AjaxDone func(data js.Value)
type AjaxFail func(jqXHR js.Value)

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

func (ajax Ajax) Send() {

	go func() {

		var (
			err    error
			data   []byte
			reqURL *url.URL
			res    *http.Response
		)
		if reqURL, err = url.Parse(ajax.URL); err != nil {
			fmt.Println(err)
			return
		}

		req := &http.Request{
			Method: ajax.Method,
			URL:    reqURL,
			Header: map[string][]string{
				"Content-Type": {"application/json; charset=UTF-8"},
			},
		}

		if len(ajax.Data) > 0 {
			req.Body = ioutil.NopCloser(strings.NewReader(string(ajax.Data)))
		}

		if res, err = http.DefaultClient.Do(req); err != nil {
			fmt.Println(err)
			return
		}

		if data, err = ioutil.ReadAll(res.Body); err != nil {
			fmt.Println(err)
			return
		}
		res.Body.Close()

		//jsResponse := js.Global().Get("JSON").Call("stringify", string(data))
		jsResponse := js.ValueOf(string(data))

		if res.StatusCode == 200 {
			ajax.Done(jsResponse)
		} else {
			if ajax.Fail != nil {
				ajax.Fail(jsResponse)
			}
		}

	}()

}
