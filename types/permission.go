package types

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/models"
)

type Permission struct {
	*models.Permission
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
