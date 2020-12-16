package common

import (
	"fmt"
	"syscall/js"

	. "github.com/tbellembois/gochimitheque-wasm/globals"
	. "github.com/tbellembois/gochimitheque-wasm/types"
	. "github.com/tbellembois/gochimitheque-wasm/utils"
	"github.com/tbellembois/gochimitheque-wasm/views/product"
	"github.com/tbellembois/gochimitheque-wasm/views/storage"
)

func Export(this js.Value, args []js.Value) interface{} {

	BSTableQueryFilter.Lock()
	BSTableQueryFilter.QueryFilter.Export = true

	if CurrentView != "storage" {
		Jq("#Product_table").Bootstraptable(nil).Refresh(nil)
	} else {
		Jq("#Storage_table").Bootstraptable(nil).Refresh(nil)
	}

	return nil

}

func SwitchProductStorageWrapper(this js.Value, args []js.Value) interface{} {

	fmt.Println(CurrentView)

	storageCallbackWrapper := func(args ...interface{}) {
		storage.Storage_listCallback(js.Null(), nil)
	}
	productCallbackWrapper := func(args ...interface{}) {
		product.Product_listCallback(js.Null(), nil)
	}

	if CurrentView != "storage" {
		LoadContent("storage", fmt.Sprintf("%sv/storages", ApplicationProxyPath), storageCallbackWrapper, nil)
	} else {
		LoadContent("product", fmt.Sprintf("%sv/products", ApplicationProxyPath), productCallbackWrapper, nil)
	}

	return nil

}
