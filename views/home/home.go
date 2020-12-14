package home

import (
	"syscall/js"

	. "github.com/tbellembois/gochimitheque-wasm/types"
	"github.com/tbellembois/gochimitheque-wasm/utils"
)

func Home_listCallback(this js.Value, args []js.Value) interface{} {

	email := js.Global().Call("readCookie", "email").String()
	Jq("#logged").SetHtml(email)

	utils.HasPermission("products", "-2", "get", func() {
		Jq("#menu_scan_qrcode").FadeIn()
		Jq("#menu_list_products").FadeIn()
		Jq("#menu_list_bookmarks").FadeIn()
	}, func() {
	})

	utils.HasPermission("products", "", "post", func() {
		Jq("#menu_create_product").FadeIn()
	}, func() {
	})

	utils.HasPermission("entities", "-2", "get", func() {
		Jq("#menu_entities").FadeIn()
	}, func() {
	})

	utils.HasPermission("entities", "", "post", func() {
		Jq("#menu_create_entity").FadeIn()
	}, func() {
	})

	utils.HasPermission("entities", "-2", "put", func() {
		Jq("#menu_update_welcomeannounce").FadeIn()
	}, func() {
	})

	utils.HasPermission("storages", "-2", "get", func() {
		Jq("#menu_storelocations").FadeIn()
	}, func() {
	})

	utils.HasPermission("storelocations", "", "post", func() {
		Jq("#menu_create_storelocation").FadeIn()
	}, func() {
	})

	utils.HasPermission("people", "-2", "get", func() {
		Jq("#menu_people").FadeIn()
	}, func() {
	})

	utils.HasPermission("people", "", "post", func() {
		Jq("#menu_create_person").FadeIn()
	}, func() {
	})

	return nil

}
