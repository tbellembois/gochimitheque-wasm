package types

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"syscall/js"
)

// Select2ResultAble represents a type that can
// fill in select2 data as defined
// https://select2.org/data-sources/ajax#transforming-response-data
type Select2ResultAble interface {
	FromJsJSONValue(js.Value) Select2ResultAble
	GetRows() []Select2ItemAble
	GetTotal() int
	IsExactMatch() bool
}

type Select2ItemAble interface {
	GetSelect2Id() int
	GetSelect2Text() string
	FromJsJSONValue(js.Value) Select2ItemAble
}

// Select2Config is a select2 parameters struct
// as defined https://select2.org/configuration
type Select2Config struct {
	Placeholder       string      `json:"placeholder"`
	AllowClear        bool        `json:"allowClear"`
	Tags              bool        `json:"tags"`
	Ajax              Select2Ajax `json:"ajax"`
	TemplateResult    interface{} `json:"templateResult,omitempty"`
	TemplateSelection interface{} `json:"templateSelection,omitempty"`
	CreateTag         interface{} `json:"createTag,omitempty"`
}

// Select2Ajax is a select2 ajax request
// as defined https://select2.org/data-sources/ajax
type Select2Ajax struct {
	URL            string      `json:"url"`
	DataType       string      `json:"datatype"`
	Data           interface{} `json:"data"`
	ProcessResults interface{} `json:"processResults"`
}

// Select2Data is a select2 data struct
// as defined https://select2.org/data-sources/formats
type Select2Data struct {
	// Results    []Select2Item     `json:"results"`
	Results    []Select2ItemAble `json:"results"`
	Pagination Select2Pagination `json:"pagination"`
}

type Select2Pagination struct {
	More bool `json:"more"`
}

type Select2Item struct {
	Id   string `json:"id"`
	Text string `json:"text"`
}

func (s *Select2Item) UnmarshalJSON(data []byte) error {

	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch reflect.ValueOf(v["id"]).Kind() {
	case reflect.String:
		s.Id = v["id"].(string)
	case reflect.Float64:
		s.Id = strconv.Itoa(int(v["id"].(float64)))
	}

	s.Text = v["text"].(string)

	return nil

}

func (s Select2Item) IDIsDigit() bool {

	r := regexp.MustCompile("^[0-9]+$")

	return r.MatchString(s.Id)

}

func (s Select2Item) IsEmpty() bool {

	return (s.Id == "0" && s.Text == "")

}

func (jq Jquery) Select2(config Select2Config) {

	configMap := StructToMap(config)
	jq.Object.Call("select2", configMap)

	// Can not use:
	// jq.Object.Call("select2", config.ToJsValue())
	// because Select2Config contains functions.

}

func (jq Jquery) Select2Data() []Select2Item {

	var (
		select2Items []Select2Item
		err          error
	)

	select2Data := jq.Object.Call("select2", "data")

	jsSelect2ItemsString := js.Global().Get("JSON").Call("stringify", select2Data).String()
	if err = json.Unmarshal([]byte(jsSelect2ItemsString), &select2Items); err != nil {
		fmt.Println(err)
	}

	return select2Items

}

func (jq Jquery) Select2IsInitialized() bool {

	return jq.HasClass("select2-hidden-accessible")

}

func (jq Jquery) Select2AppendOption(o interface{}) {

	jq.Object.Call("append", js.ValueOf(o)).Call("trigger", "change")

}

func (jq Jquery) Select2Clear() {

	jq.SetVal(nil).Object.Call("trigger", "change")
	jq.Find("option").Remove()

}

func (s Select2Data) ToJsValue() js.Value {

	var (
		marshalJson []byte
		err         error
	)

	if marshalJson, err = json.Marshal(s); err != nil {
		fmt.Println(err)
		return js.Null()
	}

	return js.Global().Get("JSON").Call("parse", string(marshalJson))

}

func (s Select2Item) ToJsValue() js.Value {

	var (
		marshalJson []byte
		err         error
	)

	if marshalJson, err = json.Marshal(s); err != nil {
		fmt.Println(err)
		return js.Null()
	}

	return js.Global().Get("JSON").Call("parse", string(marshalJson))

}

func Select2ItemFromJsJSONValue(jsvalue js.Value) Select2Item {

	var (
		select2Item Select2Item
		err         error
	)

	jsSelect2ItemString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsSelect2ItemString), &select2Item); err != nil {
		fmt.Println(err)
	}

	return select2Item

}

// Select2GenericAjaxData returns the most common request
// parameters for select2 as defined
// https://select2.org/data-sources/ajax#request-parameters
func Select2GenericAjaxData(this js.Value, args []js.Value) interface{} {

	params := args[0]

	search := ""
	if params.Get("term").Truthy() {
		search = params.Get("term").String()
	}
	page := 1
	if params.Get("page").Truthy() {
		page = params.Get("page").Int()
	}
	offset := (page - 1) * 10
	limit := 10

	if offset < 0 {
		offset = 0
	}

	return QueryFilter{
		Search: search,
		Offset: offset,
		Page:   page,
		Limit:  limit,
	}.ToJsValue()

}

func Select2GenericAjaxProcessResults(select2ResultAble Select2ResultAble) func(this js.Value, args []js.Value) interface{} {

	return func(this js.Value, args []js.Value) interface{} {

		data := args[0]
		params := args[1]
		page := 1
		if params.Get("page").Truthy() {
			page = params.Get("page").Int()
		}

		objects := select2ResultAble.FromJsJSONValue(data)

		rows := objects.GetRows()
		if len(rows) > 0 {
			row := rows[0]
			rowType := reflect.TypeOf(row).Elem()
			rowTypeName := rowType.Name()

			if objects.IsExactMatch() {
				Jq(fmt.Sprintf("input#exactMatch%s", rowTypeName)).SetVal(true)
			} else {
				Jq(fmt.Sprintf("input#exactMatch%s", rowTypeName)).SetVal(false)
			}
		}

		var select2ItemAbles []Select2ItemAble
		for _, object := range objects.GetRows() {
			select2ItemAbles = append(select2ItemAbles, object)
		}

		// Needed to avoid a Jquery exception. Don't know why.
		if len(select2ItemAbles) == 0 {
			select2ItemAbles = append(select2ItemAbles, FakeItem{})
		}

		select2Data := Select2Data{
			Results: select2ItemAbles,
			Pagination: Select2Pagination{
				More: (page * 10) < objects.GetTotal(),
			},
		}

		return select2Data.ToJsValue()

	}
}

func Select2GenericTemplateResults(select2ItemAble Select2ItemAble) func(this js.Value, args []js.Value) interface{} {

	return func(this js.Value, args []js.Value) interface{} {

		data := args[0]

		object := select2ItemAble.FromJsJSONValue(data)

		return object.GetSelect2Text()

	}

}

func Select2GenericCreateTag(select2ItemAble Select2ItemAble) func(this js.Value, args []js.Value) interface{} {

	return func(this js.Value, args []js.Value) interface{} {

		params := args[0]

		objectName := reflect.TypeOf(select2ItemAble).Name()

		if Jq(fmt.Sprintf("input#exactMatch%s", objectName)).GetVal().String() == "true" {
			return nil
		}

		return Select2Item{
			Id:   params.Get("term").String(),
			Text: params.Get("term").String(),
		}.ToJsValue()

	}

}
