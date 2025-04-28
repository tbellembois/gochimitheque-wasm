//go:build go1.24 && js && wasm

package person

import (
	"fmt"
	"syscall/js"

	. "github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/jquery"
)

func ValidatePersonEmailBeforeSend(this js.Value, args []js.Value) interface{} {

	settings := args[1]

	id := "-1"
	email := jquery.Jq("input#person_email").GetVal().String()
	pid := jquery.Jq("input#person_id")

	if pid.Object.Length() > 0 {
		id = pid.GetVal().String()
	}

	settings.Set("url", fmt.Sprintf("%svalidate/person/%s/email/%s", ApplicationProxyPath, id, email))

	return nil

}
