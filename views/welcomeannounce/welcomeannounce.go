package welcomeannounce

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/tbellembois/gochimitheque-wasm/ajax"
	. "github.com/tbellembois/gochimitheque-wasm/globals"
	"github.com/tbellembois/gochimitheque-wasm/jquery"
	"github.com/tbellembois/gochimitheque-wasm/jsutils"
	"github.com/tbellembois/gochimitheque-wasm/locales"
	. "github.com/tbellembois/gochimitheque-wasm/types"
)

func WelcomeAnnounce_listCallback(this js.Value, args []js.Value) interface{} {

	done := func(data js.Value) {

		var (
			welcomeAnnounce WelcomeAnnounce
			err             error
		)

		if err = json.Unmarshal([]byte(data.String()), &welcomeAnnounce); err != nil {
			fmt.Println(err)
		}
		jquery.Jq("#welcomeannounce_text").SetVal(welcomeAnnounce.WelcomeAnnounceText)

	}
	fail := func(data js.Value) {

		jsutils.DisplayGenericErrorMessage()

	}

	ajax.Ajax{
		Method: "get",
		URL:    fmt.Sprintf("%swelcomeannounce", ApplicationProxyPath),
		Done:   done,
		Fail:   fail,
	}.Send()

	jquery.Jq("#searchbar").Hide()
	jquery.Jq("#actions").Hide()

	return nil

}
func SaveWelcomeAnnounce(this js.Value, args []js.Value) interface{} {

	var (
		welcomeAnnounce WelcomeAnnounce
		dataBytes       []byte
		err             error
	)

	welcomeAnnounce.WelcomeAnnounceText = jquery.Jq("#welcomeannounce_text").GetVal().String()
	if dataBytes, err = json.Marshal(welcomeAnnounce); err != nil {
		fmt.Println(err)
		return nil
	}

	ajax.Ajax{
		URL:    fmt.Sprintf("%swelcomeannounce", ApplicationProxyPath),
		Method: "put",
		Data:   dataBytes,
		Done: func(data js.Value) {

			jsutils.DisplaySuccessMessage(locales.Translate("welcomeannounce_text_modificationsuccess", HTTPHeaderAcceptLanguage))

		},
		Fail: func(jqXHR js.Value) {

			jsutils.DisplayGenericErrorMessage()

		},
	}.Send()

	return nil

}
