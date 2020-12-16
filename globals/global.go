package globals

import (
	"net/url"
	"syscall/js"
)

var (
	Win js.Value
	Doc js.Value

	URLParameters                                  url.Values
	ApplicationProxyPath, HTTPHeaderAcceptLanguage string
	ConnectedUserEmail                             string
	ConnectedUserID                                int
	DisableCache                                   bool
	CurrentView                                    string
	MenuLoaded                                     bool
	SearchLoaded                                   bool
	UserLoaded                                     bool

	// permissions
	PermItems = [3]string{
		"rproducts",
		"products",
		"storages"}
)
