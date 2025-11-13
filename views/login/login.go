//go:build go1.24 && js && wasm

package login

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/ajax"
	. "github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/jquery"
	"github.com/tbellembois/gochimitheque-wasm/jsutils"
	"github.com/tbellembois/gochimitheque-wasm/locales"
	. "github.com/tbellembois/gochimitheque-wasm/types"
	"github.com/tbellembois/gochimitheque-wasm/views/menu"
	"github.com/tbellembois/gochimitheque-wasm/views/product"
	"github.com/tbellembois/gochimitheque-wasm/views/search"
)

func Login_listCallback(this js.Value, args []js.Value) interface{} {

	// Call the ping handler to see if if are already authenticated.
	ajax.Ajax{
		URL:    fmt.Sprintf("%sping", ApplicationProxyPath),
		Method: "get",
		Done: func(data js.Value) {
			AfterLogin_listCallback(js.Null(), nil)
		},
		Fail: func(jqXHR js.Value) {
			// Just do nothing.
			if jqXHR.String() == "forbidden" {
				jsutils.DisplayErrorMessage(locales.Translate("login_forbidden", HTTPHeaderAcceptLanguage))
			}
		},
	}.Send()

	return nil
}

func AfterLogin_listCallback(this js.Value, args []js.Value) interface{} {

	productCallbackWrapper := func(args ...interface{}) {
		product.Product_listCallback(js.Null(), nil)
	}

	jsutils.LoadContent("div#menu", "menu", fmt.Sprintf("%smenu", ApplicationProxyPath), menu.ShowIfAuthorizedMenuItems)
	jsutils.LoadContent("div#searchbar", "product", fmt.Sprintf("%ssearch", ApplicationProxyPath), search.Search_listCallback)
	jsutils.LoadContent("div#content", "product", fmt.Sprintf("%sv/products", ApplicationProxyPath), productCallbackWrapper)

	ajax.Ajax{
		URL:    fmt.Sprintf("%suserinfo", ApplicationProxyPath),
		Method: "get",
		Done: func(data js.Value) {

			var (
				person *Person
				err    error
			)

			if err = json.Unmarshal([]byte(data.String()), &person); err != nil {
				fmt.Println(err)
			}

			ConnectedUserEmail = person.PersonEmail
			ConnectedUserID = *person.PersonID

			jquery.Jq("#logged").SetHtml(ConnectedUserEmail)

			product.Product_listCallback(js.Null(), nil)
		},
		Fail: func(jqXHR js.Value) {

			jsutils.DisplayGenericErrorMessage()

		},
	}.Send()

	return nil

}
