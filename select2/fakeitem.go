//go:build go1.24 && js && wasm

package select2

import (
	"encoding/json"
	"fmt"
	"syscall/js"
)

type FakeItem struct {
	Id   int64  `json:"id"`
	Text string `json:"text"`
}

func (f FakeItem) FromJsJSONValue(jsvalue js.Value) Select2ItemAble {

	var (
		fakeItem FakeItem
		err      error
	)

	jsFakeItemString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsFakeItemString), &fakeItem); err != nil {
		fmt.Println(err)
	}

	return fakeItem

}

func (f FakeItem) GetSelect2Id() int64 {

	return f.Id

}

func (f FakeItem) GetSelect2Text() string {

	return f.Text

}
