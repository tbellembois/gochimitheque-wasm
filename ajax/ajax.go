//go:build go1.24 && js && wasm

package ajax

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"syscall/js"
)

// Ajax represents a JQuery Ajax method.
type Ajax struct {
	URL    string
	Method string
	Data   []byte
	Done   AjaxDone
	Fail   AjaxFail
}

// Response contains the data retrieved
// from the query above.
type Response struct {
	Rows  js.Value `json:"rows"`
	Total js.Value `json:"total"`
}

type AjaxDone func(data js.Value)
type AjaxFail func(jqXHR js.Value)

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

		window := js.Global()

		var keycloak js.Value
		keycloak = window.Get("keycloak")
		token := keycloak.Get("token").String()

		req := &http.Request{
			Method: ajax.Method,
			URL:    reqURL,
			Header: map[string][]string{
				"Content-Type":  {"application/json; charset=UTF-8"},
				"Authorization": {"Bearer " + token},
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
