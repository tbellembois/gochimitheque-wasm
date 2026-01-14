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
	pid := jquery.Jq("input#person_id")

	if pid.Object.Length() > 0 {
		id = pid.GetVal().String()
	}

	settings.Set("url", fmt.Sprintf("%svalidate/person/%s/email/", BackProxyPath, id))

	return nil

}

func ValidatePersonEmailData(this js.Value, args []js.Value) interface{} {

	// return select2.NewSelect2(jquery.Jq("input#person_email"), nil).Select2Data()[0].Text
	return jquery.Jq("input#person_email").GetVal().String()

}
