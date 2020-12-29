package menu

import (
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/utils"
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
		callbackFunc = nil
	case "WelcomeAnnounce_list":
		callbackFunc = nil
	case "About_list":
		callbackFunc = nil
	}

	utils.LoadContent(args[0].String(), args[1].String(), callbackFunc, nil)

	return nil

}
