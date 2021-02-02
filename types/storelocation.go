package types

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/select2"
)

type StoreLocations struct {
	Rows  []*StoreLocation `json:"rows"`
	Total int              `json:"total"`
}

type StoreLocation struct {
	StoreLocationID       sql.NullInt64  `json:"storelocation_id"`
	StoreLocationName     sql.NullString `json:"storelocation_name"`
	StoreLocationCanStore sql.NullBool   `json:"storelocation_canstore"`
	StoreLocationColor    sql.NullString `json:"storelocation_color"`
	Entity                `json:"entity"`
	StoreLocation         *StoreLocation `json:"storelocation"`
	StoreLocationFullPath string         `json:"storelocation_fullpath"`

	Children []*StoreLocation `json:"children"`
	Stocks   []Stock          `json:"stock"`
}

type Stock struct {
	Total   float64 `json:"total"`
	Current float64 `json:"current"`
	Unit    Unit    `json:"unit"`
}

func (elems StoreLocations) IsExactMatch() bool {

	return false

}

func (s *StoreLocation) MarshalJSON() ([]byte, error) {

	type Copy StoreLocation
	return json.Marshal(struct {
		Id   int    `json:"id"`
		Text string `json:"text"`
		Copy
	}{
		Id:   s.GetSelect2Id(),
		Text: s.GetSelect2Text(),
		Copy: (Copy)(*s),
	})

}

func (StoreLocations) FromJsJSONValue(jsvalue js.Value) select2.Select2ResultAble {

	var (
		storelocations StoreLocations
		err            error
	)

	jsStoreLocationsString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsStoreLocationsString), &storelocations); err != nil {
		fmt.Println(err)
	}

	return storelocations

}

func (s StoreLocations) GetRows() []select2.Select2ItemAble {

	var select2ItemAble []select2.Select2ItemAble = make([]select2.Select2ItemAble, len(s.Rows))

	for i, row := range s.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (s StoreLocations) GetTotal() int {

	return s.Total

}

func (s StoreLocation) FromJsJSONValue(jsvalue js.Value) select2.Select2ItemAble {

	var (
		storelocation StoreLocation
		err           error
	)

	jsStoreLocationString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsStoreLocationString), &storelocation); err != nil {
		fmt.Println(err)
	}

	return storelocation

}

func (s StoreLocation) GetSelect2Id() int {

	return int(s.StoreLocationID.Int64)

}

func (s StoreLocation) GetSelect2Text() string {

	return s.StoreLocationFullPath

}
