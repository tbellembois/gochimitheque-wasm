package globals

import (
	"net/url"

	"github.com/tbellembois/gochimitheque-wasm/ajax"
	"github.com/tbellembois/gochimitheque-wasm/localstorage"
	"github.com/tbellembois/gochimitheque-wasm/types"
)

var (
	URLParameters                                  url.Values
	ApplicationProxyPath, HTTPHeaderAcceptLanguage string
	ConnectedUserEmail                             string
	ConnectedUserID                                int
	DisableCache                                   bool
	CurrentView                                    string
	BSTableQueryFilter                             ajax.SafeQueryFilter
	LocalStorage                                   localstorage.LocalStorage
	// CurrentProduct is the current viewed or edited product, or the product of
	// the listed storages.
	CurrentProduct types.Product
	// CurrentStorage is the current viewed or edited storage.
	CurrentStorage types.Storage
	// CurrentStorages are the newly created storages.
	CurrentStorages           []types.Storage
	DBPrecautionaryStatements []types.PrecautionaryStatement // for magic selector
	DBHazardStatements        []types.HazardStatement        // for magic selector

	// permissions
	PermItems = [3]string{
		"rproducts",
		"products",
		"storages"}
)
