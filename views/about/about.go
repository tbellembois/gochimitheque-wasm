package about

import (
	"syscall/js"

	. "github.com/tbellembois/gochimitheque-wasm/types"
)

func About_listCallback(this js.Value, args []js.Value) interface{} {

	Jq("#search").Hide()
	Jq("#actions").Hide()

	return nil

}
