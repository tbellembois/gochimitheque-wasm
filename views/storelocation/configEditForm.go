package storelocation

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"syscall/js"

	. "github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/localStorage"
	. "github.com/tbellembois/gochimitheque-wasm/types"
	"github.com/tbellembois/gochimitheque-wasm/utils"
	"github.com/tbellembois/gochimitheque-wasm/widgets"
)

func FillInStoreLocationForm(s StoreLocation, id string) {

	Jq(fmt.Sprintf("#%s #storelocation_id", id)).SetVal(s.StoreLocationID.Int64)
	Jq(fmt.Sprintf("#%s #storelocation_name", id)).SetVal(s.StoreLocationName.String)
	Jq(fmt.Sprintf("#%s #storelocation_canstore", id)).SetProp("checked", s.StoreLocationCanStore.Bool)
	Jq(fmt.Sprintf("#%s #storelocation_color", id)).SetVal(s.StoreLocationColor.String)

	Jq("select#entity").Select2Clear()

	Jq("select#entity").Select2AppendOption(
		widgets.NewOption(widgets.OptionAttributes{
			Text:            s.Entity.EntityName,
			Value:           strconv.Itoa(s.Entity.EntityID),
			DefaultSelected: true,
			Selected:        true,
		}).HTMLElement.OuterHTML())

	Jq("select#storelocation").Select2Clear()

	var parentId, parentName string
	if s.StoreLocation != nil {
		parentId = strconv.Itoa(int(s.StoreLocation.StoreLocationID.Int64))
		parentName = s.StoreLocation.StoreLocationName.String
	}
	Jq("select#storelocation").Select2AppendOption(
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

	if !Jq("#storelocation").Valid() {
		return nil
	}

	storelocation = &StoreLocation{}
	if Jq("input#storelocation_id").GetVal().Truthy() {
		if storelocationId, err = strconv.Atoi(Jq("input#storelocation_id").GetVal().String()); err != nil {
			fmt.Println(err)
			return nil
		}
		storelocation.StoreLocationID = sql.NullInt64{
			Int64: int64(storelocationId),
			Valid: true,
		}
	}
	storelocation.StoreLocationCanStore = sql.NullBool{
		Bool:  Jq("input#storelocation_canstore").Prop("checked").(js.Value).Bool(),
		Valid: true,
	}
	storelocation.StoreLocationName = sql.NullString{
		String: Jq("input#storelocation_name").GetVal().String(),
		Valid:  true,
	}
	storelocation.StoreLocationColor = sql.NullString{
		String: Jq("input#storelocation_color").GetVal().String(),
		Valid:  true,
	}

	select2ItemEntity := Jq("select#entity").Select2Data()[0]
	storelocation.Entity = Entity{}
	if storelocation.Entity.EntityID, err = strconv.Atoi(select2ItemEntity.Id); err != nil {
		fmt.Println(err)
		return nil
	}
	storelocation.Entity.EntityName = select2ItemEntity.Text

	if len(Jq("select#storelocation").Select2Data()) > 0 {
		select2ItemStoreLocation := Jq("select#storelocation").Select2Data()[0]
		if !select2ItemStoreLocation.IsEmpty() {
			storelocation.StoreLocation = &StoreLocation{}
			if storelocationId, err = strconv.Atoi(select2ItemStoreLocation.Id); err != nil {
				fmt.Println(err)
				return nil
			}
			storelocation.StoreLocation.StoreLocationID = sql.NullInt64{
				Int64: int64(storelocationId),
				Valid: true,
			}
		}
	}

	if dataBytes, err = json.Marshal(storelocation); err != nil {
		fmt.Println(err)
		return nil
	}

	if Jq("form#storelocation input#storelocation_id").Object.Length() > 0 {
		ajaxURL = fmt.Sprintf("%sstorelocations/%d", ApplicationProxyPath, storelocation.StoreLocationID.Int64)
		ajaxMethod = "put"
	} else {
		ajaxURL = fmt.Sprintf("%sstorelocations", ApplicationProxyPath)
		ajaxMethod = "post"
	}

	Ajax{
		URL:    ajaxURL,
		Method: ajaxMethod,
		Data:   dataBytes,
		Done: func(data js.Value) {

			localStorage.Clear()

			var (
				storelocation StoreLocation
				err           error
			)

			if err = json.Unmarshal([]byte(data.String()), &storelocation); err != nil {
				fmt.Println(err)
				return
			}

			// TODO: use entityId for redirection
			href := fmt.Sprintf("%sv/storelocations", ApplicationProxyPath)
			utils.LoadContent("storelocation", href, StoreLocation_SaveCallback, storelocation.StoreLocationName.String)

		},
		Fail: func(jqXHR js.Value) {

			utils.DisplayGenericErrorMessage()

		},
	}.Send()

	return nil

}
