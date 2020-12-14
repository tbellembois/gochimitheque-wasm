package types

import (
	"encoding/json"
	"fmt"
	"syscall/js"
)

type Permission struct {
	PermissionID       int    `json:"permission_id"`
	PermissionPermName string `json:"permission_perm_name" js:"permission_perm_name"` // ex: r
	PermissionItemName string `json:"permission_item_name" js:"permission_item_name"` // ex: entity
	PermissionEntityID int    `json:"permission_entity_id" js:"permission_entity_id"` // ex: 8
}

func PermissionFromJsJSONValue(jsvalue js.Value) Permission {

	var (
		permission Permission
		err        error
	)

	jsPermissionString := js.Global().Get("JSON").Call("stringify", jsvalue).String()
	if err = json.Unmarshal([]byte(jsPermissionString), &permission); err != nil {
		fmt.Println(err)
	}

	return permission

}
