//go:build go1.24 && js && wasm

package validate

import (
	"reflect"
	"strings"

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
				switch v.Field(i).Type.Elem() {
				case reflect.TypeOf(ValidateRule{}):
					m := field.(map[string]ValidateRule)

					res2 := map[string]interface{}{}
					for k, v := range m {
						res2[k] = structToMap(v)
					}
					res[tag] = res2

				case reflect.TypeOf(ValidateMessage{}):
					m := field.(map[string]ValidateMessage)

					res2 := map[string]interface{}{}
					for k, v := range m {
						res2[k] = structToMap(v)
					}
					res[tag] = res2

				default:
					m := field.(map[string]interface{})

					res2 := map[string]interface{}{}
					for k, v := range m {
						res2[k] = v
					}
					res[tag] = res2
				}

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

type Validate struct {
	jquery.Jquery
	config *ValidateConfig
}

// ValidateConfig is a jQuery validate
// parameters struct as defined
// https://jqueryvalidation.org/validate/
type ValidateConfig struct {
	Ignore     string                     `json:"ignore"`
	ErrorClass string                     `json:"errorClass"`
	Rules      map[string]ValidateRule    `json:"rules"`
	Messages   map[string]ValidateMessage `json:"messages"`
}
type ValidateRule struct {
	Email    bool           `json:"email,omitempty"`
	EqualTo  string         `json:"equalTo,omitempty"`
	Required interface{}    `json:"required"`
	Remote   ValidateRemote `json:"remote"`
}
type ValidateRemote struct {
	URL        string                 `json:"url"`
	Type       string                 `json:"type"`
	BeforeSend interface{}            `json:"beforeSend"`
	Data       map[string]interface{} `json:"data"`
}

type ValidateMessage struct {
	Email     string `json:"email,omitempty"`
	EqualTo   string `json:"equalTo,omitempty"`
	Required  string `json:"required"`
	MinLength string `json:"minlength"`
}

func NewValidate(jq jquery.Jquery, config *ValidateConfig) Validate {

	return Validate{Jquery: jq, config: config}
}

func (v Validate) ValidateAdd(rule ValidateRule) {

	v.Jquery.Object.Call("rules", "add", structToMap(rule))

}

func (v Validate) ValidateAddRequired() {

	v.Jquery.Object.Call("rules", "add", "required")

}

func (v Validate) ValidateRemoveRequired() {

	v.Jquery.Object.Call("rules", "remove", "required")

}

func (v Validate) Validate() {

	configMap := structToMap(v.config)
	v.Jquery.Object.Call("validate", configMap)

}

func (v Validate) Valid() bool {

	return v.Jquery.Object.Call("valid").Bool()

}
