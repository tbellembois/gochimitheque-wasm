package jsutils

import (
	"fmt"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/ajax"
	"github.com/tbellembois/gochimitheque-wasm/globals"
	. "github.com/tbellembois/gochimitheque-wasm/globals"
)

func IsLDAPUser(email string, done, fail func()) {

	go func() {

		url := ApplicationProxyPath + "people/isldap/" + email

		ajaxDone := func(data js.Value) {

			isLDAP := js.Global().Get("JSON").Call("parse", data).Bool()

			if isLDAP {
				done()
			} else {
				fail()
			}
		}
		ajaxFail := func(js.Value) {
			fail()
		}
		ajax := ajax.Ajax{
			URL:    url,
			Method: "GET",
			Done:   ajaxDone,
			Fail:   ajaxFail,
		}

		ajax.Send()

	}()

}

func HasPermission(item, id, method string, done, fail func()) {

	go func() {

		cacheKey := fmt.Sprintf("%s:%s:%s", item, id, method)
		cachedPermission := globals.LocalStorage.GetItem(cacheKey)

		if !DisableCache && cachedPermission != "" {
			if cachedPermission == "true" {
				done()
			} else {
				fail()
			}
		} else {
			var url string
			if id != "" {
				url = ApplicationProxyPath + "f/" + item + "/" + id
			} else {
				url = ApplicationProxyPath + "f/" + item
			}

			ajaxDone := func(js.Value) {
				globals.LocalStorage.SetItem(cacheKey, "true")
				done()
			}
			ajaxFail := func(js.Value) {
				globals.LocalStorage.SetItem(cacheKey, "false")
				fail()
			}
			ajax := ajax.Ajax{
				URL:    url,
				Method: method,
				Done:   ajaxDone,
				Fail:   ajaxFail,
			}

			ajax.Send()

		}
	}()

}
