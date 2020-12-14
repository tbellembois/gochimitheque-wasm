package utils

import (
	"fmt"
	"syscall/js"

	. "github.com/tbellembois/gochimitheque-wasm/globals"
	localStorage "github.com/tbellembois/gochimitheque-wasm/localStorage"
	. "github.com/tbellembois/gochimitheque-wasm/types"
)

func HasPermission(item, id, method string, done, fail func()) {

	go func() {
		cacheKey := fmt.Sprintf("%s:%s:%s", item, id, method)
		cachedPermission := localStorage.GetItem(cacheKey)

		if !DisableCache && cachedPermission != "" {
			if cachedPermission == "true" {
				done()
			} else {
				fail()
			}
		} else {
			var url string
			if id != "" {
				url = ApplicationProxyPath + "f/" + item + "/" + id
			} else {
				url = ApplicationProxyPath + "f/" + item
			}

			ajaxDone := func(js.Value) {
				localStorage.SetItem(cacheKey, "true")
				done()
			}
			ajaxFail := func(js.Value) {
				localStorage.SetItem(cacheKey, "false")
				fail()
			}
			ajax := Ajax{
				URL:    url,
				Method: method,
				Done:   ajaxDone,
				Fail:   ajaxFail,
			}

			ajax.Send()

		}
	}()

}
