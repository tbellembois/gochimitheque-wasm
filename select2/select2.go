//go:build go1.24 && js && wasm

package select2

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/ajax"
	"github.com/tbellembois/gochimitheque-wasm/jquery"
)

func structToMap(item interface{}) map[string]interface{} {

	res := map[string]interface{}{}
	if item == nil {
		return res
	}
	v := reflect.TypeOf(item)
	reflectValue := reflect.ValueOf(item)
	reflectValue = reflect.Indirect(reflectValue)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	var (
		omitempty bool
	)

	for i := 0; i < v.NumField(); i++ {

		tag := v.Field(i).Tag.Get("json")
		if strings.Contains(tag, "omitempty") {
			omitempty = true
			tag = strings.Replace(tag, ",omitempty", "", 1)
		} else {
			omitempty = false
		}

		field := reflectValue.Field(i).Interface()

		if tag != "" && tag != "-" {
			if v.Field(i).Type.Kind() == reflect.String {
				if !omitempty && len(field.(string)) > 0 {
					res[tag] = field
				}
			} else if v.Field(i).Type.Kind() == reflect.Bool {
				res[tag] = field
			} else if v.Field(i).Type.Kind() == reflect.Struct {
				res[tag] = structToMap(field)
			} else if v.Field(i).Type.Kind() == reflect.Map {

				m := field.(map[string]interface{})

				res2 := map[string]interface{}{}
				for k, v := range m {
					res2[k] = v
				}
				res[tag] = res2
			} else if v.Field(i).Type.Kind() == reflect.Interface {
				if !reflectValue.Field(i).IsNil() {
					res[tag] = field
				}
			} else {
				res[tag] = field
			}
		}
	}
	return res
}

type Select2 struct {
	jquery.Jquery
	config *Select2Config
}

// select2.Select2ResultAble represents a type that can
// fill in select2 data as defined
// https://select2.org/data-sources/ajax#transforming-response-data
type Select2ResultAble interface {
	FromJsJSONValue(js.Value) Select2ResultAble
	GetRows() []Select2ItemAble
	GetRowConcreteTypeName() string
	GetTotal() int
	IsExactMatch() bool
}

type Select2ItemAble interface {
	GetSelect2Id() int64
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

func NewSelect2(jq jquery.Jquery, config *Select2Config) Select2 {

	// configMap := StructToMap(config)
	// jq.Object.Call("select2", configMap)

	// Can not use:
	// jq.Object.Call("select2", config.ToJsValue())
	// because Select2Config contains functions.

	return Select2{Jquery: jq, config: config}

}

func (s Select2) Select2ify() {

	configMap := structToMap(s.config)
	s.Jquery.Object.Call("select2", configMap)

}

func (s Select2) Select2Data() []Select2Item {

	var (
		select2Items []Select2Item
		err          error
	)

	select2Data := s.Jquery.Object.Call("select2", "data")

	jsSelect2ItemsString := js.Global().Get("JSON").Call("stringify", select2Data).String()
	if err = json.Unmarshal([]byte(jsSelect2ItemsString), &select2Items); err != nil {
		fmt.Println(err)
	}

	return select2Items

}

func (s Select2) Select2IsInitialized() bool {

	return s.Jquery.HasClass("select2-hidden-accessible")

}

func (s Select2) Select2AppendOption(o interface{}) {

	s.Jquery.Object.Call("append", js.ValueOf(o)).Call("trigger", "change")

}

func (s Select2) Select2Clear() {

	s.Jquery.SetVal(nil).Object.Call("trigger", "change")
	s.Jquery.Find("option").Remove()

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

	return ajax.QueryFilter{
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
		rowTypeName := objects.GetRowConcreteTypeName()

		rows := objects.GetRows()
		if len(rows) > 0 {

			// row := rows[0]
			// rowType := reflect.TypeOf(row).Elem()
			// rowTypeName := rowType.Name()

			if objects.IsExactMatch() {
				jquery.Jq(fmt.Sprintf("input#exactMatch%s", rowTypeName)).SetVal(true)
			} else {
				jquery.Jq(fmt.Sprintf("input#exactMatch%s", rowTypeName)).SetVal(false)
			}

		} else {

			jquery.Jq(fmt.Sprintf("input#exactMatch%s", rowTypeName)).SetVal(false)

		}

		var select2ItemAbles []Select2ItemAble
		select2ItemAbles = append(select2ItemAbles, objects.GetRows()...)

		// Needed to avoid a Jquery exception. Don't know why.
		if len(select2ItemAbles) == 0 {
			select2ItemAbles = append(select2ItemAbles, FakeItem{})
		}

		var pagination Select2Pagination
		if select2ResultAble.GetRowConcreteTypeName() != "" {
			pagination.More = (page * 10) < objects.GetTotal()
		}

		select2Data := Select2Data{
			Results:    select2ItemAbles,
			Pagination: pagination,
		}

		// js.Global().Get("console").Call("log", fmt.Sprintf("%#v", pagination))

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

		if jquery.Jq(fmt.Sprintf("input#exactMatch%s", objectName)).GetVal().String() == "true" {
			return js.Null()
		}

		return Select2Item{
			Id:   params.Get("term").String(),
			Text: params.Get("term").String(),
		}.ToJsValue()

	}

}
