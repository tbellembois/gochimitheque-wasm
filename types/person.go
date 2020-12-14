package types

import (
	"encoding/json"
	"fmt"
	"syscall/js"
)

type People struct {
	Rows  []*Person `json:"rows"`
	Total int       `json:"total"`
}
type Person struct {
	PersonId       int           `json:"person_id"`
	PersonEmail    string        `json:"person_email"`
	PersonPassword string        `json:"person_password"`
	Entities       []*Entity     `json:"entities"`
	Permissions    []*Permission `json:"permissions"`
	CaptchaText    string        `json:"captcha_text"`
	CaptchaUID     string        `json:"captcha_uid"`
}

func (elems People) IsExactMatch() bool {

	return false

}

func (p *Person) MarshalJSON() ([]byte, error) {
	type Copy Person
	return json.Marshal(&struct {
		Id   int    `json:"id"`
		Text string `json:"text"`

		*Copy
	}{
		Id:   p.GetSelect2Id(),
		Text: p.GetSelect2Text(),
		Copy: (*Copy)(p),
	})
}

func (People) FromJsJSONValue(jsvalue js.Value) Select2ResultAble {

	var (
		people People
		err    error
	)

	jsPeopleString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsPeopleString), &people); err != nil {
		fmt.Println(err)
	}

	return people

}

func (p Person) FromJsJSONValue(jsvalue js.Value) Select2ItemAble {

	var (
		person Person
		err    error
	)

	jsPersonString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsPersonString), &person); err != nil {
		fmt.Println(err)
	}

	return person

}

func (p People) GetRows() []Select2ItemAble {

	var select2ItemAble []Select2ItemAble = make([]Select2ItemAble, len(p.Rows))

	for i, row := range p.Rows {
		select2ItemAble[i] = row
	}

	return select2ItemAble

}

func (p People) GetTotal() int {

	return p.Total

}

func (p Person) GetSelect2Id() int {

	return p.PersonId

}

func (p Person) GetSelect2Text() string {

	return p.PersonEmail

}
