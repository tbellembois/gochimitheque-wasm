package entity

import (
	"fmt"
	"syscall/js"

	. "github.com/tbellembois/gochimitheque-wasm/globals"
	. "github.com/tbellembois/gochimitheque-wasm/types"
)

func ValidateEntityNameBeforeSend(this js.Value, args []js.Value) interface{} {

	settings := args[1]

	id := "-1"
	eid := Jq("input#entity_id")

	if eid.Object.Length() > 0 {
		id = eid.GetVal().String()
	}

	settings.Set("url", fmt.Sprintf("%svalidate/entity/%s/name/", ApplicationProxyPath, id))

	return nil

}
