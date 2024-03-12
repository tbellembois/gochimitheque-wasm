package person

import (
	"fmt"
	"strconv"
	"sync"
	"syscall/js"

	"github.com/rocketlaunchr/react/forks/encoding/json"
	"github.com/tbellembois/gochimitheque-wasm/ajax"
	"github.com/tbellembois/gochimitheque-wasm/globals"
	. "github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/jquery"
	"github.com/tbellembois/gochimitheque-wasm/jsutils"
	"github.com/tbellembois/gochimitheque-wasm/models"
	"github.com/tbellembois/gochimitheque-wasm/select2"
	"github.com/tbellembois/gochimitheque-wasm/types"
	. "github.com/tbellembois/gochimitheque-wasm/types"
	"github.com/tbellembois/gochimitheque-wasm/validate"
	"github.com/tbellembois/gochimitheque-wasm/widgets"
	"honnef.co/go/js/dom/v2"
)

// populatePermission checks the permissions checkboxes in the person edition page
func populatePermission(permissions []types.Permission, managedEntitiesIds map[int]string) {

	Doc := dom.GetWindow().Document()

	// unchecking all permissions
	for _, e := range Doc.GetElementsByClassName("perm") {
		e.(*dom.HTMLInputElement).RemoveAttribute("checked")
	}

	// setting all permissions at none by defaut
	for _, e := range Doc.GetElementsByClassName("permn") {
		e.(*dom.HTMLInputElement).SetChecked(true)
	}

	// then setting up new permissions
	for _, p := range permissions {

		pentityid := strconv.Itoa(p.PermissionEntityID)
		if _, ok := managedEntitiesIds[p.PermissionEntityID]; ok {
			continue
		}

		switch p.PermissionItemName {
		case "products":
			switch p.PermissionPermName {
			case "w", "all":
				if pentityid == "-1" {
					for _, e := range Doc.GetElementsByClassName("permwproducts") {
						e.(*dom.HTMLInputElement).SetChecked(true)
					}
				} else {
					input := Doc.GetElementByID("permw" + p.PermissionItemName + pentityid)
					if input != nil {
						input.(*dom.HTMLInputElement).SetChecked(true)
					}
				}
			case "r":
				if pentityid == "-1" {
					for _, e := range Doc.GetElementsByClassName("permrproducts") {
						// avoid selecting r if w is already selected
						eid := e.GetAttribute("entity_id")
						if !Doc.GetElementByID("permw" + p.PermissionItemName + eid).(*dom.HTMLInputElement).Checked() {
							e.(*dom.HTMLInputElement).SetChecked(true)
						}
					}
				} else {
					// avoid selecting r if w is already selected
					if !Doc.GetElementByID("permw" + p.PermissionItemName + pentityid).(*dom.HTMLInputElement).Checked() {
						Doc.GetElementByID("permr" + p.PermissionItemName + pentityid).(*dom.HTMLInputElement).SetChecked(true)
					}
				}
			case "n":
				if pentityid == "-1" {
					for _, e := range Doc.GetElementsByClassName("permnproducts") {
						e.(*dom.HTMLInputElement).SetChecked(true)
					}
				} else {
					Doc.GetElementByID("permn" + p.PermissionItemName + pentityid).(*dom.HTMLInputElement).SetChecked(true)
				}
			}
		case "rproducts":
			switch p.PermissionPermName {
			case "w", "all":
				if pentityid == "-1" {
					// for _, e := range Doc.GetElementsByClassName("permwrproducts") {
					for _, e := range Doc.GetElementsByClassName("permrrproducts") {
						e.(*dom.HTMLInputElement).SetChecked(true)
					}
				} else {
					input := Doc.GetElementByID("permw" + p.PermissionItemName + pentityid)
					if input != nil {
						input.(*dom.HTMLInputElement).SetChecked(true)
					}
				}
			case "r":
				if pentityid == "-1" {
					for _, e := range Doc.GetElementsByClassName("permrrproducts") {
						// avoid selecting r if w is already selected
						e.(*dom.HTMLInputElement).SetChecked(true)
						// eid := e.GetAttribute("entity_id")
						// if !Doc.GetElementByID("permw" + p.PermissionItemName + eid).(*dom.HTMLInputElement).Checked() {
						// 	e.(*dom.HTMLInputElement).SetChecked(true)
						// }
					}
				} else {
					// avoid selecting r if w is already selected
					Doc.GetElementByID("permr" + p.PermissionItemName + pentityid).(*dom.HTMLInputElement).SetChecked(true)
					// if !Doc.GetElementByID("permw" + p.PermissionItemName + pentityid).(*dom.HTMLInputElement).Checked() {
					// 	Doc.GetElementByID("permr" + p.PermissionItemName + pentityid).(*dom.HTMLInputElement).SetChecked(true)
					// }
				}
			case "n":
				if pentityid == "-1" {
					for _, e := range Doc.GetElementsByClassName("permnrproducts") {
						e.(*dom.HTMLInputElement).SetChecked(true)
					}
				} else {
					Doc.GetElementByID("permn" + p.PermissionItemName + pentityid).(*dom.HTMLInputElement).SetChecked(true)
				}
			}
		case "storages":
			switch p.PermissionPermName {
			case "w", "all":
				if pentityid == "-1" {
					for _, e := range Doc.GetElementsByClassName("permwstorages") {
						e.(*dom.HTMLInputElement).SetChecked(true)
					}
				} else {
					Doc.GetElementByID("permw" + p.PermissionItemName + pentityid).(*dom.HTMLInputElement).SetChecked(true)
				}
			case "r":
				if pentityid == "-1" {
					for _, e := range Doc.GetElementsByClassName("permrstorages") {
						// avoid selecting r if w is already selected
						eid := e.GetAttribute("entity_id")
						if !Doc.GetElementByID("permw" + p.PermissionItemName + eid).(*dom.HTMLInputElement).Checked() {
							e.(*dom.HTMLInputElement).SetChecked(true)
						}
					}
				} else {
					// avoid selecting r if w is already selected
					if !Doc.GetElementByID("permw" + p.PermissionItemName + pentityid).(*dom.HTMLInputElement).Checked() {
						Doc.GetElementByID("permr" + p.PermissionItemName + pentityid).(*dom.HTMLInputElement).SetChecked(true)
					}
				}
			}
		case "all":
			switch p.PermissionPermName {
			case "w", "all":
				if pentityid == "-1" {
					// super admin (if "all")
					for _, e := range Doc.GetElementsByClassName("permw") {
						e.(*dom.HTMLInputElement).SetChecked(true)
					}
				}
			case "r":
				for _, e := range Doc.GetElementsByClassName("permr") {
					e.(*dom.HTMLInputElement).SetChecked(true)
				}
			}
		}
	}

}

func SelectAllEntity(this js.Value, args []js.Value) interface{} {

	url := fmt.Sprintf("%sentities", ApplicationProxyPath)
	method := "get"

	done := func(data js.Value) {

		var (
			entities Select2Entities
			err      error
		)

		if err = json.Unmarshal([]byte(data.String()), &entities); err != nil {
			fmt.Println(err)
			return
		}

		jquery.Jq("#permissions").SetHtml("")

		select2Entities := select2.NewSelect2(jquery.Jq("select#entities"), nil)
		select2Entities.Select2Clear()
		for _, entity := range entities.Rows {

			select2Entities.Select2AppendOption(
				widgets.NewOption(widgets.OptionAttributes{
					Text:            entity.EntityName,
					Value:           strconv.Itoa(entity.EntityID),
					DefaultSelected: true,
					Selected:        true,
				}).HTMLElement.OuterHTML())

			jquery.Jq("#permissions").Append(widgets.Permission(entity.EntityID, entity.EntityName, false).OuterHTML())

		}

		jquery.Jq("button#selectAllEntity").SetProp("disabled", false)
		jquery.Jq("#selectAllPermissions").Show()

	}
	fail := func(data js.Value) {

		jsutils.DisplayGenericErrorMessage()

	}

	jquery.Jq("button#selectAllEntity").SetProp("disabled", true)

	ajax.Ajax{
		Method: method,
		URL:    url,
		Done:   done,
		Fail:   fail,
	}.Send()

	return nil

}

func FillInPersonForm(p Person, id string) {

	type Permissions []Permission
	type Entities []Entity

	var (
		wg                        sync.WaitGroup
		entities, managedEntities Entities
		managedEntitiesIds        map[int]string
		permissions               Permissions
		err                       error
	)

	managedEntitiesIds = make(map[int]string)

	select2Entities := select2.NewSelect2(jquery.Jq("select#entities"), nil)
	select2Entities.Select2Clear()

	jquery.Jq("#permissions").Empty()

	// Getting the entities the person is manager of.
	wg.Add(1)
	go func() {

		ajax.Ajax{
			URL:    fmt.Sprintf("%speople/%d/manageentities", ApplicationProxyPath, p.PersonID),
			Method: "get",
			Done: func(data js.Value) {
				if err = json.Unmarshal([]byte(data.String()), &managedEntities); err != nil {
					fmt.Println(err)
					jsutils.DisplayGenericErrorMessage()
				}
				for _, entity := range managedEntities {
					managedEntitiesIds[entity.EntityID] = entity.EntityName
				}
				wg.Done()
			},
			Fail: func(jqXHR js.Value) {
				jsutils.DisplayGenericErrorMessage()
				wg.Done()
			},
		}.Send()

	}()

	// Getting the person permissions.
	wg.Add(1)
	go func() {

		ajax.Ajax{
			URL:    fmt.Sprintf("%speople/%d/permissions", ApplicationProxyPath, p.PersonID),
			Method: "get",
			Done: func(data js.Value) {
				if err = json.Unmarshal([]byte(data.String()), &permissions); err != nil {
					fmt.Println(err)
					jsutils.DisplayGenericErrorMessage()
				}
				wg.Done()
			},
			Fail: func(jqXHR js.Value) {
				jsutils.DisplayGenericErrorMessage()
				wg.Done()
			},
		}.Send()

	}()

	// Getting the person entities.
	wg.Add(1)
	go func() {

		ajax.Ajax{
			URL:    fmt.Sprintf("%speople/%d/entities", ApplicationProxyPath, p.PersonID),
			Method: "get",
			Done: func(data js.Value) {
				if err = json.Unmarshal([]byte(data.String()), &entities); err != nil {
					fmt.Println(err)
					jsutils.DisplayGenericErrorMessage()
				}
				wg.Done()
			},
			Fail: func(jqXHR js.Value) {
				jsutils.DisplayGenericErrorMessage()
				wg.Done()
			},
		}.Send()

	}()

	wg.Wait()

	jquery.Jq(fmt.Sprintf("#%s #person_id", id)).SetVal(p.PersonID)
	jquery.Jq(fmt.Sprintf("#%s #person_email", id)).SetVal(p.PersonEmail)
	jquery.Jq(fmt.Sprintf("#%s #person_password", id)).SetVal("")

	// Appending managed entities in hidden inputs for further use.
	jquery.Jq(fmt.Sprintf("#%s option.manageentities", id)).Remove()
	for _, entity := range managedEntities {
		option := widgets.NewOption(widgets.OptionAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Classes: []string{"manageentities"},
				Attributes: map[string]string{
					"type": "hidden",
				},
			},
			Value: strconv.Itoa(entity.EntityID),
		})
		jquery.Jq("form#person").Append(option.OuterHTML())
	}

	// Populating the entities select2.
	for _, entity := range entities {
		select2Entities.Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            entity.EntityName,
				Value:           strconv.Itoa(entity.EntityID),
				DefaultSelected: true,
				Selected:        true,
			}).HTMLElement.OuterHTML())
	}

	// Adding a permission widget for each entity
	// except for managed entities.
	for _, entity := range entities {
		if _, ok := managedEntitiesIds[entity.EntityID]; !ok {
			jquery.Jq("#permissions").Append(widgets.Permission(entity.EntityID, entity.EntityName, false).OuterHTML())
		}
	}

	// Populating the permissions widget.
	populatePermission(permissions, managedEntitiesIds)

}

func SavePerson(this js.Value, args []js.Value) interface{} {

	var (
		ajaxURL, ajaxMethod string
		dataBytes           []byte
		person              *Person
		err                 error
	)

	if !validate.NewValidate(jquery.Jq("#person"), nil).Valid() {
		return nil
	}

	person = &Person{Person: &models.Person{}}

	if jquery.Jq("input#person_id").GetVal().Truthy() {
		if person.PersonID, err = strconv.Atoi(jquery.Jq("input#person_id").GetVal().String()); err != nil {
			fmt.Println(err)
			return nil
		}
	}

	person.PersonEmail = jquery.Jq("input#person_email").GetVal().String()

	select2Entities := select2.NewSelect2(jquery.Jq("select#entities"), nil)

	for _, select2Item := range select2Entities.Select2Data() {
		entity := &models.Entity{}
		if entity.EntityID, err = strconv.Atoi(select2Item.Id); err != nil {
			fmt.Println(err)
			return nil
		}
		entity.EntityName = select2Item.Text

		person.Entities = append(person.Entities, entity)
	}

	permissions := jquery.Jq("input[type=radio]:checked").Object
	for i := 0; i < permissions.Length(); i++ {
		permission := &models.Permission{}

		permission.PermissionPermName = permissions.Index(i).Call("getAttribute", "perm_name").String()
		permission.PermissionItemName = permissions.Index(i).Call("getAttribute", "item_name").String()
		entityId := permissions.Index(i).Call("getAttribute", "entity_id").String()

		if permission.PermissionEntityID, err = strconv.Atoi(entityId); err != nil {
			fmt.Println(err)
			return nil
		}

		person.Permissions = append(person.Permissions, permission)
	}

	if dataBytes, err = json.Marshal(person); err != nil {
		fmt.Println(err)
		return nil
	}

	if jquery.Jq("form#person input#person_id").Object.Length() > 0 {
		ajaxURL = fmt.Sprintf("%speople/%d", ApplicationProxyPath, person.PersonID)
		ajaxMethod = "put"
	} else {
		ajaxURL = fmt.Sprintf("%speople", ApplicationProxyPath)
		ajaxMethod = "post"
	}

	ajax.Ajax{
		URL:    ajaxURL,
		Method: ajaxMethod,
		Data:   dataBytes,
		Done: func(data js.Value) {

			globals.LocalStorage.Clear()

			var (
				person Person
				err    error
			)

			if err = json.Unmarshal([]byte(data.String()), &person); err != nil {
				fmt.Println(err)
				return
			}

			// TODO: use personId for redirection
			href := fmt.Sprintf("%sv/people", ApplicationProxyPath)
			jsutils.ClearSearch(js.Null(), nil)
			jsutils.LoadContent("div#content", "person", href, Person_SaveCallback, person.PersonEmail)

		},
		Fail: func(jqXHR js.Value) {

			jsutils.DisplayGenericErrorMessage()

		},
	}.Send()

	return nil

}
