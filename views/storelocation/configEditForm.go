//go:build go1.24 && js && wasm

package storelocation

import (
	"database/sql"
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

	jquery.Jq(fmt.Sprintf("#%s #store_location_id", id)).SetVal(s.StoreLocationID.Int64)
	jquery.Jq(fmt.Sprintf("#%s #store_location_name", id)).SetVal(s.StoreLocationName.String)
	jquery.Jq(fmt.Sprintf("#%s #store_location_can_store", id)).SetProp("checked", s.StoreLocationCanStore.Bool)
	jquery.Jq(fmt.Sprintf("#%s #store_location_color", id)).SetVal(s.StoreLocationColor.String)

	select2Entity := select2.NewSelect2(jquery.Jq("select#entity"), nil)
	select2Entity.Select2Clear()
	select2Entity.Select2AppendOption(
		widgets.NewOption(widgets.OptionAttributes{
			Text:            s.Entity.EntityName,
			Value:           strconv.Itoa(s.Entity.EntityID),
			DefaultSelected: true,
			Selected:        true,
		}).HTMLElement.OuterHTML())

	select2StoreLocation := select2.NewSelect2(jquery.Jq("select#store_location"), nil)
	select2StoreLocation.Select2Clear()

	var parentId, parentName string
	if s.StoreLocation.StoreLocation != nil {
		parentId = strconv.Itoa(int(s.StoreLocation.StoreLocation.StoreLocationID.Int64))
		parentName = s.StoreLocation.StoreLocation.StoreLocationName.String
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

		storelocation.StoreLocationID = models.NullInt64{sql.NullInt64{
			Int64: int64(storelocationId),
			Valid: true,
		}}
	}

	if jquery.Jq("input#store_location_can_store:checked").Object.Length() > 0 {
		storelocation.StoreLocationCanStore = models.NullBool{sql.NullBool{
			Bool:  true,
			Valid: true,
		}}
	} else {
		storelocation.StoreLocationCanStore = models.NullBool{sql.NullBool{
			Bool:  false,
			Valid: true,
		}}
	}

	storelocation.StoreLocationName = models.NullString{sql.NullString{
		String: jquery.Jq("input#store_location_name").GetVal().String(),
		Valid:  true,
	}}
	storelocation.StoreLocationColor = models.NullString{sql.NullString{
		String: jquery.Jq("input#store_location_color").GetVal().String(),
		Valid:  true,
	}}

	select2ItemEntity := select2.NewSelect2(jquery.Jq("select#entity"), nil)
	storelocation.Entity = models.Entity{}

	if storelocation.Entity.EntityID, err = strconv.Atoi(select2ItemEntity.Select2Data()[0].Id); err != nil {
		fmt.Println("select2ItemEntity:" + err.Error())
		return nil
	}

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

			storelocation.StoreLocation.StoreLocation.StoreLocationID = models.NullInt64{sql.NullInt64{
				Int64: int64(storelocationId),
				Valid: true,
			}}
			// FIELD_REQUIRED_BY_RUST_MODEL_BUT_NOT_USED_IN_CREATE_OR_UPDATE
			storelocation.StoreLocation.StoreLocation.StoreLocationName = models.NullString{sql.NullString{
				String: "",
				Valid:  true,
			}}
			// FIELD_REQUIRED_BY_RUST_MODEL_BUT_NOT_USED_IN_CREATE_OR_UPDATE
			storelocation.StoreLocation.StoreLocation.StoreLocationCanStore = models.NullBool{sql.NullBool{
				Bool:  false,
				Valid: true,
			}}
		}
	}

	if dataBytes, err = json.Marshal(storelocation); err != nil {
		fmt.Println(err)
		return nil
	}

	if jquery.Jq("form#store_location input#store_location_id").Object.Length() > 0 {
		ajaxURL = fmt.Sprintf("%sstore_locations/%d", ApplicationProxyPath, storelocation.StoreLocationID.Int64)
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
			jsutils.LoadContent("div#content", "store_location", href, StoreLocation_SaveCallback, storelocation.StoreLocationName.String)

		},
		Fail: func(jqXHR js.Value) {

			jsutils.DisplayGenericErrorMessage()

		},
	}.Send()

	return nil

}
