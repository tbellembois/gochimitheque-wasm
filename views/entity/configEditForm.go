package entity

import (
	"encoding/json"
	"fmt"
	"strconv"
	"syscall/js"

	. "github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/localStorage"
	. "github.com/tbellembois/gochimitheque-wasm/types"
	"github.com/tbellembois/gochimitheque-wasm/utils"
	"github.com/tbellembois/gochimitheque-wasm/views/search"
	"github.com/tbellembois/gochimitheque-wasm/widgets"
)

func FillInEntityForm(e Entity, id string) {

	Jq(fmt.Sprintf("#%s #entity_id", id)).SetVal(e.EntityID)
	Jq(fmt.Sprintf("#%s #entity_name", id)).SetVal(e.EntityName)
	Jq(fmt.Sprintf("#%s #entity_description", id)).SetVal(e.EntityDescription)
	Jq("select#managers").Select2Clear()

	for _, manager := range e.Managers {
		Jq("select#managers").Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            manager.PersonEmail,
				Value:           strconv.Itoa(manager.PersonId),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

}

func SaveEntity(this js.Value, args []js.Value) interface{} {

	var (
		ajaxURL, ajaxMethod string
		entity              *Entity
		dataBytes           []byte
		err                 error
	)

	if !Jq("#entity").Valid() {
		return nil
	}

	entity = &Entity{}
	if Jq("input#entity_id").GetVal().Truthy() {
		if entity.EntityID, err = strconv.Atoi(Jq("input#entity_id").GetVal().String()); err != nil {
			fmt.Println(err)
			return nil
		}
	}
	entity.EntityName = Jq("input#entity_name").GetVal().String()
	entity.EntityDescription = Jq("input#entity_description").GetVal().String()

	for _, select2Item := range Jq("select#managers").Select2Data() {
		person := &Person{}
		if person.PersonId, err = strconv.Atoi(select2Item.Id); err != nil {
			fmt.Println(err)
			return nil
		}
		person.PersonEmail = select2Item.Text

		entity.Managers = append(entity.Managers, person)
	}

	if dataBytes, err = json.Marshal(entity); err != nil {
		fmt.Println(err)
		return nil
	}

	if Jq("form#entity input#entity_id").Object.Length() > 0 {
		ajaxURL = fmt.Sprintf("%sentities/%d", ApplicationProxyPath, entity.EntityID)
		ajaxMethod = "put"
	} else {
		ajaxURL = fmt.Sprintf("%sentities", ApplicationProxyPath)
		ajaxMethod = "post"
	}

	Ajax{
		URL:    ajaxURL,
		Method: ajaxMethod,
		Data:   dataBytes,
		Done: func(data js.Value) {

			localStorage.Clear()

			var (
				entity Entity
				err    error
			)

			if err = json.Unmarshal([]byte(data.String()), &entity); err != nil {
				fmt.Println(err)
				return
			}

			// TODO: use entityId for redirection
			href := fmt.Sprintf("%sv/entities", ApplicationProxyPath)
			search.ClearSearch(js.Null(), nil)
			utils.LoadContent("entity", href, Entity_SaveCallback, entity.EntityName)

		},
		Fail: func(jqXHR js.Value) {

			utils.DisplayGenericErrorMessage()

		},
	}.Send()

	return nil

}
