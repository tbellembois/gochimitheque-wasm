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
	pid := jquery.Jq("input#person_id")

	if pid.GetVal().String() != "" {
		id = pid.GetVal().String()
	}

	settings.Set("url", fmt.Sprintf("%svalidate/person/%s/email/", ApplicationProxyPath, id))

	return nil

}
