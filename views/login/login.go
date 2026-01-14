//go:build go1.24 && js && wasm

package login

import (
	"encoding/json"
	"fmt"
	"syscall/js"
	"time"

	"github.com/tbellembois/gochimitheque-wasm/ajax"
	. "github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/jquery"
	"github.com/tbellembois/gochimitheque-wasm/jsutils"
	. "github.com/tbellembois/gochimitheque-wasm/types"
	"github.com/tbellembois/gochimitheque-wasm/views/menu"
	"github.com/tbellembois/gochimitheque-wasm/views/product"
	"github.com/tbellembois/gochimitheque-wasm/views/search"
)

func Login_listCallback(this js.Value, args []js.Value) interface{} {

	js.Global().Get("console").Call("log", fmt.Sprintf("Login_listCallback"))
	js.Global().Get("console").Call("log", fmt.Sprintf("%#v", args))

	// Get the global object
	window := js.Global()

	// Wait until keycloak is defined
	var keycloak js.Value
	for {
		keycloak = window.Get("keycloak")
		if !keycloak.IsUndefined() {
			break
		}
		fmt.Println("Waiting for Keycloak...")
		time.Sleep(1000 * time.Millisecond)
	}

	fmt.Println("Keycloak is ready!")

	// Now you can access token or call functions
	token := keycloak.Get("token").String()
	fmt.Println("Keycloak token:", token)

	AfterLogin_listCallback(js.Null(), nil)

	// js.Global().Get("window").Get("location").Set("href", fmt.Sprintf("%sauthenticated", BackProxyPath))

	// Call the ping handler to see if if are already authenticated.
	// ajax.Ajax{
	// 	URL:    fmt.Sprintf("%sauthenticated", BackProxyPath),
	// 	Method: "get",
	// 	Done: func(data js.Value) {
	// 		AfterLogin_listCallback(js.Null(), nil)
	// 	},
	// 	Fail: func(jqXHR js.Value) {
	// 		// Just do nothing.
	// 		if jqXHR.String() == "forbidden" {
	// 			jsutils.DisplayErrorMessage(locales.Translate("login_forbidden", HTTPHeaderAcceptLanguage))
	// 		}
	// 	},
	// }.Send()

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
		URL:    fmt.Sprintf("%sgetconnecteduser", BackProxyPath),
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
