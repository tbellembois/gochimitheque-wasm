package types

import (
	"reflect"
	"strings"
)

func StructToMap(item interface{}) map[string]interface{} {

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
				res[tag] = StructToMap(field)
			} else if v.Field(i).Type.Kind() == reflect.Map {
				switch v.Field(i).Type.Elem() {
				case reflect.TypeOf(ValidateRule{}):
					m := field.(map[string]ValidateRule)

					res2 := map[string]interface{}{}
					for k, v := range m {
						res2[k] = StructToMap(v)
					}
					res[tag] = res2

				case reflect.TypeOf(ValidateMessage{}):
					m := field.(map[string]ValidateMessage)

					res2 := map[string]interface{}{}
					for k, v := range m {
						res2[k] = StructToMap(v)
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
