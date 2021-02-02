package home

import (
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/jquery"
	"github.com/tbellembois/gochimitheque-wasm/jsutils"
)

func Home_listCallback(this js.Value, args []js.Value) interface{} {

	email := js.Global().Call("readCookie", "email").String()
	jquery.Jq("#logged").SetHtml(email)

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

	return nil

}
