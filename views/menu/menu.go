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
	entityCallbackWrapper := func(args ...interface{}) {
		entity.Entity_listCallback(js.Null(), nil)
	}
	entityCreateCallbackWrapper := func(args ...interface{}) {
		entity.Entity_createCallBack(js.Null(), nil)
	}
	storelocationCallbackWrapper := func(args ...interface{}) {
		storelocation.StoreLocation_listCallback(js.Null(), nil)
	}
	storelocationCreateCallbackWrapper := func(args ...interface{}) {
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

	var callbackFunc func(args ...interface{})

	switch args[2].String() {
	case "Product_list":
		callbackFunc = productCallbackWrapper
	case "Product_create":
		callbackFunc = productCreateCallbackWrapper
	case "Entity_list":
		callbackFunc = entityCallbackWrapper
	case "Entity_create":
		callbackFunc = entityCreateCallbackWrapper
	case "StoreLocation_list":
		callbackFunc = storelocationCallbackWrapper
	case "StoreLocation_create":
		callbackFunc = storelocationCreateCallbackWrapper
	case "Person_list":
		callbackFunc = personCallbackWrapper
	case "Person_create":
		callbackFunc = personCreateCallbackWrapper
	case "PersonPass_list":
		callbackFunc = aboutListCallbackWrapper
	case "WelcomeAnnounce_list":
		callbackFunc = aboutListCallbackWrapper
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
	}, func() {
	})

	jsutils.HasPermission("products", "", "post", func() {
		jquery.Jq("#menu_create_product").FadeIn()
	}, func() {
	})

	jsutils.HasPermission("entities", "-2", "get", func() {
		jquery.Jq("#menu_entities").FadeIn()
	}, func() {
	})

	jsutils.HasPermission("entities", "", "post", func() {
		jquery.Jq("#menu_create_entity").FadeIn()
	}, func() {
	})

	jsutils.HasPermission("entities", "-2", "put", func() {
		jquery.Jq("#menu_update_welcomeannounce").FadeIn()
	}, func() {
	})

	jsutils.HasPermission("storages", "-2", "get", func() {
		jquery.Jq("#menu_storelocations").FadeIn()
	}, func() {
	})

	jsutils.HasPermission("storelocations", "", "post", func() {
		jquery.Jq("#menu_create_storelocation").FadeIn()
	}, func() {
	})

	jsutils.HasPermission("people", "-2", "get", func() {
		jquery.Jq("#menu_people").FadeIn()
	}, func() {
	})

	jsutils.HasPermission("people", "", "post", func() {
		jquery.Jq("#menu_create_person").FadeIn()
	}, func() {
	})

}
