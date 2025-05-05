//go:build go1.24 && js && wasm

package menu

import (
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/jquery"
	"github.com/tbellembois/gochimitheque-wasm/jsutils"
	"github.com/tbellembois/gochimitheque-wasm/views/about"
	"github.com/tbellembois/gochimitheque-wasm/views/entity"
	"github.com/tbellembois/gochimitheque-wasm/views/person"
	"github.com/tbellembois/gochimitheque-wasm/views/product"
	"github.com/tbellembois/gochimitheque-wasm/views/storelocation"
	"github.com/tbellembois/gochimitheque-wasm/views/welcomeannounce"
)

// LoadContentWrapper is only called in the menu.jade template
// renderer by the server part.
func LoadContentWrapper(this js.Value, args []js.Value) interface{} {

	productCallbackWrapper := func(args ...interface{}) {
		product.Product_listCallback(js.Null(), nil)
	}
	productCreateCallbackWrapper := func(args ...interface{}) {
		product.Product_createCallback(js.Null(), nil)
	}
	productPubchemCallbackWrapper := func(args ...interface{}) {
		product.Product_pubchemCallback(js.Null(), nil)
	}
	entityCallbackWrapper := func(args ...interface{}) {
		entity.Entity_listCallback(js.Null(), nil)
	}
	entityCreateCallbackWrapper := func(args ...interface{}) {
		entity.Entity_createCallBack(js.Null(), nil)
	}
	store_locationCallbackWrapper := func(args ...interface{}) {
		storelocation.StoreLocation_listCallback(js.Null(), nil)
	}
	store_locationCreateCallbackWrapper := func(args ...interface{}) {
		storelocation.StoreLocation_createCallBack(js.Null(), nil)
	}
	personCallbackWrapper := func(args ...interface{}) {
		person.Person_listCallback(js.Null(), nil)
	}
	personCreateCallbackWrapper := func(args ...interface{}) {
		person.Person_createCallBack(js.Null(), nil)
	}
	aboutListCallbackWrapper := func(args ...interface{}) {
		about.About_listCallback(js.Null(), nil)
	}
	welcomeannounceListCallbackWrapper := func(args ...interface{}) {
		welcomeannounce.WelcomeAnnounce_listCallback(js.Null(), nil)
	}

	var callbackFunc func(args ...interface{})

	switch args[2].String() {
	case "Product_list":
		callbackFunc = productCallbackWrapper
	case "Product_create":
		callbackFunc = productCreateCallbackWrapper
	case "Product_pubchem":
		callbackFunc = productPubchemCallbackWrapper
	case "Entity_list":
		callbackFunc = entityCallbackWrapper
	case "Entity_create":
		callbackFunc = entityCreateCallbackWrapper
	case "StoreLocation_list":
		callbackFunc = store_locationCallbackWrapper
	case "StoreLocation_create":
		callbackFunc = store_locationCreateCallbackWrapper
	case "Person_list":
		callbackFunc = personCallbackWrapper
	case "Person_create":
		callbackFunc = personCreateCallbackWrapper
	case "WelcomeAnnounce_list":
		callbackFunc = welcomeannounceListCallbackWrapper
	case "About_list":
		callbackFunc = aboutListCallbackWrapper
	}

	jsutils.LoadContent("div#content", args[0].String(), args[1].String(), callbackFunc, nil)

	return nil

}

func ShowIfAuthorizedMenuItems(args ...interface{}) {

	jsutils.HasPermission("products", "-2", "get", func() {
		jquery.Jq("#menu_scan_qrcode").FadeIn()
		jquery.Jq("#menu_list_products").FadeIn()
		jquery.Jq("#menu_list_bookmarks").FadeIn()
		jquery.Jq("#menu_pubchem").FadeIn()
	}, func() {
	})

	jsutils.HasPermission("products", "-2", "put", func() {
		jquery.Jq("#menu_create_product").FadeIn()
	}, func() {
	})

	jsutils.HasPermission("entities", "-2", "get", func() {
		jquery.Jq("#menu_entities").FadeIn()
		jquery.Jq("#menu_management").FadeIn()
	}, func() {
	})

	jsutils.HasPermission("entities", "-2", "put", func() {
		jquery.Jq("#menu_create_entity").FadeIn()
		jquery.Jq("#menu_management").FadeIn()
	}, func() {
	})

	jsutils.HasPermission("entities", "-2", "put", func() {
		jquery.Jq("#menu_update_welcomeannounce").FadeIn()
		jquery.Jq("#menu_settings").FadeIn()
		jquery.Jq("#menu_management").FadeIn()
	}, func() {
	})

	jsutils.HasPermission("storages", "-2", "get", func() {
		jquery.Jq("#menu_store_locations").FadeIn()
		jquery.Jq("#menu_management").FadeIn()
	}, func() {
	})

	jsutils.HasPermission("store_locations", "-2", "put", func() {
		jquery.Jq("#menu_create_store_location").FadeIn()
		jquery.Jq("#menu_management").FadeIn()
	}, func() {
	})

	jsutils.HasPermission("people", "-2", "get", func() {
		jquery.Jq("#menu_people").FadeIn()
		jquery.Jq("#menu_management").FadeIn()
	}, func() {
	})

}
