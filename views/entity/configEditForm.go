//go:build go1.24 && js && wasm

package entity

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

func FillInEntityForm(e Entity, id string) {

	// js.Global().Get("console").Call("log", fmt.Sprintf("%#v", e.Entity))

	jquery.Jq(fmt.Sprintf("#%s #entity_id", id)).SetVal(*e.EntityID)
	jquery.Jq(fmt.Sprintf("#%s #entity_name", id)).SetVal(e.EntityName)
	jquery.Jq(fmt.Sprintf("#%s #entity_description", id)).SetVal(e.EntityDescription)

	select2Managers := select2.NewSelect2(jquery.Jq("select#managers"), nil)

	select2Managers.Select2Clear()
	if e.Managers != nil {
		for _, manager := range *e.Managers {
			select2Managers.Select2AppendOption(
				widgets.NewOption(widgets.OptionAttributes{
					Text:            manager.PersonEmail,
					Value:           strconv.Itoa(int(*manager.PersonID)),
					DefaultSelected: true,
					Selected:        true,
				}).HTMLElement.OuterHTML())
		}
	}

	// select2LDAPGroups := select2.NewSelect2(jquery.Jq("select#ldapgroups"), nil)

	// select2LDAPGroups.Select2Clear()
	// for _, group := range e.LDAPGroups {
	// 	select2LDAPGroups.Select2AppendOption(
	// 		widgets.NewOption(widgets.OptionAttributes{
	// 			Text:            group,
	// 			Value:           group,
	// 			DefaultSelected: true,
	// 			Selected:        true,
	// 		}).HTMLElement.OuterHTML())
	// }
}

func SaveEntity(this js.Value, args []js.Value) interface{} {

	var (
		ajaxURL, ajaxMethod string
		entity              *Entity
		dataBytes           []byte
		err                 error
	)

	if !validate.NewValidate(jquery.Jq("#entity"), nil).Valid() {
		return nil
	}

	entity = &Entity{Entity: &models.Entity{}}

	if jquery.Jq("input#entity_id").GetVal().Truthy() {
		var _entity_id int
		if _entity_id, err = strconv.Atoi(jquery.Jq("input#entity_id").GetVal().String()); err != nil {
			fmt.Println(err)
			return nil
		}
		_entity_id_64 := int64(_entity_id)
		entity.Entity.EntityID = &_entity_id_64
	}

	entity.EntityName = jquery.Jq("input#entity_name").GetVal().String()
	entity.EntityDescription = jquery.Jq("input#entity_description").GetVal().String()

	var managers []models.Person

	select2Managers := select2.NewSelect2(jquery.Jq("select#managers"), nil)
	for _, select2Item := range select2Managers.Select2Data() {
		person := models.Person{}
		var _person_id int
		if _person_id, err = strconv.Atoi(select2Item.Id); err != nil {
			fmt.Println(err)
			return nil
		}
		_person_id_64 := int64(_person_id)
		person.PersonID = &_person_id_64
		person.PersonEmail = select2Item.Text

		// entity.Managers = append(entity.Managers, person)
		managers = append(managers, person)
	}
	entity.Managers = &managers

	// select2LDAPGroups := select2.NewSelect2(jquery.Jq("select#ldapgroups"), nil)
	// for _, select2Item := range select2LDAPGroups.Select2Data() {
	// 	entity.LDAPGroups = append(entity.LDAPGroups, select2Item.Text)
	// }

	if dataBytes, err = json.Marshal(entity); err != nil {
		fmt.Println(err)
		return nil
	}

	if jquery.Jq("form#entity input#entity_id").Object.Length() > 0 {
		ajaxURL = fmt.Sprintf("%sentities", ApplicationProxyPath)
		ajaxMethod = "put"
	} else {
		ajaxURL = fmt.Sprintf("%sentities", ApplicationProxyPath)
		ajaxMethod = "post"
	}

	ajax.Ajax{
		URL:    ajaxURL,
		Method: ajaxMethod,
		Data:   dataBytes,
		Done: func(data js.Value) {

			globals.LocalStorage.Clear()

			// var (
			// 	entity Entity
			// 	err    error
			// )

			// if err = json.Unmarshal([]byte(data.String()), &entity); err != nil {
			// 	fmt.Println(err)
			// 	return
			// }

			// TODO: use entityId for redirection
			href := fmt.Sprintf("%sv/entities", ApplicationProxyPath)
			jsutils.ClearSearch(js.Null(), nil)
			jsutils.LoadContent("div#content", "entity", href, Entity_SaveCallback, entity.EntityName)

		},
		Fail: func(jqXHR js.Value) {

			jsutils.DisplayGenericErrorMessage()

		},
	}.Send()

	return nil

}
