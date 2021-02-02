package about

import (
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/jquery"
)

func About_listCallback(this js.Value, args []js.Value) interface{} {

	jquery.Jq("#search").Hide()
	jquery.Jq("#actions").Hide()

	return nil

}
