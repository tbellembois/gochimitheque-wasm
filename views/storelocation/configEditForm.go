//go:build go1.24 && js && wasm

package storelocation

import (
	"encoding/json"
	"fmt"
	"strconv"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/ajax"
	"github.com/tbellembois/gochimitheque-wasm/globals"
	. "github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/jquery"
	"github.com/tbellembois/gochimitheque-wasm/jsutils"
	"github.com/tbellembois/gochimitheque-wasm/select2"
	. "github.com/tbellembois/gochimitheque-wasm/types"
	"github.com/tbellembois/gochimitheque-wasm/validate"
	"github.com/tbellembois/gochimitheque-wasm/widgets"
	"github.com/tbellembois/gochimitheque/models"
)

func FillInStoreLocationForm(s StoreLocation, id string) {

	jquery.Jq(fmt.Sprintf("#%s #store_location_id", id)).SetVal(*s.StoreLocationID)
	jquery.Jq(fmt.Sprintf("#%s #store_location_name", id)).SetVal(s.StoreLocationName)
	jquery.Jq(fmt.Sprintf("#%s #store_location_can_store", id)).SetProp("checked", s.StoreLocationCanStore)
	jquery.Jq(fmt.Sprintf("#%s #store_location_color", id)).SetVal(*s.StoreLocationColor)

	select2Entity := select2.NewSelect2(jquery.Jq("select#entity"), nil)
	select2Entity.Select2Clear()
	select2Entity.Select2AppendOption(
		widgets.NewOption(widgets.OptionAttributes{
			Text:            s.Entity.EntityName,
			Value:           strconv.Itoa(int(*s.Entity.EntityID)),
			DefaultSelected: true,
			Selected:        true,
		}).HTMLElement.OuterHTML())

	select2StoreLocation := select2.NewSelect2(jquery.Jq("select#store_location"), nil)
	select2StoreLocation.Select2Clear()

	var parentId, parentName string
	if s.StoreLocation.StoreLocation != nil {
		parentId = strconv.Itoa(int(*s.StoreLocation.StoreLocation.StoreLocationID))
		parentName = s.StoreLocation.StoreLocation.StoreLocationName
	}
	select2StoreLocation.Select2AppendOption(
		widgets.NewOption(widgets.OptionAttributes{
			Text:            parentName,
			Value:           parentId,
			DefaultSelected: true,
			Selected:        true,
		}).HTMLElement.OuterHTML())

}

func SaveStoreLocation(this js.Value, args []js.Value) interface{} {

	var (
		ajaxURL, ajaxMethod string
		storelocation       *StoreLocation
		storelocationId     int
		dataBytes           []byte
		err                 error
	)

	if !validate.NewValidate(jquery.Jq("#store_location"), nil).Valid() {
		return nil
	}

	storelocation = &StoreLocation{StoreLocation: &models.StoreLocation{}}

	if jquery.Jq("input#store_location_id").GetVal().Truthy() {
		if storelocationId, err = strconv.Atoi(jquery.Jq("input#store_location_id").GetVal().String()); err != nil {
			fmt.Println(err)
			return nil
		}

		var _id64 int64 = int64(storelocationId)
		storelocation.StoreLocationID = &_id64
	}

	if jquery.Jq("input#store_location_can_store:checked").Object.Length() > 0 {
		storelocation.StoreLocationCanStore = true
	} else {
		storelocation.StoreLocationCanStore = false
	}

	var _color string = jquery.Jq("input#store_location_color").GetVal().String()
	storelocation.StoreLocationName = jquery.Jq("input#store_location_name").GetVal().String()
	storelocation.StoreLocationColor = &_color

	select2ItemEntity := select2.NewSelect2(jquery.Jq("select#entity"), nil)
	storelocation.Entity = &models.Entity{}

	var _id int
	if _id, err = strconv.Atoi(select2ItemEntity.Select2Data()[0].Id); err != nil {
		fmt.Println("select2ItemEntity:" + err.Error())
		return nil
	}
	var _id64 = int64(_id)
	storelocation.Entity.EntityID = &_id64
	storelocation.Entity.EntityName = select2ItemEntity.Select2Data()[0].Text

	select2StoreLocation := select2.NewSelect2(jquery.Jq("select#store_location"), nil)

	if len(select2StoreLocation.Select2Data()) > 0 && select2StoreLocation.Select2Data()[0].Id != "" {
		select2ItemStoreLocation := select2StoreLocation.Select2Data()[0]
		if !select2ItemStoreLocation.IsEmpty() {
			storelocation.StoreLocation.StoreLocation = &models.StoreLocation{}

			if storelocationId, err = strconv.Atoi(select2ItemStoreLocation.Id); err != nil {
				fmt.Println("select2ItemStoreLocation:" + err.Error())
				return nil
			}

			var _id64 = int64(storelocationId)
			storelocation.StoreLocation.StoreLocation.StoreLocationID = &_id64

			// FIELD_REQUIRED_BY_RUST_MODEL_BUT_NOT_USED_IN_CREATE_OR_UPDATE
			storelocation.StoreLocation.StoreLocation.StoreLocationName = ""
			// FIELD_REQUIRED_BY_RUST_MODEL_BUT_NOT_USED_IN_CREATE_OR_UPDATE
			storelocation.StoreLocation.StoreLocation.StoreLocationCanStore = false
		}
	}

	if dataBytes, err = json.Marshal(storelocation); err != nil {
		fmt.Println(err)
		return nil
	}

	if jquery.Jq("form#store_location input#store_location_id").Object.Length() > 0 {
		ajaxURL = fmt.Sprintf("%sstore_locations", ApplicationProxyPath)
		ajaxMethod = "put"
	} else {
		ajaxURL = fmt.Sprintf("%sstore_locations", ApplicationProxyPath)
		ajaxMethod = "post"
	}

	ajax.Ajax{
		URL:    ajaxURL,
		Method: ajaxMethod,
		Data:   dataBytes,
		Done: func(data js.Value) {

			globals.LocalStorage.Clear()

			// var (
			// 	storelocation StoreLocation
			// 	err           error
			// )
			//
			// if err = json.Unmarshal([]byte(data.String()), &storelocation); err != nil {
			// 	fmt.Println(err)
			// 	return
			// }

			// TODO: use entityId for redirection
			href := fmt.Sprintf("%sv/store_locations", ApplicationProxyPath)
			jsutils.ClearSearch(js.Null(), nil)
			jsutils.LoadContent("div#content", "store_location", href, StoreLocation_SaveCallback, storelocation.StoreLocationName)

		},
		Fail: func(jqXHR js.Value) {

			jsutils.DisplayGenericErrorMessage()

		},
	}.Send()

	return nil

}
